package router

import "github.com/labstack/echo/v5"

// Controller defines interface for route-aware controller.
type Controller interface {
	RegisterRoute(g *echo.Group)
}

// RegisterAPIV1 calls Register func and map them to `/api/v1`.
func RegisterAPIV1(e *echo.Echo, controllers []Controller) {
	apiv1 := e.Group("/api/v1")

	for _, con := range controllers {
		con.RegisterRoute(apiv1)
	}
}
