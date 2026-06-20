package middleware

import (
	"fmt"
	"net/http"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/config"
)

// AlgorithmES256 equals to ES256.
const AlgorithmES256 = "ES256"

// NewJwtMiddleware initializes the JWKS cache and returns the Echo JWT middleware.
func NewJwtMiddleware(cfg *config.Config) (echo.MiddlewareFunc, error) {
	kf, err := keyfunc.NewDefault([]string{cfg.Supabase.JwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS: %v", err)
	}

	jwtConfig := echojwt.Config{
		KeyFunc:        kf.Keyfunc,
		SuccessHandler: jwtSuccessHandler,
	}

	return echojwt.WithConfig(jwtConfig), nil
}

// jwtSuccessHandler extracts the parsed JWT, validates claims, and attaches CurrentUser to context.
func jwtSuccessHandler(c *echo.Context) error {
	// extract token parsed by echo-jwt (stored under context key "user" by default)
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token instance")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims structure")
	}

	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)

	if sub == "" || email == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token payload")
	}

	id, err := uuid.Parse(sub)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
	}

	currentUser := &entity.CurrentUser{
		ID:    id,
		Email: email,
	}
	c.Set(entity.ContextKeyCurrentUser, currentUser)

	return nil
}

func getCurrentUser(c *echo.Context) (*entity.CurrentUser, error) {
	u, ok := c.Get(entity.ContextKeyCurrentUser).(*entity.CurrentUser)
	if !ok || u == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "No current user in context")
	}
	return u, nil
}

// RequireUser extracts user from current context and pass it to handler.
func RequireUser(handler func(c *echo.Context, user *entity.CurrentUser) error) func(c *echo.Context) error {
	return func(c *echo.Context) error {
		user, err := getCurrentUser(c)
		if err != nil {
			return err
		}
		return handler(c, user)
	}
}
