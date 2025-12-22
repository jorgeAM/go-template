package mailer

import "context"

type MailerPayload struct {
	From    string   `json:"from"`
	CC      []string `json:"cc"`
	To      string   `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

//go:generate mockgen -source=./mailer.go -destination=./mocks/mailer.go -package=mock -mock_names=Mailer=MockMailer
type Mailer interface {
	Send(ctx context.Context, payload *MailerPayload) error
}
