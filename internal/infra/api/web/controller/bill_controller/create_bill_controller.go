package bill_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/validation"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/bill_usecase"
)

type BillController struct {
	billUseCase bill_usecase.BillUseCaseInterface
}

func NewBillController(billUseCase bill_usecase.BillUseCaseInterface) *BillController {
	return &BillController{
		billUseCase: billUseCase,
	}
}

func (u *BillController) CreateBill(c *gin.Context) {
	var billInputDTO bill_usecase.BillInputDTO

	if err := c.ShouldBindJSON(&billInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.billUseCase.CreateBill(context.Background(), billInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
