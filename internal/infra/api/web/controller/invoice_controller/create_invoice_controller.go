package invoice_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/validation"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/invoice_usecase"
)

type InvoiceController struct {
	invoiceUseCase invoice_usecase.InvoiceUseCaseInterface
}

func NewInvoiceController(invoiceUseCase invoice_usecase.InvoiceUseCaseInterface) *InvoiceController {
	return &InvoiceController{
		invoiceUseCase: invoiceUseCase,
	}
}

func (u *InvoiceController) CreateInvoice(c *gin.Context) {
	var invoiceInputDTO invoice_usecase.InvoiceInputDTO

	if err := c.ShouldBindJSON(&invoiceInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.invoiceUseCase.CreateInvoice(context.Background(), invoiceInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
