package bid

import (
	"context"
	"fmt"
	"lab_fullcyle-auction_go/cmd/internal/entity/bid_entity"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
	"lab_fullcyle-auction_go/configuration/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (bd *BidRepository) FindBidByAuctionID(ctx context.Context, auctionID string) ([]bid_entity.Bid, *internal_error.InternalError) {
	var bidEntitiesMongo []BidEntityMongo
	filter := bson.M{"auction_id": auctionID}
	cursor, err := bd.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to find bids for auction with ID = %s", auctionID), err)
		return nil, internal_error.NewInternalServerError("failed to find bids for auction")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var bidEntityMongo BidEntityMongo
		if err := cursor.Decode(&bidEntityMongo); err != nil {
			logger.Error("Failed to decode bid entity from MongoDB", err)
			return nil, internal_error.NewInternalServerError("failed to decode bid entity")
		}
		bidEntitiesMongo = append(bidEntitiesMongo, bidEntityMongo)
	}

	var bidEntities []bid_entity.Bid
	for _, bidMongo := range bidEntitiesMongo {
		bidEntity := bid_entity.Bid{
			ID:        bidMongo.ID,
			AuctionID: bidMongo.AuctionID,
			UserID:    bidMongo.UserID,
			Amount:    bidMongo.Amount,
			Timestamp: time.Unix(bidMongo.Timestamp, 0), // Assuming you want to convert this back to time.Time if needed
		}
		bidEntities = append(bidEntities, bidEntity)
	}

	return bidEntities, nil
}

func (bd *BidRepository) FindWinningBidByAuctionID(ctx context.Context, auctionID string) (*bid_entity.Bid, *internal_error.InternalError) {
	var winningBidMongo BidEntityMongo
	filter := bson.M{"auction_id": auctionID}
	opts := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}}) // Sort by amount in descending order to get the highest bid
	err := bd.Collection.FindOne(ctx, filter, opts).Decode(&winningBidMongo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error(fmt.Sprintf("No bids found for auction with ID = %s", auctionID), err)
			return nil, internal_error.NewNotFoundError("no bids found for auction with id: " + auctionID)
		}
		logger.Error(fmt.Sprintf("Failed to find winning bid for auction with ID = %s", auctionID), err)
		return nil, internal_error.NewInternalServerError("failed to find winning bid for auction")
	}

	winningBid := &bid_entity.Bid{
		ID:        winningBidMongo.ID,
		AuctionID: winningBidMongo.AuctionID,
		UserID:    winningBidMongo.UserID,
		Amount:    winningBidMongo.Amount,
		Timestamp: time.Unix(winningBidMongo.Timestamp, 0), // Assuming you want to convert this back to time.Time if needed
	}

	return winningBid, nil
}
