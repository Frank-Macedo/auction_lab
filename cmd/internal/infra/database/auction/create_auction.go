package auction

import (
	"context"
	"lab_fullcyle-auction_go/cmd/internal/entity/auction_entity"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
	"lab_fullcyle-auction_go/configuration/logger"
	"os"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	ID          string                          `json:"id" bson:"_id,omitempty"`
	ProductName string                          `json:"product_name" bson:"product_name"`
	Category    string                          `json:"category" bson:"category"`
	Condition   auction_entity.ProductCondition `json:"condition" bson:"condition"`
	Status      auction_entity.AuctionStatus    `json:"status" bson:"status"`
	Timestamp   int64                           `json:"timestamp" bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
	Timers     map[string]*time.Timer
	mu         sync.Mutex
	EndTime    int
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	endtimeEnv, err := strconv.Atoi(os.Getenv("AUCTION_END_TIME"))
	if err != nil {
		logger.Error("Failed to parse AUCTION_END_TIME environment variable", err)
		return nil
	}

	if err != nil {
		logger.Error("Failed to parse end time", err)
		return nil
	}
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
		Timers:     make(map[string]*time.Timer),
		EndTime:    endtimeEnv,
	}
}

func (au *AuctionRepository) CreateAuction(ctx context.Context, auction auction_entity.AuctionEntity) *internal_error.InternalError {

	auctionEntityMongo := &AuctionEntityMongo{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   auction.Timestamp.Unix(),
	}

	_, err := au.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Failed to create auction", err)
		return internal_error.NewInternalServerError("failed to create auction")
	}
	au.monitorAuction(auction.ID, au.EndTime)

	return nil
}

func (au *AuctionRepository) monitorAuction(auctionID string, endTime int) {

	ctx := context.Background()
	duration := time.Duration(endTime) * time.Second
	if duration <= 0 {
		duration = time.Second
	}

	timer := time.NewTimer(duration)

	au.mu.Lock()
	au.Timers[auctionID] = timer
	au.mu.Unlock()

	go func() {

		select {

		case <-timer.C:

			_, err := au.Collection.UpdateOne(
				ctx,
				bson.M{"_id": auctionID},
				bson.M{"$set": bson.M{"status": auction_entity.Completed}},
			)

			if err != nil {
				logger.Error("Failed to update auction status", err)
				return
			}

			logger.Info("Auction finished: " + auctionID)

			au.removeTimer(auctionID)

		case <-ctx.Done():
			return
		}

	}()
}

func (au *AuctionRepository) removeTimer(auctionID string) {

	au.mu.Lock()
	defer au.mu.Unlock()

	if timer, exists := au.Timers[auctionID]; exists {
		timer.Stop()
		delete(au.Timers, auctionID)
	}
}
