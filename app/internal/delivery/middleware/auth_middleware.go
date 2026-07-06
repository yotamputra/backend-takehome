package middleware

import (
	"app/internal/model"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type AuthMiddleware struct {
	Config *viper.Viper
	Log    *zerolog.Logger
}

func NewAuthMiddleware(config *viper.Viper, log *zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		Config: config,
		Log:    log,
	}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			m.unauthorized(w)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.Config.GetString("JWT_ACCESS_SECRET")), nil
		})

		if err != nil || !token.Valid {
			m.Log.Warn().Err(err).Msg("Failed to parse or validate token")
			m.unauthorized(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.unauthorized(w)
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			m.unauthorized(w)
			return
		}

		name, ok := claims["name"].(string)
		if !ok {
			name = "Unknown"
		}

		auth := &model.Auth{
			ID:   int(sub),
			Name: name,
		}

		ctx := context.WithValue(r.Context(), "auth", auth)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) unauthorized(w http.ResponseWriter) {
	msg := "Unauthorized"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(model.WebResponse[any]{
		Errors:  "Unauthorized",
		Message: &msg,
	})
}

func GetUser(r *http.Request) *model.Auth {
	auth, ok := r.Context().Value("auth").(*model.Auth)
	if !ok {
		return nil
	}
	return auth
}
