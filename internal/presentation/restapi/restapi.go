package restapi

import (
	"encoding/json"
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/bank_account"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"io"
	"net/http"
	"reflect"
)

type restApi struct {
	validator          *validator.Validate
	authService        auth.Service
	consumerService    consumer.Service
	bankAccountService bank_account.Service
}

func New(dependency *service.Dependency) *restApi {
	v := validator.New()
	v.SetTagName("binding")
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("json")
	})

	return &restApi{
		validator:          v,
		authService:        dependency.AuthService,
		consumerService:    dependency.ConsumerService,
		bankAccountService: dependency.BankAccountService,
	}
}

func (h *restApi) bodyRequestBindToStruct(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return false
	}
	defer r.Body.Close()

	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		Error(w, r, http.StatusUnprocessableEntity, err, err.Error())
		return false
	}

	err = h.validator.Struct(v)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return false
	}
	return true
}

func (h *restApi) bindUploadFileRequest(w http.ResponseWriter, r *http.Request, input api.FileUploadRequest) (output primitive.PresignedFileUpload, ok bool) {
	fileUploadInput := primitive.NewPresignedFileUploadInput{
		Identifier:       input.Identifier,
		OriginalFileName: input.OriginalFilename,
		MimeType:         primitive.MimeType(input.MimeType),
		Size:             input.Size,
		ChecksumSHA256:   input.ChecksumSha256,
	}
	output, err := primitive.NewPresignedFileUpload(fileUploadInput)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err, err.Error())
		return output, false
	}
	return output, true
}

func (h *restApi) bindUploadFileResponse(input primitive.PresignedFileUploadOutput) (output api.FileUploadResponse) {
	return api.FileUploadResponse{
		Identifier:      input.Identifier,
		UploadExpiredAt: input.UploadExpiredAt,
		UploadUrl:       input.UploadURL,
		MinioFormData:   input.MinioFormData,
	}
}

func (h *restApi) writeJson(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	respByte, err := json.Marshal(v)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respByte)
}

var decoder = schema.NewDecoder()

func (h *restApi) queryParamBindToStruct(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := r.ParseForm(); err != nil {
		Error(w, r, http.StatusBadRequest, err, err.Error())
		return false
	}

	decoder.SetAliasTag("json")
	if err := decoder.Decode(v, r.Form); err != nil {
		Error(w, r, http.StatusBadRequest, err, err.Error())
		return false
	}

	err := h.validator.Struct(v)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return false
	}
	return true
}

func (h *restApi) getUserID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	userIDAny := r.Context().Value(primitive.UserIDKey)

	if userIDAny == nil {
		Error(w, r, http.StatusUnauthorized, errors.New("no user id in context"))
		return 0, false
	}

	userID, ok := userIDAny.(int64)
	if !ok {
		Error(w, r, http.StatusUnauthorized, errors.New("invalid user id, cannot type assertion to int64"))
		return 0, false
	}

	if userID == 0 {
		Error(w, r, http.StatusUnauthorized, errors.New("invalid user id"))
		return 0, false
	}

	return userID, true
}

func (h *restApi) bindToPaginationInput(page int64, pageSize int64) pagination.PaginationInput {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 20 {
		pageSize = 20
	}

	return pagination.PaginationInput{
		Page:     page,
		PageSize: pageSize,
	}
}

func (h *restApi) bindToPaginationResponse(input pagination.PaginationOutput) api.PaginationResponse {
	return api.PaginationResponse{
		Page:      input.Page,
		PageCount: input.PageCount,
		PageSize:  input.PageSize,
		TotalData: input.TotalData,
	}
}
