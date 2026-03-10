package auction

import (
	"context"
	"lab_fullcyle-auction_go/cmd/internal/entity/auction_entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	return client.Database("auction_test")
}

func TestCreateAuctionAndComplete(t *testing.T) {
	db := setupTestDB(t)
	repo := &AuctionRepository{
		Collection: db.Collection("auctions"),
		Timers:     make(map[string]*time.Timer),
		EndTime:    5,
	}

	auction := auction_entity.AuctionEntity{
		ID:          "test-auction-1",
		ProductName: "Test Product",
		Category:    "Electronics",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}
	err := repo.CreateAuction(context.Background(), auction)
	assert.Nil(t, err)

	completed := false
	for completed == false {
		time.Sleep(1 * time.Second)

		var result AuctionEntityMongo
		err := repo.Collection.FindOne(context.Background(), map[string]interface{}{"_id": "test-auction-1"}).Decode(&result)
		require.NoError(t, err)

		if result.Status == auction_entity.Completed {
			completed = true
		}

	}

	assert.Nil(t, err)

	db.Collection("auctions").DeleteOne(context.Background(), map[string]interface{}{"_id": "test-auction-1"})
}

func TestCreateAuction(t *testing.T) {
	db := setupTestDB(t)
	repo := &AuctionRepository{
		Collection: db.Collection("auctions"),
		Timers:     make(map[string]*time.Timer),
		EndTime:    1,
	}

	auction := auction_entity.AuctionEntity{
		ID:          "test-auction-1",
		ProductName: "Test Product",
		Category:    "Electronics",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	err := repo.CreateAuction(context.Background(), auction)
	assert.Nil(t, err)

	db.Collection("auctions").DeleteOne(context.Background(), map[string]interface{}{"_id": "test-auction-1"})
}

func TestMonitorAuction(t *testing.T) {
	db := setupTestDB(t)
	repo := &AuctionRepository{
		Collection: db.Collection("auctions"),
		Timers:     make(map[string]*time.Timer),
		EndTime:    1,
	}

	repo.monitorAuction("test-auction-2", 1)
	time.Sleep(2 * time.Second)

	repo.mu.Lock()
	_, exists := repo.Timers["test-auction-2"]
	repo.mu.Unlock()

	assert.False(t, exists)
}

func TestRemoveTimer(t *testing.T) {
	repo := &AuctionRepository{
		Timers: make(map[string]*time.Timer),
	}

	timer := time.NewTimer(10 * time.Second)
	repo.Timers["test-auction-3"] = timer

	repo.removeTimer("test-auction-3")

	repo.mu.Lock()
	_, exists := repo.Timers["test-auction-3"]
	repo.mu.Unlock()

	assert.False(t, exists)
}
