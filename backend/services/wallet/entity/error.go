package entity

import "net/http"

var (
	// ErrEmptyWallet occurs when wallet is nil or empty.
	ErrEmptyWallet = NewWalletError(http.StatusBadRequest, "Wallet is empty or nil")
	// ErrInvalidUser occurs when user is invalid.
	ErrInvalidUser = NewWalletError(http.StatusBadRequest, "User is invalid")
	// ErrInvalidWallet occurs when wallet is invalid.
	ErrInvalidWallet = NewWalletError(http.StatusBadRequest, "Wallet is invalid")
	// ErrInvalidIdempotencyKey occurs when idempotency key is invalid.
	ErrInvalidIdempotencyKey = NewWalletError(http.StatusBadRequest, "Idempotency key is invalid")
	// ErrInvalidTopupAmount occurs when topup amount is invalid (less than or equal to zero).
	ErrInvalidTopupAmount = NewWalletError(http.StatusBadRequest, "Invalid topup amount")
	// ErrInvalidCurrency occurs when currency is invalid.
	// It uses https://github.com/bojanz/currency as source of truth.
	ErrInvalidCurrency = NewWalletError(http.StatusBadRequest, "Invalid currency")
	// ErrNilCustomer occurs when customer is not found.
	ErrNilCustomer = NewWalletError(http.StatusNotFound, "Customer not found")
	// ErrNilTransaction occurs when transaction is not found.
	ErrNilTransaction = NewWalletError(http.StatusNotFound, "Transaction not found")
	// ErrEmptyInput occurs when input is empty.
	ErrEmptyInput = NewWalletError(http.StatusBadRequest, "Input is empty")

	// ErrBadRequest occurs when request is not as expected.
	ErrBadRequest = NewWalletError(http.StatusBadRequest, "Bad request in body or param")

	// ErrInternal occurs for any unknown or when server crashes.
	ErrInternal = NewWalletError(http.StatusInternalServerError, "Internal error")
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
