package main

import (
	"app/internal/config"
	"fmt"
	"net/http"
	"time"

	_ "app/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-playground/validator/v10"
)

// @title           Backend Takehome API
// @version         1.0
// @description     This is a sample blog platform API.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, logger)
	validate := validator.New()

	mux := http.NewServeMux()

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Log:      logger,
		Validate: validate,
		Config:   viperConfig,
		Mux:      mux,
	})

	port := viperConfig.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info().Msgf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Msgf("failed to start server: %v", err)
	}
}
