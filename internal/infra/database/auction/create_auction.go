package auction

import (
	"context"

	"os"
	"sync"
	"time"

	"github.com/Frank-Macedo/lab_auction/configuration/logger"
	"github.com/Frank-Macedo/lab_auction/internal/entity/auction_entity"
	"github.com/Frank-Macedo/lab_auction/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
	Timers     map[string]*time.Timer
	mu         sync.Mutex
	EndTime    string
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {

	endtimeEnv := os.Getenv("AUCTION_END_TIME")

	if endtimeEnv == "" {
		endtimeEnv = "10s"
	}

	return &AuctionRepository{
		Collection: database.Collection("auctions"),
		Timers:     make(map[string]*time.Timer),
		EndTime:    endtimeEnv,
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	ar.monitorAuction(auctionEntityMongo.Id, ar.EndTime)

	return nil
}

func (au *AuctionRepository) monitorAuction(auctionID string, endTime string) {

	ctx := context.Background()
	duration, err := time.ParseDuration(endTime)
	if err != nil {
		logger.Error("Failed to parse end time", err)
		return
	}
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
