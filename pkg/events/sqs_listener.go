package events

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jorgeAM/go-template/pkg/log"
)

var _ Listener = (*SQSListener)(nil)

type SQSListener struct {
	client   *sqs.Client
	queueURL string
	handlers map[Topic]Handler
	waitTime int32
}

func NewSQSListener(client *sqs.Client, queueURL string, handlers map[Topic]Handler, waitTime int32) *SQSListener {
	return &SQSListener{
		client:   client,
		queueURL: queueURL,
		handlers: handlers,
		waitTime: waitTime,
	}
}

func (s *SQSListener) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			output, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            &s.queueURL,
				WaitTimeSeconds:     s.waitTime,
				MaxNumberOfMessages: 10,
			})
			if err != nil {
				log.Error(ctx, "error receiving message from SQS", log.WithError(err))
				continue
			}

			for _, msg := range output.Messages {
				var event Event
				if err := json.Unmarshal([]byte(*msg.Body), &event); err != nil {
					log.Error(ctx, "invalid message format from SQS", log.WithError(err))
					continue
				}

				handler, ok := s.handlers[event.Topic]
				if !ok {
					log.Warn(ctx, "no handler for event type", log.WithString("topic", event.Topic.String()))
					continue
				}

				if err := handler.Handle(ctx, &event); err != nil {
					log.Error(ctx, "error handling event from SQS", log.WithString("topic", event.Topic.String()), log.WithError(err))
					continue
				}

				_, _ = s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &s.queueURL,
					ReceiptHandle: msg.ReceiptHandle,
				})
			}
		}
	}
}
