package dto

import (
	"net/http"

	"github.com/labstack/echo/v5"

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

// CreateWalletRequest defines JSON form of create wallet API.
type CreateWalletRequest struct {
	Currency string `json:"currency" validate:"required,alpha,len=3,uppercase"`
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
