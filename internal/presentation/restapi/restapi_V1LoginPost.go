package restapi

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"net/http"
)

func (h *restApi) V1LoginPost(w http.ResponseWriter, r *http.Request) {
	req := api.V1LoginPostRequestBody{}

	if !h.bodyRequestBindToStruct(w, r, &req) {
		return
	}

	loginOutput, err := h.authService.Login(r.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			Error(w, r, http.StatusBadRequest, err, "invalid email or password")
		} else if errors.Is(err, auth.ErrInvalidPassword) {
			Error(w, r, http.StatusBadRequest, err, "invalid email or password")
		} else {
			Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	refreshTokenCookie := &http.Cookie{
		Name:     primitive.RefreshTokenCookieKey,
		Value:    loginOutput.RefreshToken.Token,
		Expires:  loginOutput.RefreshToken.ExpiredAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/v1/refresh-token",
	}

	http.SetCookie(w, refreshTokenCookie)

	resp := api.V1LoginPostResponseBody{
		AccessToken: api.V1TokenJwtResponse{
			ExpiredAt: loginOutput.AccessToken.ExpiredAt,
			Token:     loginOutput.AccessToken.Token,
		},
		Email:  loginOutput.Email,
		UserId: loginOutput.UserID,
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
