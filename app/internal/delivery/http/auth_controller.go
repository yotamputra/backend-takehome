package http

import (
	"app/internal/model"
	"app/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

type AuthController struct {
	AuthUseCase *usecase.AuthUseCase
	Log         *zerolog.Logger
}

func NewAuthController(authUseCase *usecase.AuthUseCase, log *zerolog.Logger) *AuthController {
	return &AuthController{
		AuthUseCase: authUseCase,
		Log:         log,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.RegisterRequest true "Register Request"
// @Success      200 {object} model.WebResponse[model.UserResponse]
// @Failure      400 {object} model.WebResponse[any]
// @Router       /api/auth/register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var request model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid request body"})
		return
	}

	response, err := c.AuthUseCase.Register(&request)
	if err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: err.Error()})
		return
	}

	c.jsonResponse(w, http.StatusOK, model.WebResponse[*model.UserResponse]{Data: response})
}

// Login godoc
// @Summary      Login
// @Description  Login and receive token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequest true "Login Request"
// @Success      200 {object} model.WebResponse[model.LoginResponse]
// @Failure      400 {object} model.WebResponse[any]
// @Router       /api/auth/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var request model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid request body"})
		return
	}

	response, err := c.AuthUseCase.Login(&request)
	if err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: err.Error()})
		return
	}

	c.jsonResponse(w, http.StatusOK, model.WebResponse[*model.LoginResponse]{Data: response})
}

func (c *AuthController) jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}