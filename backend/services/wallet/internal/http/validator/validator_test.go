package validator_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/validator"
)

var (
	testName  = "Test User"
	testEmail = "test.user@mail.com"
	testAge   = 30
)

type testStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}

func TestNew(t *testing.T) {
	t.Run("success create an instance of Validator", func(t *testing.T) {
		v := validator.New()

		assert.NotNil(t, v)
	})
}

func TestValidator_Validate(t *testing.T) {
	t.Run("should return error when required field is missing", func(t *testing.T) {
		v := validator.New()
		input := testStruct{
			Email: testEmail,
			Age:   testAge,
		}

		err := v.Validate(input)

		assert.Error(t, err)
	})

	t.Run("should return echo bad request error when validation fails", func(t *testing.T) {
		v := validator.New()
		input := testStruct{
			Name:  testName,
			Email: "not-an-email",
			Age:   testAge,
		}

		err := v.Validate(input)
		httpErr, ok := err.(*echo.HTTPError)

		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("should return error when age is out of range", func(t *testing.T) {
		v := validator.New()
		input := testStruct{
			Name:  testName,
			Email: testEmail,
			Age:   200,
		}

		err := v.Validate(input)

		assert.Error(t, err)
	})

	t.Run("should return error when input is not a struct", func(t *testing.T) {
		v := validator.New()

		err := v.Validate("not a struct")

		assert.Error(t, err)
	})

	t.Run("should return error when required struct field is nil", func(t *testing.T) {
		v := validator.New()
		type wrapper struct {
			Inner *testStruct `validate:"required"`
		}
		input := wrapper{}

		err := v.Validate(input)

		assert.Error(t, err)
	})

	t.Run("should return nil when struct is valid", func(t *testing.T) {
		v := validator.New()
		input := testStruct{
			Name:  testName,
			Email: testEmail,
			Age:   testAge,
		}

		err := v.Validate(input)

		assert.NoError(t, err)
	})
}
