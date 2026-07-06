package route

import (
	"app/internal/delivery/http"
	"app/internal/delivery/middleware"
	"app/internal/model"
	"encoding/json"
	netHttp "net/http"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	Mux    *netHttp.ServeMux
	Config *viper.Viper
	Log    *zerolog.Logger

	AuthController    *http.AuthController
	BlogController    *http.BlogController
	CommentController *http.CommentController
}

func (c *RouteConfig) Setup() {
	apiPrefix := "/api"

	c.SetupAuthRoute(apiPrefix)
	c.SetupPostRoute(apiPrefix)

	c.Mux.HandleFunc("/health", HealthCheck)

	c.Mux.HandleFunc("/", func(w netHttp.ResponseWriter, r *netHttp.Request) {
		msg := "Error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(netHttp.StatusNotFound)
		json.NewEncoder(w).Encode(model.WebResponse[any]{
			Message: &msg,
			Errors:  "Not Found",
		})
	})
}

func (c *RouteConfig) SetupAuthRoute(apiPrefix string) {
	authPrefix := apiPrefix + "/auth"

	c.Mux.HandleFunc(authPrefix+"/register", c.AuthController.Register)
	c.Mux.HandleFunc(authPrefix+"/login", c.AuthController.Login)
}

func (c *RouteConfig) SetupPostRoute(apiPrefix string) {
	postPrefix := apiPrefix + "/posts"
	authMiddleware := middleware.NewAuthMiddleware(c.Config, c.Log)

	// Public routes
	c.Mux.HandleFunc("GET "+postPrefix, c.BlogController.GetAll)
	c.Mux.HandleFunc("GET "+postPrefix+"/{id}", c.BlogController.GetById)

	// Protected routes
	c.Mux.Handle("POST "+postPrefix, authMiddleware.Handle(netHttp.HandlerFunc(c.BlogController.Create)))
	c.Mux.Handle("PUT "+postPrefix+"/{id}", authMiddleware.Handle(netHttp.HandlerFunc(c.BlogController.Update)))
	c.Mux.Handle("DELETE "+postPrefix+"/{id}", authMiddleware.Handle(netHttp.HandlerFunc(c.BlogController.Delete)))

	// Comment routes
	c.Mux.HandleFunc("GET "+postPrefix+"/{id}/comments", c.CommentController.GetByPostId)
	c.Mux.Handle("POST "+postPrefix+"/{id}/comments", authMiddleware.Handle(netHttp.HandlerFunc(c.CommentController.Create)))
}

// HealthCheck godoc
// @Summary      Health Check
// @Description  Check if the server is running
// @Tags         health
// @Produce      json
// @Router       /health [get]
func HealthCheck(w netHttp.ResponseWriter, r *netHttp.Request) {
	if r.Method != netHttp.MethodGet {
		netHttp.Error(w, "Method Not Allowed", netHttp.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(netHttp.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
