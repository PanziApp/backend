package mail

import (
	"context"
	"fmt"
	"github.com/PanziApp/backend/internal/domain"
	"net/smtp"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

type Gmail struct {
	from string
	auth smtp.Auth
}

func NewGmailSender(username, password string) Mailer {
	return Gmail{
		from: username,
		auth: smtp.PlainAuth("", username, password, smtpHost),
	}
}

func (m Gmail) Send(ctx context.Context, receiver, name, subject, messageInHtml string) error {

	data := []byte(fmt.Sprintf(
		`To: %s <%s>
Subject: %s
Mime-Version: 1.0
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: 8bit
X-Auto-Response-Suppress: All

%s`, name, receiver, subject, messageInHtml,
	))

	err := smtp.SendMail(smtpHost+":"+smtpPort, m.auth, m.from, []string{receiver}, data)
	if err != nil {
		return domain.ServiceError{
			Name: "gmail",
			Err:  err,
		}
	}
	return nil
}
