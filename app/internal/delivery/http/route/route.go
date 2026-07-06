package route

import (
	"app/internal/delivery/http"
	"app/internal/model"
	"encoding/json"
	netHttp "net/http"

	"github.com/spf13/viper"
)

type RouteConfig struct {
	Mux    *netHttp.ServeMux
	Config *viper.Viper

	AuthController *http.AuthController
}

func (c *RouteConfig) Setup() {
	apiPrefix := "/api"

	c.SetupAuthRoute(apiPrefix)

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
