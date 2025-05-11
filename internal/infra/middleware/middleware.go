package middleware

import (
	"net/http"
)

type Middlewares struct {
	JwtAuthenticator func(http.Handler) http.Handler
}
