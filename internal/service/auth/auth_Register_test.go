package auth_test

import (
	"context"
	"database/sql"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	extra "github.com/oxyno-zeta/gomock-extra-matcher"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"testing"
	"time"
)

func Test_service_Register(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	conf.Init()

	mockS3Repository := s3.NewMockRepository(mock)
	mockDBTx := db.NewMockSqlxTransaction(mock)
	mockUserRepository := users.NewMockRepository(mock)
	mockConsumerRepository := consumers.NewMockRepository(mock)
	mockLimitRepository := limits.NewMockRepository(mock)

	s := auth.NewService(auth.NewServiceOpts{
		UserRepository:     mockUserRepository,
		ConsumerRepository: mockConsumerRepository,
		LimitRepository:    mockLimitRepository,
		S3Repository:       mockS3Repository,
		DBTx:               mockDBTx,
	})
	ctx := context.TODO()

	t.Run("should be return correct", func(t *testing.T) {
		expectedBucket := conf.GetConfig().Minio.PrivateBucket
		expectedInput := auth.RegisterInput{
			DateOfBirth: time.Now().UTC(),
			Email:       faker.Email(),
			Nik:         faker.UUIDDigit(),
			FullName:    faker.Name(),
			LegalName:   faker.Name(),
			PhotoKtp: primitive.PresignedFileUpload{
				Identifier:        faker.UUIDDigit(),
				OriginalFileName:  faker.Name(),
				MimeType:          primitive.MimeTypeJpeg,
				Size:              rand.Int63(),
				ChecksumSHA256:    faker.UUIDDigit(),
				GeneratedFileName: faker.UUIDDigit(),
				Extension:         "image/jpeg",
			},
			PhotoSelfie: primitive.PresignedFileUpload{
				Identifier:        faker.UUIDDigit(),
				OriginalFileName:  faker.Name(),
				MimeType:          primitive.MimeTypeJpeg,
				Size:              rand.Int63(),
				ChecksumSHA256:    faker.UUIDDigit(),
				GeneratedFileName: faker.UUIDDigit(),
				Extension:         "image/jpeg",
			},
			PlaceOfBirth: faker.Name(),
			Salary:       rand.Float64(),
			Password:     faker.Password(),
		}

		expectedOutput := auth.RegisterOutput{
			UserID:     rand.Int63(),
			ConsumerID: rand.Int63(),
			PhotoKtpPresignedFileUploadOutput: primitive.PresignedFileUploadOutput{
				Identifier:      expectedInput.PhotoKtp.Identifier,
				UploadURL:       faker.URL(),
				UploadExpiredAt: time.Now().UTC(),
				MinioFormData:   make(map[string]string),
			},
			PhotoSelfiePresignedFileUploadOutput: primitive.PresignedFileUploadOutput{
				Identifier:      expectedInput.PhotoSelfie.Identifier,
				UploadURL:       faker.URL(),
				UploadExpiredAt: time.Now().UTC(),
				MinioFormData:   make(map[string]string),
			},
		}

		mockS3Repository.EXPECT().
			CreatePresignedUrl(ctx, s3.CreatePresignedUrlInput{
				BucketName: expectedBucket,
				Path:       expectedInput.PhotoKtp.GeneratedFileName,
				MimeType:   string(expectedInput.PhotoKtp.MimeType),
				Checksum:   expectedInput.PhotoKtp.ChecksumSHA256,
			}).
			Return(s3.CreatePresignedUrlOutput{
				URL:           expectedOutput.PhotoKtpPresignedFileUploadOutput.UploadURL,
				ExpiredAt:     expectedOutput.PhotoKtpPresignedFileUploadOutput.UploadExpiredAt,
				MinioFormData: expectedOutput.PhotoKtpPresignedFileUploadOutput.MinioFormData,
			}, nil)

		mockS3Repository.EXPECT().
			CreatePresignedUrl(ctx, s3.CreatePresignedUrlInput{
				BucketName: expectedBucket,
				Path:       expectedInput.PhotoSelfie.GeneratedFileName,
				MimeType:   string(expectedInput.PhotoSelfie.MimeType),
				Checksum:   expectedInput.PhotoSelfie.ChecksumSHA256,
			}).
			Return(s3.CreatePresignedUrlOutput{
				URL:           expectedOutput.PhotoSelfiePresignedFileUploadOutput.UploadURL,
				ExpiredAt:     expectedOutput.PhotoSelfiePresignedFileUploadOutput.UploadExpiredAt,
				MinioFormData: expectedOutput.PhotoSelfiePresignedFileUploadOutput.MinioFormData,
			}, nil)

		mockUserRepository.EXPECT().
			CheckExisting(ctx, users.CheckExistingInput{
				ByEmail: null.StringFrom(expectedInput.Email),
			}).
			Return(users.CheckExistingOutput{
				Existing: false,
			}, nil)

		mockConsumerRepository.EXPECT().
			CheckExisting(ctx, consumers.CheckExistingInput{
				ByNIK: null.StringFrom(expectedInput.Nik),
			}).
			Return(consumers.CheckExistingOutput{
				Existing: false,
			}, nil)

		mockDBTx.EXPECT().DoTransaction(ctx, &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		}, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx *sql.TxOptions, fn func(tx *db.SqlxWrapper) error) error {
				mockSqlxWrapper := &db.SqlxWrapper{}

				createInputStructMatcher := extra.StructMatcher().
					Field("Transaction", mockSqlxWrapper).
					Field("Email", expectedInput.Email).
					Field("Password", gomock.Any())

				mockUserRepository.EXPECT().
					Create(ctx, createInputStructMatcher).
					Return(users.CreateOutput{
						ID: expectedOutput.UserID,
					}, nil)

				mockConsumerRepository.EXPECT().
					Create(ctx, consumers.CreateInput{
						Transaction:  mockSqlxWrapper,
						UserID:       expectedOutput.UserID,
						Nik:          expectedInput.Nik,
						FullName:     expectedInput.FullName,
						LegalName:    expectedInput.LegalName,
						PlaceOfBirth: expectedInput.PlaceOfBirth,
						DateOfBirth:  expectedInput.DateOfBirth,
						Salary:       expectedInput.Salary,
						PhotoKTP:     expectedInput.PhotoKtp.GeneratedFileName,
						PhotoSelfie:  expectedInput.PhotoSelfie.GeneratedFileName,
					}).
					Return(consumers.CreateOutput{
						ID: expectedOutput.ConsumerID,
					}, nil)

				mockLimitRepository.EXPECT().
					Creates(ctx, limits.CreatesInput{
						Transaction: mockSqlxWrapper,
						ConsumerID:  expectedOutput.ConsumerID,
						Items:       limits.DefaultLimitData(),
					}).Return(nil)

				return fn(mockSqlxWrapper)
			},
		).Return(nil)

		output, err := s.Register(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return error nik is available", func(t *testing.T) {
		expectedInput := auth.RegisterInput{
			DateOfBirth: time.Now().UTC(),
			Email:       faker.Email(),
			Nik:         faker.UUIDDigit(),
			FullName:    faker.Name(),
			LegalName:   faker.Name(),
			PhotoKtp: primitive.PresignedFileUpload{
				Identifier:        faker.UUIDDigit(),
				OriginalFileName:  faker.Name(),
				MimeType:          primitive.MimeTypeJpeg,
				Size:              rand.Int63(),
				ChecksumSHA256:    faker.UUIDDigit(),
				GeneratedFileName: faker.UUIDDigit(),
				Extension:         "image/jpeg",
			},
			PhotoSelfie: primitive.PresignedFileUpload{
				Identifier:        faker.UUIDDigit(),
				OriginalFileName:  faker.Name(),
				MimeType:          primitive.MimeTypeJpeg,
				Size:              rand.Int63(),
				ChecksumSHA256:    faker.UUIDDigit(),
				GeneratedFileName: faker.UUIDDigit(),
				Extension:         "image/jpeg",
			},
			PlaceOfBirth: faker.Name(),
			Salary:       rand.Float64(),
			Password:     faker.Password(),
		}

		mockUserRepository.EXPECT().
			CheckExisting(ctx, users.CheckExistingInput{
				ByEmail: null.StringFrom(expectedInput.Email),
			}).
			Return(users.CheckExistingOutput{
				Existing: false,
			}, nil)

		mockConsumerRepository.EXPECT().
			CheckExisting(ctx, consumers.CheckExistingInput{
				ByNIK: null.StringFrom(expectedInput.Nik),
			}).
			Return(consumers.CheckExistingOutput{
				Existing: true,
			}, nil)

		_, err := s.Register(ctx, expectedInput)
		require.Error(t, err, auth.ErrNikIsAvailable)
	})

	t.Run("should be return error email is available", func(t *testing.T) {
		expectedInput := auth.RegisterInput{
			DateOfBirth: time.Now().UTC(),
			Email:       faker.Email(),
			Nik:         faker.UUIDDigit(),
			FullName:    faker.Name(),
			LegalName:   faker.Name(),
			PhotoKtp: primitive.PresignedFileUpload{
				Identifier:        faker.UUIDDigit(),
				OriginalFileName:  faker.Name(),
				MimeType:          primitive.MimeTypeJpeg,
				Size:              rand.Int63(),
				ChecksumSHA256:    faker.UUIDDigit(),
				GeneratedFileName: faker.UUIDDigit(),
				Extension:         "image/jpeg",
			},
			PhotoSelfie: primitive.PresignedFileUpload{
				Identifier:        faker.UUIDDigit(),
				OriginalFileName:  faker.Name(),
				MimeType:          primitive.MimeTypeJpeg,
				Size:              rand.Int63(),
				ChecksumSHA256:    faker.UUIDDigit(),
				GeneratedFileName: faker.UUIDDigit(),
				Extension:         "image/jpeg",
			},
			PlaceOfBirth: faker.Name(),
			Salary:       rand.Float64(),
			Password:     faker.Password(),
		}

		mockUserRepository.EXPECT().
			CheckExisting(ctx, users.CheckExistingInput{
				ByEmail: null.StringFrom(expectedInput.Email),
			}).
			Return(users.CheckExistingOutput{
				Existing: true,
			}, nil)

		mockConsumerRepository.EXPECT().
			CheckExisting(ctx, consumers.CheckExistingInput{
				ByNIK: null.StringFrom(expectedInput.Nik),
			}).
			Return(consumers.CheckExistingOutput{
				Existing: false,
			}, nil)

		_, err := s.Register(ctx, expectedInput)
		require.Error(t, err, auth.ErrEmailIsAvailable)
	})
}
