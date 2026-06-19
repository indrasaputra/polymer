package controller

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/dto"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/http/middleware"
)

// Wallet handles HTTP request for wallet.
type Wallet struct {
	creator service.CreateWallet
}

// NewWallet creates an instance of Wallet.
func NewWallet(c service.CreateWallet) *Wallet {
	return &Wallet{creator: c}
}

// RegisterRoute registers all routes in wallet controller.
func (w *Wallet) RegisterRoute(g *echo.Group) {
	wg := g.Group("/wallets")

	wg.POST("", middleware.RequireUser(w.Create))
}

// Create creates a new wallet with idempotency in mind.
func (w *Wallet) Create(c *echo.Context, currentUser *entity.CurrentUser) error {
	var request dto.CreateWalletRequest
	if err := c.Bind(&request); err != nil {
		log.Printf("bind err: %v\n", err)
		return dto.SendResponse(c, nil, entity.ErrEmptyWallet, 0)
	}

	if err := c.Validate(request); err != nil {
		log.Printf("validate err: %v\n", err)
		return dto.SendResponse(c, nil, entity.ErrEmptyWallet, 0)
	}

	log.Printf("ddcdcef")

	input := &entity.CreateWalletInput{UserID: currentUser.ID, Currency: request.Currency}
	result, err := w.creator.Create(c.Request().Context(), input)
	log.Printf("service err %v\n", err)
	return dto.SendResponse(c, result, err, http.StatusCreated)
}
