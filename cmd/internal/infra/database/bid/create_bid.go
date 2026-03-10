package bid

import (
	"context"
	"fmt"
	"lab_fullcyle-auction_go/cmd/internal/entity/auction_entity"
	"lab_fullcyle-auction_go/cmd/internal/entity/bid_entity"
	"lab_fullcyle-auction_go/cmd/internal/infra/database/auction"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
	"lab_fullcyle-auction_go/configuration/logger"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

type BidEntityMongo struct {
	ID        string  `json:"id" bson:"_id,omitempty"`
	AuctionID string  `json:"auction_id" bson:"auction_id"`
	UserID    string  `json:"user_id" bson:"user_id"`
	Amount    float64 `json:"amount" bson:"amount"`
	Timestamp int64   `json:"timestamp" bson:"timestamp"`
}

type BidRepository struct {
	Collection        *mongo.Collection
	AuctionRepository *auction.AuctionRepository
}

func NewBidRepository(database *mongo.Database, auctionRepo *auction.AuctionRepository) *BidRepository {
	return &BidRepository{
		Collection:        database.Collection("bids"),
		AuctionRepository: auctionRepo,
	}
}

func (bd *BidRepository) CreateBid(ctx context.Context, bidEntities []bid_entity.Bid) *internal_error.InternalError {
	var wg sync.WaitGroup

	auctionCloseTimes := make(map[string]auction_entity.AuctionStatus)
	for _, bid := range bidEntities {
		if _, exists := auctionCloseTimes[bid.AuctionID]; !exists {
			auctionEntity, err := bd.AuctionRepository.FindAuctionByID(ctx, bid.AuctionID)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to find auction with ID = %s for bid creation", bid.AuctionID), err)
				continue
			}
			auctionCloseTimes[bid.AuctionID] = auctionEntity.Status
		}
	}

	for _, bid := range bidEntities {
		wg.Add(1)
		go func(ctx context.Context, bidValue bid_entity.Bid) {
			defer wg.Done()
			auctionStatus, exists := auctionCloseTimes[bidValue.AuctionID]
			if !exists {
				logger.Error(fmt.Sprintf("Auction with ID = %s not found for bid creation", bidValue.AuctionID), nil)
				return
			}
			if auctionStatus != auction_entity.Active {
				logger.Error(fmt.Sprintf("Auction with ID = %s is not active for bid creation", bidValue.AuctionID), nil)
				return
			}

			bidEntityMongo := &BidEntityMongo{
				ID:        bidValue.ID,
				AuctionID: bidValue.AuctionID,
				UserID:    bidValue.UserID,
				Amount:    bidValue.Amount,
				Timestamp: bidValue.Timestamp.Unix(),
			}

			if _, err := bd.Collection.InsertOne(ctx, bidEntityMongo); err != nil {
				logger.Error(fmt.Sprintf("Failed to create bid with ID = %s", bidValue.ID), err)
				return
			}

			logger.Info(fmt.Sprintf("Bid with ID = %s created successfully", bidValue.ID))

		}(ctx, bid)

	}

	wg.Wait()

	return nil
}
