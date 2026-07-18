package controller

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/dto"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/http/middleware"
)

// Wallet handles HTTP request for wallet.
type Wallet struct {
	creator service.CreateWallet
	topup   service.TopupWallet
}

// NewWallet creates an instance of Wallet.
func NewWallet(c service.CreateWallet, t service.TopupWallet) *Wallet {
	return &Wallet{creator: c, topup: t}
}

// RegisterRoute registers all routes in wallet controller.
func (w *Wallet) RegisterRoute(g *echo.Group) {
	wg := g.Group("/wallets")

	wg.POST("", middleware.RequireUser(w.Create))
	wg.POST("/topup", middleware.RequireUser(w.Topup))
}

// Create creates a new wallet with idempotency in mind.
func (w *Wallet) Create(c *echo.Context, currentUser *entity.CurrentUser) error {
	var request dto.CreateWalletRequest
	if err := c.Bind(&request); err != nil {
		return dto.SendResponse(c, nil, entity.ErrEmptyWallet, 0)
	}

	if err := c.Validate(request); err != nil {
		log.Println(err)
		return dto.SendResponse(c, nil, entity.ErrEmptyWallet, 0)
	}

	input := &entity.CreateWalletInput{UserID: currentUser.ID, Currency: request.Currency, Email: currentUser.Email}
	result, err := w.creator.Create(c.Request().Context(), input)
	response := convertCreateWalletOutputToCreateWalletResponse(result)
	return dto.SendResponse(c, response, err, http.StatusCreated)
}

// Topup top-ups wallet.
func (w *Wallet) Topup(c *echo.Context, currentUser *entity.CurrentUser) error {
	var request dto.TopupWalletRequest
	if err := echo.BindHeaders(c, &request); err != nil {
		return dto.SendResponse(c, nil, entity.ErrEmptyTopup, 0)
	}
	if err := c.Bind(&request); err != nil {
		return dto.SendResponse(c, nil, entity.ErrEmptyTopup, 0)
	}

	if err := c.Validate(request); err != nil {
		log.Println(err)
		return dto.SendResponse(c, nil, entity.ErrEmptyTopup, 0)
	}

	amount, err := decimal.NewFromString(request.Amount)
	if err != nil {
		return dto.SendResponse(c, nil, entity.ErrEmptyTopup, http.StatusBadRequest)
	}

	input := &entity.TopupWalletInput{
		Amount:         amount,
		WalletID:       request.WalletID,
		IdempotencyKey: request.IdempotencyKey,
		UserID:         currentUser.ID,
	}
	result, err := w.topup.Topup(c.Request().Context(), input)
	response := convertTopupWalletOutpuToTopupWalletResponse(result)
	return dto.SendResponse(c, response, err, http.StatusCreated)
}

func convertCreateWalletOutputToCreateWalletResponse(res *entity.Wallet) *dto.CreateWalletResponse {
	if res == nil {
		return nil
	}
	return &dto.CreateWalletResponse{
		ID:       res.ID,
		Balance:  res.Balance,
		UserID:   res.UserID,
		Currency: res.Currency,
		Timestamp: dto.Timestamp{
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
	}
}

func convertTopupWalletOutpuToTopupWalletResponse(res *entity.TopupWalletOutput) *dto.TopupWalletResponse {
	if res == nil {
		return nil
	}
	return &dto.TopupWalletResponse{
		CheckoutSessionURL: res.CheckoutSessionURL,
	}
}
