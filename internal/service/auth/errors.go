package auth

import "errors"

var ErrNikIsAvailable = errors.New("Nik is available")
var ErrEmailIsAvailable = errors.New("email is available")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInvalidTokenType = errors.New("invalid token type")
var ErrInvalidToken = errors.New("invalid token")
var ErrTokenIsExpired = errors.New("token is expired")
