package middleware

// deliberately not using _test suffix because I want to test private function.
// changing private func to public func makes test too complex.

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/config"
)

var (
	testUserID = uuid.Must(uuid.NewV7())
)

func TestNewJwtMiddleware(t *testing.T) {
	t.Run("invalid jwks url returns error", func(t *testing.T) {
		cfg := &config.Config{
			Supabase: config.Supabase{
				JwksURL: "invalid",
			},
		}

		res, err := NewJwtMiddleware(cfg)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success create jwt middleware", func(t *testing.T) {
		cfg := &config.Config{
			Supabase: config.Supabase{
				JwksURL: "http://localhost:54321/auth/v1/.well-known/jwks.json",
			},
		}

		res, err := NewJwtMiddleware(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestJwtSuccessHandler(t *testing.T) {
	e := echo.New()

	t.Run("should return unauthorized when user key is missing", func(t *testing.T) {
		c := createTestContext(e)

		err := jwtSuccessHandler(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	t.Run("should return unauthorized when user is not a *jwt.Token", func(t *testing.T) {
		c := createTestContext(e)
		c.Set("user", "not-a-token")

		err := jwtSuccessHandler(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
		assert.Equal(t, "Invalid token instance", httpErr.Message)
	})

	t.Run("should return unauthorized when claims are not jwt.MapClaims", func(t *testing.T) {
		c := createTestContext(e)
		token := &jwt.Token{
			Claims: &jwt.RegisteredClaims{},
		}
		c.Set("user", token)

		err := jwtSuccessHandler(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
		assert.Equal(t, "Invalid token claims structure", httpErr.Message)
	})

	t.Run("should return unauthorized when sub is missing", func(t *testing.T) {
		c := createTestContext(e)
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"email": "test.user@mail.com",
			},
		}
		c.Set("user", token)

		err := jwtSuccessHandler(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, "Invalid token payload", httpErr.Message)
	})

	t.Run("should return unauthorized when email is missing", func(t *testing.T) {
		c := createTestContext(e)
		userID := testUserID
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"sub": userID.String(),
			},
		}
		c.Set("user", token)

		err := jwtSuccessHandler(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, "Invalid token payload", httpErr.Message)
	})

	t.Run("should return unauthorized when sub is not a valid UUID", func(t *testing.T) {
		c := createTestContext(e)
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"sub":   "not-a-uuid",
				"email": "test.user@mail.com",
			},
		}
		c.Set("user", token)

		err := jwtSuccessHandler(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, "Invalid user ID", httpErr.Message)
	})

	t.Run("should set CurrentUser when token is valid", func(t *testing.T) {
		c := createTestContext(e)
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"sub":   testUserID.String(),
				"email": "test.user@mail.com",
			},
		}
		c.Set("user", token)

		err := jwtSuccessHandler(c)

		assert.NoError(t, err)
		currentUser, ok := c.Get(entity.ContextKeyCurrentUser).(*entity.CurrentUser)
		assert.True(t, ok)
		assert.Equal(t, testUserID, currentUser.ID)
		assert.Equal(t, "test.user@mail.com", currentUser.Email)
	})
}

func TestGetCurrentUser(t *testing.T) {
	e := echo.New()

	t.Run("should return error when user is not in context", func(t *testing.T) {
		c := createTestContext(e)

		user, err := getCurrentUser(c)

		assert.Error(t, err)
		assert.Nil(t, user)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	t.Run("should return error when user is nil", func(t *testing.T) {
		c := createTestContext(e)
		var nilUser *entity.CurrentUser
		c.Set(entity.ContextKeyCurrentUser, nilUser)

		user, err := getCurrentUser(c)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return error when context value is wrong type", func(t *testing.T) {
		c := createTestContext(e)
		c.Set(entity.ContextKeyCurrentUser, "not-a-user")

		user, err := getCurrentUser(c)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return user when present in context", func(t *testing.T) {
		c := createTestContext(e)
		expected := &entity.CurrentUser{
			ID:    testUserID,
			Email: "test.user@mail.com",
		}
		c.Set(entity.ContextKeyCurrentUser, expected)

		user, err := getCurrentUser(c)

		assert.NoError(t, err)
		assert.Equal(t, expected, user)
	})

}

func TestRequireUser(t *testing.T) {
	e := echo.New()

	t.Run("should return error and not call handler when user missing", func(t *testing.T) {
		c := createTestContext(e)

		called := false
		handler := RequireUser(func(c *echo.Context, user *entity.CurrentUser) error {
			called = true
			return nil
		})

		err := handler(c)

		assert.Error(t, err)
		assert.False(t, called)
	})

	t.Run("should call handler with current user when present", func(t *testing.T) {
		c := createTestContext(e)
		expected := &entity.CurrentUser{
			ID:    testUserID,
			Email: "test.user@mail.com",
		}
		c.Set(entity.ContextKeyCurrentUser, expected)

		called := false
		var receivedUser *entity.CurrentUser
		handler := RequireUser(func(c *echo.Context, user *entity.CurrentUser) error {
			called = true
			receivedUser = user
			return nil
		})

		err := handler(c)

		assert.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, expected, receivedUser)
	})
}

func createTestContext(e *echo.Echo) *echo.Context {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}
