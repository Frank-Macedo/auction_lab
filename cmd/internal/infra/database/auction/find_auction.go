package auction

import (
	"context"
	"errors"
	"fmt"
	"lab_fullcyle-auction_go/cmd/internal/entity/auction_entity"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
	"lab_fullcyle-auction_go/configuration/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ar *AuctionRepository) FindAuctionByID(ctx context.Context, auctionID string) (*auction_entity.AuctionEntity, *internal_error.InternalError) {
	var auctionEntityMongo AuctionEntityMongo
	filter := bson.M{"_id": auctionID}
	err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("Auction not found with this ID = %s", auctionID), err)
			return nil, internal_error.NewNotFoundError("auction not found with id: " + auctionID)
		}
		logger.Error(fmt.Sprintf("Error while finding auction with ID = %s", auctionID), err)
		return nil, internal_error.NewInternalServerError("failed to find auction")
	}

	auctionEntity := &auction_entity.AuctionEntity{
		ID:          auctionEntityMongo.ID,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
		// Assuming you want to convert this back to time.Time if needed
	}
	return auctionEntity, nil
}

func (ar *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category, productName string) ([]*auction_entity.AuctionEntity, *internal_error.InternalError) {
	filter := bson.M{}
	if status != 0 {
		filter["status"] = status
	}
	if category != "" {
		filter["category"] = category
	}
	if productName != "" {
		filter["product_name"] = bson.M{"$regex": productName, "$options": "i"}
	}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error while finding auctions", err)
		return nil, internal_error.NewInternalServerError("failed to find auctions")
	}
	defer cursor.Close(ctx)

	var auctionsMongo []*AuctionEntityMongo
	if err := cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error while decoding auctions", err)
		return nil, internal_error.NewInternalServerError("failed to decode auctions")
	}

	auctions := make([]*auction_entity.AuctionEntity, len(auctionsMongo))
	for i, auctionMongo := range auctionsMongo {
		auctions[i] = &auction_entity.AuctionEntity{
			ID:          auctionMongo.ID,
			ProductName: auctionMongo.ProductName,
			Category:    auctionMongo.Category,
			Condition:   auctionMongo.Condition,
			Status:      auctionMongo.Status,
			Timestamp:   time.Unix(auctionMongo.Timestamp, 0),
		}
	}

	return auctions, nil

}
