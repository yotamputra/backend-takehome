package config

import (
	"app/internal/delivery/http"
	"app/internal/delivery/http/route"
	"app/internal/repository"
	"app/internal/usecase"
	netHttp "net/http"

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
	Mux      *netHttp.ServeMux
}

func Bootstrap(config *BootstrapConfig) {
	// Add health check route
	config.Mux.HandleFunc("/health", func(w netHttp.ResponseWriter, r *netHttp.Request) {
		if r.Method != netHttp.MethodGet {
			netHttp.Error(w, "Method Not Allowed", netHttp.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(netHttp.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Repositories
	userRepository := repository.NewUserRepository(config.Log)

	// UseCases
	authUseCase := usecase.NewAuthUseCase(config.DB, config.Log, config.Validate, userRepository, config.Config)

	// Controllers
	authController := http.NewAuthController(authUseCase, config.Log)

	// Routes
	routeConfig := route.RouteConfig{
		Mux:            config.Mux,
		Config:         config.Config,
		AuthController: authController,
	}
	routeConfig.Setup()
}
