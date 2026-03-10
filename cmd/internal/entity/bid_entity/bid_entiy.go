package bid_entity

import (
	"context"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	AuctionID string    `json:"auction_id" bson:"auction_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Amount    float64   `json:"amount" bson:"amount"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

func CreateBid(ctx context.Context, auctionID, userID string, amount float64) (*Bid, *internal_error.InternalError) {
	bid := &Bid{
		ID:        uuid.New().String(),
		AuctionID: auctionID,
		UserID:    userID,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	if err := bid.Validate(); err != nil {
		return nil, internal_error.NewInternalServerError("invalid bid data: " + err.Error())
	}
	return bid, nil
}

type BidRepositoryInterface interface {
	CreateBid(ctx context.Context, bids []Bid) *internal_error.InternalError
	FindBidByAuctionID(ctx context.Context, auctionID string) ([]Bid, *internal_error.InternalError)
	FindWinningBidByAuctionID(ctx context.Context, auctionID string) (*Bid, *internal_error.InternalError)
}

func (b *Bid) Validate() error {
	if err := uuid.Validate(b.ID); err != nil {
		return internal_error.NewInternalServerError("invalid bid ID")
	}
	if err := uuid.Validate(b.AuctionID); err != nil {
		return internal_error.NewInternalServerError("invalid auction ID")
	}
	if err := uuid.Validate(b.UserID); err != nil {
		return internal_error.NewInternalServerError("invalid user ID")
	}
	if b.Amount <= 0 {
		return internal_error.NewInternalServerError("invalid bid amount")
	}
	return nil
}
