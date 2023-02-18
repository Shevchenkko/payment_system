// Package api implements external api calls.
package emails

import (
	"context"
	"fmt"
	"os"

	"github.com/Shevchenkko/payment_system/internal/service"
	gomail "gopkg.in/mail.v2"
)

// Emails - represents api which is used for emails.
type Emails struct{}

// NewEmails - creates new instance of emails api.
func New() *Emails {
	return &Emails{}
}

func (e *Emails) SendEmail(ctx context.Context, inp service.SendEmailInput) error {
	// build email
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {os.Getenv("MAIL_USERNAME")},
		"To":      {inp.To},
		"Subject": {inp.Subject},
	})
	m.SetBody(inp.ContentType, inp.Body)

	// sending email
	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), 587, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_APP_PASSWORD"))
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println("ERROR", err)
		return err
	}

	return nil
}
