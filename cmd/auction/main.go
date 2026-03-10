package main

import (
	"context"
	"fmt"
	"lab_fullcyle-auction_go/cmd/internal/infra/api/web/controller/auction_controller"
	"lab_fullcyle-auction_go/cmd/internal/infra/api/web/controller/bid_controller"
	"lab_fullcyle-auction_go/cmd/internal/infra/api/web/controller/user_controller"
	"lab_fullcyle-auction_go/cmd/internal/infra/database/auction"
	"lab_fullcyle-auction_go/cmd/internal/infra/database/bid"
	"lab_fullcyle-auction_go/cmd/internal/infra/database/user"
	"lab_fullcyle-auction_go/cmd/internal/usecase/auction_usecase"
	"lab_fullcyle-auction_go/cmd/internal/usecase/bid_usecase"
	"lab_fullcyle-auction_go/cmd/internal/usecase/user_usecase"
	"lab_fullcyle-auction_go/configuration/database/mongodb"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	fmt.Println(os.Getenv("MONGO_URL"))
	_, err := mongodb.NewMongoDBConnection(ctx)

	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
		return
	}

	userController, auctionController, bidController := initDependencies(databaseConnection)

	router := gin.Default()

	router.GET("/auctions", auctionController.FindAuctions)
	router.POST("/auctions", auctionController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionController.FindWinningBidByAuctionID)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidsByAuctionID)
	router.GET("/users/:id", userController.FindUserByID)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (userController *user_controller.UserController,
	auctionController *auction_controller.AuctionController,
	bidController *bid_controller.BidController) {
	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	return user_controller.NewUserController(user_usecase.NewUserUseCase(userRepository)),
		auction_controller.NewAuctionController(auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository)),
		bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

}
