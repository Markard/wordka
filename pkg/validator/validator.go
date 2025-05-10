package validator

import (
	"github.com/go-playground/validator/v10"
	"unicode"
)

type Validator struct {
	Validator *validator.Validate
}

func NewValidator() (*Validator, error) {
	v := &Validator{Validator: validator.New()}
	err := v.Validator.RegisterValidation("validate_password", validatePassword)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func validatePassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, ch := range value {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}
