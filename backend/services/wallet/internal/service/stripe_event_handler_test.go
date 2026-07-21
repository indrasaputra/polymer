package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go/v86"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	mockservice "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/service"
)

type StripeEventHandlerSuite struct {
	eventHandler *service.StripeEventHandler
	repo         *mockservice.MockTransactionRepository
}

func TestNewStripeEventHandler(t *testing.T) {
	t.Run("successfully create an instance of StripeEventHandler", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)

		assert.NotNil(t, st.eventHandler)
	})
}

func TestStripeEventHandler_HandleCheckoutSessionCompleted(t *testing.T) {
	t.Run("payment status is not paid, does nothing", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestCheckoutSession(stripe.CheckoutSessionPaymentStatusUnpaid)

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.NoError(t, err)
	})

	t.Run("payment status is no payment required, does nothing", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestCheckoutSession(stripe.CheckoutSessionPaymentStatusNoPaymentRequired)

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.NoError(t, err)
	})

	t.Run("payment status is paid, but transaction is already gone (no pending transaction)", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestCheckoutSession(stripe.CheckoutSessionPaymentStatusPaid)
		st.repo.EXPECT().UpdatePendingTransactionByCheckoutSessionIDToCompleted(testCtx, session.ID).
			Return(entity.ErrNilTransaction)

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.NoError(t, err)
	})

	t.Run("payment status is paid, update pending transaction returns error", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestCheckoutSession(stripe.CheckoutSessionPaymentStatusPaid)
		st.repo.EXPECT().UpdatePendingTransactionByCheckoutSessionIDToCompleted(testCtx, session.ID).
			Return(assert.AnError)

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("payment status is paid, success update pending transaction to completed", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestCheckoutSession(stripe.CheckoutSessionPaymentStatusPaid)
		st.repo.EXPECT().UpdatePendingTransactionByCheckoutSessionIDToCompleted(testCtx, session.ID).
			Return(nil)

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.NoError(t, err)
	})
}

func createStripeEventHandlerSuite(t *testing.T) *StripeEventHandlerSuite {
	r := mockservice.NewMockTransactionRepository(t)
	h := service.NewStripeEventHandler(r)

	return &StripeEventHandlerSuite{
		eventHandler: h,
		repo:         r,
	}
}

func createTestCheckoutSession(status stripe.CheckoutSessionPaymentStatus) *stripe.CheckoutSession {
	return &stripe.CheckoutSession{
		ID:            testCheckoutSessionID,
		PaymentStatus: status,
	}
}
