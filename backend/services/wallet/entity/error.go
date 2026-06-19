package entity

import "net/http"

var (
	// ErrEmptyWallet occurs when wallet is nil or empty.
	ErrEmptyWallet = NewWalletError(http.StatusBadRequest, "wallet is empty or nil")
	// ErrInvalidUser occurs when user is invalid.
	ErrInvalidUser = NewWalletError(http.StatusBadRequest, "user is invalid")

	// ErrBadRequest occurs when request is not as expected.
	ErrBadRequest = NewWalletError(http.StatusBadRequest, "bad request in body or param")

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
