package email_value_source_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/validation"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/email_value_source_usecase"
)

func (u *EmailValueSourceController) UpdateEmailValueSource(c *gin.Context) {
	emailValueSourceId := c.Param("id")

	var emailValueSourceInputDTO email_value_source_usecase.UpdateEmailValueSourceInputDTO

	if err := c.ShouldBindJSON(&emailValueSourceInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.EmailValueSourceUseCase.UpdateEmailValueSource(context.Background(), emailValueSourceId, emailValueSourceInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
