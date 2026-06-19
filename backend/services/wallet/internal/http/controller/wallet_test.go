package controller_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/controller"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/validator"
	mockservice "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/service"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testUserID      = uuid.Must(uuid.NewV7())
	testUserEmail   = "test.user@mail.com"
	testCurrency    = "USD"
	testCurrentUser = &entity.CurrentUser{
		ID:    testUserID,
		Email: testUserEmail,
	}
	testWallet = &entity.Wallet{
		ID:       uuid.Must(uuid.NewV7()),
		UserID:   testUserID,
		Currency: testCurrency,
		Balance:  decimal.Zero,
		Auditable: entity.Auditable{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			CreatedBy: testUserID,
			UpdatedBy: testUserID,
		},
	}
)

type WalletSuite struct {
	walletController *controller.Wallet
	walletService    *mockservice.MockCreateWallet
}

func TestNewWallet(t *testing.T) {
	t.Run("success create an instace of Wallet", func(t *testing.T) {
		st := createWalletSuite(t)

		assert.NotNil(t, st)
	})
}

func TestWallet_RegisterRoute(t *testing.T) {
	t.Run("success register route", func(t *testing.T) {
		st := createWalletSuite(t)
		e := echo.New()

		assert.NotPanics(t, func() { st.walletController.RegisterRoute(e.Group("/api/v1")) })
	})
}

func TestWallet_Create(t *testing.T) {
	t.Run("binding fail due to invalid json body", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
			},
			JSONBody: []byte(`{"bad":"json"`),
		}.ToContextRecorder(t)

		st := createWalletSuite(t)

		err := st.walletController.Create(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation fail due to non-exist required field", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
			},
			JSONBody: []byte(`{"good":"json"}`),
		}.ToContextRecorder(t)
		(*c).Echo().Validator = validator.New()

		st := createWalletSuite(t)

		err := st.walletController.Create(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("wallet service returns error", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
			},
			JSONBody: []byte(`{"currency":"` + testCurrency + `"}`),
		}.ToContextRecorder(t)
		(*c).Echo().Validator = validator.New()

		st := createWalletSuite(t)
		st.walletService.EXPECT().Create(c.Request().Context(), mock.MatchedBy(func(input *entity.CreateWalletInput) bool {
			return input.UserID == testCurrentUser.ID && input.Currency == testCurrency
		})).Return(nil, assert.AnError)

		err := st.walletController.Create(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success create a walelt", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
			},
			JSONBody: []byte(`{"currency":"` + testCurrency + `"}`),
		}.ToContextRecorder(t)
		(*c).Echo().Validator = validator.New()

		st := createWalletSuite(t)
		st.walletService.EXPECT().Create(c.Request().Context(), mock.MatchedBy(func(input *entity.CreateWalletInput) bool {
			return input.UserID == testCurrentUser.ID && input.Currency == testCurrency
		})).Return(testWallet, nil)

		err := st.walletController.Create(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func createWalletSuite(t *testing.T) *WalletSuite {
	s := mockservice.NewMockCreateWallet(t)
	w := controller.NewWallet(s)
	return &WalletSuite{
		walletController: w,
		walletService:    s,
	}
}
