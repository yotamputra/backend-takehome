package config

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	Log      *zerolog.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Mux      *http.ServeMux
}

func Bootstrap(config *BootstrapConfig) {
}
