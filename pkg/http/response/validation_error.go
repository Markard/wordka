package response

import (
	"encoding/json"
	"net/http"
)

type FieldValidationError struct {
	Field   string
	Message string
}

func NewFieldValidationError(field string, message string) *FieldValidationError {
	return &FieldValidationError{Field: field, Message: message}
}

type ValidationError struct {
	*HttpError
	FieldErrors []*FieldValidationError
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		HttpError:   NewHttpError(http.StatusBadRequest, "Validation error"),
		FieldErrors: []*FieldValidationError{},
	}
}

func (valErr *ValidationError) AddFieldError(field, message string) *ValidationError {
	valErr.FieldErrors = append(valErr.FieldErrors, NewFieldValidationError(field, message))
	return valErr
}

type JsonFieldValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type JsonValidationErrorResponse struct {
	Message     string                              `json:"message"`
	FieldErrors []*JsonFieldValidationErrorResponse `json:"field_errors,omitempty"`
}

func (valErr *ValidationError) ErrValidation(w http.ResponseWriter) {
	r := &JsonValidationErrorResponse{Message: valErr.HttpError.Message}
	for _, fieldErr := range valErr.FieldErrors {
		r.FieldErrors = append(r.FieldErrors, &JsonFieldValidationErrorResponse{
			Field:   fieldErr.Field,
			Message: fieldErr.Message,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(valErr.HttpError.StatusCode)
	_ = json.NewEncoder(w).Encode(r)
}
