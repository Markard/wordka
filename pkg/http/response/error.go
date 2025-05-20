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
