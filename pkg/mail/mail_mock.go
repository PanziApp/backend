package mail

import (
	"context"
	"log"
)

type MailMock struct {
}

func (m MailMock) Send(
	ctx context.Context,
	receiver, name, subject, messageInHtml string,
) error {
	log.Printf("Email for %s (%s): %s => %s", name, receiver, subject, messageInHtml)
	return nil
}
