package auction_controller

import (
	"lab_fullcyle-auction_go/cmd/internal/rest_err"
	"lab_fullcyle-auction_go/cmd/internal/usecase/auction_usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ac *AuctionController) FindAuctionByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	auctionData, err := ac.AuctionUseCase.FindAuctionByID(ctx, id)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	if auctionData == nil {
		c.JSON(http.StatusNotFound, rest_err.NewnotfoundError("auction not found"))
		return
	}

	c.JSON(http.StatusOK, auctionData)
}

func (ac *AuctionController) FindAuctions(c *gin.Context) {
	ctx := c.Request.Context()
	status := c.Query("status")

	statusInt, err := strconv.Atoi(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, rest_err.NewBadRequestError("invalid status"))
		return
	}
	category := c.Query("category")
	productName := c.Query("product_name")

	auctions, errRest := ac.AuctionUseCase.FindAuctions(ctx, auction_usecase.AuctionStatus(statusInt), category, productName)
	if errRest != nil {
		errRest := rest_err.ConvertError(errRest)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (ac *AuctionController) FindWinningBidByAuctionID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	winningInfo, err := ac.AuctionUseCase.FindWinningBidByAuctionID(ctx, id)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, winningInfo)
}
