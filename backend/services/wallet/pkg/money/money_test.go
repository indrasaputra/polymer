package money_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/money"
)

func TestToSubunits(t *testing.T) {
	t.Run("USD", func(t *testing.T) {
		amount, err := decimal.NewFromString("123.45")
		assert.NoError(t, err)

		res, err := money.ToSubunits(amount, "USD")

		assert.NoError(t, err)
		assert.Equal(t, int64(12345), res)
	})

	t.Run("JPY", func(t *testing.T) {
		amount, err := decimal.NewFromString("123")
		assert.NoError(t, err)

		res, err := money.ToSubunits(amount, "JPY")

		assert.NoError(t, err)
		assert.Equal(t, int64(123), res)
	})

	t.Run("BHD", func(t *testing.T) {
		amount, err := decimal.NewFromString("123.456")
		assert.NoError(t, err)

		res, err := money.ToSubunits(amount, "BHD")

		assert.NoError(t, err)
		assert.Equal(t, int64(123456), res)
	})

	t.Run("invalid", func(t *testing.T) {
		amount, err := decimal.NewFromString("123.45")
		assert.NoError(t, err)

		res, err := money.ToSubunits(amount, "XXX")

		assert.Error(t, err)
		assert.Zero(t, res)
	})
}
