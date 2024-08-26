package restapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/transaction"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/go-chi/chi/v5"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_restApi_V1TransactionPost(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()
	conf.Init()

	mockTransactionService := transaction.NewMockService(mock)
	mockConsumerService := consumer.NewMockService(mock)

	h := restapi.New(&service.Dependency{
		TransactionService: mockTransactionService,
		ConsumerService:    mockConsumerService,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedReq := api.V1TransactionPostRequestBody{
			LimitId: rand.Int63(),
			Products: []api.V1TransactionPostRequestBodyProductItem{
				{
					ProductId: rand.Int63(),
					Qty:       rand.Int63(),
				},
			},
		}
		expectedTransactionID := rand.Int63()
		expectedConsumerID := rand.Int63()
		expectedUserID := rand.Int63()
		expectedReqBytes, err := json.Marshal(expectedReq)
		require.NoError(t, err)

		r := chi.NewRouter()
		r.Post("/api/v1/transaction", h.V1TransactionPost)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/transaction", io.NopCloser(bytes.NewReader(expectedReqBytes)))
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), primitive.UserIDKey, expectedUserID))

		mockConsumerService.EXPECT().
			Get(req.Context(), consumer.GetInput{
				UserID: null.IntFrom(expectedUserID),
			}).
			Return(consumer.GetOutput{
				ID: expectedConsumerID,
			}, nil)

		mockTransactionService.EXPECT().
			Create(req.Context(), transaction.CreateInput{
				Products: []transaction.CreateInputProductItem{
					{
						ID:  expectedReq.Products[0].ProductId,
						Qty: expectedReq.Products[0].Qty,
					},
				},
				LimitID:    expectedReq.LimitId,
				ConsumerID: expectedConsumerID,
			}).Return(transaction.CreateOutput{
			ID: expectedTransactionID,
		}, nil)

		rr := httptest.NewRecorder()
		h.V1TransactionPost(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp api.V1TransactionPostResponseBody
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, api.V1TransactionPostResponseBody{
			TransactionId: expectedTransactionID,
		}, resp)
	})
}
