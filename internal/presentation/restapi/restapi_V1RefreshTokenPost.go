package restapi

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	jwtutil "github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/jwt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"net/http"
	"time"
)

func (h *restApi) V1RefreshTokenPost(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(primitive.RefreshTokenCookieKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			Error(w, r, http.StatusUnauthorized, err)
		} else {
			Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	verifyTokenOutput, err := h.authService.VerifyToken(r.Context(), auth.VerifyTokenInput{
		Token:     c.Value,
		TokenType: primitive.TokenTypeRefreshToken,
	})
	if err != nil {
		if errors.Is(err, auth.ErrTokenIsExpired) {
			Error(w, r, http.StatusUnauthorized, err, "your token is expired, you can log in again")
		} else if errors.Is(err, auth.ErrInvalidToken) {
			Error(w, r, http.StatusUnauthorized, err, "your token is invalid, you can log in again")
		} else {
			Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	accessTokenConfDefault := jwtutil.HS256AccessTokenDefault(verifyTokenOutput.UserID)
	accessToken, err := jwtutil.GenerateHS256(accessTokenConfDefault)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := api.V1RefreshTokenPostResponseBody{
		AccessToken: api.V1TokenJwtResponse{
			ExpiredAt: time.Now().UTC().Add(accessTokenConfDefault.Exp),
			Token:     accessToken,
		},
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
