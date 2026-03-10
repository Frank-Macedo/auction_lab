package auction_usecase

import (
	"context"
	"lab_fullcyle-auction_go/cmd/internal/entity/auction_entity"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
	"lab_fullcyle-auction_go/cmd/internal/usecase/bid_usecase"
)

func (au *AuctionUseCase) FindAuctionByID(ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntity, err := au.AuctionRepositoryInterface.FindAuctionByID(ctx, id)
	if err != nil {
		return nil, internal_error.NewInternalServerError("failed to find auction by id: " + id)
	}

	return &AuctionOutputDTO{
		ID:          auctionEntity.ID,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Condition:   ProductCondition(auctionEntity.Condition),
		Status:      AuctionStatus(auctionEntity.Status),
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}, nil
}

func (au *AuctionUseCase) FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]*AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntities, err := au.AuctionRepositoryInterface.FindAuctions(ctx, auction_entity.AuctionStatus(status), category, productName)
	if err != nil {
		return nil, internal_error.NewInternalServerError("failed to find auctions")
	}

	var auctionOutputDTOs []*AuctionOutputDTO
	for _, auctionEntity := range auctionEntities {
		auctionOutputDTOs = append(auctionOutputDTOs, &AuctionOutputDTO{
			ID:          auctionEntity.ID,
			ProductName: auctionEntity.ProductName,
			Category:    auctionEntity.Category,
			Condition:   ProductCondition(auctionEntity.Condition),
			Status:      AuctionStatus(auctionEntity.Status),
			Timestamp:   auctionEntity.Timestamp.Unix(),
		})
	}

	return auctionOutputDTOs, nil
}

func (au *AuctionUseCase) FindWinningBidByAuctionID(ctx context.Context, id string) (*WinningInfoOutputDTO, *internal_error.InternalError) {
	auctionEntity, err := au.AuctionRepositoryInterface.FindAuctionByID(ctx, id)
	if err != nil {
		return nil, internal_error.NewInternalServerError("failed to find auction by id: " + id)
	}

	if auctionEntity == nil {
		return nil, nil
	}
	auctionOutputDTO := &AuctionOutputDTO{
		ID:          auctionEntity.ID,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Condition:   ProductCondition(auctionEntity.Condition),
		Status:      AuctionStatus(auctionEntity.Status),
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	bidEntity, err := au.BidRepositoryInterface.FindWinningBidByAuctionID(ctx, id)

	if err != nil {
		return nil, internal_error.NewInternalServerError("failed to find winning bid by auction id: " + id)
	}

	if bidEntity == nil {
		return &WinningInfoOutputDTO{
			Auction:    auctionOutputDTO,
			WinningBid: nil,
		}, nil
	}

	return &WinningInfoOutputDTO{
		Auction: auctionOutputDTO,
		WinningBid: &bid_usecase.BidOutputDTO{
			ID:        bidEntity.ID,
			AuctionID: bidEntity.AuctionID,
			UserID:    bidEntity.UserID,
			Amount:    bidEntity.Amount,
			Timestamp: bidEntity.Timestamp.Unix(),
		},
	}, nil

}

type WinningInfoOutputDTO struct {
	Auction    *AuctionOutputDTO         `json:"auction"`
	WinningBid *bid_usecase.BidOutputDTO `json:"winning_bid"`
}
