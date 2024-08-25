package auth

import (
	"context"
	"database/sql"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/guregu/null/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
)

func (s *service) Register(ctx context.Context, input RegisterInput) (output RegisterOutput, err error) {

	err = s.registerValidateInput(ctx, input)
	if err != nil {
		return output, tracer.Error(err)
	}

	var erg errgroup.Group
	erg.Go(func() (err error) {
		presignedUrlOutputs, err := s.s3Repository.CreatePresignedUrl(ctx, s3.CreatePresignedUrlInput{
			BucketName: conf.GetConfig().Minio.PrivateBucket,
			Path:       input.PhotoKtp.GeneratedFileName,
			MimeType:   string(input.PhotoKtp.MimeType),
			Checksum:   input.PhotoKtp.ChecksumSHA256,
		})
		if err != nil {
			return tracer.Error(err)
		}

		output.PhotoKtpPresignedFileUploadOutput = primitive.PresignedFileUploadOutput{
			Identifier:      input.PhotoKtp.Identifier,
			UploadURL:       presignedUrlOutputs.URL,
			UploadExpiredAt: presignedUrlOutputs.ExpiredAt,
			MinioFormData:   presignedUrlOutputs.MinioFormData,
		}

		return nil
	})

	erg.Go(func() (err error) {
		presignedUrlOutputs, err := s.s3Repository.CreatePresignedUrl(ctx, s3.CreatePresignedUrlInput{
			BucketName: conf.GetConfig().Minio.PrivateBucket,
			Path:       input.PhotoSelfie.GeneratedFileName,
			MimeType:   string(input.PhotoSelfie.MimeType),
			Checksum:   input.PhotoSelfie.ChecksumSHA256,
		})
		if err != nil {
			return tracer.Error(err)
		}

		output.PhotoSelfiePresignedFileUploadOutput = primitive.PresignedFileUploadOutput{
			Identifier:      input.PhotoSelfie.Identifier,
			UploadURL:       presignedUrlOutputs.URL,
			UploadExpiredAt: presignedUrlOutputs.ExpiredAt,
			MinioFormData:   presignedUrlOutputs.MinioFormData,
		}

		return nil
	})

	if err = erg.Wait(); err != nil {
		return output, tracer.Error(err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return output, tracer.Error(err)
	}

	err = s.dbTx.DoTransaction(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false},
		func(tx db.Rdbms) (err error) {
			createUserOutput, err := s.userRepository.Create(ctx, users.CreateInput{
				Transaction: tx,
				Email:       input.Email,
				Password:    string(passwordHash),
			})
			if err != nil {
				return tracer.Error(err)
			}

			createConsumerOutput, err := s.consumerRepository.Create(ctx, consumers.CreateInput{
				Transaction:  tx,
				UserID:       createUserOutput.ID,
				Nik:          input.Nik,
				FullName:     input.FullName,
				LegalName:    input.LegalName,
				PlaceOfBirth: input.PlaceOfBirth,
				DateOfBirth:  input.DateOfBirth,
				Salary:       input.Salary,
				PhotoKTP:     input.PhotoKtp.GeneratedFileName,
				PhotoSelfie:  input.PhotoSelfie.GeneratedFileName,
			})
			if err != nil {
				return tracer.Error(err)
			}

			err = s.limitRepository.Creates(ctx, limits.CreatesInput{
				Transaction: tx,
				ConsumerID:  createConsumerOutput.ID,
				Items:       limits.DefaultLimitData(),
			})
			if err != nil {
				return tracer.Error(err)
			}

			output.UserID = createUserOutput.ID
			output.ConsumerID = createConsumerOutput.ID
			return
		},
	)
	if err != nil {
		return output, tracer.Error(err)
	}

	return
}

func (s *service) registerValidateInput(ctx context.Context, input RegisterInput) (err error) {
	var erg errgroup.Group

	erg.Go(func() (err error) {
		outputExisting, err := s.consumerRepository.CheckExisting(ctx, consumers.CheckExistingInput{
			ByNIK: null.StringFrom(input.Nik),
		})
		if err != nil {
			return tracer.Error(err)
		}

		if outputExisting.Existing {
			return tracer.Error(ErrNikIsAvailable)
		}
		return nil
	})

	erg.Go(func() (err error) {
		outputExisting, err := s.userRepository.CheckExisting(ctx, users.CheckExistingInput{
			ByEmail: null.StringFrom(input.Email),
		})
		if err != nil {
			return tracer.Error(err)
		}

		if outputExisting.Existing {
			return tracer.Error(ErrEmailIsAvailable)
		}
		return nil
	})

	if err = erg.Wait(); err != nil {
		return tracer.Error(err)
	}
	return
}
