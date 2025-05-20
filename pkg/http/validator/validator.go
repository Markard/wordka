package validator

import (
	"github.com/Markard/wordka/pkg/http/response"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"strings"
	"unicode"
)

type ProjectValidator interface {
	Struct(s interface{}) []*response.ValidationErr
}

type Validator struct {
	Validate *validator.Validate
}

func NewValidator() (*Validator, error) {
	v := &Validator{Validate: validator.New()}
	err := v.Validate.RegisterValidation("validate_password", validatePassword)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Validator) Struct(s interface{}) []*response.ValidationErr {
	if errVal := v.Validate.Struct(s); errVal != nil {
		validationErrors := errVal.(validator.ValidationErrors)
		errResult := make([]*response.ValidationErr, 0, len(validationErrors))
		for _, validationError := range validationErrors {
			errResult = append(
				errResult,
				response.NewValidationErr(
					validationError.Tag(),
					formatFieldForMsg(validationError),
					validationError.Param(),
					validationError.Error(),
					validationError.Field(),
				),
			)
		}
		return errResult
	}
	return nil
}

func formatFieldForMsg(fieldError validator.FieldError) string {
	fields := strings.Split(fieldError.Namespace(), ".")
	fields = fields[1:]
	snakeCasedFields := make([]string, 0, len(fields))
	for _, field := range fields {
		snakeCasedFields = append(snakeCasedFields, strcase.ToSnake(field))
	}

	return strings.Join(snakeCasedFields, ".")
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
