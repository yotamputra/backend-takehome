package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func NewLogger(v *viper.Viper) *zerolog.Logger {
	env := v.GetString("APP_ENV")

	var log zerolog.Logger
	if env == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		log = zerolog.New(os.Stderr).
			With().
			Timestamp().
			Logger().
			Level(zerolog.ErrorLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).
			With().
			Timestamp().
			Logger()
	}

	return &log
}
