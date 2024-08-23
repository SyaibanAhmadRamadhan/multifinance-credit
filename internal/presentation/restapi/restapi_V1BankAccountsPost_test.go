package restapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/bank_account"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/go-chi/chi/v5"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_restApi_V1BankAccountsPost(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mockConsumerService := consumer.NewMockService(mock)
	mockBankAccountService := bank_account.NewMockService(mock)

	h := restapi.New(&service.Dependency{
		ConsumerService:    mockConsumerService,
		BankAccountService: mockBankAccountService,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedUserID := rand.Int63()
		expectedConsumerID := rand.Int63()
		expectedRequestBody := api.V1BankAccountsPostRequestBody{
			Items: []api.V1BankAccountsPostRequestBodyItem{
				{
					AccountHolderName: faker.Name(),
					AccountNumber:     faker.Name(),
					Name:              faker.Name(),
				},
				{
					AccountHolderName: faker.Name(),
					AccountNumber:     faker.Name(),
					Name:              faker.Name(),
				},
			},
		}
		expectedReqJson, err := json.Marshal(expectedRequestBody)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Post("/api/v1/bank-account", h.V1BankAccountsPost)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/bank-account", io.NopCloser(bytes.NewBuffer(expectedReqJson)))
		require.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(), primitive.UserIDKey, expectedUserID))

		mockConsumerService.EXPECT().
			Get(req.Context(), consumer.GetInput{
				UserID: null.IntFrom(expectedUserID),
			}).
			Return(consumer.GetOutput{
				ID: expectedConsumerID,
			}, nil)

		mockBankAccountService.EXPECT().
			Creates(req.Context(), bank_account.CreatesInput{
				ConsumerID: expectedConsumerID,
				Items: []bank_account.CreatesInputItem{
					{
						Name:              expectedRequestBody.Items[0].Name,
						AccountNumber:     expectedRequestBody.Items[0].AccountNumber,
						AccountHolderName: expectedRequestBody.Items[0].AccountHolderName,
					},
					{
						Name:              expectedRequestBody.Items[1].Name,
						AccountNumber:     expectedRequestBody.Items[1].AccountNumber,
						AccountHolderName: expectedRequestBody.Items[1].AccountHolderName,
					},
				},
			}).Return(nil)

		rr := httptest.NewRecorder()
		h.V1BankAccountsPost(rr, req)

		require.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("should be return error consumer not found", func(t *testing.T) {
		expectedUserID := rand.Int63()
		expectedRequestBody := api.V1BankAccountsPostRequestBody{
			Items: []api.V1BankAccountsPostRequestBodyItem{
				{
					AccountHolderName: faker.Name(),
					AccountNumber:     faker.Name(),
					Name:              faker.Name(),
				},
				{
					AccountHolderName: faker.Name(),
					AccountNumber:     faker.Name(),
					Name:              faker.Name(),
				},
			},
		}
		expectedReqJson, err := json.Marshal(expectedRequestBody)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Post("/api/v1/bank-account", h.V1BankAccountsPost)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/bank-account", io.NopCloser(bytes.NewBuffer(expectedReqJson)))
		require.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(), primitive.UserIDKey, expectedUserID))

		mockConsumerService.EXPECT().
			Get(req.Context(), consumer.GetInput{
				UserID: null.IntFrom(expectedUserID),
			}).
			Return(consumer.GetOutput{}, consumer.ErrConsumerNotFound)

		rr := httptest.NewRecorder()
		h.V1BankAccountsPost(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
