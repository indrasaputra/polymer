package money

import (
	"github.com/bojanz/currency"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

const (
	ten = 10
)

var (
	decimalTen = decimal.NewFromInt(ten)
)

// ToSubunits converts amount to its smallest currency.
// E.g: USD to cents.
func ToSubunits(amount decimal.Decimal, currencyCode string) (int64, error) {
	digit, ok := currency.GetDigits(currencyCode)
	if !ok {
		return 0, entity.ErrInvalidCurrency
	}

	su := amount.Mul(decimalTen.Pow(decimal.NewFromInt32(int32(digit)))).IntPart()
	return su, nil
}
