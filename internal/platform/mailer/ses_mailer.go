package mailer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

var _ Mailer = (*SESMailer)(nil)

type SESMailer struct {
	client *ses.Client
}

func NewSESMailer(client *ses.Client) *SESMailer {
	return &SESMailer{client}
}

func (s *SESMailer) Send(ctx context.Context, payload *MailerPayload) error {
	_, err := s.client.SendEmail(ctx, &ses.SendEmailInput{
		FromEmailAddress: aws.String(payload.From),
		Destination: &types.Destination{
			ToAddresses: []string{payload.To},
			CcAddresses: payload.CC,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(payload.Subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(payload.Body),
					},
				},
			},
		},
	})
	return err
}
