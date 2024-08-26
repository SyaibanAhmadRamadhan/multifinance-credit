package middleware

import (
	"context"
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"strings"
)

var ErrMissingAuthorizationHeader = errors.New("missing authorization header")
var ErrInvalidAuthorizationHeader = errors.New("invalid token")
var ErrAuthorizationHeaderNotBearer = errors.New("authorization header not bearer")

func parseAuthorizationHeader(header string) (token string, err error) {
	if header == "" {
		err = ErrMissingAuthorizationHeader
		return
	}
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		err = ErrInvalidAuthorizationHeader
		return
	}

	if parts[0] != "Bearer" {
		err = ErrAuthorizationHeaderNotBearer
		return
	}

	token = parts[1]
	return
}

func (m *middleware) AuthUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otelTracer.Start(r.Context(), "checking authentication user")
		r = r.WithContext(ctx)
		defer span.End()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			restapi.Error(w, r, http.StatusUnauthorized, errors.New("header Authorization not exists"), "invalid authorization")
			return
		}

		token, err := parseAuthorizationHeader(authHeader)
		if err != nil {
			restapi.Error(w, r, http.StatusUnauthorized, err, "invalid authorization")
			return
		}

		verifyTokenOutput, err := m.authService.VerifyToken(r.Context(), auth.VerifyTokenInput{
			Token:     token,
			TokenType: primitive.TokenTypeAccessToken,
		})
		if err != nil {
			if errors.Is(err, auth.ErrTokenIsExpired) {
				restapi.Error(w, r, http.StatusUnauthorized, err, "your token is expired, you can log in again")
			} else if errors.Is(err, auth.ErrInvalidToken) {
				restapi.Error(w, r, http.StatusUnauthorized, err, "your token is invalid, you can log in again")
			} else {
				restapi.Error(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		span.SetAttributes(attribute.String("status", "user is authentication"))
		span.SetAttributes(attribute.Int64("user_id", verifyTokenOutput.UserID))

		ctx = context.WithValue(r.Context(), primitive.UserIDKey, verifyTokenOutput.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
