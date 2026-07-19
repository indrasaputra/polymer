package controller_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/controller"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/validator"
	mockservice "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/service"
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
	walletCreator    *mockservice.MockCreateWallet
	walletTopup      *mockservice.MockTopupWallet
}

func TestNewWallet(t *testing.T) {
	t.Run("success create an instance of Wallet", func(t *testing.T) {
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
		st.walletCreator.EXPECT().Create(c.Request().Context(), mock.MatchedBy(func(input *entity.CreateWalletInput) bool {
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
		st.walletCreator.EXPECT().Create(c.Request().Context(), mock.MatchedBy(func(input *entity.CreateWalletInput) bool {
			return input.UserID == testCurrentUser.ID && input.Currency == testCurrency
		})).Return(testWallet, nil)

		err := st.walletController.Create(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestWallet_Topup(t *testing.T) {
	testWalletID := uuid.Must(uuid.NewV7())
	testIdempotencyKey := uuid.Must(uuid.NewV7())

	t.Run("binding header fail due to invalid idempotency key", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				"x-idempotency-key":    {"not-a-uuid"},
			},
			JSONBody: []byte(`{"amount":"100","wallet_id":"` + testWalletID.String() + `"}`),
		}.ToContextRecorder(t)

		st := createWalletSuite(t)

		err := st.walletController.Topup(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("binding body fail due to invalid json body", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				"x-idempotency-key":    {testIdempotencyKey.String()},
			},
			JSONBody: []byte(`{"bad":"json"`),
		}.ToContextRecorder(t)

		st := createWalletSuite(t)

		err := st.walletController.Topup(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation fail due to non-exist required field", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				"x-idempotency-key":    {testIdempotencyKey.String()},
			},
			JSONBody: []byte(`{"good":"json"}`),
		}.ToContextRecorder(t)
		(*c).Echo().Validator = validator.New()

		st := createWalletSuite(t)

		err := st.walletController.Topup(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("amount is not a valid numeric string", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				"x-idempotency-key":    {testIdempotencyKey.String()},
			},
			JSONBody: []byte(`{"amount":"abc","wallet_id":"` + testWalletID.String() + `"}`),
		}.ToContextRecorder(t)
		(*c).Echo().Validator = validator.New()

		st := createWalletSuite(t)

		err := st.walletController.Topup(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("wallet service returns error", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				"x-idempotency-key":    {testIdempotencyKey.String()},
			},
			JSONBody: []byte(`{"amount":"100","wallet_id":"` + testWalletID.String() + `"}`),
		}.ToContextRecorder(t)
		(*c).Echo().Validator = validator.New()

		st := createWalletSuite(t)
		amount, _ := decimal.NewFromString("100")
		st.walletTopup.EXPECT().Topup(c.Request().Context(), mock.MatchedBy(func(input *entity.TopupWalletInput) bool {
			return input.UserID == testCurrentUser.ID &&
				input.WalletID == testWalletID &&
				input.IdempotencyKey == testIdempotencyKey &&
				input.Amount.Equal(amount)
		})).Return(nil, assert.AnError)

		err := st.walletController.Topup(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success topup wallet", func(t *testing.T) {
		c, rec := echotest.ContextConfig{
			Headers: map[string][]string{
				echo.HeaderContentType: {echo.MIMEApplicationJSON},
				"x-idempotency-key":    {testIdempotencyKey.String()},
			},
			JSONBody: []byte(`{"amount":"100","wallet_id":"` + testWalletID.String() + `"}`),
		}.ToContextRecorder(t)
		(*c).Echo().Validator = validator.New()

		st := createWalletSuite(t)
		amount, _ := decimal.NewFromString("100")
		st.walletTopup.EXPECT().Topup(c.Request().Context(), mock.MatchedBy(func(input *entity.TopupWalletInput) bool {
			return input.UserID == testCurrentUser.ID &&
				input.WalletID == testWalletID &&
				input.IdempotencyKey == testIdempotencyKey &&
				input.Amount.Equal(amount)
		})).Return(&entity.TopupWalletOutput{CheckoutSessionURL: "https://checkout.example.com/session/123"}, nil)

		err := st.walletController.Topup(c, testCurrentUser)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func createWalletSuite(t *testing.T) *WalletSuite {
	c := mockservice.NewMockCreateWallet(t)
	tp := mockservice.NewMockTopupWallet(t)

	w := controller.NewWallet(c, tp)
	return &WalletSuite{
		walletController: w,
		walletCreator:    c,
		walletTopup:      tp,
	}
}
