package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	jwtutil "github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/jwt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/golang-jwt/jwt/v5"
)

func (s *service) VerifyToken(ctx context.Context, input VerifyTokenInput) (output VerifyTokenOutput, err error) {
	tokenKey := ""
	switch input.TokenType {
	case primitive.TokenTypeAccessToken:
		tokenKey = conf.GetConfig().Jwt.HS256.AccessToken.Key
	case primitive.TokenTypeRefreshToken:
		tokenKey = conf.GetConfig().Jwt.HS256.RefreshToken.Key
	default:
		return output, tracer.Error(ErrInvalidTokenType)
	}

	claims, err := jwtutil.ClaimsHS256(input.Token, tokenKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = ErrTokenIsExpired
		} else {
			err = fmt.Errorf("err from jwt: %s. %w", err.Error(), ErrInvalidToken)
		}
		return output, tracer.Error(err)
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return output, tracer.Error(fmt.Errorf("cannot type assertion sub to int64. %w", ErrInvalidToken))
	}

	output = VerifyTokenOutput{
		UserID: int64(userID),
	}
	return
}
