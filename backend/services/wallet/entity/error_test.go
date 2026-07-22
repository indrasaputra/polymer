package entity_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

func TestNewWalletError(t *testing.T) {
	t.Run("success create an instance of WalletError", func(t *testing.T) {
		err := entity.NewWalletError(500, "error")

		assert.NotNil(t, err)
	})
}

func TestWalletError_Error(t *testing.T) {
	t.Run("error returns error message", func(t *testing.T) {
		err := entity.NewWalletError(500, "error message")

		assert.Equal(t, "error message", err.Error())
	})
}

func TestErrWalletEmpty(t *testing.T) {
	t.Run("empty wallet error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrWalletEmpty.Code)
		assert.Equal(t, "wallet is empty or nil", entity.ErrWalletEmpty.Error())
	})
}

func TestErrUserEmpty(t *testing.T) {
	t.Run("empty user error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrUserEmpty.Code)
		assert.Equal(t, "user is empty or nil", entity.ErrUserEmpty.Error())
	})
}

func TestErrGeneralInvalid(t *testing.T) {
	t.Run("general invalid returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrGeneralInvalid.Code)
		assert.Equal(t, "request, param, or instance's value is invalid", entity.ErrGeneralInvalid.Error())
	})
}

func TestErrInternal(t *testing.T) {
	t.Run("internal error returns 500 code", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, entity.ErrInternal.Code)
		assert.Equal(t, "internal error", entity.ErrInternal.Error())
	})
}

func TestErrIdempotencyKeyEmpty(t *testing.T) {
	t.Run("invalid idempotency key error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrIdempotencyKeyEmpty.Code)
		assert.Equal(t, "idempotency key is empty or nil or not UUID", entity.ErrIdempotencyKeyEmpty.Error())
	})
}

func TestErrTopupAmountInvalid(t *testing.T) {
	t.Run("invalid topup amount error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrTopupAmountInvalid.Code)
		assert.Equal(t, "topup amount is invalid", entity.ErrTopupAmountInvalid.Error())
	})
}

func TestErrCurrencyInvalid(t *testing.T) {
	t.Run("invalid currency error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrCurrencyInvalid.Code)
		assert.Equal(t, "currency is invalid", entity.ErrCurrencyInvalid.Error())
	})
}

func TestErrCustomerNotFound(t *testing.T) {
	t.Run("nil customer error returns 404 code", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, entity.ErrCustomerNotFound.Code)
		assert.Equal(t, "customer not found", entity.ErrCustomerNotFound.Error())
	})
}

func TestErrCustomerEmpty(t *testing.T) {
	t.Run("empty customer error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrCustomerEmpty.Code)
		assert.Equal(t, "customer is empty or nil", entity.ErrCustomerEmpty.Error())
	})
}

func TestErrTransactionNotFound(t *testing.T) {
	t.Run("nil transaction error returns 404 code", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, entity.ErrTransactionNotFound.Code)
		assert.Equal(t, "transaction not found", entity.ErrTransactionNotFound.Error())
	})
}

func TestErrTransactionEmpty(t *testing.T) {
	t.Run("empty transaction error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrTransactionEmpty.Code)
		assert.Equal(t, "transaction is empty or nil", entity.ErrTransactionEmpty.Error())
	})
}

func TestErrWalletNotFound(t *testing.T) {
	t.Run("nil wallet error returns 404 code", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, entity.ErrWalletNotFound.Code)
		assert.Equal(t, "wallet not found", entity.ErrWalletNotFound.Error())
	})
}

func TestErrTopupEmpty(t *testing.T) {
	t.Run("empty topup error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrTopupEmpty.Code)
		assert.Equal(t, "topup is empty or nil", entity.ErrTopupEmpty.Error())
	})
}

func TestErrStripeEventInvalid(t *testing.T) {
	t.Run("invalid stripe event error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrStripeEventInvalid.Code)
		assert.Equal(t, "stripe event is invalid", entity.ErrStripeEventInvalid.Error())
	})
}

func TestErrTransactionUnprocessable(t *testing.T) {
	t.Run("invalid transaction error returns 422 code", func(t *testing.T) {
		assert.Equal(t, http.StatusUnprocessableEntity, entity.ErrTransactionUnprocessable.Code)
		assert.Equal(t, "transaction is unprocessable", entity.ErrTransactionUnprocessable.Error())
	})
}
