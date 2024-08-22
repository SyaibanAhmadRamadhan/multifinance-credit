package restapi

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/guregu/null/v5"
	"io"
	"net/http"
)

func (h *restApi) ImageDisplay(w http.ResponseWriter, r *http.Request) {
	queryParam := api.V1ImageDisplayParams{}

	if !h.queryParamBindToStruct(w, r, &queryParam) {
		return
	}

	output, err := h.consumerService.GetPrivateImage(r.Context(), consumer.GetPrivateImageInput{
		ConsumerID:  null.IntFromPtr(queryParam.ConsumerId),
		UserID:      queryParam.UserId,
		ImageKtp:    null.BoolFromPtr(queryParam.ImageKtp),
		ImageSelfie: null.BoolFromPtr(queryParam.ImageSelfie),
	})
	if err != nil {
		if errors.Is(err, consumer.ErrConsumerNotFound) {
			Error(w, r, http.StatusNotFound, err, consumer.ErrConsumerNotFound.Error())
		} else {
			Error(w, r, http.StatusInternalServerError, err, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	if _, err = io.Copy(w, output.Object); err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
}
