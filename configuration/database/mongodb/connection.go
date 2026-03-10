package mongodb

import (
	"context"
	"fmt"
	"lab_fullcyle-auction_go/configuration/logger"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_URL = "MONGO_URL"
	MONGODB_DB  = "MONGO_DB"
)

func NewMongoDBConnection(ctx context.Context) (*mongo.Database, error) {
	mongoDatabase := os.Getenv(MONGODB_DB)

	uri := strings.TrimSpace(os.Getenv(MONGODB_URL))

	fmt.Println("URL:", os.Getenv(MONGODB_URL))
	fmt.Println("DB:", os.Getenv(MONGODB_DB))

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
	)

	if err != nil {
		logger.Error("Failed to connect to MongoDB", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("Failed to ping MongoDB", err)
		return nil, err
	}
	return client.Database(mongoDatabase), nil
}
