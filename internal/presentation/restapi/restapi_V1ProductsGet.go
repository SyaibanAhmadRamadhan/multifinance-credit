package restapi

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/product"
	"github.com/guregu/null/v5"
	"net/http"
)

func (h *restApi) V1ProductsGet(w http.ResponseWriter, r *http.Request) {
	params := api.V1ProductsGetParams{}

	if !h.queryParamBindToStruct(w, r, &params) {
		return
	}

	productsOutput, err := h.productService.GetAll(r.Context(), product.GetAllInput{
		MerchantID: null.IntFromPtr(params.MerchantId),
		IDs:        params.Ids,
		Pagination: h.bindToPaginationInput(params.Page, params.PageSize),
	})
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := api.V1ProductsGetResponseBody{
		Items:      make([]api.V1ProductsGetResponseBodyItem, 0),
		Pagination: h.bindToPaginationResponse(productsOutput.Pagination),
	}

	for _, item := range productsOutput.Items {
		resp.Items = append(resp.Items, api.V1ProductsGetResponseBodyItem{
			Id:         item.ID,
			Image:      item.Image,
			MerchantId: item.MerchantID,
			Name:       item.Name,
			Price:      item.Price,
			Qty:        item.Qty,
		})
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
