package restapi

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/product"
	"net/http"
)

func (h *restApi) V1ProductPost(w http.ResponseWriter, r *http.Request) {
	req := api.V1ProductPostRequestBody{}

	if !h.bodyRequestBindToStruct(w, r, &req) {
		return
	}

	fileUpload, ok := h.bindUploadFileRequest(w, r, req.Image)
	if !ok {
		return
	}

	createProductOutput, err := h.productService.Create(r.Context(), product.CreateInput{
		MerchantID: req.MerchantId,
		Name:       req.Name,
		Image:      fileUpload,
		Qty:        req.Qty,
		Price:      req.Price,
	})
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := api.V1ProductPostResponseBody{
		PresignedImageUpload: h.bindUploadFileResponse(createProductOutput.Image),
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
