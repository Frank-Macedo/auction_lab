package auction_entity

import (
	"context"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
	"time"

	"github.com/google/uuid"
)

func CreateAuction(ctx context.Context, productName, category, description string, condition ProductCondition) (*AuctionEntity, error) {
	auction := AuctionEntity{
		ID:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	err := auction.Validate()
	if err != nil {
		return nil, err
	}

	return &auction, nil
}

func (au *AuctionEntity) Validate() error {
	if len(au.ProductName) <= 1 ||
		len(au.Category) <= 2 ||
		len(au.Description) <= 10 && (au.Condition < New || au.Condition > Refurbished) {
		return internal_error.NewInternalServerError("invalid auction data")
	}
	return nil
}

type AuctionEntity struct {
	ID          string           `json:"id" bson:"_id,omitempty"`
	ProductName string           `json:"product_name" bson:"product_name"`
	Category    string           `json:"category" bson:"category"`
	Condition   ProductCondition `json:"condition" bson:"condition"`
	Status      AuctionStatus    `json:"status" bson:"status"`
	Timestamp   time.Time        `json:"timestamp" bson:"timestamp"`
	Description string           `json:"description" bson:"description"`
}

type ProductCondition int
type AuctionStatus int

const (
	Active AuctionStatus = iota
	Inactive
	Completed
)

const (
	New ProductCondition = iota
	Used
	Refurbished
)

type AuctionRepositoryInterface interface {
	CreateAuction(ctx context.Context, auction AuctionEntity) *internal_error.InternalError
	FindAuctionByID(ctx context.Context, id string) (*AuctionEntity, *internal_error.InternalError)
	FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]*AuctionEntity, *internal_error.InternalError)
}
