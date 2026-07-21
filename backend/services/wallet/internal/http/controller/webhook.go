package controller

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/dto"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
)

const (
	headerStripeSignature = "Stripe-Signature"
	maxBodyBytes          = int64(65536) // 64KB
)

// Webhook handles HTTP request for webhook.
type Webhook struct {
	handler service.HandleStripeWebhook
}

// NewWebhook creates an instance of Webhook.
func NewWebhook(h service.HandleStripeWebhook) *Webhook {
	return &Webhook{handler: h}
}

// RegisterRoute registers all routes in webhook controller.
func (w *Webhook) RegisterRoute(g *echo.Group, _ echo.MiddlewareFunc) {
	wg := g.Group("/webhooks")

	wg.POST("/stripe", w.Stripe)
}

// Stripe receives webhook from Stripe.
func (w *Webhook) Stripe(c *echo.Context) error {
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxBodyBytes)
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "[WebhookController-Stripe] fail read payload")
		return dto.SendResponse(c, nil, entity.ErrBadRequest, 0)
	}

	header := c.Request().Header.Get(headerStripeSignature)
	if strings.TrimSpace(header) == "" {
		slog.ErrorContext(c.Request().Context(), "[WebhookController-Stripe] fail get header")
		return dto.SendResponse(c, nil, entity.ErrBadRequest, 0)
	}

	event := &entity.StripeEvent{Payload: payload, Header: header}
	err = w.handler.Receive(c.Request().Context(), event)
	return dto.SendResponse(c, nil, err, http.StatusOK)
}
