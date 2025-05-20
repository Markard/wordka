package response

import (
	"encoding/json"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Message    string
}

func NewHttpError(statusCode int, message string) *HttpError {
	return &HttpError{StatusCode: statusCode, Message: message}
}

type JsonResponse struct {
	Message string `json:"error"`
	Details string `json:"details,omitempty"`
}

func replyAsJson(w http.ResponseWriter, httpError *HttpError) {
	r := JsonResponse{Message: httpError.Message}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpError.StatusCode)
	_ = json.NewEncoder(w).Encode(r)
}

func ErrNotFound(w http.ResponseWriter, err error) {
	httpError := NewHttpError(http.StatusNotFound, err.Error())
	replyAsJson(w, httpError)
}

func ErrConflict(w http.ResponseWriter, err error) {
	httpError := NewHttpError(http.StatusConflict, err.Error())
	replyAsJson(w, httpError)
}

func ErrInternalServer(w http.ResponseWriter) {
	httpError := NewHttpError(http.StatusInternalServerError, "Internal server error.")
	replyAsJson(w, httpError)
}

func ErrHttpError(w http.ResponseWriter, statusCode int, msg string) {
	httpError := NewHttpError(statusCode, msg)
	replyAsJson(w, httpError)
}

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

func (err *ValidationError) AddFieldError(field, message string) *ValidationError {
	err.FieldErrors = append(err.FieldErrors, NewFieldValidationError(field, message))
	return err
}

type JsonFieldValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type JsonValidationErrorResponse struct {
	Message     string                              `json:"message"`
	FieldErrors []*JsonFieldValidationErrorResponse `json:"field_errors,omitempty"`
}

func ErrValidation(w http.ResponseWriter, valErr *ValidationError) {
	replyValErrAsJson(w, valErr)
}

func replyValErrAsJson(w http.ResponseWriter, valError *ValidationError) {
	r := &JsonValidationErrorResponse{Message: valError.HttpError.Message}
	for _, fieldErr := range valError.FieldErrors {
		r.FieldErrors = append(r.FieldErrors, &JsonFieldValidationErrorResponse{
			Field:   fieldErr.Field,
			Message: fieldErr.Message,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(valError.HttpError.StatusCode)
	_ = json.NewEncoder(w).Encode(r)
}
