package response

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun/driver/pgdriver"
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

func ErrInvalidJson(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid JSON",
		ErrorText:      err.Error(),
	}
}

func ErrInvalidRequest(err error) render.Renderer {
	var msg string
	switch err.(type) {
	case *json.UnmarshalTypeError:
		jsonErr := err.(*json.UnmarshalTypeError)
		msg = fmt.Sprint("Field with error: ", jsonErr.Field)
		break
	case pgdriver.Error:
		if pgErr, ok := err.(pgdriver.Error); ok {
			if pgErr.IntegrityViolation() && pgErr.Field('C') == pgerrcode.UniqueViolation {
				msg = fmt.Sprintf("%s", pgErr.Field('D'))
			}
		} else {
			msg = err.Error()
		}
	default:
		msg = err.Error()
		break
	}

	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      msg,
	}
}

func ErrIncorrectCredentials(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Authorization failed.",
		ErrorText:      "The credentials provided are incorrect.",
	}
}

func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error.",
		ErrorText:      "Internal Server Error.",
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     "Resource not found.",
		ErrorText:      err.Error(),
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
