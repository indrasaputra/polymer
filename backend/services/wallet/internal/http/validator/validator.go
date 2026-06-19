package validator

import (
	"log"

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
		// Optionally return the error to let each route control the status code.
		log.Printf("validate %v\n", err)
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}
