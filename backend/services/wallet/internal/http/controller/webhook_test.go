package controller_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/labstack/echo/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/controller"
	mockservice "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/service"
)

const (
	headerStripeSignature = "Stripe-Signature"
)

func TestNewWebhook(t *testing.T) {
	t.Run("success create an instance of Webhook", func(t *testing.T) {
		st := createWebhookSuite(t)

		assert.NotNil(t, st)
	})
}

func TestWebhook_RegisterRoute(t *testing.T) {
	t.Run("success register route", func(t *testing.T) {
		st := createWebhookSuite(t)
		e := echo.New()

		assert.NotPanics(t, func() { st.webhookController.RegisterRoute(e.Group("/api/v1"), middleware.RequestID()) })
	})
}

func TestWebhook_Stripe(t *testing.T) {
	testPayload := []byte(`{"id":"evt_test_123","type":"checkout.session.completed"}`)
	testSignature := "t=123,v1=abcdef"

	t.Run("fail read payload due to body exceeding max bytes", func(t *testing.T) {
		oversized := make([]byte, 70000)
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				headerStripeSignature:  {testSignature},
			},
			JSONBody: oversized,
		}.ToContextRecorder(t)

		st := createWebhookSuite(t)

		err := st.webhookController.Stripe(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("fail due to missing stripe signature header", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
			},
			JSONBody: testPayload,
		}.ToContextRecorder(t)

		st := createWebhookSuite(t)

		err := st.webhookController.Stripe(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("fail due to blank stripe signature header", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				headerStripeSignature:  {"   "},
			},
			JSONBody: testPayload,
		}.ToContextRecorder(t)

		st := createWebhookSuite(t)

		err := st.webhookController.Stripe(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("receiver returns error", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				headerStripeSignature:  {testSignature},
			},
			JSONBody: testPayload,
		}.ToContextRecorder(t)

		st := createWebhookSuite(t)
		st.webhookReceiver.EXPECT().Receive(c.Request().Context(), mock.MatchedBy(func(event *entity.StripeEvent) bool {
			return string(event.Payload) == string(testPayload) && event.Header == testSignature
		})).Return(assert.AnError)

		err := st.webhookController.Stripe(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success receive stripe webhook", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				headerStripeSignature:  {testSignature},
			},
			JSONBody: testPayload,
		}.ToContextRecorder(t)

		st := createWebhookSuite(t)
		st.webhookReceiver.EXPECT().Receive(c.Request().Context(), mock.MatchedBy(func(event *entity.StripeEvent) bool {
			return string(event.Payload) == string(testPayload) && event.Header == testSignature
		})).Return(nil)

		err := st.webhookController.Stripe(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

type WebhookSuite struct {
	webhookController *controller.Webhook
	webhookReceiver   *mockservice.MockReceiveStripeWebhook
}

func createWebhookSuite(t *testing.T) *WebhookSuite {
	r := mockservice.NewMockReceiveStripeWebhook(t)
	w := controller.NewWebhook(r)
	return &WebhookSuite{
		webhookController: w,
		webhookReceiver:   r,
	}
}
