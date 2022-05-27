package mail

//import (
//	"context"
//	"github.com/fundever/backend/domain"
//	"github.com/mailjet/mailjet-apiv3-go/v3"
//)
//
//func NewMailJet(
//	apiKey string,
//	apiSecret string,
//) Mailer {
//	m := mailjet.NewMailjetClient(apiKey, apiSecret)
//	return mailJet{m}
//}
//
//type mailJet struct {
//	*mailjet.Client
//}
//
//func (m mailJet) Send(
//	ctx context.Context,
//	receiver,
//	name,
//	subject,
//	messageInHtml string,
//) error {
//	messagesInfo := []mailjet.InfoMessagesV31{
//		{
//			From: &mailjet.RecipientV31{
//				Email: "no-reply@fundever.com",
//				//Email: "chet@rail.town",
//				Name: "Fundever Team",
//			},
//			To: &mailjet.RecipientsV31{
//				mailjet.RecipientV31{
//					Email: receiver,
//					Name:  name,
//				},
//			},
//			Subject:  subject,
//			HTMLPart: messageInHtml,
//		},
//	}
//	messages := mailjet.MessagesV31{Info: messagesInfo}
//
//	var err error
//	done := make(chan struct{})
//	go func() {
//		_, err = m.SendMailV31(&messages)
//		done <- struct{}{}
//	}()
//
//	select {
//	case <-ctx.Done():
//	case <-done:
//		if err != nil {
//			return domain.ServiceError{
//				Name: "mail-jet",
//				Err:  err,
//			}
//		}
//	}
//
//	return nil
//}
