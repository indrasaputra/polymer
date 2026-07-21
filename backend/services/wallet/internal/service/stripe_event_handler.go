package service

import (
	"context"
	"log/slog"

	"github.com/stripe/stripe-go/v86"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

// TransactionRepository defines the interface for transaction-related table.
type TransactionRepository interface {
	// UpdatePendingTransactionByCheckoutSessionIDToCompleted updates the pending transaction to paid.
	UpdatePendingTransactionByCheckoutSessionIDToCompleted(ctx context.Context, id string) error
}

// StripeEventHandler is responsible to handle Stripe event.
type StripeEventHandler struct {
	repo TransactionRepository
}

// NewStripeEventHandler creates an instance of StripeEventHandler.
func NewStripeEventHandler() *StripeEventHandler {
	return &StripeEventHandler{}
}

// HandleCheckoutSessionCompleted handles checkout.session.completed.
func (s *StripeEventHandler) HandleCheckoutSessionCompleted(ctx context.Context, session *stripe.CheckoutSession) error {
	if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		err := s.repo.UpdatePendingTransactionByCheckoutSessionIDToCompleted(ctx, session.ID)
		if err == entity.ErrNilTransaction {
			return nil
		}
		if err != nil {
			slog.ErrorContext(ctx, "[StripeEventHandler-HandleCheckoutSessionCompleted] fail update pending transaction", "error", err)
			return entity.ErrInternal
		}
	}
	return nil
}
