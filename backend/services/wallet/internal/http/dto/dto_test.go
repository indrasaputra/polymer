package dto_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/dto"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestSendResponse(t *testing.T) {
	t.Run("wallet error entity returns code as defined", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		err := dto.SendResponse(c, nil, entity.ErrEmptyWallet, 0)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("non wallet error returns internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		err := dto.SendResponse(c, nil, assert.AnError, 0)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success returns http status ok by default", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		err := dto.SendResponse(c, nil, nil, 0)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("success returns custom code if provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		err := dto.SendResponse(c, nil, nil, http.StatusCreated)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}
