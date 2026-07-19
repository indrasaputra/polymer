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

func TestErrEmptyWallet(t *testing.T) {
	t.Run("empty wallet error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrEmptyWallet.Code)
		assert.Equal(t, "Wallet is empty or nil", entity.ErrEmptyWallet.Error())
	})
}

func TestErrInvalidUser(t *testing.T) {
	t.Run("invalid user error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrInvalidUser.Code)
		assert.Equal(t, "User is invalid", entity.ErrInvalidUser.Error())
	})
}

func TestErrBadRequest(t *testing.T) {
	t.Run("bad request returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrBadRequest.Code)
		assert.Equal(t, "Bad request in body or param", entity.ErrBadRequest.Error())
	})
}

func TestErrInternal(t *testing.T) {
	t.Run("internal error returns 500 code", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, entity.ErrInternal.Code)
		assert.Equal(t, "Internal error", entity.ErrInternal.Error())
	})
}

func TestErrInvalidWallet(t *testing.T) {
	t.Run("invalid wallet error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrInvalidWallet.Code)
		assert.Equal(t, "Wallet is invalid", entity.ErrInvalidWallet.Error())
	})
}

func TestErrInvalidIdempotencyKey(t *testing.T) {
	t.Run("invalid idempotency key error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrInvalidIdempotencyKey.Code)
		assert.Equal(t, "Idempotency key is invalid", entity.ErrInvalidIdempotencyKey.Error())
	})
}

func TestErrInvalidTopupAmount(t *testing.T) {
	t.Run("invalid topup amount error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrInvalidTopupAmount.Code)
		assert.Equal(t, "Invalid topup amount", entity.ErrInvalidTopupAmount.Error())
	})
}

func TestErrInvalidCurrency(t *testing.T) {
	t.Run("invalid currency error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrInvalidCurrency.Code)
		assert.Equal(t, "Invalid currency", entity.ErrInvalidCurrency.Error())
	})
}

func TestErrNilCustomer(t *testing.T) {
	t.Run("nil customer error returns 404 code", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, entity.ErrNilCustomer.Code)
		assert.Equal(t, "Customer not found", entity.ErrNilCustomer.Error())
	})
}

func TestErrNilTransaction(t *testing.T) {
	t.Run("nil transaction error returns 404 code", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, entity.ErrNilTransaction.Code)
		assert.Equal(t, "Transaction not found", entity.ErrNilTransaction.Error())
	})
}

func TestErrNilWallet(t *testing.T) {
	t.Run("nil wallet error returns 404 code", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, entity.ErrNilWallet.Code)
		assert.Equal(t, "Wallet not found", entity.ErrNilWallet.Error())
	})
}

func TestErrEmptyInput(t *testing.T) {
	t.Run("empty input error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrEmptyInput.Code)
		assert.Equal(t, "Input is empty", entity.ErrEmptyInput.Error())
	})
}

func TestErrEmptyTopup(t *testing.T) {
	t.Run("empty topup error returns 400 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrEmptyTopup.Code)
		assert.Equal(t, "Empty topup", entity.ErrEmptyTopup.Error())
	})
}
