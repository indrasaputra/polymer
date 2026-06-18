package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	mockservice "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/service"
)

var (
	testCtx      = context.Background()
	testUserID   = uuid.Must(uuid.NewV7())
	testCurrency = "USD"
)

type WalletCreatorSuite struct {
	walletService *service.WalletCreator
	walletRepo    *mockservice.MockCreateWalletRepository
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

		err := st.walletService.Create(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet, err)
	})

	t.Run("user id is invalid", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		input.UserID = uuid.Nil

		err := st.walletService.Create(testCtx, input)

		assert.Error(t, err)
	})

	t.Run("wallet repo insert returns error", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		st.walletRepo.EXPECT().Insert(testCtx, mock.MatchedBy(func(wallet *entity.Wallet) bool {
			return wallet.UserID == input.UserID && wallet.Currency == input.Currency && wallet.ID.String() != ""
		})).Return(assert.AnError)

		err := st.walletService.Create(testCtx, input)

		assert.Error(t, err)
	})

	t.Run("success create a wallet", func(t *testing.T) {
		st := createWalletCreatorSuite(t)
		input := createCreateWalletInput()
		st.walletRepo.EXPECT().Insert(testCtx, mock.MatchedBy(func(wallet *entity.Wallet) bool {
			return wallet.UserID == input.UserID && wallet.Currency == input.Currency && wallet.ID.String() != ""
		})).Return(nil)

		err := st.walletService.Create(testCtx, input)

		assert.NoError(t, err)
	})
}

func createWalletCreatorSuite(t *testing.T) *WalletCreatorSuite {
	r := mockservice.NewMockCreateWalletRepository(t)
	w := service.NewWalletCreator(r)
	return &WalletCreatorSuite{
		walletService: w,
		walletRepo:    r,
	}
}

func createCreateWalletInput() *entity.CreateWalletInput {
	return &entity.CreateWalletInput{
		UserID:   testUserID,
		Currency: testCurrency,
	}
}
