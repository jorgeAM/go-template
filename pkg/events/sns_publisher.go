package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/jorgeAM/go-template/pkg/collections"
	"golang.org/x/sync/errgroup"
)

var _ Publisher = (*SNSPublisher)(nil)

const MAX_BATCH_SIZE = 10

type snsMessage struct {
	ID        string          `json:"id"`
	Topic     string          `json:"topic"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

type SNSPublisher struct {
	client   *sns.Client
	topicArn string
}

func NewSNSPublisher(
	client *sns.Client,
	topicArn string,
) *SNSPublisher {
	return &SNSPublisher{
		client:   client,
		topicArn: topicArn,
	}
}

func (s *SNSPublisher) Publish(ctx context.Context, events ...*Event) error {
	batchEvents := collections.Chunks(events, MAX_BATCH_SIZE)
	gr, ctx := errgroup.WithContext(ctx)

	for _, events := range batchEvents {
		events := events

		gr.Go(func() error {
			return s.batchPublish(ctx, events)
		})
	}

	return gr.Wait()
}

func (s *SNSPublisher) batchPublish(ctx context.Context, events []*Event) error {
	requests := make([]types.PublishBatchRequestEntry, 0, len(events))

	for _, event := range events {
		payload, err := event.MarshalPayload()
		if err != nil {
			return err
		}

		message := &snsMessage{
			ID:        event.ID,
			Topic:     string(event.Topic),
			Payload:   payload,
			Timestamp: event.Timestamp,
		}

		msgJson, err := json.Marshal(message)
		if err != nil {
			return err
		}

		attrs := map[string]types.MessageAttributeValue{
			"topic": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.Topic.String()),
			},
		}

		requests = append(requests, types.PublishBatchRequestEntry{
			Id:                aws.String(event.ID),
			Message:           aws.String(string(msgJson)),
			MessageAttributes: attrs,
		})
	}

	_, err := s.client.PublishBatch(
		ctx,
		&sns.PublishBatchInput{
			TopicArn:                   &s.topicArn,
			PublishBatchRequestEntries: requests,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
