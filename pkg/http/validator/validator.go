package validator

import (
	"fmt"
	"github.com/Markard/wordka/pkg/http/response"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"strings"
	"unicode"
)

type ProjectValidator interface {
	Struct(s interface{}) *response.ValidationError
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

func (v *Validator) Struct(s interface{}) *response.ValidationError {
	if errVal := v.Validate.Struct(s); errVal != nil {
		validationErrors := errVal.(validator.ValidationErrors)
		errResult := response.NewValidationError()
		for _, validationError := range validationErrors {
			errResult.FieldErrors = append(
				errResult.FieldErrors,
				response.NewFieldValidationError(
					validationError.Field(),
					msgForTag(
						validationError.Tag(),
						formatFieldForMsg(validationError),
						validationError.Param(),
						validationError.Error(),
					),
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

func msgForTag(tag, fieldForErrMsg, param, originErrMessage string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("The '%s' field is required.", fieldForErrMsg)
	case "min":
		return fmt.Sprintf("The '%s' field must be at least %v.", fieldForErrMsg, param)
	case "max":
		return fmt.Sprintf("The '%s' field may not be greater than %v.", fieldForErrMsg, param)
	case "len":
		return fmt.Sprintf("The '%s' field must be %v characters.", fieldForErrMsg, param)
	case "email":
		return fmt.Sprintf("The '%s' field must be a valid email address.", fieldForErrMsg)
	case "validate_password":
		return fmt.Sprintf(
			"The '%s' field must contains at least one uppercase letter, one lowercase letter, one number and one special character.",
			fieldForErrMsg,
		)
	}
	return originErrMessage
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
