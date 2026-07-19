package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/stripe/stripe-go/v86"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

// ReceiveStripeWebhook defines interface to receive webhook from Stripe.
type ReceiveStripeWebhook interface {
	// Receive receives a webhook.
	Receive(ctx context.Context, incoming *entity.StripeEvent) error
}

// Producer defines the interface to produce event.
type Producer interface {
	// Produce produces event to message queue.
	Produce(ctx context.Context, event *entity.Event) error
}

// StripeWebhookReceiver is responsible to receive webhook from Stripe.
type StripeWebhookReceiver struct {
	producer    Producer
	constructor StripeEventConstructor
	secret      string
	topic       string
}

// StripeEventConstructor defines interface to construct Stripe event.
type StripeEventConstructor interface {
	// ConstructEvent construct incoming payload to be a Stripe event.
	ConstructEvent(payload []byte, header string, secret string) (*stripe.Event, error)
}

// StripeWebhookReceiverConfig defines config for Stripe webhook receiver.
// All are required.
type StripeWebhookReceiverConfig struct {
	Producer         Producer
	EventConstructor StripeEventConstructor
	Secret           string
	Topic            string
}

// NewStripeWebhookReceiver creates an instance of StripeWebhookReceiver.
func NewStripeWebhookReceiver(c StripeWebhookReceiverConfig) (*StripeWebhookReceiver, error) {
	if c.Producer == nil {
		return nil, errors.New("producer is required")
	}
	if c.EventConstructor == nil {
		return nil, errors.New("event constructor is required")
	}
	if strings.TrimSpace(c.Secret) == "" {
		return nil, errors.New("secret is required")
	}
	if strings.TrimSpace(c.Topic) == "" {
		return nil, errors.New("topic is required")
	}

	return &StripeWebhookReceiver{
		producer:    c.Producer,
		constructor: c.EventConstructor,
		secret:      strings.TrimSpace(c.Secret),
		topic:       strings.TrimSpace(c.Topic),
	}, nil
}

// Receive receives a webhook payload, validates the payload, and sends to message queue for further processing.
func (s *StripeWebhookReceiver) Receive(ctx context.Context, incoming *entity.StripeEvent) error {
	se, err := s.constructor.ConstructEvent(incoming.Payload, incoming.Header, s.secret)
	if err != nil {
		slog.ErrorContext(ctx, "[StripeWebhookReceiver-Receive] incoming event is invalid", "error", err)
		return entity.ErrInvalidStripeEvent
	}

	event := &entity.Event{
		Key:     []byte(se.ID),
		Payload: incoming.Payload,
		Topic:   s.topic,
	}
	if err := s.producer.Produce(ctx, event); err != nil {
		slog.ErrorContext(ctx, "[StripeWebhookReceiver-Receive] fail produce event", "error", err)
		return entity.ErrInternal
	}
	return nil
}
