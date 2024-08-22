package consumer_test

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"testing"
)

func Test_service_Get(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mockConsumerRepository := consumers.NewMockRepository(mock)
	ctx := context.Background()

	s := consumer.NewService(consumer.NewServiceOpts{
		ConsumerRepository: mockConsumerRepository,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := consumer.GetInput{
			UserID:     null.IntFrom(rand.Int63()),
			ConsumerID: null.IntFrom(rand.Int63()),
		}
		expectedOutput := consumer.GetOutput{
			ID:        rand.Int63(),
			UserID:    rand.Int63(),
			FullName:  faker.Name(),
			LegalName: faker.Name(),
		}

		mockConsumerRepository.EXPECT().
			Get(ctx, consumers.GetInput{
				ID:     expectedInput.ConsumerID,
				UserID: expectedInput.UserID,
			}).
			Return(consumers.GetOutput{
				ID:        expectedOutput.ID,
				UserID:    expectedOutput.UserID,
				FullName:  expectedOutput.FullName,
				LegalName: expectedOutput.LegalName,
			}, nil)

		output, err := s.Get(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return error consumer not found", func(t *testing.T) {
		expectedInput := consumer.GetInput{
			UserID:     null.IntFrom(rand.Int63()),
			ConsumerID: null.IntFrom(rand.Int63()),
		}

		mockConsumerRepository.EXPECT().
			Get(ctx, consumers.GetInput{
				ID:     expectedInput.ConsumerID,
				UserID: expectedInput.UserID,
			}).
			Return(consumers.GetOutput{}, datastore.ErrRecordNotFound)

		output, err := s.Get(ctx, expectedInput)
		require.ErrorIs(t, err, consumer.ErrConsumerNotFound)
		require.Empty(t, output)
	})
}
