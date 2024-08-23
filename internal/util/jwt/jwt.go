package jwt

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"time"
)

type Jwt struct {
	UserID int64
	Key    string
	Exp    time.Duration
}

func HS256AccessTokenDefault(userID int64) Jwt {
	accessToken := conf.GetConfig().Jwt.HS256.AccessToken
	return Jwt{
		UserID: userID,
		Key:    accessToken.Key,
		Exp:    accessToken.Expired,
	}
}

func HS256RefreshTokenDefault(userID int64) Jwt {
	refreshToken := conf.GetConfig().Jwt.HS256.RefreshToken
	return Jwt{
		UserID: userID,
		Key:    refreshToken.Key,
		Exp:    refreshToken.Expired,
	}
}
