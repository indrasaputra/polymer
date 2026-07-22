package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/stripe/stripe-go/v86"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

// HandleStripeWebhook defines interface to receive webhook from Stripe.
type HandleStripeWebhook interface {
	// Receive receives a webhook.
	Receive(ctx context.Context, incoming *entity.StripeEvent) error
}

// Producer defines the interface to produce event.
type Producer interface {
	// Produce produces event to message queue.
	Produce(ctx context.Context, event *entity.Event) error
}

// StripeWebhookHandler is responsible to handle webhook from Stripe.
type StripeWebhookHandler struct {
	producer    Producer
	constructor StripeEventConstructor
	handler     HandleStripeEvent
	topic       string
}

// StripeEventConstructor defines interface to construct Stripe event.
type StripeEventConstructor interface {
	// ConstructEvent construct incoming payload to be a Stripe event.
	ConstructEvent(ctx context.Context, payload []byte, header string) (*stripe.Event, error)
}

// HandleStripeEvent defines interface to handle Stripe event.
type HandleStripeEvent interface {
	// HandleCheckoutSessionCompleted handles checkout.session.completed.
	HandleCheckoutSessionCompleted(ctx context.Context, event *stripe.CheckoutSession) error
}

// StripeWebhookHandlerConfig defines config for Stripe webhook receiver.
// All are required.
type StripeWebhookHandlerConfig struct {
	Producer         Producer
	EventConstructor StripeEventConstructor
	EventHandler     HandleStripeEvent
	Topic            string
}

// NewStripeWebhookHandler creates an instance of StripeWebhookHandler.
func NewStripeWebhookHandler(c StripeWebhookHandlerConfig) (*StripeWebhookHandler, error) {
	if c.Producer == nil {
		return nil, errors.New("producer is required")
	}
	if c.EventConstructor == nil {
		return nil, errors.New("event constructor is required")
	}
	if c.EventHandler == nil {
		return nil, errors.New("event handler is required")
	}
	if strings.TrimSpace(c.Topic) == "" {
		return nil, errors.New("topic is required")
	}

	return &StripeWebhookHandler{
		producer:    c.Producer,
		handler:     c.EventHandler,
		constructor: c.EventConstructor,
		topic:       strings.TrimSpace(c.Topic),
	}, nil
}

// Receive receives a webhook payload, validates the payload, and sends to message queue for further processing.
func (s *StripeWebhookHandler) Receive(ctx context.Context, incoming *entity.StripeEvent) error {
	se, err := s.constructor.ConstructEvent(ctx, incoming.Payload, incoming.Header)
	if err != nil {
		slog.ErrorContext(ctx, "[StripeWebhookHandler-Receive] incoming event is invalid", "error", err)
		return entity.ErrStripeEventInvalid
	}

	event := &entity.Event{
		Key:     []byte(se.ID),
		Payload: incoming.Payload,
		Topic:   s.topic,
	}
	if err := s.producer.Produce(ctx, event); err != nil {
		slog.ErrorContext(ctx, "[StripeWebhookHandler-Receive] fail produce event", "error", err)
		return entity.ErrInternal
	}
	return nil
}

// Handle handles incoming payload, which comes from Kafka consumer, and do the business process.
func (s *StripeWebhookHandler) Handle(ctx context.Context, payload []byte) error {
	var event stripe.Event
	if err := json.Unmarshal(payload, &event); err != nil {
		slog.ErrorContext(ctx, "[StripeWebhookHandler-Handle] fail unmarshal event payload", "error", err)
		return entity.ErrGeneralInvalid
	}

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var se stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &se); err != nil {
			slog.ErrorContext(ctx, "[StripeWebhookHandler-Handle] fail unmarshal checkout.session.completed payload", "error", err)
			return entity.ErrGeneralInvalid
		}
		return s.handler.HandleCheckoutSessionCompleted(ctx, &se)
	default:
		slog.InfoContext(ctx, "[StripeWebhookHandler-Handle] ignore event", "type", event.Type)
	}

	return nil
}
