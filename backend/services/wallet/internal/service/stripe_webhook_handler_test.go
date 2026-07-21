package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v86"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	mockservice "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/service"
)

var (
	testStripeTopic       = "wallet.stripe.events"
	testStripeHeader      = "t=123,v1=abcdef"
	testStripePayload     = []byte(`{"id":"evt_test_123","type":"checkout.session.completed"}`)
	testStripeEventID     = "evt_test_123"
	testCheckoutSessionID = "cs_test_123"
)

type StripeWebhookHandlerSuite struct {
	handler      *service.StripeWebhookHandler
	producer     *mockservice.MockProducer
	constructor  *mockservice.MockStripeEventConstructor
	eventHandler *mockservice.MockHandleStripeEvent
}

func TestNewStripeWebhookHandler(t *testing.T) {
	t.Run("producer is required", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)

		r, err := service.NewStripeWebhookHandler(service.StripeWebhookHandlerConfig{
			Producer:         nil,
			EventConstructor: st.constructor,
			Topic:            testStripeTopic,
			EventHandler:     st.eventHandler,
		})

		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("event constructor is required", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)

		r, err := service.NewStripeWebhookHandler(service.StripeWebhookHandlerConfig{
			Producer:         st.producer,
			EventConstructor: nil,
			Topic:            testStripeTopic,
			EventHandler:     st.eventHandler,
		})

		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("event handler is required", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)

		r, err := service.NewStripeWebhookHandler(service.StripeWebhookHandlerConfig{
			Producer:         st.producer,
			EventConstructor: st.constructor,
			Topic:            testStripeTopic,
			EventHandler:     nil,
		})

		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("topic is required", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)

		r, err := service.NewStripeWebhookHandler(service.StripeWebhookHandlerConfig{
			Producer:         st.producer,
			EventConstructor: st.constructor,
			Topic:            "",
			EventHandler:     st.eventHandler,
		})

		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("topic is blank spaces", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)

		r, err := service.NewStripeWebhookHandler(service.StripeWebhookHandlerConfig{
			Producer:         st.producer,
			EventConstructor: st.constructor,
			Topic:            "   ",
		})

		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("successfully create an instance of StripeWebhookHandler", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)

		r, err := service.NewStripeWebhookHandler(service.StripeWebhookHandlerConfig{
			Producer:         st.producer,
			EventConstructor: st.constructor,
			Topic:            testStripeTopic,
			EventHandler:     st.eventHandler,
		})

		assert.NoError(t, err)
		assert.NotNil(t, r)
	})
}

func TestStripeWebhookHandler_Receive(t *testing.T) {
	t.Run("construct event returns error", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)
		incoming := createTestStripeEvent()
		st.constructor.EXPECT().ConstructEvent(testCtx, incoming.Payload, incoming.Header).
			Return(nil, assert.AnError)

		err := st.handler.Receive(testCtx, incoming)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidStripeEvent, err)
	})

	t.Run("produce event returns error", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)
		incoming := createTestStripeEvent()
		se := createTestStripeGoEvent()
		st.constructor.EXPECT().ConstructEvent(testCtx, incoming.Payload, incoming.Header).
			Return(se, nil)
		st.producer.EXPECT().Produce(testCtx, mock.MatchedBy(func(event *entity.Event) bool {
			return string(event.Key) == se.ID &&
				string(event.Payload) == string(incoming.Payload) &&
				event.Topic == testStripeTopic
		})).Return(assert.AnError)

		err := st.handler.Receive(testCtx, incoming)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("success receive webhook", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)
		incoming := createTestStripeEvent()
		se := createTestStripeGoEvent()
		st.constructor.EXPECT().ConstructEvent(testCtx, incoming.Payload, incoming.Header).
			Return(se, nil)
		st.producer.EXPECT().Produce(testCtx, mock.MatchedBy(func(event *entity.Event) bool {
			return string(event.Key) == se.ID &&
				string(event.Payload) == string(incoming.Payload) &&
				event.Topic == testStripeTopic
		})).Return(nil)

		err := st.handler.Receive(testCtx, incoming)

		assert.NoError(t, err)
	})
}

func TestStripeWebhookHandler_Handle(t *testing.T) {
	t.Run("fail unmarshal event payload", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)
		payload := []byte(`{"bad":"json"`)

		err := st.handler.Handle(testCtx, payload)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrBadRequest, err)
	})

	t.Run("fail unmarshal checkout session payload for checkout.session.completed event", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)
		payload := []byte(`{"id":"evt_test_123","type":"checkout.session.completed","data":{"object":"not-an-object"}}`)

		err := st.handler.Handle(testCtx, payload)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrBadRequest, err)
	})

	t.Run("checkout.session.completed handler returns error", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)
		payload := createTestCheckoutSessionCompletedPayload()
		st.eventHandler.EXPECT().HandleCheckoutSessionCompleted(testCtx, mock.MatchedBy(func(se *stripe.CheckoutSession) bool {
			return se.ID == testCheckoutSessionID
		})).Return(assert.AnError)

		err := st.handler.Handle(testCtx, payload)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("success handle checkout.session.completed", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)
		payload := createTestCheckoutSessionCompletedPayload()
		st.eventHandler.EXPECT().HandleCheckoutSessionCompleted(testCtx, mock.MatchedBy(func(se *stripe.CheckoutSession) bool {
			return se.ID == testCheckoutSessionID
		})).Return(nil)

		err := st.handler.Handle(testCtx, payload)

		assert.NoError(t, err)
	})

	t.Run("ignore unhandled event type", func(t *testing.T) {
		st := createStripeWebhookHandlerSuite(t)
		payload := []byte(`{"id":"evt_test_456","type":"payment_intent.succeeded","data":{"object":{"id":"pi_test_123"}}}`)

		err := st.handler.Handle(testCtx, payload)

		assert.NoError(t, err)
	})
}

func createStripeWebhookHandlerSuite(t *testing.T) *StripeWebhookHandlerSuite {
	p := mockservice.NewMockProducer(t)
	c := mockservice.NewMockStripeEventConstructor(t)
	e := mockservice.NewMockHandleStripeEvent(t)
	r, err := service.NewStripeWebhookHandler(service.StripeWebhookHandlerConfig{
		Producer:         p,
		EventConstructor: c,
		Topic:            testStripeTopic,
		EventHandler:     e,
	})
	assert.NoError(t, err)

	return &StripeWebhookHandlerSuite{
		handler:      r,
		producer:     p,
		constructor:  c,
		eventHandler: e,
	}
}

func createTestStripeEvent() *entity.StripeEvent {
	return &entity.StripeEvent{
		Payload: testStripePayload,
		Header:  testStripeHeader,
	}
}

func createTestStripeGoEvent() *stripe.Event {
	return &stripe.Event{
		ID: testStripeEventID,
	}
}

func createTestCheckoutSessionCompletedPayload() []byte {
	return []byte(`{"id":"evt_test_123","type":"checkout.session.completed","data":{"object":{"id":"` + testCheckoutSessionID + `"}}}`)
}
