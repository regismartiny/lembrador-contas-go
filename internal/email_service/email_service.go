package email_service

import (
	"log"

	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
	"google.golang.org/api/gmail/v1"
)

type EmailServiceInterface interface {
	FindMessages(subject string, address string, startDate string, endDate string) ([]*EmailServiceMessage, *internal_error.InternalError)
	GetMessage(messageId string) (*EmailServiceMessage, *internal_error.InternalError)
}

type EmailServiceMessage struct {
	HistoryId       uint64                   `json:"historyId,omitempty,string"`
	Id              string                   `json:"id,omitempty"`
	InternalDate    int64                    `json:"internalDate,omitempty,string"`
	LabelIds        []string                 `json:"labelIds,omitempty"`
	Payload         *EmailServiceMessagePart `json:"payload,omitempty"`
	Raw             string                   `json:"raw,omitempty"`
	SizeEstimate    int64                    `json:"sizeEstimate,omitempty"`
	Snippet         string                   `json:"snippet,omitempty"`
	ThreadId        string                   `json:"threadId,omitempty"`
	ForceSendFields []string                 `json:"-"`
	NullFields      []string                 `json:"-"`
}

type EmailServiceMessagePart struct {
	Body            *MessagePartBody           `json:"body,omitempty"`
	Filename        string                     `json:"filename,omitempty"`
	Headers         []*MessagePartHeader       `json:"headers,omitempty"`
	MimeType        string                     `json:"mimeType,omitempty"`
	PartId          string                     `json:"partId,omitempty"`
	Parts           []*EmailServiceMessagePart `json:"parts,omitempty"`
	ForceSendFields []string                   `json:"-"`
	NullFields      []string                   `json:"-"`
}

type MessagePartBody struct {
	AttachmentId    string   `json:"attachmentId,omitempty"`
	Data            string   `json:"data,omitempty"`
	Size            int64    `json:"size,omitempty"`
	ForceSendFields []string `json:"-"`
	NullFields      []string `json:"-"`
}

type MessagePartHeader struct {
	Name            string   `json:"name,omitempty"`
	Value           string   `json:"value,omitempty"`
	ForceSendFields []string `json:"-"`
	NullFields      []string `json:"-"`
}

func NewGmailEmailService(gmailService *gmail.Service) EmailServiceInterface {
	return &GmailEmailService{
		gmailService: gmailService,
	}
}

type GmailEmailService struct {
	gmailService *gmail.Service
}

func (g *GmailEmailService) FindMessages(subject string, address string, startDate string, endDate string) ([]*EmailServiceMessage, *internal_error.InternalError) {

	messagesFound := make([]*EmailServiceMessage, 0)

	query := "from:" + address + " subject:\"" + subject + "\" after:" + startDate + " before:" + endDate

	mes, err := g.gmailService.Users.Messages.List("me").Q(query).Do()
	if err != nil {
		log.Printf("Error Listing emails: %v", err)
		return nil, internal_error.NewInternalServerError("Error listing emails")
	}

	log.Printf("Found %d messages", len(mes.Messages))

	for _, message := range mes.Messages {
		messagesFound = append(messagesFound, &EmailServiceMessage{
			Id: message.Id,
		})
	}

	return messagesFound, nil
}

func (g *GmailEmailService) GetMessage(messageId string) (*EmailServiceMessage, *internal_error.InternalError) {
	msg, err := g.gmailService.Users.Messages.Get("me", messageId).Do()
	if err != nil {
		log.Printf("Error getting message: %v", err)
		return &EmailServiceMessage{}, internal_error.NewInternalServerError("Error getting message")
	}

	headers := make([]*MessagePartHeader, 0)

	for _, h := range msg.Payload.Headers {
		headers = append(headers, &MessagePartHeader{
			Name:            h.Name,
			Value:           h.Value,
			ForceSendFields: h.ForceSendFields,
			NullFields:      h.NullFields,
		})
	}

	parts := make([]*EmailServiceMessagePart, 0)

	for _, p := range msg.Payload.Parts {
		parts = append(parts, &EmailServiceMessagePart{
			Body: &MessagePartBody{
				AttachmentId:    p.Body.AttachmentId,
				Data:            p.Body.Data,
				Size:            p.Body.Size,
				ForceSendFields: p.Body.ForceSendFields,
				NullFields:      p.Body.NullFields,
			},
			Filename:        p.Filename,
			Headers:         nil,
			MimeType:        p.MimeType,
			PartId:          p.PartId,
			Parts:           nil,
			ForceSendFields: p.ForceSendFields,
			NullFields:      p.NullFields,
		})
	}

	return &EmailServiceMessage{
		Id:           messageId,
		HistoryId:    msg.HistoryId,
		InternalDate: msg.InternalDate,
		LabelIds:     msg.LabelIds,
		Payload: &EmailServiceMessagePart{
			Body: &MessagePartBody{
				AttachmentId:    msg.Payload.Body.AttachmentId,
				Data:            msg.Payload.Body.Data,
				Size:            msg.Payload.Body.Size,
				ForceSendFields: msg.Payload.Body.ForceSendFields,
				NullFields:      msg.Payload.Body.NullFields,
			},
			Filename:        msg.Payload.Filename,
			Headers:         headers,
			MimeType:        msg.Payload.MimeType,
			PartId:          msg.Payload.PartId,
			Parts:           parts,
			ForceSendFields: msg.Payload.ForceSendFields,
			NullFields:      msg.Payload.NullFields,
		},
		Raw:          msg.Raw,
		SizeEstimate: msg.SizeEstimate,
		Snippet:      msg.Snippet,
		ThreadId:     msg.ThreadId,
	}, nil
}
