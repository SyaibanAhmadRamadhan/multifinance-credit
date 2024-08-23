package conf

import "time"

type Config struct {
	AppPort       int                 `mapstructure:"APP_PORT"`
	AppMode       string              `mapstructure:"APP_MODE"`
	OpenTelemetry ConfigOpenTelemetry `mapstructure:"OPEN_TELEMETRY"`
	Minio         ConfigMinio         `mapstructure:"MINIO"`
	DatabaseDSN   string              `mapstructure:"DATABASE_DSN"`
	Jwt           ConfigJwt           `mapstructure:"JWT"`
}

type ConfigOpenTelemetry struct {
	Password   string `mapstructure:"PASSWORD"`
	Username   string `mapstructure:"USERNAME"`
	Endpoint   string `mapstructure:"ENDPOINT"`
	TracerName string `mapstructure:"TRACER_NAME"`
}

type ConfigMinio struct {
	Endpoint        string `mapstructure:"ENDPOINT"`
	AccessID        string `mapstructure:"ACCESS_ID"`
	SecretAccessKey string `mapstructure:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `mapstructure:"USE_SSL"`
	PrivateBucket   string `mapstructure:"PRIVATE_BUCKET"`
}

type ConfigJwt struct {
	HS256 ConfigJwtHS256 `mapstructure:"HS256"`
}
type ConfigJwtHS256 struct {
	AccessToken struct {
		Expired time.Duration `mapstructure:"EXPIRED"`
		Key     string        `mapstructure:"KEY"`
	} `mapstructure:"ACCESS_TOKEN"`

	RefreshToken struct {
		Expired time.Duration `mapstructure:"EXPIRED"`
		Key     string        `mapstructure:"KEY"`
	} `mapstructure:"REFRESH_TOKEN"`
}
