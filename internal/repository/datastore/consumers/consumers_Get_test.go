package consumers_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func Test_repository_Get(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewRdbms(sqlxDB)

	r := consumers.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := consumers.GetInput{
			ID:     null.IntFrom(rand.Int63()),
			UserID: null.IntFrom(rand.Int63()),
		}
		expectedOutput := consumers.GetOutput{
			ID:           expectedInput.ID.Int64,
			UserID:       expectedInput.UserID.Int64,
			Nik:          faker.UUIDDigit(),
			FullName:     faker.Name(),
			LegalName:    faker.Name(),
			PlaceOfBirth: faker.UUIDDigit(),
			DateOfBirth:  time.Now().UTC(),
			Salary:       rand.Float64(),
			PhotoKTP:     faker.UUIDDigit(),
			PhotoSelfie:  faker.UUIDDigit(),
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, user_id, nik, full_name, legal_name, place_of_birth, date_of_birth, salary, photo_ktp, photo_selfie 
					FROM consumers WHERE user_id = ? AND id = ?`,
		)).WithArgs(expectedInput.UserID.Int64, expectedInput.ID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "user_id", "nik", "full_name", "legal_name", "place_of_birth", "date_of_birth", "salary", "photo_ktp", "photo_selfie",
			}).AddRow(
				expectedOutput.ID, expectedOutput.UserID, expectedOutput.Nik, expectedOutput.FullName, expectedOutput.LegalName,
				expectedOutput.PlaceOfBirth, expectedOutput.DateOfBirth, expectedOutput.Salary, expectedOutput.PhotoKTP,
				expectedOutput.PhotoSelfie,
			))

		output, err := r.Get(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should be return error no record found", func(t *testing.T) {
		expectedInput := consumers.GetInput{
			ID:     null.IntFrom(rand.Int63()),
			UserID: null.IntFrom(rand.Int63()),
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, user_id, nik, full_name, legal_name, place_of_birth, date_of_birth, salary, photo_ktp, photo_selfie 
					FROM consumers WHERE user_id = ? AND id = ?`,
		)).WithArgs(expectedInput.UserID.Int64, expectedInput.ID.Int64).
			WillReturnError(sql.ErrNoRows)

		output, err := r.Get(ctx, expectedInput)
		require.Error(t, err, datastore.ErrRecordNotFound)
		require.Empty(t, output)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
