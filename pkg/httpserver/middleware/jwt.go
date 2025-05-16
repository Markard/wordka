package middleware

import (
	"context"
	"errors"
	"github.com/Markard/wordka/pkg/httpserver/response"
	"github.com/go-chi/render"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	Sub int64
	Exp time.Time
	Iat time.Time
}

func NewToken(sub int64, iat time.Time, exp time.Time) (*Token, error) {
	return &Token{
		Sub: sub,
		Iat: iat,
		Exp: exp,
	}, nil
}

type TokenVerifier interface {
	VerifyTokenStringWithES256(tokenString string) (*Token, error)
}

type contextKey struct {
	name string
}

var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

var ErrNoTokenFound = errors.New("no token found")

// JwtVerifier http middleware handler will verify a JWT string from a http request.
//
// JwtVerifier will search for a JWT token in a http request, in the order:
//  1. 'Authorization: BEARER T' request header
//  2. Cookie 'jwt' value
//
// The first JWT string that is found as a query parameter, authorization header
// or cookie header is then decoded by the `jwt-go` library and a *jwt.Token
// object is set on the request context. In the case of a signature decoding error
// the JwtVerifier will also set the error on the request context.
func JwtVerifier(tv TokenVerifier) func(http.Handler) http.Handler {
	return verify(tv, tokenFromHeader, tokenFromCookie, tokenFromQuery)
}

func verify(tv TokenVerifier, findTokenFns ...func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := verifyRequest(tv, r, findTokenFns...)
			ctx = newContext(ctx, token, err)

			if err != nil {
				_ = render.Render(w, r, response.ErrUnauthorized())
				return
			}

			if token == nil {
				_ = render.Render(w, r, response.ErrUnauthorized())
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func verifyRequest(tv TokenVerifier, r *http.Request, findTokenFns ...func(r *http.Request) string) (*Token, error) {
	var tokenString string

	for _, fn := range findTokenFns {
		tokenString = fn(r)
		if tokenString != "" {
			break
		}
	}
	if tokenString == "" {
		return nil, ErrNoTokenFound
	}

	return tv.VerifyTokenStringWithES256(tokenString)
}

const (
	headerName   = "Authorization"
	headerPrefix = "BEARER "
	cookieName   = "jwt"
	queryName    = "jwt"
)

func tokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get(headerName)
	if len(bearer) > 7 && strings.ToUpper(bearer[0:7]) == headerPrefix {
		return bearer[7:]
	}
	return ""
}

func tokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func tokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get(queryName)
}

func newContext(ctx context.Context, t *Token, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}
