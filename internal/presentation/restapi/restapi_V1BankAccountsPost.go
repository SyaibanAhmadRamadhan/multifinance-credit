package restapi

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/bank_account"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/guregu/null/v5"
	"net/http"
)

func (h *restApi) V1BankAccountsPost(w http.ResponseWriter, r *http.Request) {
	req := api.V1BankAccountsPostRequestBody{}

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

	createBankAccountsInput := bank_account.CreatesInput{
		ConsumerID: consumerOutput.ID,
		Items:      make([]bank_account.CreatesInputItem, 0),
	}

	for _, item := range req.Items {
		createBankAccountsInput.Items = append(createBankAccountsInput.Items, bank_account.CreatesInputItem{
			Name:              item.Name,
			AccountNumber:     item.AccountNumber,
			AccountHolderName: item.AccountHolderName,
		})
	}

	err = h.bankAccountService.Creates(r.Context(), createBankAccountsInput)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
