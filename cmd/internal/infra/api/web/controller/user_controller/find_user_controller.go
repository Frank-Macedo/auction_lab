package user_controller

import (
	"lab_fullcyle-auction_go/cmd/internal/rest_err"
	"lab_fullcyle-auction_go/cmd/internal/usecase/user_usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	usecase user_usecase.UserUseCaseInterface
}

func NewUserController(usecase user_usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		usecase: usecase,
	}
}
func (uc *UserController) FindUserByID(c *gin.Context) {

	ctx := c.Request.Context()
	userId := c.Param("id")

	if err := uuid.Validate(userId); err != nil {
		errRest := rest_err.NewBadRequestError("invalid user id", rest_err.Causes{
			Field:   "id",
			Message: "must be a valid UUID",
		})
		c.JSON(errRest.Code, errRest)

		return
	}

	userData, err := uc.usecase.FindUserByID(ctx, userId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)

		return
	}

	c.JSON(200, userData)

}
