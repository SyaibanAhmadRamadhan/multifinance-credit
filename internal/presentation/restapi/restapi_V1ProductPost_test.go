package restapi_test

import (
	"bytes"
	"encoding/json"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/product"
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

func Test_restApi_V1ProductPost(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()
	conf.Init()

	mockProductService := product.NewMockService(mock)
	h := restapi.New(&service.Dependency{
		ProductService: mockProductService,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedReq := api.V1ProductPostRequestBody{
			Image: api.FileUploadRequest{
				ChecksumSha256:   faker.UUIDDigit(),
				Identifier:       faker.UUIDDigit(),
				MimeType:         "image/jpeg",
				OriginalFilename: "image.jpeg",
				Size:             5000,
			},
			MerchantId: 1,
			Name:       faker.Name(),
			Price:      rand.Float64(),
			Qty:        rand.Int63(),
		}
		expectedPresignedOutput := primitive.PresignedFileUploadOutput{
			Identifier:      expectedReq.Image.Identifier,
			UploadURL:       faker.URL(),
			UploadExpiredAt: time.Now().UTC(),
			MinioFormData:   make(map[string]string),
		}
		expectedReqBodyByte, err := json.Marshal(expectedReq)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Post("/api/v1/product", h.V1ProductPost)

		req, err := http.NewRequest("POST", "/api/v1/product", io.NopCloser(bytes.NewBuffer(expectedReqBodyByte)))
		assert.NoError(t, err)

		mockProductService.EXPECT().
			Create(req.Context(), product.CreateInput{
				MerchantID: expectedReq.MerchantId,
				Name:       expectedReq.Name,
				Image: primitive.PresignedFileUpload{
					Identifier:        expectedReq.Image.Identifier,
					OriginalFileName:  expectedReq.Image.OriginalFilename,
					MimeType:          primitive.MimeType(expectedReq.Image.MimeType),
					Size:              expectedReq.Image.Size,
					ChecksumSHA256:    expectedReq.Image.ChecksumSha256,
					GeneratedFileName: expectedReq.Image.OriginalFilename,
					Extension:         ".jpeg",
				},
				Qty:   expectedReq.Qty,
				Price: expectedReq.Price,
			}).
			Return(product.CreateOutput{
				Image: expectedPresignedOutput,
			}, nil)

		expectedResp := api.V1ProductPostResponseBody{
			PresignedImageUpload: api.FileUploadResponse{
				Identifier:      expectedPresignedOutput.Identifier,
				MinioFormData:   expectedPresignedOutput.MinioFormData,
				UploadExpiredAt: expectedPresignedOutput.UploadExpiredAt,
				UploadUrl:       expectedPresignedOutput.UploadURL,
			},
		}

		rr := httptest.NewRecorder()
		h.V1ProductPost(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var resp api.V1ProductPostResponseBody
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResp, resp)
	})
}
