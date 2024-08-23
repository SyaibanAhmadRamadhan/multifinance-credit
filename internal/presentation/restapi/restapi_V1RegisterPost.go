package restapi

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	"net/http"
)

func (h *restApi) V1RegisterPost(w http.ResponseWriter, r *http.Request) {

	req := api.V1RegisterPostRequestBody{}

	if !h.bodyRequestBindToStruct(w, r, &req) {
		return
	}

	photoKtpUploadFile, ok := h.bindUploadFileRequest(w, r, req.PhotoKtp)
	if !ok {
		return
	}

	photoSelfieUploadFile, ok := h.bindUploadFileRequest(w, r, req.PhotoSelfie)
	if !ok {
		return
	}

	output, err := h.authService.Register(r.Context(), auth.RegisterInput{
		DateOfBirth:  req.DateOfBirth,
		Email:        req.Email,
		FullName:     req.FullName,
		LegalName:    req.LegalName,
		PhotoKtp:     photoKtpUploadFile,
		PhotoSelfie:  photoSelfieUploadFile,
		PlaceOfBirth: req.PlaceOfBirth,
		Salary:       req.Salary,
		Password:     req.Password,
		Nik:          req.Nik,
	})
	if err != nil {
		if errors.Is(err, auth.ErrNikIsAvailable) {
			Error(w, r, http.StatusBadRequest, err, auth.ErrNikIsAvailable.Error())
		} else if errors.Is(err, auth.ErrEmailIsAvailable) {
			Error(w, r, http.StatusBadRequest, err, auth.ErrEmailIsAvailable.Error())
		} else {
			Error(w, r, http.StatusInternalServerError, err, err.Error())
		}
		return
	}

	resp := api.V1RegisterPost200Response{
		ConsumerId:        output.ConsumerID,
		UploadPhotoKtp:    h.bindUploadFileResponse(output.PhotoKtpPresignedFileUploadOutput),
		UploadPhotoSelfie: h.bindUploadFileResponse(output.PhotoSelfiePresignedFileUploadOutput),
		UserId:            output.UserID,
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
