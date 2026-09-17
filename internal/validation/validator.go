package validation

import (
	"gopkg.in/go-playground/validator.v9"
)




type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(value any) error {
	return v.validate.Struct(value)
}