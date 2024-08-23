package restapi_test

import (
	"context"
	"encoding/json"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/go-chi/chi/v5"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func Test_restApi_ImageDisplay(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()
	conf.Init()

	mockConsumerService := consumer.NewMockService(mock)

	h := restapi.New(&service.Dependency{
		ConsumerService: mockConsumerService,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedOutput := consumer.GetPrivateImageOutput{
			Object: io.NopCloser(strings.NewReader("image data")),
		}
		expectedParams := api.V1ImagePrivateParams{
			ImageKtp:    null.BoolFrom(true).Ptr(),
			ImageSelfie: nil,
		}
		expectedUserID := rand.Int63()

		r := chi.NewRouter()
		r.Get("/api/v1/image-display", h.V1ImagePrivateGet)

		params := url.Values{}
		params.Add("image_ktp", strconv.FormatBool(*expectedParams.ImageKtp))

		req, err := http.NewRequest("GET", "/api/v1/image-display?"+params.Encode(), nil)
		assert.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(), primitive.UserIDKey, expectedUserID))

		mockConsumerService.EXPECT().
			GetPrivateImage(req.Context(), consumer.GetPrivateImageInput{
				UserID:      expectedUserID,
				ImageKtp:    null.BoolFromPtr(expectedParams.ImageKtp),
				ImageSelfie: null.Bool{},
			}).
			Return(expectedOutput, nil)

		rr := httptest.NewRecorder()
		h.V1ImagePrivateGet(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "image/png", rr.Header().Get("Content-Type"))
		assert.Equal(t, "image data", rr.Body.String())
	})

	t.Run("should be return error unauthorized", func(t *testing.T) {
		expectedParams := api.V1ImagePrivateParams{
			ImageKtp:    null.BoolFrom(true).Ptr(),
			ImageSelfie: nil,
		}

		r := chi.NewRouter()
		r.Get("/api/v1/image-display", h.V1ImagePrivateGet)

		params := url.Values{}
		params.Add("image_ktp", strconv.FormatBool(*expectedParams.ImageKtp))

		req, err := http.NewRequest("GET", "/api/v1/image-display?"+params.Encode(), nil)
		assert.NoError(t, err)

		expectedResp := api.Error{
			Message: "Unauthorized",
		}
		rr := httptest.NewRecorder()
		h.V1ImagePrivateGet(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		var resp api.Error
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResp, resp)
	})
}
