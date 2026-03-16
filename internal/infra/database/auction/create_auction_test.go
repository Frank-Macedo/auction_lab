package auction

import (
	"context"
	"testing"
	"time"

	"github.com/Frank-Macedo/lab_auction/internal/entity/auction_entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewAuctionRepository(t *testing.T) {
	t.Run("should create repository with default end time", func(t *testing.T) {
		t.Setenv("AUCTION_END_TIME", "")
		client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
		db := client.Database("test")

		repo := NewAuctionRepository(db)

		assert.NotNil(t, repo)
		assert.Equal(t, "10s", repo.EndTime)
		assert.NotNil(t, repo.Timers)
	})

	t.Run("should create repository with custom end time", func(t *testing.T) {
		t.Setenv("AUCTION_END_TIME", "120s")
		client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
		db := client.Database("test")

		repo := NewAuctionRepository(db)

		assert.Equal(t, "120s", repo.EndTime)
	})
}

func TestCreateAuction(t *testing.T) {
	t.Run("should create auction successfully", func(t *testing.T) {
		client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
		db := client.Database("test")
		repo := NewAuctionRepository(db)

		auction := &auction_entity.Auction{
			Id:          "test-id-99",
			ProductName: "Test Product",
			Category:    "Electronics",
			Description: "Test Description",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := repo.CreateAuction(context.Background(), auction)

		assert.Nil(t, err)
		db.Collection("auctions").DeleteOne(context.Background(), map[string]interface{}{"_id": "test-id-99"})

	})
}

func TestCreateAuctionAndWhaitToComplete(t *testing.T) {

	t.Run("should create auction successfully", func(t *testing.T) {
		client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
		db := client.Database("test")
		repo := NewAuctionRepository(db)

		auction := &auction_entity.Auction{
			Id:          "test-auction-1",
			ProductName: "Test Product",
			Category:    "Electronics",
			Description: "Test Description",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := repo.CreateAuction(context.Background(), auction)

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

		assert.Nil(t, err)
	})
}

func TestMonitorAuction(t *testing.T) {
	t.Run("should create and store timer", func(t *testing.T) {
		client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
		db := client.Database("test")
		repo := NewAuctionRepository(db)

		repo.monitorAuction("test-auction", "1s")
		time.Sleep(100 * time.Millisecond)

		repo.mu.Lock()
		_, exists := repo.Timers["test-auction"]
		repo.mu.Unlock()

		assert.True(t, exists)
	})
}

func TestRemoveTimer(t *testing.T) {
	t.Run("should remove timer from map", func(t *testing.T) {
		client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
		db := client.Database("test")
		repo := NewAuctionRepository(db)

		repo.monitorAuction("test-auction", "10s")
		time.Sleep(50 * time.Millisecond)

		repo.removeTimer("test-auction")

		repo.mu.Lock()
		_, exists := repo.Timers["test-auction"]
		repo.mu.Unlock()

		assert.False(t, exists)
	})
}
