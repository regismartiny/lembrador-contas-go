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

type CpflEmailDataExtractor struct {
	emailService email_service.EmailServiceInterface
}

type CpflEmailData struct {
	Instalacao    string
	Vencimento    time.Time
	MesReferencia string
	Valor         float64
}

func NewCpflEmailDataExtractor(emailService email_service.EmailServiceInterface) *CpflEmailDataExtractor {
	return &CpflEmailDataExtractor{
		emailService: emailService,
	}
}

func (x *CpflEmailDataExtractor) Extract(request EmailDataExtractorRequest) (*EmailDataExtractorResponse, *internal_error.InternalError) {

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

func (x *CpflEmailDataExtractor) parse(msg *email_service.EmailServiceMessage) (*CpflEmailData, *internal_error.InternalError) {
	log.Println("Parsing cpfl message body")

	snippet := msg.Snippet
	if snippet == "" {
		return nil, internal_error.NewInternalServerError("No data found in message snippet")
	}

	STR_NUMERO_INSTALACAO := "Número da instalação:"
	indexInstalacao := strings.Index(snippet, STR_NUMERO_INSTALACAO)
	STR_DATA_VENCIMENTO := "Data de vencimento:"
	indexDataVencimento := strings.Index(snippet, STR_DATA_VENCIMENTO)

	valorInstalacao := snippet[indexInstalacao+len(STR_NUMERO_INSTALACAO)+1 : indexDataVencimento-1]

	STR_VALOR_A_PAGAR := "Valor a pagar:"
	indexValorAPAgar := strings.Index(snippet, STR_VALOR_A_PAGAR)
	STR_MES_REFERENCIA := "Mês de referência:"
	indexMesReferencia := strings.Index(snippet, STR_MES_REFERENCIA)

	valorDataVencimento := snippet[indexDataVencimento+len(STR_DATA_VENCIMENTO)+1 : indexMesReferencia-1]
	vencimento, err := time.Parse("2/1/2006", valorDataVencimento)
	if err != nil {
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error parsing date %s", valorDataVencimento))
	}

	valorMesReferencia := snippet[indexMesReferencia+len(STR_MES_REFERENCIA)+1 : indexValorAPAgar-1]

	STR_PARA_ABRIR := "Para abrir"
	indexParaAbrir := strings.Index(snippet, STR_PARA_ABRIR)
	valorValorAPagarStr := snippet[indexValorAPAgar+len(STR_VALOR_A_PAGAR)+1 : indexParaAbrir-1]
	valorValorAPagarStrNumber := valorValorAPagarStr[3 : len(valorValorAPagarStr)-1]
	valorValorAPagarStrNumber = strings.ReplaceAll(valorValorAPagarStrNumber, ",", ".")
	valorValorAPagar, err := strconv.ParseFloat(valorValorAPagarStrNumber, 64)
	if err != nil {
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error parsing value %s", valorValorAPagarStrNumber))
	}
	valorValorAPagar = math.Ceil(valorValorAPagar*100) / 100

	return &CpflEmailData{
		Instalacao:    valorInstalacao,
		Vencimento:    vencimento,
		MesReferencia: valorMesReferencia,
		Valor:         valorValorAPagar,
	}, nil
}
