package response

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	HTTPStatusCode int `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

type ValidationErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidationErr(tag, fieldForErrMsg, param, message, field string) *ValidationErr {
	return &ValidationErr{
		Field:   field,
		Message: msgForTag(tag, fieldForErrMsg, param, message),
	}
}

type ValidationErrResponse struct {
	*ErrResponse

	ErrorTexts []*ValidationErr `json:"errors,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Not found",
		ErrorText:      err.Error(),
	}
}

func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Conflict",
		ErrorText:      err.Error(),
	}
}

func ErrInvalidJson(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid JSON",
		ErrorText:      err.Error(),
	}
}

func ErrUnauthorized() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Authentication required",
		ErrorText: "Access to this resource requires authentication. Please provide a valid JWT token in the " +
			"Authorization header (Bearer {token}), in the 'jwt' cookie, or as the 'jwt' query parameter.",
	}
}

func ErrIncorrectCredentials() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Authorization failed",
		ErrorText:      "The credentials provided are incorrect.",
	}
}

func ErrInternalServer() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error",
		ErrorText:      "Internal Server Error.",
	}
}

func ErrValidation(errors []*ValidationErr) render.Renderer {
	return &ValidationErrResponse{
		ErrResponse: &ErrResponse{
			HTTPStatusCode: 400,
			StatusText:     "Validation failed.",
			ErrorText:      "",
		},
		ErrorTexts: errors,
	}
}

func msgForTag(tag, fieldForErrMsg, param, originErrMessage string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("The '%s' field is required.", fieldForErrMsg)
	case "min":
		return fmt.Sprintf("The '%s' field must be at least %v.", fieldForErrMsg, param)
	case "max":
		return fmt.Sprintf("The '%s' field may not be greater than %v.", fieldForErrMsg, param)
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
