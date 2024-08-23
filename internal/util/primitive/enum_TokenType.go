package primitive

type TokenType int8

const (
	TokenTypeUnknown TokenType = iota
	TokenTypeAccessToken
	TokenTypeRefreshToken
)
