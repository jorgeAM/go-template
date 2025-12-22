package mailer

import (
	"context"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var _ Mailer = (*SendgridMailer)(nil)

type SendgridMailer struct {
	client *sendgrid.Client
}

func NewSendgridMailer(client *sendgrid.Client) *SendgridMailer {
	return &SendgridMailer{client}
}

func (s *SendgridMailer) Send(ctx context.Context, payload *MailerPayload) error {
	message := mail.NewV3MailInit(
		&mail.Email{
			Address: payload.From,
		},
		payload.Subject,
		&mail.Email{
			Address: payload.To,
		},
		&mail.Content{
			Type:  "text/html",
			Value: payload.Body,
		},
	)

	_, err := s.client.SendWithContext(ctx, message)

	return err
}
