package restapi_test

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/go-chi/chi/v5"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/assert"
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

	mockConsumerService := consumer.NewMockService(mock)

	h := restapi.New(&service.Dependency{
		ConsumerService: mockConsumerService,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedOutput := consumer.GetPrivateImageOutput{
			Object: io.NopCloser(strings.NewReader("image data")),
		}
		expectedParams := api.V1ImageDisplayParams{
			ConsumerId:  nil,
			UserId:      rand.Int63(),
			ImageKtp:    null.BoolFrom(true).Ptr(),
			ImageSelfie: nil,
		}

		r := chi.NewRouter()
		r.Get("/api/v1/image-display", h.ImageDisplay)

		params := url.Values{}
		params.Add("image_ktp", strconv.FormatBool(*expectedParams.ImageKtp))
		params.Add("user_id", strconv.Itoa(int(expectedParams.UserId)))

		req, err := http.NewRequest("GET", "/api/v1/image-display?"+params.Encode(), nil)
		assert.NoError(t, err)

		mockConsumerService.EXPECT().
			GetPrivateImage(req.Context(), consumer.GetPrivateImageInput{
				UserID:      expectedParams.UserId,
				ConsumerID:  null.Int{},
				ImageKtp:    null.BoolFromPtr(expectedParams.ImageKtp),
				ImageSelfie: null.Bool{},
			}).
			Return(expectedOutput, nil)

		rr := httptest.NewRecorder()
		h.ImageDisplay(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "image/png", rr.Header().Get("Content-Type"))
		assert.Equal(t, "image data", rr.Body.String())
	})
}
