package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	"github.com/indrasaputra/polymer/backend/services/wallet/test/mock/pkg/sdk/uow"
	mockservice "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/service"
)

type ctxKey string

var (
	testCtx        = context.Background()
	testCtxTx      = context.WithValue(testCtx, ctxKey("tx"), true)
	testUserID     = uuid.Must(uuid.NewV7())
	testCurrency   = "USD"
	testCustomerID = "cus_NffrFeUfNV2Hib"
)

type WalletCreatorSuite struct {
	walletService  *service.WalletCreator
	txManager      *uow.MockTxManager
	walletRepo     *mockservice.MockCreateWalletRepository
	customerClient *mockservice.MockCreateCustomerClient
}

func TestNewWalletCreator(t *testing.T) {
	t.Run("successfully create an instance of Wallet", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		assert.NotNil(t, st.walletService)
	})
}

func TestWalletCreator_Create(t *testing.T) {
	t.Run("empty wallet is prohibited", func(t *testing.T) {
		st := createWalletCreatorSuite(t)

		res, err := st.walletService.Create(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet, err)
		assert.Nil(t, res)
	})

	t.Run("user id is invalid", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		input.UserID = uuid.Nil

		res, err := st.walletService.Create(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("currency is invalid", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		input.Currency = "XXX"

		res, err := st.walletService.Create(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("get customer returns error", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).Return(nil, assert.AnError)

		res, err := st.walletService.Create(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("create customer client returns error", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).Return(nil, entity.ErrNilCustomer)
		st.customerClient.EXPECT().CreateCustomer(testCtx, input.Email).Return("", assert.AnError)

		res, err := st.walletService.Create(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("insert customer returns error", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).Return(nil, entity.ErrNilCustomer)
		st.customerClient.EXPECT().CreateCustomer(testCtx, input.Email).Return(testCustomerID, nil)
		st.walletRepo.EXPECT().InsertCustomer(testCtxTx, mock.MatchedBy(func(customer *entity.Customer) bool {
			return customer.StripeCustomerID == testCustomerID
		})).Return(nil, assert.AnError)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return assert.AnError
			})

		res, err := st.walletService.Create(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("insert wallet returns error", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).Return(nil, entity.ErrNilCustomer)
		st.customerClient.EXPECT().CreateCustomer(testCtx, input.Email).Return(testCustomerID, nil)
		st.walletRepo.EXPECT().InsertCustomer(testCtxTx, mock.MatchedBy(func(customer *entity.Customer) bool {
			return customer.StripeCustomerID == testCustomerID
		})).Return(nil, nil)
		st.walletRepo.EXPECT().InsertWallet(testCtxTx, mock.MatchedBy(func(wallet *entity.Wallet) bool {
			return wallet.UserID == input.UserID && wallet.Currency == input.Currency && wallet.ID.String() != ""
		})).Return(nil, assert.AnError)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return assert.AnError
			})

		res, err := st.walletService.Create(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("tx manager returns error", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		output := createWallet()
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).Return(nil, entity.ErrNilCustomer)
		st.customerClient.EXPECT().CreateCustomer(testCtx, input.Email).Return(testCustomerID, nil)
		st.walletRepo.EXPECT().InsertCustomer(testCtxTx, mock.MatchedBy(func(customer *entity.Customer) bool {
			return customer.StripeCustomerID == testCustomerID
		})).Return(nil, nil)
		st.walletRepo.EXPECT().InsertWallet(testCtxTx, mock.MatchedBy(func(wallet *entity.Wallet) bool {
			return wallet.UserID == input.UserID && wallet.Currency == input.Currency && wallet.ID.String() != ""
		})).Return(output, nil)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.NoError(t, fn(testCtxTx))
				return assert.AnError
			})

		res, err := st.walletService.Create(testCtx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success create wallet", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		output := createWallet()
		st.walletRepo.EXPECT().GetActiveCustomerByUserID(testCtx, input.UserID).Return(nil, entity.ErrNilCustomer)
		st.customerClient.EXPECT().CreateCustomer(testCtx, input.Email).Return(testCustomerID, nil)
		st.walletRepo.EXPECT().InsertCustomer(testCtxTx, mock.MatchedBy(func(customer *entity.Customer) bool {
			return customer.StripeCustomerID == testCustomerID
		})).Return(nil, nil)
		st.walletRepo.EXPECT().InsertWallet(testCtxTx, mock.MatchedBy(func(wallet *entity.Wallet) bool {
			return wallet.UserID == input.UserID && wallet.Currency == input.Currency && wallet.ID.String() != ""
		})).Return(output, nil)
		st.txManager.EXPECT().Do(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.NoError(t, fn(testCtxTx))
				return nil
			})

		res, err := st.walletService.Create(testCtx, input)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, output.UserID, res.UserID)
		assert.Equal(t, output.Currency, res.Currency)
	})
}

func createWalletCreatorSuite(t *testing.T) *WalletCreatorSuite {
	m := uow.NewMockTxManager(t)
	r := mockservice.NewMockCreateWalletRepository(t)
	c := mockservice.NewMockCreateCustomerClient(t)
	w := service.NewWalletCreator(m, r, c)
	return &WalletCreatorSuite{
		walletService:  w,
		txManager:      m,
		walletRepo:     r,
		customerClient: c,
	}
}

func createCreateWalletInput() *entity.CreateWalletInput {
	return &entity.CreateWalletInput{
		UserID:   testUserID,
		Currency: testCurrency,
	}
}

func createWallet() *entity.Wallet {
	now := time.Now().UTC()

	return &entity.Wallet{
		ID:       uuid.Must(uuid.NewV7()),
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
