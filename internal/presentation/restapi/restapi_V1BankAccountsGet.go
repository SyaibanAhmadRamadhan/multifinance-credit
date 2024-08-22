package restapi

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/bank_account"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/guregu/null/v5"
	"net/http"
)

func (h *restApi) V1BankAccountsGet(w http.ResponseWriter, r *http.Request) {
	params := api.V1BankAccountsGetParams{}
	if !h.queryParamBindToStruct(w, r, &params) {
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

	bankAccountsOutput, err := h.bankAccountService.GetAll(r.Context(), bank_account.GetAllInput{
		ConsumerID: consumerOutput.ID,
		Pagination: h.bindToPaginationInput(params.Page, params.PageSize),
	})
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := api.V1BankAccountsGetResponseBody{
		Items:      make([]api.V1BankAccountsGetResponseBodyItem, 0),
		Pagination: h.bindToPaginationResponse(bankAccountsOutput.Pagination),
	}

	for _, item := range bankAccountsOutput.Items {
		resp.Items = append(resp.Items, api.V1BankAccountsGetResponseBodyItem{
			AccountHolderName: item.AccountHolderName,
			AccountNumber:     item.AccountNumber,
			ConsumerId:        item.ConsumerID,
			Id:                item.ID,
			Name:              item.Name,
		})
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
