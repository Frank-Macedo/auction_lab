package bid_usecase

import (
	"context"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
)

func (bu *BidUseCase) FindBidByAuctionID(ctx context.Context, id string) ([]*BidOutputDTO, *internal_error.InternalError) {
	bidEntities, err := bu.BidRepository.FindBidByAuctionID(ctx, id)
	if err != nil {
		return nil, internal_error.NewInternalServerError("failed to find bid by auction id: " + id)
	}

	bidOutputDTOs := make([]*BidOutputDTO, len(bidEntities))
	for i, bidEntity := range bidEntities {
		bidOutputDTOs[i] = &BidOutputDTO{
			ID:        bidEntity.ID,
			AuctionID: bidEntity.AuctionID,
			UserID:    bidEntity.UserID,
			Amount:    bidEntity.Amount,
			Timestamp: bidEntity.Timestamp.Unix(),
		}
	}

	return bidOutputDTOs, nil
}

func (bu *BidUseCase) FindWinningBidByAuctionID(ctx context.Context, id string) (*BidOutputDTO, *internal_error.InternalError) {
	bidEntity, err := bu.BidRepository.FindWinningBidByAuctionID(ctx, id)
	if err != nil {
		return nil, internal_error.NewInternalServerError("failed to find winning bid by auction id: " + id)
	}

	if bidEntity == nil {
		return nil, nil
	}

	bidOutputDTO := &BidOutputDTO{
		ID:        bidEntity.ID,
		AuctionID: bidEntity.AuctionID,
		UserID:    bidEntity.UserID,
		Amount:    bidEntity.Amount,
		Timestamp: bidEntity.Timestamp.Unix(),
	}

	return bidOutputDTO, nil
}
