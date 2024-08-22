package restapi_test

import (
	"bytes"
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
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_restApi_V1LoginPost(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()
	conf.Init()

	mockAuthService := auth.NewMockService(mock)

	h := restapi.New(&service.Dependency{
		AuthService: mockAuthService,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedReq := api.V1LoginPostRequestBody{
			Email:    faker.Email(),
			Password: faker.Password(),
		}
		expectedAuthLoginOutput := auth.LoginOutput{
			AccessToken: auth.LoginOutputToken{
				ExpiredAt: time.Now().UTC(),
				Token:     faker.Jwt(),
			},
			RefreshToken: auth.LoginOutputToken{
				ExpiredAt: time.Now().UTC(),
				Token:     faker.Jwt(),
			},
			UserID: rand.Int63(),
			Email:  faker.Email(),
		}

		expectedReqBodyByte, err := json.Marshal(expectedReq)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Post("/api/v1/login", h.V1LoginPost)

		req, err := http.NewRequest("POST", "/api/v1/login", io.NopCloser(bytes.NewBuffer(expectedReqBodyByte)))
		assert.NoError(t, err)

		mockAuthService.EXPECT().
			Login(req.Context(), auth.LoginInput{
				Email:    expectedReq.Email,
				Password: expectedReq.Password,
			}).
			Return(expectedAuthLoginOutput, nil)

		expectedResp := api.V1LoginPostResponseBody{
			AccessToken: api.V1TokenJwtResponse{
				ExpiredAt: expectedAuthLoginOutput.AccessToken.ExpiredAt,
				Token:     expectedAuthLoginOutput.AccessToken.Token,
			},
			Email:  expectedAuthLoginOutput.Email,
			UserId: expectedAuthLoginOutput.UserID,
		}

		rr := httptest.NewRecorder()
		h.V1LoginPost(rr, req)

		cookie := rr.Result().Cookies()
		require.NotEmpty(t, cookie)

		foundCookie := false
		for _, c := range cookie {
			if c.Name == primitive.RefreshTokenCookieKey {
				foundCookie = true
				assert.Equal(t, expectedAuthLoginOutput.RefreshToken.Token, c.Value)
				assert.WithinDuration(t, expectedAuthLoginOutput.RefreshToken.ExpiredAt, c.Expires, time.Second)
				assert.True(t, c.HttpOnly)
				assert.Equal(t, "/api/v1/refresh-token", c.Path)
				assert.Equal(t, http.SameSiteStrictMode, c.SameSite)
				break
			}
		}
		assert.True(t, foundCookie, "Expected refresh token cookie to be set")

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp api.V1LoginPostResponseBody
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResp, resp)
	})

	t.Run("should be return error user not found", func(t *testing.T) {
		expectedReq := api.V1LoginPostRequestBody{
			Email:    faker.Email(),
			Password: faker.Password(),
		}

		expectedReqBodyByte, err := json.Marshal(expectedReq)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Post("/api/v1/login", h.V1LoginPost)

		req, err := http.NewRequest("POST", "/api/v1/login", io.NopCloser(bytes.NewBuffer(expectedReqBodyByte)))
		assert.NoError(t, err)

		mockAuthService.EXPECT().
			Login(req.Context(), auth.LoginInput{
				Email:    expectedReq.Email,
				Password: expectedReq.Password,
			}).
			Return(auth.LoginOutput{}, auth.ErrUserNotFound)

		expectedResp := api.Error{
			Message: "invalid email or password",
		}

		rr := httptest.NewRecorder()
		h.V1LoginPost(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var resp api.Error
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResp, resp)
	})

	t.Run("should be return error invalid password", func(t *testing.T) {
		expectedReq := api.V1LoginPostRequestBody{
			Email:    faker.Email(),
			Password: faker.Password(),
		}

		expectedReqBodyByte, err := json.Marshal(expectedReq)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Post("/api/v1/login", h.V1LoginPost)

		req, err := http.NewRequest("POST", "/api/v1/login", io.NopCloser(bytes.NewBuffer(expectedReqBodyByte)))
		assert.NoError(t, err)

		mockAuthService.EXPECT().
			Login(req.Context(), auth.LoginInput{
				Email:    expectedReq.Email,
				Password: expectedReq.Password,
			}).
			Return(auth.LoginOutput{}, auth.ErrInvalidPassword)

		expectedResp := api.Error{
			Message: "invalid email or password",
		}

		rr := httptest.NewRecorder()
		h.V1LoginPost(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var resp api.Error
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResp, resp)
	})
}
