package client

import (
	"context"
	"log/slog"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/stripe/stripe-go/v86"
)

// Stripe is responsible to connect with Stripe API.
type Stripe struct {
	client *stripe.Client
}

// NewStripe creates an instance of Stripe.
func NewStripe(apiKey string) *Stripe {
	c := stripe.NewClient(apiKey)
	return &Stripe{client: c}
}

// CreateCustomer creates a customer in client side.
// It returns ID from client.
func (s *Stripe) CreateCustomer(ctx context.Context, email string) (string, error) {
	param := &stripe.CustomerCreateParams{
		Email: stripe.String(email),
	}

	cust, err := s.client.V1Customers.Create(ctx, param)
	if err != nil {
		slog.ErrorContext(ctx, "[Stripe-CreateCustomer] fail create customer", "error", err)
		return "", entity.ErrInternal
	}
	return cust.ID, nil
}

// GetCheckoutSessionURL gets checkout session by ID.
func (s *Stripe) GetCheckoutSessionURL(ctx context.Context, id string) (string, error) {
	param := &stripe.CheckoutSessionRetrieveParams{}
	session, err := s.client.V1CheckoutSessions.Retrieve(ctx, id, param)
	if err != nil {
		slog.ErrorContext(ctx, "[Stripe-GetCheckoutSessionURL] fail get checkout session", "error", err)
		return "", entity.ErrInternal
	}
	return session.URL, nil
}
