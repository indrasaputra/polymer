package router

import "github.com/labstack/echo/v5"

// RouteRegistrar defines interface for route-aware controller.
type RouteRegistrar interface {
	RegisterRoute(g *echo.Group, jwtmid echo.MiddlewareFunc)
}

// RegisterAPIV1 calls Register func and map them to `/api/v1`.
func RegisterAPIV1(e *echo.Echo, jwtmid echo.MiddlewareFunc, controllers ...RouteRegistrar) {
	apiv1 := e.Group("/api/v1")

	for _, ctrl := range controllers {
		ctrl.RegisterRoute(apiv1, jwtmid)
	}
}
