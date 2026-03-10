package bid_usecase

import (
	"context"
	"lab_fullcyle-auction_go/cmd/internal/entity/bid_entity"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
	"lab_fullcyle-auction_go/configuration/logger"
	"os"
	"strconv"
	"time"
)

type BidOutputDTO struct {
	ID        string  `json:"id"`
	AuctionID string  `json:"auction_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
}

type BidUseCase struct {
	BidRepository bid_entity.BidRepositoryInterface

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bid_entity.Bid
}

type BidInputDTO struct {
	AuctionID string
	UserID    string
	Amount    float64
}

type BidUseCaseInterface interface {
	CreateBid(ctx context.Context, input BidInputDTO) (*BidOutputDTO, *internal_error.InternalError)
	FindBidByAuctionID(ctx context.Context, id string) ([]*BidOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionID(ctx context.Context, id string) (*BidOutputDTO, *internal_error.InternalError)
}

func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}
	return value
}

func getMaxBatchInsertInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Second
	}
	return duration
}

func NewBidUseCase(bidRepository bid_entity.BidRepositoryInterface) *BidUseCase {

	maxSizeInterval := getMaxBatchInsertInterval()
	maxBatchSize := getMaxBatchSize()

	buildUseCase := &BidUseCase{
		BidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxSizeInterval,
		timer:               time.NewTimer(maxSizeInterval),
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
	}
	buildUseCase.triggerCreateRoutine(context.Background())
	return buildUseCase
}

func (bu *BidUseCase) CreateBid(ctx context.Context, input BidInputDTO) (*BidOutputDTO, *internal_error.InternalError) {

	bidEntity, err := bid_entity.CreateBid(ctx, input.AuctionID, input.UserID, input.Amount)
	if err != nil {
		return nil, internal_error.NewInternalServerError("failed to create bid")
	}

	bu.bidChannel <- *bidEntity

	return &BidOutputDTO{
		ID:        bidEntity.ID,
		AuctionID: bidEntity.AuctionID,
		UserID:    bidEntity.UserID,
		Amount:    bidEntity.Amount,
		Timestamp: bidEntity.Timestamp.Unix(),
	}, nil
}

func (bu *BidUseCase) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)

		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
							logger.Error("Failed to create bid batch", err)
						}
					}
					return
				}

				bidBatch = append(bidBatch, bidEntity)

				if len(bidBatch) >= bu.maxBatchSize {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("Failed to create bid batch", err)
					}
					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			case <-bu.timer.C:
				if len(bidBatch) > 0 {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("Failed to create bid batch", err)
					}
					bidBatch = nil
				}
				bu.timer.Reset(bu.batchInsertInterval)

			}
		}

	}()
}

var bidBatch []bid_entity.Bid
