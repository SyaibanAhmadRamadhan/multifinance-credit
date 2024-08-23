package auth

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"time"
)

type RegisterInput struct {
	DateOfBirth  time.Time
	Email        string
	Nik          string
	FullName     string
	LegalName    string
	PhotoKtp     primitive.PresignedFileUpload
	PhotoSelfie  primitive.PresignedFileUpload
	PlaceOfBirth string
	Salary       float64
	Password     string
}

type RegisterOutput struct {
	UserID                               int64
	ConsumerID                           int64
	PhotoKtpPresignedFileUploadOutput    primitive.PresignedFileUploadOutput
	PhotoSelfiePresignedFileUploadOutput primitive.PresignedFileUploadOutput
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  LoginOutputToken
	RefreshToken LoginOutputToken
	UserID       int64
	Email        string
}

type LoginOutputToken struct {
	ExpiredAt time.Time
	Token     string
}

type VerifyTokenInput struct {
	Token     string
	TokenType primitive.TokenType
}

type VerifyTokenOutput struct {
	UserID int64
}
