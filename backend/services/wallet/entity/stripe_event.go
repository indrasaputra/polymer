package entity

// StripeEventType defines Stripe event type.
type StripeEventType string

const (
	// StripeEventTypeCheckoutSessionCompleted derived from Stripe API docs.
	StripeEventTypeCheckoutSessionCompleted = "checkout.session.completed"
)

// StripeEvent defines Stripe event.
type StripeEvent struct {
	Header  string
	Payload []byte
}
