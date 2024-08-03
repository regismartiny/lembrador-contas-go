package bill_processing_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/regismartiny/lembrador-contas-go/configuration/rest_err"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_processing_entity"
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

func (u *BillProcessingController) FindBillProcessings(c *gin.Context) {
	status := c.Query("status")

	billStatus, err := bill_processing_entity.GetBillProcessingStatusByName(status)
	if err != nil {
		errRest := rest_err.NewBadRequestError("Error trying to validate billProcessing status param")
		c.JSON(errRest.Code, errRest)
		return
	}

	billProcessings, err := u.billProcessingUseCase.FindBillProcessings(context.Background(), billStatus)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, billProcessings)
}
