package dto

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

// SuccessResponse defines JSON form of success response.
type SuccessResponse struct {
	Data any `json:"data"`
}

// ErrorResponse defines JSON form of error response.
type ErrorResponse struct {
	Message string `json:"message"`
}

// CreateWalletRequest defines JSON form of create wallet API request.
type CreateWalletRequest struct {
	Currency string `json:"currency" validate:"required,alpha,len=3,uppercase"`
}

// CreateWalletResponse defines JSON form of create wallet API response.
type CreateWalletResponse struct {
	Timestamp
	Balance  decimal.Decimal `json:"balance"`
	Currency string          `json:"currency"`
	ID       uuid.UUID       `json:"id"`
	UserID   uuid.UUID       `json:"user_id"`
}

// TopupWalletRequest defines JSON form of topup wallet API request.
type TopupWalletRequest struct {
	Amount         string    `json:"amount" validate:"required,numeric"`
	WalletID       uuid.UUID `json:"wallet_id" validate:"required,uuid"`
	IdempotencyKey uuid.UUID `header:"x-idempotency-key" validate:"required,uuid"`
}

// TopupWalletResponse defines JSON form of topup wallet API response.
type TopupWalletResponse struct {
	CheckoutSessionURL string `json:"checkout_session_url"`
}

// Timestamp defines logical data related to timestamp.
type Timestamp struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SendResponse sends JSON response using Echo framework.
// It will decide whether to return error response or success response.
func SendResponse(c *echo.Context, data any, err error, customCode int) error {
	if err != nil {
		if we, ok := err.(*entity.WalletError); ok {
			return c.JSON(we.Code, ErrorResponse{Message: we.Error()})
		}
		return c.JSON(http.StatusInternalServerError, entity.ErrInternal)
	}

	code := http.StatusOK
	if customCode > 0 {
		code = customCode
	}
	return c.JSON(code, SuccessResponse{Data: data})
}
