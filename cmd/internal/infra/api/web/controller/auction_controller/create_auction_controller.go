package auction_controller

import (
	"lab_fullcyle-auction_go/cmd/internal/infra/api/web/validation"
	"lab_fullcyle-auction_go/cmd/internal/rest_err"
	"lab_fullcyle-auction_go/cmd/internal/usecase/auction_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuctionController struct {
	AuctionUseCase auction_usecase.AuctionUseCaseInterface
}

func NewAuctionController(auctionUseCase auction_usecase.AuctionUseCaseInterface) *AuctionController {
	return &AuctionController{
		AuctionUseCase: auctionUseCase,
	}
}

type AuctionControllerInterface interface {
	FindAuctionByID(c *gin.Context)
	FindAuctions(c *gin.Context)
	FindWinningBidByAuctionID(c *gin.Context)
	CreateAuction(c *gin.Context)
}

func (ac *AuctionController) CreateAuction(c *gin.Context) {
	var input auction_usecase.AuctionInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errRest := validation.ValidateErr(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	ctx := c.Request.Context()
	_, err := ac.AuctionUseCase.CreateAuction(ctx, input)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.Status(http.StatusCreated)
}
