package user_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/user_entity"
)

func (u *UserController) FindUserById(c *gin.Context) {
	userId := c.Param("id")

	if err := uuid.Validate(userId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "id",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	userData, err := u.userUseCase.FindUserById(context.Background(), userId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, userData)
}

func (u *UserController) FindUsers(c *gin.Context) {
	status := c.Query("status")
	name := c.Query("name")
	email := c.Query("email")

	userStatus, err := user_entity.GetUserStatusByName(status)
	if err != nil {
		errRest := rest_err.NewBadRequestError("Error trying to validate user status param")
		c.JSON(errRest.Code, errRest)
		return
	}

	auctions, err := u.userUseCase.FindUsers(context.Background(),
		userStatus, name, email)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}
