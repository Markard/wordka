package jwt

import (
	"context"
	"errors"
	"github.com/Markard/wordka/internal/entity"
	serviceJwt "github.com/Markard/wordka/internal/infra/service/jwt"
	"github.com/Markard/wordka/pkg/http/response"
	"github.com/Markard/wordka/pkg/slogext"
	"net/http"
	"strings"
)

type TokenVerifier interface {
	VerifyTokenStringWithES256(tokenString string) (*serviceJwt.Token, error)
}

type UserProvider interface {
	FindById(id int64) (*entity.User, error)
}

type contextKey struct {
	name string
}

var (
	TokenCtxKey       = &contextKey{"Token"}
	ErrorCtxKey       = &contextKey{"Error"}
	CurrentUserCtxKey = &contextKey{"CurrentUser"}
)

var ErrNoTokenFound = errors.New("no token found")

// Authenticator http middleware handler will verify a JWT string from a http request.
//
// Authenticator will search for a JWT token in a http request, in the order:
//  1. 'Authorization: BEARER T' request header
//  2. Cookie 'jwt' value
//
// The first JWT string that is found as a query parameter, authorization header
// or cookie header is then decoded by the `jwt-go` library and a *jwt.Token
// object is set on the request context. In the case of a signature decoding error
// the Authenticator will also set the error on the request context.
func Authenticator(tv TokenVerifier, up UserProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := verifyRequest(tv, r, tokenFromHeader, tokenFromCookie, tokenFromQuery)
			errMsg := "Access to this resource requires authentication. Please provide a valid JWT token in the " +
				"Authorization header (Bearer {token}), in the 'jwt' cookie, or as the 'jwt' query parameter."
			if err != nil {
				response.ErrHttpError(w, http.StatusUnauthorized, errMsg)
				return
			}

			if token == nil {
				response.ErrHttpError(w, http.StatusUnauthorized, errMsg)
				return
			}

			user, errUserFindById := up.FindById(token.Sub)
			if errUserFindById != nil {
				response.ErrHttpError(w, http.StatusUnauthorized, errMsg)
				return
			}

			ctx = newTokenContext(ctx, token, user, err)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func verifyRequest(tv TokenVerifier, r *http.Request, findTokenFns ...func(r *http.Request) string) (*serviceJwt.Token, error) {
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

func newTokenContext(ctx context.Context, t *serviceJwt.Token, u *entity.User, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	ctx = context.WithValue(ctx, CurrentUserCtxKey, u)
	ctx = slogext.WithLogUserID(ctx, u.Id)
	return ctx
}
