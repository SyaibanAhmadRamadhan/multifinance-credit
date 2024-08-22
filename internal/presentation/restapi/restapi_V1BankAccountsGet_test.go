package restapi_test

import (
	"context"
	"encoding/json"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/bank_account"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/go-chi/chi/v5"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func Test_restApi_V1BankAccountsGet(t *testing.T) {
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
		expectedPagination := pagination.PaginationInput{
			Page:     1,
			PageSize: 2,
		}
		expectedTotalData := int64(2)
		expectedBankAccountOutput := bank_account.GetAllOutput{
			Pagination: pagination.CreatePaginationOutput(expectedPagination, expectedTotalData),
			Items: []bank_account.GetAllOutputItem{
				{
					ID:                rand.Int63(),
					ConsumerID:        expectedConsumerID,
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
				{
					ID:                rand.Int63(),
					ConsumerID:        expectedConsumerID,
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
			},
		}

		params := url.Values{}
		params.Add("page", strconv.FormatInt(1, 10))
		params.Add("page_size", strconv.FormatInt(2, 10))

		r := chi.NewRouter()
		r.Get("/api/v1/bank-account", h.V1BankAccountsGet)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/bank-account?"+params.Encode(), nil)
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
			GetAll(req.Context(), bank_account.GetAllInput{
				ConsumerID: expectedConsumerID,
				Pagination: expectedPagination,
			}).Return(expectedBankAccountOutput, nil)

		expectedResp := api.V1BankAccountsGetResponseBody{
			Items: []api.V1BankAccountsGetResponseBodyItem{
				{
					AccountHolderName: expectedBankAccountOutput.Items[0].AccountHolderName,
					AccountNumber:     expectedBankAccountOutput.Items[0].AccountNumber,
					ConsumerId:        expectedConsumerID,
					Id:                expectedBankAccountOutput.Items[0].ID,
					Name:              expectedBankAccountOutput.Items[0].Name,
				},
				{
					AccountHolderName: expectedBankAccountOutput.Items[1].AccountHolderName,
					AccountNumber:     expectedBankAccountOutput.Items[1].AccountNumber,
					ConsumerId:        expectedConsumerID,
					Id:                expectedBankAccountOutput.Items[1].ID,
					Name:              expectedBankAccountOutput.Items[1].Name,
				},
			},
			Pagination: api.PaginationResponse{
				Page:      1,
				PageCount: 1,
				PageSize:  2,
				TotalData: 2,
			},
		}
		rr := httptest.NewRecorder()
		h.V1BankAccountsGet(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp api.V1BankAccountsGetResponseBody
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Equal(t, expectedResp, resp)
	})

	t.Run("should be return error consumer not found", func(t *testing.T) {
		expectedUserID := rand.Int63()

		r := chi.NewRouter()
		r.Get("/api/v1/bank-account", h.V1BankAccountsGet)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/bank-account", nil)
		require.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(), primitive.UserIDKey, expectedUserID))

		mockConsumerService.EXPECT().
			Get(req.Context(), consumer.GetInput{
				UserID: null.IntFrom(expectedUserID),
			}).
			Return(consumer.GetOutput{}, consumer.ErrConsumerNotFound)

		rr := httptest.NewRecorder()
		h.V1BankAccountsGet(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
