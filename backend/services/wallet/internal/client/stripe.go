package client

import (
	"context"
	"log/slog"

	"github.com/bojanz/currency"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v86"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

const (
	defaultCurrencyDigit = uint8(2)
	ten                  = 10
)

var (
	decimalTen = decimal.NewFromInt(ten)
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

// GetCheckoutSession gets checkout session by ID.
func (s *Stripe) GetCheckoutSession(ctx context.Context, id string) (*entity.CheckoutSession, error) {
	param := &stripe.CheckoutSessionRetrieveParams{}
	session, err := s.client.V1CheckoutSessions.Retrieve(ctx, id, param)
	if err != nil {
		slog.ErrorContext(ctx, "[Stripe-GetCheckoutSessionURL] fail get checkout session", "error", err)
		return nil, entity.ErrInternal
	}
	return &entity.CheckoutSession{ID: session.ID, URL: session.URL}, nil
}

// CreateCheckoutSession creates a checkout.
func (s *Stripe) CreateCheckoutSession(ctx context.Context, input *entity.CheckoutInput) (*entity.CheckoutSession, error) {
	amount := toSmallestUnitCurrency(input.Amount, input.Currency)

	param := &stripe.CheckoutSessionCreateParams{
		Customer:   stripe.String(input.StripeCustomerID),
		SuccessURL: stripe.String(input.SuccessURL),
		Mode:       stripe.String(stripe.CheckoutSessionModePayment),
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{
				Quantity: stripe.Int64(int64(input.Quantity)),
				PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
					Currency:   stripe.String(input.Currency),
					UnitAmount: stripe.Int64(amount),
					ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
						Name:        stripe.String(input.Purpose),
						Description: stripe.String(input.Purpose),
					},
				},
			},
		},
	}

	session, err := s.client.V1CheckoutSessions.Create(ctx, param)
	if err != nil {
		slog.ErrorContext(ctx, "[Stripe-CreateCheckoutSessionURL] fail create checkout session", "error", err)
		return nil, entity.ErrInternal
	}
	return &entity.CheckoutSession{ID: session.ID, URL: session.URL}, nil
}

// ConstructEvent constructs an incoming payload to be an event.
// Prior to constructing, it will validate the header and secret.
func (s *Stripe) ConstructEvent(ctx context.Context, payload []byte, header string, secret string) (*stripe.Event, error) {
	event, err := s.client.ConstructEvent(payload, header, secret)
	if err != nil {
		slog.ErrorContext(ctx, "[Stripe-CreateCheckoutSessionURL] fail create checkout session", "error", err)
		return nil, entity.ErrInvalidStripeEvent
	}
	return &event, nil
}

func toSmallestUnitCurrency(amout decimal.Decimal, currencyCode string) int64 {
	digit, ok := currency.GetDigits(currencyCode)
	if !ok {
		digit = defaultCurrencyDigit
	}

	cents := amout.Mul(decimalTen.Pow(decimal.NewFromInt32(int32(digit)))).IntPart()
	return cents
}
