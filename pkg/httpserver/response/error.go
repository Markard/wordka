package response

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"net/http"
	"strings"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

type ValidationErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrResponse struct {
	*ErrResponse

	ErrorTexts []*ValidationErr `json:"errors,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Conflict",
		ErrorText:      err.Error(),
	}
}

func ErrInvalidJson(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid JSON",
		ErrorText:      err.Error(),
	}
}

func ErrUnauthorized(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Authentication required",
		ErrorText: "Access to this resource requires authentication. Please provide a valid JWT token in the " +
			"Authorization header (Bearer {token}), in the 'jwt' cookie, or as the 'jwt' query parameter.",
	}
}

func ErrIncorrectCredentials(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Authorization failed",
		ErrorText:      "The credentials provided are incorrect.",
	}
}

func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error",
		ErrorText:      "Internal Server Error.",
	}
}

func ErrValidation(err error) render.Renderer {
	validationErrors := err.(validator.ValidationErrors)
	errorsAsSlice := make([]*ValidationErr, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		errorsAsSlice = append(
			errorsAsSlice,
			&ValidationErr{
				Field:   validationError.Field(),
				Message: msgForTag(validationError),
			},
		)
	}

	return &ValidationErrResponse{
		ErrResponse: &ErrResponse{
			Err:            err,
			HTTPStatusCode: 400,
			StatusText:     "Validation failed.",
			ErrorText:      "",
		},
		ErrorTexts: errorsAsSlice,
	}
}

func msgForTag(fieldError validator.FieldError) string {
	field := formatFieldForMsg(fieldError)

	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("The '%s' field is required.", field)
	case "min":
		return fmt.Sprintf("The '%s' field must be at least %v.", field, fieldError.Param())
	case "max":
		return fmt.Sprintf("The '%s' field may not be greater than %v.", field, fieldError.Param())
	case "email":
		return fmt.Sprintf("The '%s' field must be a valid email address.", field)
	case "validate_password":
		return fmt.Sprintf(
			"The '%s' field must contains at least one uppercase letter, one lowercase letter, one number and one special character.",
			field,
		)
	}
	return fieldError.Error()
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
