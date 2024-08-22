package auth_test

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	jwtutil "github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/jwt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_service_VerifyToken(t *testing.T) {
	s := auth.NewService(auth.NewServiceOpts{})
	ctx := context.TODO()
	conf.Init()
	conf.GetConfig().Jwt.HS256.AccessToken.Expired = 5 * time.Minute

	t.Run("should be return correct", func(t *testing.T) {
		expectedUserID := int64(1)
		expectedToken, err := jwtutil.GenerateHS256(jwtutil.HS256AccessTokenDefault(expectedUserID))
		require.NoError(t, err)

		output, err := s.VerifyToken(ctx, auth.VerifyTokenInput{
			Token:     expectedToken,
			TokenType: primitive.TokenTypeAccessToken,
		})
		require.NoError(t, err)
		require.Equal(t, expectedUserID, output.UserID)
	})

	t.Run("should be return error token expired", func(t *testing.T) {
		conf.GetConfig().Jwt.HS256.AccessToken.Expired = -5 * time.Minute
		expectedUserID := int64(1)
		expectedToken, err := jwtutil.GenerateHS256(jwtutil.HS256AccessTokenDefault(expectedUserID))
		require.NoError(t, err)

		output, err := s.VerifyToken(ctx, auth.VerifyTokenInput{
			Token:     expectedToken,
			TokenType: primitive.TokenTypeAccessToken,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrTokenIsExpired)
		require.Empty(t, output)
	})

	t.Run("should be return error invalid token type", func(t *testing.T) {
		conf.GetConfig().Jwt.HS256.AccessToken.Expired = -5 * time.Minute
		expectedUserID := int64(1)
		expectedToken, err := jwtutil.GenerateHS256(jwtutil.HS256AccessTokenDefault(expectedUserID))
		require.NoError(t, err)

		output, err := s.VerifyToken(ctx, auth.VerifyTokenInput{
			Token:     expectedToken,
			TokenType: primitive.TokenTypeUnknown,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrInvalidTokenType)
		require.Empty(t, output)
	})

	t.Run("should be return error invalid key token", func(t *testing.T) {
		conf.GetConfig().Jwt.HS256.AccessToken.Expired = 5 * time.Minute
		expectedUserID := int64(1)
		expectedToken, err := jwtutil.GenerateHS256(jwtutil.HS256AccessTokenDefault(expectedUserID))
		require.NoError(t, err)

		output, err := s.VerifyToken(ctx, auth.VerifyTokenInput{
			Token:     expectedToken,
			TokenType: primitive.TokenTypeRefreshToken,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrInvalidToken)
		require.Empty(t, output)
	})
}
