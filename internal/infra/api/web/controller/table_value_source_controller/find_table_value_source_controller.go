package table_value_source_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/table_value_source_entity"
)

func (u *TableValueSourceController) FindTableValueSourceById(c *gin.Context) {
	tableValueSourceId := c.Param("id")

	if err := uuid.Validate(tableValueSourceId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "id",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	tableValueSourceData, err := u.TableValueSourceUseCase.FindTableValueSourceById(context.Background(), tableValueSourceId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, tableValueSourceData)
}

func (u *TableValueSourceController) FindTableValueSources(c *gin.Context) {
	status := c.Query("status")
	name := c.Query("name")

	tableValueSourceStatus, err := table_value_source_entity.GetTableValueSourceStatusByName(status)
	if err != nil {
		errRest := rest_err.NewBadRequestError("Error trying to validate tableValueSource status param")
		c.JSON(errRest.Code, errRest)
		return
	}

	auctions, err := u.TableValueSourceUseCase.FindTableValueSources(context.Background(),
		tableValueSourceStatus, name)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}
