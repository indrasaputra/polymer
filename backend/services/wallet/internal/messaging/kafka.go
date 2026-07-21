package messaging

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

// Kafka is responsible for producing event to kafka.
type Kafka struct {
	client *kgo.Client
}

// NewKafka creates an instance of Producer.
func NewKafka(bs []string) (*Kafka, error) {
	c, err := kgo.NewClient(
		kgo.SeedBrokers(bs...),
	)
	if err != nil {
		return nil, fmt.Errorf("fail instantiate kafka client: %v", err)
	}
	return &Kafka{client: c}, nil
}

// Close closes kafka's client.
func (k *Kafka) Close() {
	k.client.Close()
}

// Produce produces event to kafka.
// It is synchronous process.
func (k *Kafka) Produce(ctx context.Context, event *entity.Event) error {
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
