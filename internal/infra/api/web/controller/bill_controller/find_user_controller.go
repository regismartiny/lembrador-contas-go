package bill_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_entity"
)

func (u *BillController) FindBillById(c *gin.Context) {
	billId := c.Param("billId")

	if err := uuid.Validate(billId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "billId",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	billData, err := u.billUseCase.FindBillById(context.Background(), billId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, billData)
}

func (u *BillController) FindBills(c *gin.Context) {
	status := c.Query("status")
	name := c.Query("name")
	company := c.Query("company")

	billStatus, err := bill_entity.GetBillStatusByName(status)
	if err != nil {
		errRest := rest_err.NewBadRequestError("Error trying to validate bill status param")
		c.JSON(errRest.Code, errRest)
		return
	}

	auctions, err := u.billUseCase.FindBills(context.Background(),
		billStatus, name, company)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}
