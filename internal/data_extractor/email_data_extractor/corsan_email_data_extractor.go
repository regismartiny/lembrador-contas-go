package email_data_extractor

import (
	"cmp"
	"fmt"
	"log"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/email_service"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type CorsanEmailDataExtractor struct {
	emailService email_service.EmailServiceInterface
}

type CorsanEmailData struct {
	CodigoImovel  string
	Vencimento    time.Time
	MesReferencia string
	Valor         float64
}

type MonthOfTheYear uint8

const (
	JANEIRO MonthOfTheYear = iota + 1
	FEVEREIRO
	MARCO
	ABRIL
	MAIO
	JUNHO
	JULHO
	AGOSTO
	SETEMBRO
	OUTUBRO
	NOVEMBRO
	DEZEMBRO
)

func (m MonthOfTheYear) Name() string {
	return monthsOfTheYearNames[m]
}

var monthsOfTheYearNames = []string{
	"",
	"JANEIRO",
	"FEVEREIRO",
	"MARÇO",
	"ABRIL",
	"MAIO",
	"JUNHO",
	"JULHO",
	"AGOSTO",
	"SETEMBRO",
	"OUTUBRO",
	"NOVEMBRO",
	"DEZEMBRO",
}

func GetMonthOfTheYearByName(name string) (MonthOfTheYear, *internal_error.InternalError) {
	for k, v := range monthsOfTheYearNames {
		if v == name {
			return MonthOfTheYear(k), nil
		}
	}

	return MonthOfTheYear(0), internal_error.NewBadRequestError("invalid monthOfTheYear name")
}

func NewCorsanEmailDataExtractor(emailService email_service.EmailServiceInterface) *CorsanEmailDataExtractor {
	return &CorsanEmailDataExtractor{
		emailService: emailService,
	}
}

func (x *CorsanEmailDataExtractor) Extract(request EmailDataExtractorRequest) (*EmailDataExtractorResponse, *internal_error.InternalError) {

	log.Println("Extracting email data. Request:", request)

	startDate := request.StartDate.Format("2006/01/02")
	endDate := request.EndDate.Format("2006/01/02")

	messages, err := x.emailService.FindMessages(request.Subject, request.Address, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, internal_error.NewNotFoundError("No messages found")
	}

	slices.SortStableFunc(messages, func(a, b *email_service.EmailServiceMessage) int {
		return cmp.Compare(a.Id, b.Id)
	})

	log.Println("Ordered Messages", messages)

	lastMessage := messages[0]

	lastMessage, err = x.emailService.GetMessage(lastMessage.Id)
	if err != nil {
		return nil, err
	}

	parsedData, err := x.parse(lastMessage)

	log.Println("Parsed data", parsedData)

	if err != nil {
		return nil, err
	}

	return &EmailDataExtractorResponse{
		Amount: parsedData.Valor,
	}, nil
}

func (x *CorsanEmailDataExtractor) parse(msg *email_service.EmailServiceMessage) (*CorsanEmailData, *internal_error.InternalError) {
	log.Println("Parsing Corsan message body")

	snippet := msg.Snippet
	if snippet == "" {
		return nil, internal_error.NewInternalServerError("No data found in message snippet")
	}

	STR_CODIGO_IMOVEL := "Código do Imóvel:"
	indexCodigoImovel := strings.Index(snippet, STR_CODIGO_IMOVEL)
	STR_VENCIMENTO := "Vencimento:"
	indexVencimento := strings.Index(snippet, STR_VENCIMENTO)

	valorCodigoImovel := snippet[indexCodigoImovel+len(STR_CODIGO_IMOVEL)+1 : indexVencimento-1]

	STR_VALOR := "Valor:"
	indexValorAPAgar := strings.Index(snippet, STR_VALOR)

	valorVencimento := snippet[indexVencimento+len(STR_VENCIMENTO)+1 : indexValorAPAgar-1]
	vencimento, err := time.Parse("2/1/2006", valorVencimento)
	if err != nil {
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error parsing date %s", valorVencimento))
	}

	STR_MES_REFERENCIA := "referente ao mês de"
	indexMesReferencia := strings.Index(snippet, STR_MES_REFERENCIA)
	indexValorMesReferencia := indexMesReferencia + len(STR_MES_REFERENCIA) + 1
	indexFimMesReferencia := indexValorMesReferencia + strings.Index(snippet[indexValorMesReferencia:], " ")
	valorMesReferencia := snippet[indexValorMesReferencia:indexFimMesReferencia]
	mesReferencia := getReferenceMonth(valorMesReferencia, vencimento)

	STR_AGRADECEMOS := "Agradecemos"
	indexAgradecemos := strings.Index(snippet, STR_AGRADECEMOS)
	valorValorAPagarStr := snippet[indexValorAPAgar+len(STR_VALOR)+1 : indexAgradecemos-1]
	valorValorAPagarStrNumber := valorValorAPagarStr[0 : len(valorValorAPagarStr)-1]
	valorValorAPagarStrNumber = strings.ReplaceAll(valorValorAPagarStrNumber, ",", ".")
	valorValorAPagar, err := strconv.ParseFloat(valorValorAPagarStrNumber, 64)
	if err != nil {
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error parsing value %s", valorValorAPagarStrNumber))
	}
	valorValorAPagar = math.Ceil(valorValorAPagar*100) / 100

	return &CorsanEmailData{
		CodigoImovel:  valorCodigoImovel,
		Vencimento:    vencimento,
		MesReferencia: mesReferencia,
		Valor:         valorValorAPagar,
	}, nil
}

func getReferenceMonth(referenceMonthName string, dueDate time.Time) string {

	monthOfTheYear, err := GetMonthOfTheYearByName(referenceMonthName)
	if err != nil {
		return ""
	}

	for i := monthOfTheYear; i > 1; i-- {
		dueDate = dueDate.AddDate(0, -1, 0)

		if uint8(dueDate.Month()) == uint8(i) {
			return fmt.Sprintf("%02d/%04d", int(i), dueDate.Year())
		}
	}

	return ""
}
