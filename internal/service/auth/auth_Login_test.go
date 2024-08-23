package auth_test

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"testing"
)

func Test_service_Login(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.TODO()

	conf.Init()

	mockUserRepository := users.NewMockRepository(mock)

	s := auth.NewService(auth.NewServiceOpts{
		UserRepository: mockUserRepository,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := auth.LoginInput{
			Email:    faker.Email(),
			Password: faker.Password(),
		}
		expectedUserID := rand.Int63()
		expectedHashedPassword, err := bcrypt.GenerateFromPassword([]byte(expectedInput.Password), 10)
		require.NoError(t, err)

		mockUserRepository.EXPECT().
			Get(ctx, users.GetInput{
				Email: null.StringFrom(expectedInput.Email),
			}).
			Return(users.GetOutput{
				ID:       expectedUserID,
				Email:    expectedInput.Email,
				Password: string(expectedHashedPassword),
			}, nil)

		output, err := s.Login(ctx, expectedInput)
		require.NoError(t, err)
		require.NotEmpty(t, output)
		require.Equal(t, expectedUserID, output.UserID)
		require.Equal(t, expectedInput.Email, output.Email)
	})

	t.Run("should be return error invalid password", func(t *testing.T) {
		expectedInput := auth.LoginInput{
			Email:    faker.Email(),
			Password: faker.Password(),
		}
		expectedUserID := rand.Int63()
		expectedHashedPassword, err := bcrypt.GenerateFromPassword([]byte("invalid pw"), 10)
		require.NoError(t, err)

		mockUserRepository.EXPECT().
			Get(ctx, users.GetInput{
				Email: null.StringFrom(expectedInput.Email),
			}).
			Return(users.GetOutput{
				ID:       expectedUserID,
				Email:    expectedInput.Email,
				Password: string(expectedHashedPassword),
			}, nil)

		output, err := s.Login(ctx, expectedInput)
		require.ErrorIs(t, err, auth.ErrInvalidPassword)
		require.Empty(t, output)
	})

	t.Run("should be return error user not found", func(t *testing.T) {
		expectedInput := auth.LoginInput{
			Email:    faker.Email(),
			Password: faker.Password(),
		}

		mockUserRepository.EXPECT().
			Get(ctx, users.GetInput{
				Email: null.StringFrom(expectedInput.Email),
			}).
			Return(users.GetOutput{}, datastore.ErrRecordNotFound)

		output, err := s.Login(ctx, expectedInput)
		require.ErrorIs(t, err, auth.ErrUserNotFound)
		require.Empty(t, output)
	})
}
