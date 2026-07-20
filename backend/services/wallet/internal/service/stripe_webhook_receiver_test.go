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
	testStripeTopic   = "wallet.stripe.events"
	testStripeHeader  = "t=123,v1=abcdef"
	testStripePayload = []byte(`{"id":"evt_test_123","type":"checkout.session.completed"}`)
	testStripeEventID = "evt_test_123"
)

type StripeWebhookReceiverSuite struct {
	receiver    *service.StripeWebhookReceiver
	producer    *mockservice.MockProducer
	constructor *mockservice.MockStripeEventConstructor
}

func TestNewStripeWebhookReceiver(t *testing.T) {
	t.Run("producer is required", func(t *testing.T) {
		st := createStripeWebhookReceiverSuite(t)

		r, err := service.NewStripeWebhookReceiver(service.StripeWebhookReceiverConfig{
			Producer:         nil,
			EventConstructor: st.constructor,
			Topic:            testStripeTopic,
		})

		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("event constructor is required", func(t *testing.T) {
		st := createStripeWebhookReceiverSuite(t)

		r, err := service.NewStripeWebhookReceiver(service.StripeWebhookReceiverConfig{
			Producer:         st.producer,
			EventConstructor: nil,
			Topic:            testStripeTopic,
		})

		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("topic is required", func(t *testing.T) {
		st := createStripeWebhookReceiverSuite(t)

		r, err := service.NewStripeWebhookReceiver(service.StripeWebhookReceiverConfig{
			Producer:         st.producer,
			EventConstructor: st.constructor,
			Topic:            "",
		})

		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("topic is blank spaces", func(t *testing.T) {
		st := createStripeWebhookReceiverSuite(t)

		r, err := service.NewStripeWebhookReceiver(service.StripeWebhookReceiverConfig{
			Producer:         st.producer,
			EventConstructor: st.constructor,
			Topic:            "   ",
		})

		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("successfully create an instance of StripeWebhookReceiver", func(t *testing.T) {
		st := createStripeWebhookReceiverSuite(t)

		r, err := service.NewStripeWebhookReceiver(service.StripeWebhookReceiverConfig{
			Producer:         st.producer,
			EventConstructor: st.constructor,
			Topic:            testStripeTopic,
		})

		assert.NoError(t, err)
		assert.NotNil(t, r)
	})
}

func TestStripeWebhookReceiver_Receive(t *testing.T) {
	t.Run("construct event returns error", func(t *testing.T) {
		st := createStripeWebhookReceiverSuite(t)
		incoming := createTestStripeEvent()
		st.constructor.EXPECT().ConstructEvent(testCtx, incoming.Payload, incoming.Header).
			Return(nil, assert.AnError)

		err := st.receiver.Receive(testCtx, incoming)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidStripeEvent, err)
	})

	t.Run("produce event returns error", func(t *testing.T) {
		st := createStripeWebhookReceiverSuite(t)
		incoming := createTestStripeEvent()
		se := createTestStripeGoEvent()
		st.constructor.EXPECT().ConstructEvent(testCtx, incoming.Payload, incoming.Header).
			Return(se, nil)
		st.producer.EXPECT().Produce(testCtx, mock.MatchedBy(func(event *entity.Event) bool {
			return string(event.Key) == se.ID &&
				string(event.Payload) == string(incoming.Payload) &&
				event.Topic == testStripeTopic
		})).Return(assert.AnError)

		err := st.receiver.Receive(testCtx, incoming)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("success receive webhook", func(t *testing.T) {
		st := createStripeWebhookReceiverSuite(t)
		incoming := createTestStripeEvent()
		se := createTestStripeGoEvent()
		st.constructor.EXPECT().ConstructEvent(testCtx, incoming.Payload, incoming.Header).
			Return(se, nil)
		st.producer.EXPECT().Produce(testCtx, mock.MatchedBy(func(event *entity.Event) bool {
			return string(event.Key) == se.ID &&
				string(event.Payload) == string(incoming.Payload) &&
				event.Topic == testStripeTopic
		})).Return(nil)

		err := st.receiver.Receive(testCtx, incoming)

		assert.NoError(t, err)
	})
}

func createStripeWebhookReceiverSuite(t *testing.T) *StripeWebhookReceiverSuite {
	p := mockservice.NewMockProducer(t)
	c := mockservice.NewMockStripeEventConstructor(t)
	r, err := service.NewStripeWebhookReceiver(service.StripeWebhookReceiverConfig{
		Producer:         p,
		EventConstructor: c,
		Topic:            testStripeTopic,
	})
	assert.NoError(t, err)

	return &StripeWebhookReceiverSuite{
		receiver:    r,
		producer:    p,
		constructor: c,
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
