package bid_controller

import (
	"lab_fullcyle-auction_go/cmd/internal/rest_err"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (bc *BidController) FindBidsByAuctionID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	bids, err := bc.BidUseCase.FindBidByAuctionID(ctx, id)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	if bids == nil || len(bids) == 0 {
		c.JSON(http.StatusNotFound, rest_err.NewnotfoundError("bids not found for this auction"))
		return
	}

	c.JSON(http.StatusOK, bids)
}

func (bc *BidController) FindWinningBidByAuctionID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	winningBid, err := bc.BidUseCase.FindWinningBidByAuctionID(ctx, id)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	if winningBid == nil {
		c.JSON(http.StatusNotFound, rest_err.NewnotfoundError("winning bid not found for this auction"))
		return
	}

	c.JSON(http.StatusOK, winningBid)
}
