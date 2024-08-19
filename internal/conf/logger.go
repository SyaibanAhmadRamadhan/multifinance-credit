package conf

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func newLogger() {
	var level zerolog.Level
	switch conf.AppMode {
	case "prod":
		level = zerolog.InfoLevel
	default:
		level = zerolog.DebugLevel
	}

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
		FormatLevel: func(i interface{}) string {
			return "[" + i.(string) + "]"
		},
		FormatMessage: func(i interface{}) string {
			return i.(string)
		},
		FormatFieldName: func(i interface{}) string {
			return i.(string) + ":"
		},
		FormatFieldValue: func(i interface{}) string {
			return i.(string)
		},
	}).With().
		Caller().
		Timestamp().
		Logger()

	zerolog.SetGlobalLevel(level)
}
