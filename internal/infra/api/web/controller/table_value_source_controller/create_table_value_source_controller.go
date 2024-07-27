package table_value_source_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/validation"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/table_value_source_usecase"
)

type TableValueSourceController struct {
	TableValueSourceUseCase table_value_source_usecase.TableValueSourceUseCaseInterface
}

func NewTableValueSourceController(TableValueSourceUseCase table_value_source_usecase.TableValueSourceUseCaseInterface) *TableValueSourceController {
	return &TableValueSourceController{
		TableValueSourceUseCase: TableValueSourceUseCase,
	}
}

func (u *TableValueSourceController) CreateTableValueSource(c *gin.Context) {
	var tableValueSourceInputDTO table_value_source_usecase.TableValueSourceInputDTO

	if err := c.ShouldBindJSON(&tableValueSourceInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.TableValueSourceUseCase.CreateTableValueSource(context.Background(), tableValueSourceInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}