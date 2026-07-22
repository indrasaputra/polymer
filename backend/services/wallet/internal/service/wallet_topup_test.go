package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	mockservice "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/service"
)

var (
	testWalletID       = uuid.Must(uuid.NewV7())
	testIdempotencyKey = uuid.Must(uuid.NewV7())
	testAmount         = decimal.NewFromInt(100)
	testSuccessURL     = "https://example.com/success"
	testCheckoutURL    = "https://checkout.example.com/session/123"
	testSessionID      = "cs_test_123"
)

type WalletTopupSuite struct {
	walletTopup   *service.WalletTopup
	walletRepo    *mockservice.MockTopupWalletRepository
	paymentClient *mockservice.MockPaymentClient
}

func TestNewWalletTopup(t *testing.T) {
	t.Run("successfully create an instance of WalletTopup", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		assert.NotNil(t, st.walletTopup)
	})
}

func TestWalletTopup_Topup(t *testing.T) {
	t.Run("nil input is prohibited", func(t *testing.T) {
		st := createWalletTopupSuite(t)

		res, err := st.walletTopup.Topup(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyInput, err)
		assert.Nil(t, res)
	})

	t.Run("user id is invalid", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		input.UserID = uuid.Nil

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidUser, err)
		assert.Nil(t, res)
	})

	t.Run("wallet id is invalid", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		input.WalletID = uuid.Nil

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidWallet, err)
		assert.Nil(t, res)
	})

	t.Run("idempotency key is invalid", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		input.IdempotencyKey = uuid.Nil

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidIdempotencyKey, err)
		assert.Nil(t, res)
	})

	t.Run("amount is zero", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		input.Amount = decimal.Zero

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidTopupAmount, err)
		assert.Nil(t, res)
	})

	t.Run("amount is negative", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		input.Amount = decimal.NewFromInt(-1)

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidTopupAmount, err)
		assert.Nil(t, res)
	})

	t.Run("get pending transaction returns unexpected error", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		st.walletRepo.EXPECT().GetActivePendingTransactionByIdempotencyKey(testCtx, input.IdempotencyKey).
			Return(nil, assert.AnError)

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("pending transaction exists with checkout session, but get checkout session returns error", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		trx := createTestTransaction()
		st.walletRepo.EXPECT().GetActivePendingTransactionByIdempotencyKey(testCtx, input.IdempotencyKey).
			Return(trx, nil)
		st.paymentClient.EXPECT().GetCheckoutSession(testCtx, *trx.CheckoutSessionID).
			Return(nil, assert.AnError)

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("pending transaction exists with checkout session, success get checkout session", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		trx := createTestTransaction()
		session := createCheckoutSession()
		st.walletRepo.EXPECT().GetActivePendingTransactionByIdempotencyKey(testCtx, input.IdempotencyKey).
			Return(trx, nil)
		st.paymentClient.EXPECT().GetCheckoutSession(testCtx, *trx.CheckoutSessionID).
			Return(session, nil)

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, session.URL, res.CheckoutSessionURL)
	})

	t.Run("get user wallet returns error", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		st.walletRepo.EXPECT().GetActivePendingTransactionByIdempotencyKey(testCtx, input.IdempotencyKey).
			Return(nil, entity.ErrNilTransaction)
		st.walletRepo.EXPECT().GetActiveWalletByIDAndUserID(testCtx, input.WalletID, input.UserID).
			Return(nil, assert.AnError)

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("get customer returns error", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		wallet := createTestWalletForTopup()
		st.walletRepo.EXPECT().GetActivePendingTransactionByIdempotencyKey(testCtx, input.IdempotencyKey).
			Return(nil, entity.ErrNilTransaction)
		st.walletRepo.EXPECT().GetActiveWalletByIDAndUserID(testCtx, input.WalletID, input.UserID).
			Return(wallet, nil)
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).
			Return(nil, assert.AnError)

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("create checkout session returns error", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		wallet := createTestWalletForTopup()
		customer := createTestCustomerForTopup()
		st.walletRepo.EXPECT().GetActivePendingTransactionByIdempotencyKey(testCtx, input.IdempotencyKey).
			Return(nil, entity.ErrNilTransaction)
		st.walletRepo.EXPECT().GetActiveWalletByIDAndUserID(testCtx, input.WalletID, input.UserID).
			Return(wallet, nil)
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).
			Return(customer, nil)
		st.paymentClient.EXPECT().CreateCheckoutSession(testCtx, mock.MatchedBy(func(ci *entity.CheckoutInput) bool {
			return ci.WalletID == input.WalletID &&
				ci.Currency == wallet.Currency &&
				ci.StripeCustomerID == customer.StripeCustomerID &&
				ci.Amount.Equal(input.Amount)
		})).Return(nil, assert.AnError)

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("insert transaction returns error", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		wallet := createTestWalletForTopup()
		customer := createTestCustomerForTopup()
		session := createCheckoutSession()
		st.walletRepo.EXPECT().GetActivePendingTransactionByIdempotencyKey(testCtx, input.IdempotencyKey).
			Return(nil, entity.ErrNilTransaction)
		st.walletRepo.EXPECT().GetActiveWalletByIDAndUserID(testCtx, input.WalletID, input.UserID).
			Return(wallet, nil)
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).
			Return(customer, nil)
		st.paymentClient.EXPECT().CreateCheckoutSession(testCtx, mock.MatchedBy(func(ci *entity.CheckoutInput) bool {
			return ci.WalletID == input.WalletID &&
				ci.Currency == wallet.Currency &&
				ci.StripeCustomerID == customer.StripeCustomerID &&
				ci.Amount.Equal(input.Amount)
		})).Return(session, nil)
		st.walletRepo.EXPECT().InsertTransaction(testCtx, mock.MatchedBy(func(trx *entity.Transaction) bool {
			return trx.UserID == input.UserID &&
				trx.IdempotencyKey == input.IdempotencyKey &&
				trx.Type == entity.TransactionTypeTopup &&
				trx.Status == entity.TransactionStatusPending &&
				trx.Currency == wallet.Currency &&
				trx.CheckoutSessionID != nil && *trx.CheckoutSessionID == session.ID &&
				trx.Amount.Equal(input.Amount)
		})).Return(nil, assert.AnError)

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success topup wallet", func(t *testing.T) {
		st := createWalletTopupSuite(t)
		input := createTopupWalletInput()
		wallet := createTestWalletForTopup()
		customer := createTestCustomerForTopup()
		session := createCheckoutSession()
		insertedTrx := createTestTransaction()
		st.walletRepo.EXPECT().GetActivePendingTransactionByIdempotencyKey(testCtx, input.IdempotencyKey).
			Return(nil, entity.ErrNilTransaction)
		st.walletRepo.EXPECT().GetActiveWalletByIDAndUserID(testCtx, input.WalletID, input.UserID).
			Return(wallet, nil)
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).
			Return(customer, nil)
		st.paymentClient.EXPECT().CreateCheckoutSession(testCtx, mock.MatchedBy(func(ci *entity.CheckoutInput) bool {
			return ci.WalletID == input.WalletID &&
				ci.Currency == wallet.Currency &&
				ci.StripeCustomerID == customer.StripeCustomerID &&
				ci.Amount.Equal(input.Amount)
		})).Return(session, nil)
		st.walletRepo.EXPECT().InsertTransaction(testCtx, mock.MatchedBy(func(trx *entity.Transaction) bool {
			return trx.UserID == input.UserID &&
				trx.IdempotencyKey == input.IdempotencyKey &&
				trx.Type == entity.TransactionTypeTopup &&
				trx.Status == entity.TransactionStatusPending &&
				trx.Currency == wallet.Currency &&
				trx.CheckoutSessionID != nil && *trx.CheckoutSessionID == session.ID &&
				trx.Amount.Equal(input.Amount)
		})).Return(insertedTrx, nil)

		res, err := st.walletTopup.Topup(testCtx, input)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, session.URL, res.CheckoutSessionURL)
	})
}

