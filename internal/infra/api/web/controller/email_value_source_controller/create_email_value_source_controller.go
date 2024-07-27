package email_value_source_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/validation"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/email_value_source_usecase"
)

type EmailValueSourceController struct {
	EmailValueSourceUseCase email_value_source_usecase.EmailValueSourceUseCaseInterface
}

func NewEmailValueSourceController(EmailValueSourceUseCase email_value_source_usecase.EmailValueSourceUseCaseInterface) *EmailValueSourceController {
	return &EmailValueSourceController{
		EmailValueSourceUseCase: EmailValueSourceUseCase,
	}
}

func (u *EmailValueSourceController) CreateEmailValueSource(c *gin.Context) {
	var emailValueSourceInputDTO email_value_source_usecase.EmailValueSourceInputDTO

	if err := c.ShouldBindJSON(&emailValueSourceInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.EmailValueSourceUseCase.CreateEmailValueSource(context.Background(), emailValueSourceInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
