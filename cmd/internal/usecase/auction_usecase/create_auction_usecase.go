package auction_usecase

import (
	"context"
	"lab_fullcyle-auction_go/cmd/internal/entity/auction_entity"
	"lab_fullcyle-auction_go/cmd/internal/entity/bid_entity"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
)

func NewAuctionUseCase(auctionRepo auction_entity.AuctionRepositoryInterface, bidRepo bid_entity.BidRepositoryInterface) *AuctionUseCase {
	buildUseCase := &AuctionUseCase{
		AuctionRepositoryInterface: auctionRepo,
		BidRepositoryInterface:     bidRepo,
	}

	return buildUseCase
}

type AuctionOutputDTO struct {
	ID          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   int64            `json:"timestamp"`
}

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=1"`
	Category    string           `json:"category" binding:"required,min=2"`
	Description string           `json:"description" binding:"required,min=10,max=200"`
	Condition   ProductCondition `json:"condition"`
}

type ProductCondition int64
type AuctionStatus int64

const (
	New ProductCondition = iota
	Used
)

const (
	Active AuctionStatus = iota
	Completed
)

type AuctionUseCaseInterface interface {
	CreateAuction(ctx context.Context, input AuctionInputDTO) (*AuctionOutputDTO, *internal_error.InternalError)
	FindAuctionByID(ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError)
	FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]*AuctionOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionID(ctx context.Context, id string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

type AuctionUseCase struct {
	auction_entity.AuctionRepositoryInterface
	bid_entity.BidRepositoryInterface
}

func (au *AuctionUseCase) CreateAuction(ctx context.Context, input AuctionInputDTO) (*AuctionOutputDTO, *internal_error.InternalError) {
	auction, err := auction_entity.CreateAuction(
		ctx,
		input.ProductName,
		input.Category, "Default description",
		auction_entity.ProductCondition(Active),
	)
	if err != nil {
		return nil, internal_error.NewInternalServerError("failed to create auction")
	}

	if err := au.AuctionRepositoryInterface.CreateAuction(ctx, *auction); err != nil {
		return nil, internal_error.NewInternalServerError("failed to create auction")
	}
	return &AuctionOutputDTO{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Condition:   ProductCondition(auction.Condition),
		Status:      AuctionStatus(auction.Status),
		Timestamp:   auction.Timestamp.Unix(),
	}, nil
}
