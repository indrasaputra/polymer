package validator

import (
	goval "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// Validator wraps go-playground validator to satisfy Echo interface.
type Validator struct {
	validator *goval.Validate
}

// New creates an instace of validator.
func New() *Validator {
	return &Validator{
		validator: goval.New(goval.WithRequiredStructEnabled()),
	}
}

// Validate validates input against rules.
func (v *Validator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}
