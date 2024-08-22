package jwt

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"time"
)

func GenerateHS256(jwtModel Jwt) (string, error) {
	timeNow := time.Now()
	timeExp := timeNow.Add(jwtModel.Exp).Unix()

	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": timeExp,
		"sub": jwtModel.UserID,
	})

	tokenStr, err := tokenParse.SignedString([]byte(jwtModel.Key))
	if err != nil {
		return "", tracer.Error(err)
	}

	return tokenStr, nil
}

func ClaimsHS256(tokenStr, key string) (map[string]any, error) {
	tokenParse, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Debug().Msgf("unexpected signing method : %v", t.Header["alg"])
			return nil, errors.New("invalid Token")
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, tracer.Error(err)
	}

	claims, _ := tokenParse.Claims.(jwt.MapClaims)

	return claims, nil
}
