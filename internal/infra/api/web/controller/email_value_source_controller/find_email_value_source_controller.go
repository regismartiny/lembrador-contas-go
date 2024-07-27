package email_value_source_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
)

func (u *EmailValueSourceController) FindEmailValueSourceById(c *gin.Context) {
	emailValueSourceId := c.Param("id")

	if err := uuid.Validate(emailValueSourceId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "id",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	emailValueSourceData, err := u.EmailValueSourceUseCase.FindEmailValueSourceById(context.Background(), emailValueSourceId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, emailValueSourceData)
}

func (u *EmailValueSourceController) FindEmailValueSources(c *gin.Context) {
	address := c.Query("address")
	subject := c.Query("subject")

	auctions, err := u.EmailValueSourceUseCase.FindEmailValueSources(context.Background(), address, subject)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}
