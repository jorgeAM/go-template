package mailer

import (
	"context"

	"github.com/jorgeAM/go-template/pkg/log"
)

var _ Mailer = (*InMemoryMailer)(nil)

type InMemoryMailer struct {
}

func NewInMemoryMailer() *InMemoryMailer {
	return &InMemoryMailer{}
}

func (s *InMemoryMailer) Send(ctx context.Context, payload *MailerPayload) error {
	log.Debug(
		ctx,
		"sending email using mock version",
		log.WithString("from", payload.From),
		log.WithString("to", payload.To),
		log.WithString("subject", payload.Subject),
	)

	return nil
}
