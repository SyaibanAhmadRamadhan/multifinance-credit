package restapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
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

func Test_restApi_Register(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mockAuthService := auth.NewMockService(mock)

	h := restapi.New(&service.Dependency{
		AuthService: mockAuthService,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedReqBody := api.V1RegisterPostRequestBody{
			DateOfBirth: time.Now().UTC(),
			Email:       faker.Email(),
			FullName:    faker.Name(),
			LegalName:   faker.Name(),
			Nik:         fmt.Sprintf("%016d", rand.Int63n(1e16)),
			Password:    "password",
			PhotoKtp: api.FileUploadRequest{
				ChecksumSha256:   faker.UUIDDigit(),
				Identifier:       faker.UUIDDigit(),
				MimeType:         "image/jpeg",
				OriginalFilename: "image.jpeg",
				Size:             5000,
			},
			PhotoSelfie: api.FileUploadRequest{
				ChecksumSha256:   faker.UUIDDigit(),
				Identifier:       faker.UUIDDigit(),
				MimeType:         "image/jpeg",
				OriginalFilename: "image2.jpeg",
				Size:             5000,
			},
			PlaceOfBirth: "Jakarta",
			RePassword:   "password",
			Salary:       5000000,
		}
		expectedRegisterAuthOutput := auth.RegisterOutput{
			UserID:     rand.Int63(),
			ConsumerID: rand.Int63(),
			PhotoKtpPresignedFileUploadOutput: primitive.PresignedFileUploadOutput{
				Identifier:      expectedReqBody.PhotoKtp.Identifier,
				UploadURL:       faker.URL(),
				UploadExpiredAt: time.Now().UTC(),
				MinioFormData:   make(map[string]string),
			},
			PhotoSelfiePresignedFileUploadOutput: primitive.PresignedFileUploadOutput{
				Identifier:      expectedReqBody.PhotoSelfie.Identifier,
				UploadURL:       faker.URL(),
				UploadExpiredAt: time.Now().UTC(),
				MinioFormData:   make(map[string]string),
			},
		}

		expectedReqBodyByte, err := json.Marshal(expectedReqBody)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Post("/api/v1/register", h.V1RegisterPost)

		req, err := http.NewRequest("POST", "/api/v1/register", io.NopCloser(bytes.NewBuffer(expectedReqBodyByte)))
		assert.NoError(t, err)

		mockAuthService.EXPECT().
			Register(req.Context(), auth.RegisterInput{
				DateOfBirth: expectedReqBody.DateOfBirth,
				Email:       expectedReqBody.Email,
				Nik:         expectedReqBody.Nik,
				FullName:    expectedReqBody.FullName,
				LegalName:   expectedReqBody.LegalName,
				PhotoKtp: primitive.PresignedFileUpload{
					Identifier:        expectedReqBody.PhotoKtp.Identifier,
					OriginalFileName:  expectedReqBody.PhotoKtp.OriginalFilename,
					MimeType:          primitive.MimeType(expectedReqBody.PhotoKtp.MimeType),
					Size:              expectedReqBody.PhotoKtp.Size,
					ChecksumSHA256:    expectedReqBody.PhotoKtp.ChecksumSha256,
					GeneratedFileName: expectedReqBody.PhotoKtp.OriginalFilename,
					Extension:         ".jpeg",
				},
				PhotoSelfie: primitive.PresignedFileUpload{
					Identifier:        expectedReqBody.PhotoSelfie.Identifier,
					OriginalFileName:  expectedReqBody.PhotoSelfie.OriginalFilename,
					MimeType:          primitive.MimeType(expectedReqBody.PhotoSelfie.MimeType),
					Size:              expectedReqBody.PhotoSelfie.Size,
					ChecksumSHA256:    expectedReqBody.PhotoSelfie.ChecksumSha256,
					GeneratedFileName: expectedReqBody.PhotoSelfie.OriginalFilename,
					Extension:         ".jpeg",
				},
				PlaceOfBirth: expectedReqBody.PlaceOfBirth,
				Salary:       expectedReqBody.Salary,
				Password:     expectedReqBody.Password,
			}).
			Return(expectedRegisterAuthOutput, nil)

		rr := httptest.NewRecorder()
		h.V1RegisterPost(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		expectedResponse := api.V1RegisterPost200Response{
			ConsumerId: expectedRegisterAuthOutput.ConsumerID,
			UploadPhotoKtp: api.FileUploadResponse{
				Identifier:      expectedRegisterAuthOutput.PhotoKtpPresignedFileUploadOutput.Identifier,
				MinioFormData:   expectedRegisterAuthOutput.PhotoKtpPresignedFileUploadOutput.MinioFormData,
				UploadExpiredAt: expectedRegisterAuthOutput.PhotoKtpPresignedFileUploadOutput.UploadExpiredAt,
				UploadUrl:       expectedRegisterAuthOutput.PhotoKtpPresignedFileUploadOutput.UploadURL,
			},
			UploadPhotoSelfie: api.FileUploadResponse{
				Identifier:      expectedRegisterAuthOutput.PhotoSelfiePresignedFileUploadOutput.Identifier,
				MinioFormData:   expectedRegisterAuthOutput.PhotoSelfiePresignedFileUploadOutput.MinioFormData,
				UploadExpiredAt: expectedRegisterAuthOutput.PhotoSelfiePresignedFileUploadOutput.UploadExpiredAt,
				UploadUrl:       expectedRegisterAuthOutput.PhotoSelfiePresignedFileUploadOutput.UploadURL,
			},
			UserId: expectedRegisterAuthOutput.UserID,
		}
		var resp api.V1RegisterPost200Response
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResponse, resp)
	})

	t.Run("should be return error nik available", func(t *testing.T) {
		expectedReqBody := api.V1RegisterPostRequestBody{
			DateOfBirth: time.Now().UTC(),
			Email:       faker.Email(),
			FullName:    faker.Name(),
			LegalName:   faker.Name(),
			Nik:         fmt.Sprintf("%016d", rand.Int63n(1e16)),
			Password:    "password",
			PhotoKtp: api.FileUploadRequest{
				ChecksumSha256:   faker.UUIDDigit(),
				Identifier:       faker.UUIDDigit(),
				MimeType:         "image/jpeg",
				OriginalFilename: "image.jpeg",
				Size:             5000,
			},
			PhotoSelfie: api.FileUploadRequest{
				ChecksumSha256:   faker.UUIDDigit(),
				Identifier:       faker.UUIDDigit(),
				MimeType:         "image/jpeg",
				OriginalFilename: "image2.jpeg",
				Size:             5000,
			},
			PlaceOfBirth: "Jakarta",
			RePassword:   "password",
			Salary:       5000000,
		}

		expectedReqBodyByte, err := json.Marshal(expectedReqBody)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Get("/api/v1/register", h.V1RegisterPost)

		req, err := http.NewRequest("GET", "/api/v1/register", io.NopCloser(bytes.NewBuffer(expectedReqBodyByte)))
		assert.NoError(t, err)

		mockAuthService.EXPECT().
			Register(req.Context(), auth.RegisterInput{
				DateOfBirth: expectedReqBody.DateOfBirth,
				Email:       expectedReqBody.Email,
				Nik:         expectedReqBody.Nik,
				FullName:    expectedReqBody.FullName,
				LegalName:   expectedReqBody.LegalName,
				PhotoKtp: primitive.PresignedFileUpload{
					Identifier:        expectedReqBody.PhotoKtp.Identifier,
					OriginalFileName:  expectedReqBody.PhotoKtp.OriginalFilename,
					MimeType:          primitive.MimeType(expectedReqBody.PhotoKtp.MimeType),
					Size:              expectedReqBody.PhotoKtp.Size,
					ChecksumSHA256:    expectedReqBody.PhotoKtp.ChecksumSha256,
					GeneratedFileName: expectedReqBody.PhotoKtp.OriginalFilename,
					Extension:         ".jpeg",
				},
				PhotoSelfie: primitive.PresignedFileUpload{
					Identifier:        expectedReqBody.PhotoSelfie.Identifier,
					OriginalFileName:  expectedReqBody.PhotoSelfie.OriginalFilename,
					MimeType:          primitive.MimeType(expectedReqBody.PhotoSelfie.MimeType),
					Size:              expectedReqBody.PhotoSelfie.Size,
					ChecksumSHA256:    expectedReqBody.PhotoSelfie.ChecksumSha256,
					GeneratedFileName: expectedReqBody.PhotoSelfie.OriginalFilename,
					Extension:         ".jpeg",
				},
				PlaceOfBirth: expectedReqBody.PlaceOfBirth,
				Salary:       expectedReqBody.Salary,
				Password:     expectedReqBody.Password,
			}).
			Return(auth.RegisterOutput{}, auth.ErrNikIsAvailable)

		rr := httptest.NewRecorder()
		h.V1RegisterPost(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		expectedResponse := api.Error{
			Message: auth.ErrNikIsAvailable.Error(),
		}
		var resp api.Error
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResponse, resp)
	})

	t.Run("should be return error email available", func(t *testing.T) {
		expectedReqBody := api.V1RegisterPostRequestBody{
			DateOfBirth: time.Now().UTC(),
			Email:       faker.Email(),
			FullName:    faker.Name(),
			LegalName:   faker.Name(),
			Nik:         fmt.Sprintf("%016d", rand.Int63n(1e16)),
			Password:    "password",
			PhotoKtp: api.FileUploadRequest{
				ChecksumSha256:   faker.UUIDDigit(),
				Identifier:       faker.UUIDDigit(),
				MimeType:         "image/jpeg",
				OriginalFilename: "image.jpeg",
				Size:             5000,
			},
			PhotoSelfie: api.FileUploadRequest{
				ChecksumSha256:   faker.UUIDDigit(),
				Identifier:       faker.UUIDDigit(),
				MimeType:         "image/jpeg",
				OriginalFilename: "image2.jpeg",
				Size:             5000,
			},
			PlaceOfBirth: "Jakarta",
			RePassword:   "password",
			Salary:       5000000,
		}

		expectedReqBodyByte, err := json.Marshal(expectedReqBody)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Get("/api/v1/register", h.V1RegisterPost)

		req, err := http.NewRequest("GET", "/api/v1/register", io.NopCloser(bytes.NewBuffer(expectedReqBodyByte)))
		assert.NoError(t, err)

		mockAuthService.EXPECT().
			Register(req.Context(), auth.RegisterInput{
				DateOfBirth: expectedReqBody.DateOfBirth,
				Email:       expectedReqBody.Email,
				Nik:         expectedReqBody.Nik,
				FullName:    expectedReqBody.FullName,
				LegalName:   expectedReqBody.LegalName,
				PhotoKtp: primitive.PresignedFileUpload{
					Identifier:        expectedReqBody.PhotoKtp.Identifier,
					OriginalFileName:  expectedReqBody.PhotoKtp.OriginalFilename,
					MimeType:          primitive.MimeType(expectedReqBody.PhotoKtp.MimeType),
					Size:              expectedReqBody.PhotoKtp.Size,
					ChecksumSHA256:    expectedReqBody.PhotoKtp.ChecksumSha256,
					GeneratedFileName: expectedReqBody.PhotoKtp.OriginalFilename,
					Extension:         ".jpeg",
				},
				PhotoSelfie: primitive.PresignedFileUpload{
					Identifier:        expectedReqBody.PhotoSelfie.Identifier,
					OriginalFileName:  expectedReqBody.PhotoSelfie.OriginalFilename,
					MimeType:          primitive.MimeType(expectedReqBody.PhotoSelfie.MimeType),
					Size:              expectedReqBody.PhotoSelfie.Size,
					ChecksumSHA256:    expectedReqBody.PhotoSelfie.ChecksumSha256,
					GeneratedFileName: expectedReqBody.PhotoSelfie.OriginalFilename,
					Extension:         ".jpeg",
				},
				PlaceOfBirth: expectedReqBody.PlaceOfBirth,
				Salary:       expectedReqBody.Salary,
				Password:     expectedReqBody.Password,
			}).
			Return(auth.RegisterOutput{}, auth.ErrEmailIsAvailable)

		rr := httptest.NewRecorder()
		h.V1RegisterPost(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		expectedResponse := api.Error{
			Message: auth.ErrEmailIsAvailable.Error(),
		}
		var resp api.Error
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, expectedResponse, resp)
	})
}
