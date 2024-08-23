package restapi

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/transaction"
	"github.com/guregu/null/v5"
	"net/http"
)

func (h *restApi) V1TransactionPost(w http.ResponseWriter, r *http.Request) {
	req := api.V1TransactionPostRequestBody{}

	if !h.bodyRequestBindToStruct(w, r, &req) {
		return
	}

	userID, ok := h.getUserID(w, r)
	if !ok {
		return
	}

	consumerOutput, err := h.consumerService.Get(r.Context(), consumer.GetInput{
		UserID: null.IntFrom(userID),
	})
	if err != nil {
		if errors.Is(err, consumer.ErrConsumerNotFound) {
			Error(w, r, http.StatusUnauthorized, err)
		} else {
			Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	inputTransaction := transaction.CreateInput{
		Products:   make([]transaction.CreateInputProductItem, 0),
		LimitID:    req.LimitId,
		ConsumerID: consumerOutput.ID,
	}

	for _, product := range req.Products {
		inputTransaction.Products = append(inputTransaction.Products, transaction.CreateInputProductItem{
			ID:  product.ProductId,
			Qty: product.Qty,
		})
	}

	outputCreateTransaction, err := h.transactionService.Create(r.Context(), inputTransaction)
	if err != nil {
		if errors.Is(err, transaction.ErrLimitNotFound) {
			Error(w, r, http.StatusBadRequest, err, "data limit not found")
		} else if errors.Is(err, transaction.ErrProductNotFound) {
			Error(w, r, http.StatusBadRequest, err, "data product not found")
		} else if errors.Is(err, transaction.ErrStockProductNotAvailable) {
			Error(w, r, http.StatusBadRequest, err, "data stock product not available")
		} else if errors.Is(err, transaction.ErrExceedLimit) {
			Error(w, r, http.StatusBadRequest, err, "exceed limit")
		} else {
			Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	resp := api.V1TransactionPostResponseBody{
		TransactionId: outputCreateTransaction.ID,
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
