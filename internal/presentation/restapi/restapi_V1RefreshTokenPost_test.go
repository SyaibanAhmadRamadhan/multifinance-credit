package restapi_test

import (
	"encoding/json"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/go-chi/chi/v5"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_restApi_V1RefreshTokenPost(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()
	conf.Init()

	mockAuthService := auth.NewMockService(mock)

	h := restapi.New(&service.Dependency{
		AuthService: mockAuthService,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedTokenInCookie := faker.Jwt()
		expectedUserID := rand.Int63()

		r := chi.NewRouter()
		r.Post("/api/v1/refresh-token", h.V1LoginPost)

		req, err := http.NewRequest("POST", "/api/v1/refresh-token", nil)
		assert.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:     primitive.RefreshTokenCookieKey,
			Value:    expectedTokenInCookie,
			Expires:  time.Now().UTC().Add(5 * time.Minute),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/api/v1/refresh-token",
		})

		mockAuthService.EXPECT().
			VerifyToken(req.Context(), auth.VerifyTokenInput{
				Token:     expectedTokenInCookie,
				TokenType: primitive.TokenTypeRefreshToken,
			}).
			Return(auth.VerifyTokenOutput{
				UserID: expectedUserID,
			}, nil)

		rr := httptest.NewRecorder()
		h.V1RefreshTokenPost(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp api.V1RefreshTokenPostResponseBody

		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))

		require.NotEmpty(t, resp)
	})

	t.Run("should be return error token expired", func(t *testing.T) {
		expectedTokenInCookie := faker.Jwt()

		r := chi.NewRouter()
		r.Post("/api/v1/refresh-token", h.V1LoginPost)

		req, err := http.NewRequest("POST", "/api/v1/refresh-token", nil)
		assert.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:     primitive.RefreshTokenCookieKey,
			Value:    expectedTokenInCookie,
			Expires:  time.Now().UTC().Add(5 * time.Minute),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/api/v1/refresh-token",
		})

		mockAuthService.EXPECT().
			VerifyToken(req.Context(), auth.VerifyTokenInput{
				Token:     expectedTokenInCookie,
				TokenType: primitive.TokenTypeRefreshToken,
			}).
			Return(auth.VerifyTokenOutput{}, auth.ErrTokenIsExpired)

		expectedResp := api.Error{
			Message: "your token is expired, you can log in again",
		}
		rr := httptest.NewRecorder()
		h.V1RefreshTokenPost(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)

		var resp api.Error

		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResp, resp)
	})

	t.Run("should be return error token invalid", func(t *testing.T) {
		expectedTokenInCookie := faker.Jwt()

		r := chi.NewRouter()
		r.Post("/api/v1/refresh-token", h.V1LoginPost)

		req, err := http.NewRequest("POST", "/api/v1/refresh-token", nil)
		assert.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:     primitive.RefreshTokenCookieKey,
			Value:    expectedTokenInCookie,
			Expires:  time.Now().UTC().Add(5 * time.Minute),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/api/v1/refresh-token",
		})

		mockAuthService.EXPECT().
			VerifyToken(req.Context(), auth.VerifyTokenInput{
				Token:     expectedTokenInCookie,
				TokenType: primitive.TokenTypeRefreshToken,
			}).
			Return(auth.VerifyTokenOutput{}, auth.ErrInvalidToken)

		expectedResp := api.Error{
			Message: "your token is invalid, you can log in again",
		}
		rr := httptest.NewRecorder()
		h.V1RefreshTokenPost(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)

		var resp api.Error

		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResp, resp)
	})
}
