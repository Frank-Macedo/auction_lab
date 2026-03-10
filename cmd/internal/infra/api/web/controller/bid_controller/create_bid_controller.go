package bid_controller

import (
	"lab_fullcyle-auction_go/cmd/internal/infra/api/web/validation"
	"lab_fullcyle-auction_go/cmd/internal/rest_err"
	"lab_fullcyle-auction_go/cmd/internal/usecase/bid_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BidController struct {
	BidUseCase bid_usecase.BidUseCaseInterface
}

func NewBidController(bidUseCase bid_usecase.BidUseCaseInterface) *BidController {
	return &BidController{
		BidUseCase: bidUseCase,
	}
}

type BidControllerInterface interface {
	FindBidsByAuctionID(c *gin.Context)
}

func (bc *BidController) CreateBid(c *gin.Context) {
	var input bid_usecase.BidInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errRest := validation.ValidateErr(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	ctx := c.Request.Context()
	_, err := bc.BidUseCase.CreateBid(ctx, input)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.Status(http.StatusCreated)
}
