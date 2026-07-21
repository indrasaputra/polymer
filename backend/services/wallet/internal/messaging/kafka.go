package messaging

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

// KafkaProducer is responsible for producing event to kafka.
type KafkaProducer struct {
	client *kgo.Client
}

// NewKafkaProducer creates an instance of KafkaProducer.
func NewKafkaProducer(bs []string) (*KafkaProducer, error) {
	c, err := kgo.NewClient(
		kgo.SeedBrokers(bs...),
	)
	if err != nil {
		return nil, fmt.Errorf("fail instantiate kafka producer client: %v", err)
	}
	return &KafkaProducer{client: c}, nil
}

// Close closes kafka's client.
func (k *KafkaProducer) Close() {
	k.client.Close()
}

// Produce produces event to kafka.
// It is synchronous process.
func (k *KafkaProducer) Produce(ctx context.Context, event *entity.Event) error {
	record := &kgo.Record{
		Key:   event.Key,
		Value: event.Payload,
		Topic: event.Topic,
	}
	res := k.client.ProduceSync(ctx, record)
	if err := res.FirstErr(); err != nil {
		slog.ErrorContext(ctx, "[Producer-Produce] fail produce event", "error", err)
		return err
	}
	return nil
}

// RecordHandler defines interface to handle incoming record.
type RecordHandler interface {
	// Handle handles payload and process according to business process.
	Handle(ctx context.Context, payload []byte) error
}

// KafkaStripeWebhookConsumer is responsible to consume Stripe webhook event from kafka.
type KafkaStripeWebhookConsumer struct {
	client  *kgo.Client
	handler RecordHandler
}

// NewKafkaStripeWebhookConsumer creates an instance of KafkaStripeWebhookConsumer.
func NewKafkaStripeWebhookConsumer(h RecordHandler, bs []string, topic string, cgid string) (*KafkaStripeWebhookConsumer, error) {
	c, err := kgo.NewClient(
		kgo.SeedBrokers(bs...),
		kgo.ConsumerGroup(cgid),
		kgo.ConsumeTopics(topic),
		kgo.AutoCommitMarks(), // balance between auto-commit and manual-commit
	)
	if err != nil {
		return nil, fmt.Errorf("fail instantiate kafka client: %v", err)
	}
	return &KafkaStripeWebhookConsumer{client: c, handler: h}, nil
}

// Close closes kafka's client.
func (k *KafkaStripeWebhookConsumer) Close() {
	k.client.Close()
}

// Consume consumes event.
func (k *KafkaStripeWebhookConsumer) Consume(ctx context.Context) {
	for {
		fetches := k.client.PollFetches(ctx)
		if err := fetches.Err(); err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.ErrorContext(ctx, "[KafkaStripeWebhookConsumer-Consume] fetch error", "error", err)
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			if err := k.handler.Handle(ctx, record.Value); err == nil {
				k.client.MarkCommitRecords(record)
			}
		}
	}
}
