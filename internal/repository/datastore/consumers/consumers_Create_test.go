package consumers_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/go-faker/faker/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func Test_repository_Create(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewRdbms(sqlxDB)

	r := consumers.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := consumers.CreateInput{
			UserID:       rand.Int63(),
			Nik:          faker.UUIDDigit(),
			FullName:     faker.Name(),
			LegalName:    faker.Name(),
			PlaceOfBirth: faker.Name(),
			DateOfBirth:  time.Now().UTC(),
			Salary:       rand.Float64(),
			PhotoKTP:     faker.UUIDDigit(),
			PhotoSelfie:  faker.UUIDDigit(),
		}

		expectedID := rand.Int63()

		mock.ExpectExec(regexp.QuoteMeta(
			`INSERT INTO consumers (user_id,nik,full_name,legal_name,place_of_birth,date_of_birth,salary,photo_ktp,photo_selfie,created_at,updated_at) 
					VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		)).WithArgs(
			expectedInput.UserID, expectedInput.Nik, expectedInput.FullName, expectedInput.LegalName,
			expectedInput.PlaceOfBirth, expectedInput.DateOfBirth, expectedInput.Salary, expectedInput.PhotoKTP, expectedInput.PhotoSelfie,
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(expectedID, 1))

		output, err := r.Create(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedID, output.ID)
	})
}
