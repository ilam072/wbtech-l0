package validator

import "github.com/go-playground/validator/v10"

type OrderValidator struct {
	validate *validator.Validate
}

func New() *OrderValidator {
	return &OrderValidator{validate: validator.New()}
}
func (v *OrderValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}
