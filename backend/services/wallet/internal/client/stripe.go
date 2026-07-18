package client

import (
	"context"

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
		return "", nil
	}
	return cust.ID, nil
}