func createWalletTopupSuite(t *testing.T) *WalletTopupSuite {
	r := mockservice.NewMockTopupWalletRepository(t)
	p := mockservice.NewMockPaymentClient(t)
	w := service.NewWalletTopup(r, p, testSuccessURL)
	return &WalletTopupSuite{
		walletTopup:   w,
		walletRepo:    r,
		paymentClient: p,
	}
}

func createTopupWalletInput() *entity.TopupWalletInput {
	return &entity.TopupWalletInput{
		Amount:         testAmount,
		WalletID:       testWalletID,
		IdempotencyKey: testIdempotencyKey,
		UserID:         testUserID,
	}
}

func createTestTransaction() *entity.Transaction {
	now := time.Now().UTC()
	sessionID := testSessionID
	return &entity.Transaction{
		ID:                uuid.Must(uuid.NewV7()),
		UserID:            testUserID,
		Type:              entity.TransactionTypeTopup,
		Status:            entity.TransactionStatusPending,
		IdempotencyKey:    testIdempotencyKey,
		Amount:            testAmount,
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

func createTestWalletForTopup() *entity.Wallet {
	now := time.Now().UTC()
	return &entity.Wallet{
		ID:       testWalletID,
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

func createTestCustomerForTopup() *entity.Customer {
	now := time.Now().UTC()
	return &entity.Customer{
		ID:               uuid.Must(uuid.NewV7()),
		UserID:           testUserID,
		StripeCustomerID: testCustomerID,
		Auditable: entity.Auditable{
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: testUserID,
			UpdatedBy: testUserID,
		},
	}
}

func createCheckoutSession() *entity.CheckoutSession {
	return &entity.CheckoutSession{
		ID:  testSessionID,
		URL: testCheckoutURL,
	}
}
