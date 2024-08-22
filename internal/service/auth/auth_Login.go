package auth

import (
	"context"
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/jwt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/guregu/null/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"time"
)

func (s *service) Login(ctx context.Context, input LoginInput) (output LoginOutput, err error) {
	userOutput, err := s.userRepository.Get(ctx, users.GetInput{
		Email: null.StringFrom(input.Email),
	})
	if err != nil {
		if errors.Is(err, datastore.ErrRecordNotFound) {
			err = ErrUserNotFound
		}
		return output, tracer.Error(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userOutput.Password), []byte(input.Password))
	if err != nil {
		return output, tracer.Error(ErrInvalidPassword)
	}

	var erg errgroup.Group

	erg.Go(func() (err error) {
		accessTokenConfigDefault := jwt.HS256AccessTokenDefault(userOutput.ID)
		accessToken, err := jwt.GenerateHS256(accessTokenConfigDefault)
		if err != nil {
			return tracer.Error(err)
		}

		output.AccessToken = LoginOutputToken{
			ExpiredAt: time.Now().UTC().Add(accessTokenConfigDefault.Exp),
			Token:     accessToken,
		}

		return nil
	})

	erg.Go(func() (err error) {
		refreshTokenConfigDefault := jwt.HS256RefreshTokenDefault(userOutput.ID)
		refreshToken, err := jwt.GenerateHS256(refreshTokenConfigDefault)
		if err != nil {
			return tracer.Error(err)
		}

		output.RefreshToken = LoginOutputToken{
			ExpiredAt: time.Now().UTC().Add(refreshTokenConfigDefault.Exp),
			Token:     refreshToken,
		}

		return nil
	})

	if err = erg.Wait(); err != nil {
		return output, tracer.Error(err)
	}

	output.UserID = userOutput.ID
	output.Email = userOutput.Email

	return
}
