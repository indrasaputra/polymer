package entity

import "net/http"

var (
	// ErrWalletEmpty occurs when wallet instance is empty or nil.
	ErrWalletEmpty = NewWalletError(http.StatusBadRequest, "wallet is empty or nil")
	// ErrWalletNotFound occurs when wallet is not found.
	ErrWalletNotFound = NewWalletError(http.StatusNotFound, "wallet not found")

	// ErrTransactionEmpty occurs when transaction instance is empty or nil.
	ErrTransactionEmpty = NewWalletError(http.StatusBadRequest, "transaction is empty or nil")
	// ErrTransactionNotFound occurs when transaction is not found.
	ErrTransactionNotFound = NewWalletError(http.StatusNotFound, "transaction not found")
	// ErrTransactionUnprocessable occurs when transaction can't be processed further.
	ErrTransactionUnprocessable = NewWalletError(http.StatusUnprocessableEntity, "transaction is unprocessable")

	// ErrCustomerEmpty occurs when customer instance is empty or nil.
	ErrCustomerEmpty = NewWalletError(http.StatusBadRequest, "customer is empty or nil")
	// ErrCustomerNotFound occurs when customer is not found.
	ErrCustomerNotFound = NewWalletError(http.StatusNotFound, "customer not found")

	// ErrTopupEmpty occurs when topup is empty or nil.
	ErrTopupEmpty = NewWalletError(http.StatusBadRequest, "topup is empty or nil")
	// ErrTopupAmountInvalid occurs when topup amount is invalid (less than or equal to zero).
	ErrTopupAmountInvalid = NewWalletError(http.StatusBadRequest, "topup amount is invalid")

	// ErrIdempotencyKeyEmpty occurs when idempotency key is empty, nil, or not uuid.
	ErrIdempotencyKeyEmpty = NewWalletError(http.StatusBadRequest, "idempotency key is empty or nil or not UUID")

	// ErrCurrencyInvalid occurs when currency is invalid.
	// It uses https://github.com/bojanz/currency as source of truth.
	ErrCurrencyInvalid = NewWalletError(http.StatusBadRequest, "currency is invalid")

	// ErrUserEmpty occurs when user is empty or nil.
	ErrUserEmpty = NewWalletError(http.StatusBadRequest, "user is empty or nil")

	// ErrStripeEventInvalid occurs when Stripe event is invalid.
	ErrStripeEventInvalid = NewWalletError(http.StatusBadRequest, "stripe event is invalid")

	// ErrGeneralInvalid occurs when request or param is not as expected.
	ErrGeneralInvalid = NewWalletError(http.StatusBadRequest, "request, param, or instance's value is invalid")

	// ErrInternal occurs for any unknown or when server crashes.
	ErrInternal = NewWalletError(http.StatusInternalServerError, "internal error")
)

// WalletError represents wallet-related error.
type WalletError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewWalletError creates an instance of WalletError.
func NewWalletError(code int, message string) *WalletError {
	return &WalletError{
		Code:    code,
		Message: message,
	}
}

// Error returns error's message.
func (w *WalletError) Error() string {
	return w.Message
}
