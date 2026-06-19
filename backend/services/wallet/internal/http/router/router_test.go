package router_test

import (
	"testing"

	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/router"
	mockrouter "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/http/router"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterAPIV1(t *testing.T) {
	t.Run("success register routes", func(t *testing.T) {
		e := echo.New()
		r := mockrouter.NewMockRouteRegistrar(t)
		r.EXPECT().RegisterRoute(mock.Anything).Return()

		assert.NotPanics(t, func() { router.RegisterAPIV1(e, r) })
	})
}
