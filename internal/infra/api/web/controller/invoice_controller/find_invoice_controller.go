package invoice_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/invoice_entity"
)

func (u *InvoiceController) FindInvoiceById(c *gin.Context) {
	invoiceId := c.Param("id")

	if err := uuid.Validate(invoiceId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "id",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	invoiceData, err := u.invoiceUseCase.FindInvoiceById(context.Background(), invoiceId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, invoiceData)
}

func (u *InvoiceController) FindInvoices(c *gin.Context) {
	billId := c.Query("billId")
	status := c.Query("status")

	invoiceStatus, err := invoice_entity.GetInvoiceStatusByName(status)
	if err != nil {
		errRest := rest_err.NewBadRequestError("Error trying to validate invoice status param")
		c.JSON(errRest.Code, errRest)
		return
	}

	auctions, err := u.invoiceUseCase.FindInvoices(context.Background(),
		billId, invoiceStatus)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}
