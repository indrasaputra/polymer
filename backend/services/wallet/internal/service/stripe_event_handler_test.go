package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v86"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	"github.com/indrasaputra/polymer/backend/services/wallet/test/mock/pkg/sdk/uow"
	mockservice "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/service"
)

type StripeEventHandlerSuite struct {
	eventHandler *service.StripeEventHandler
	txManager    *uow.MockTxManager
	repo         *mockservice.MockHandleStripeEventRepository
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

	t.Run("client reference id is not a valid uuid", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession("not-a-uuid")
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				err := fn(testCtxTx)
				assert.Error(t, err)
				assert.Equal(t, entity.ErrGeneralInvalid, err)
				return err
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("get active wallet returns not found", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(nil, entity.ErrWalletNotFound)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				err := fn(testCtxTx)
				assert.Error(t, err)
				assert.Equal(t, entity.ErrWalletNotFound, err)
				return err
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("get active wallet returns unexpected error", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(nil, assert.AnError)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				err := fn(testCtxTx)
				assert.Error(t, err)
				assert.Equal(t, entity.ErrInternal, err)
				return err
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("get active transaction returns not found", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		wallet := createTestActiveWallet(testWalletID)
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(wallet, nil)
		st.repo.EXPECT().GetActiveTransactionByCheckoutSessionIDForUpdate(testCtxTx, session.ID).
			Return(nil, entity.ErrTransactionNotFound)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				err := fn(testCtxTx)
				assert.Error(t, err)
				assert.Equal(t, entity.ErrTransactionNotFound, err)
				return err
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("get active transaction returns unexpected error", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		wallet := createTestActiveWallet(testWalletID)
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(wallet, nil)
		st.repo.EXPECT().GetActiveTransactionByCheckoutSessionIDForUpdate(testCtxTx, session.ID).
			Return(nil, assert.AnError)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				err := fn(testCtxTx)
				assert.Error(t, err)
				assert.Equal(t, entity.ErrInternal, err)
				return err
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("transaction is already completed, idempotent no-op", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		wallet := createTestActiveWallet(testWalletID)
		trx := createTestActiveTransaction(entity.TransactionStatusCompleted)
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(wallet, nil)
		st.repo.EXPECT().GetActiveTransactionByCheckoutSessionIDForUpdate(testCtxTx, session.ID).
			Return(trx, nil)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				return fn(testCtxTx)
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.NoError(t, err)
	})

	t.Run("transaction status is neither pending nor completed", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		wallet := createTestActiveWallet(testWalletID)
		trx := createTestActiveTransaction(entity.TransactionStatus("failed"))
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(wallet, nil)
		st.repo.EXPECT().GetActiveTransactionByCheckoutSessionIDForUpdate(testCtxTx, session.ID).
			Return(trx, nil)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				err := fn(testCtxTx)
				assert.Error(t, err)
				assert.Equal(t, entity.ErrTransactionUnprocessable, err)
				return err
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("update transaction to completed returns error", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		wallet := createTestActiveWallet(testWalletID)
		trx := createTestActiveTransaction(entity.TransactionStatusPending)
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(wallet, nil)
		st.repo.EXPECT().GetActiveTransactionByCheckoutSessionIDForUpdate(testCtxTx, session.ID).
			Return(trx, nil)
		st.repo.EXPECT().UpdateActiveTransactionToCompletedByCheckoutSessionID(testCtxTx, session.ID).
			Return(assert.AnError)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				err := fn(testCtxTx)
				assert.Error(t, err)
				assert.Equal(t, entity.ErrInternal, err)
				return err
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("add wallet balance returns error", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		wallet := createTestActiveWallet(testWalletID)
		trx := createTestActiveTransaction(entity.TransactionStatusPending)
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(wallet, nil)
		st.repo.EXPECT().GetActiveTransactionByCheckoutSessionIDForUpdate(testCtxTx, session.ID).
			Return(trx, nil)
		st.repo.EXPECT().UpdateActiveTransactionToCompletedByCheckoutSessionID(testCtxTx, session.ID).
			Return(nil)
		st.repo.EXPECT().AddActiveWalletBalance(testCtxTx, wallet.ID, trx.Amount).
			Return(nil, assert.AnError)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				err := fn(testCtxTx)
				assert.Error(t, err)
				assert.Equal(t, entity.ErrInternal, err)
				return err
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("tx manager fails even though business logic succeeded", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		wallet := createTestActiveWallet(testWalletID)
		trx := createTestActiveTransaction(entity.TransactionStatusPending)
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(wallet, nil)
		st.repo.EXPECT().GetActiveTransactionByCheckoutSessionIDForUpdate(testCtxTx, session.ID).
			Return(trx, nil)
		st.repo.EXPECT().UpdateActiveTransactionToCompletedByCheckoutSessionID(testCtxTx, session.ID).
			Return(nil)
		st.repo.EXPECT().AddActiveWalletBalance(testCtxTx, wallet.ID, trx.Amount).
			Return(wallet, nil)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.NoError(t, fn(testCtxTx))
				return assert.AnError
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("success complete pending transaction and add wallet balance", func(t *testing.T) {
		st := createStripeEventHandlerSuite(t)
		session := createTestPaidCheckoutSession(testWalletID.String())
		wallet := createTestActiveWallet(testWalletID)
		trx := createTestActiveTransaction(entity.TransactionStatusPending)
		st.repo.EXPECT().GetActiveWalletByIDForUpdate(testCtxTx, testWalletID).
			Return(wallet, nil)
		st.repo.EXPECT().GetActiveTransactionByCheckoutSessionIDForUpdate(testCtxTx, session.ID).
			Return(trx, nil)
		st.repo.EXPECT().UpdateActiveTransactionToCompletedByCheckoutSessionID(testCtxTx, session.ID).
			Return(nil)
		st.repo.EXPECT().AddActiveWalletBalance(testCtxTx, wallet.ID, trx.Amount).
			Return(wallet, nil)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				return fn(testCtxTx)
			})

		err := st.eventHandler.HandleCheckoutSessionCompleted(testCtx, session)

		assert.NoError(t, err)
	})
}

func createStripeEventHandlerSuite(t *testing.T) *StripeEventHandlerSuite {
	m := uow.NewMockTxManager(t)
	r := mockservice.NewMockHandleStripeEventRepository(t)
	h := service.NewStripeEventHandler(m, r)

	return &StripeEventHandlerSuite{
		eventHandler: h,
		txManager:    m,
		repo:         r,
	}
}

func createTestCheckoutSession(status stripe.CheckoutSessionPaymentStatus) *stripe.CheckoutSession {
	return &stripe.CheckoutSession{
		ID:            testCheckoutSessionID,
		PaymentStatus: status,
	}
}

func createTestPaidCheckoutSession(clientReferenceID string) *stripe.CheckoutSession {
	return &stripe.CheckoutSession{
		ID:                testCheckoutSessionID,
		PaymentStatus:     stripe.CheckoutSessionPaymentStatusPaid,
		ClientReferenceID: clientReferenceID,
	}
}

func createTestActiveWallet(id uuid.UUID) *entity.Wallet {
	now := time.Now().UTC()
	return &entity.Wallet{
		ID:       id,
		UserID:   testUserID,
		Currency: testCurrency,
		Balance:  decimal.Zero,
		Auditable: entity.Auditable{
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: testUserID,
			UpdatedBy: testUserID,
		},
	}
}

func createTestActiveTransaction(status entity.TransactionStatus) *entity.Transaction {
	now := time.Now().UTC()
	sessionID := testCheckoutSessionID
	return &entity.Transaction{
		ID:                uuid.Must(uuid.NewV7()),
		UserID:            testUserID,
		Type:              entity.TransactionTypeTopup,
		Status:            status,
		IdempotencyKey:    uuid.Must(uuid.NewV7()),
		Amount:            decimal.NewFromInt(100),
		Currency:          testCurrency,
		CheckoutSessionID: &sessionID,
		Auditable: entity.Auditable{
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: testUserID,
			UpdatedBy: testUserID,
		},
	}
}
