package entity

// StripeEvent defines Stripe event.
type StripeEvent struct {
	Header  string
	Payload []byte
}
