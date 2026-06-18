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
	t.Run("empty wallet error returns 401 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrEmptyWallet.Code)
		assert.Equal(t, "wallet is empty or nil", entity.ErrEmptyWallet.Error())
	})
}

func TestErrInvalidUser(t *testing.T) {
	t.Run("invalid user error returns 401 code", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, entity.ErrInvalidUser.Code)
		assert.Equal(t, "user is invalid", entity.ErrInvalidUser.Error())
	})
}
