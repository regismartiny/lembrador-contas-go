package bill_processing_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/validation"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/bill_processing_usecase"
)

type BillProcessingController struct {
	billProcessingUseCase bill_processing_usecase.BillProcessingUseCaseInterface
}

func NewBillProcessingController(billProcessingUseCase bill_processing_usecase.BillProcessingUseCaseInterface) *BillProcessingController {
	return &BillProcessingController{
		billProcessingUseCase: billProcessingUseCase,
	}
}

func (u *BillProcessingController) StartBillProcessing(c *gin.Context) {
	var billProcessingInputDTO bill_processing_usecase.BillProcessingInputDTO

	if err := c.ShouldBindJSON(&billProcessingInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	billProcessingOutput, err := u.billProcessingUseCase.StartBillProcessing(context.Background(), billProcessingInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusCreated, billProcessingOutput)
}

func (u *BillProcessingController) GetBillProcessingStatus(c *gin.Context) {
	billProcessingId := c.Param("id")

	status, err := u.billProcessingUseCase.GetBillProcessingStatus(context.Background(), billProcessingId)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, status)
}
