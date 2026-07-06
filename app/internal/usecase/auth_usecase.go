package usecase

import (
	"app/internal/entity"
	"app/internal/model"
	"app/internal/repository"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB             *gorm.DB
	Log            *zerolog.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	Config         *viper.Viper
}

func NewAuthUseCase(db *gorm.DB, log *zerolog.Logger, validate *validator.Validate, userRepository *repository.UserRepository, config *viper.Viper) *AuthUseCase {
	return &AuthUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		Config:         config,
	}
}

func (c *AuthUseCase) Register(request *model.RegisterRequest) (*model.UserResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		return nil, err
	}

	_, err := c.UserRepository.FindByEmail(c.DB, request.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Error().Err(err).Msg("failed to hash password")
		return nil, err
	}

	user := &entity.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
	}

	if err := c.UserRepository.Create(c.DB, user); err != nil {
		c.Log.Error().Err(err).Msg("failed to create user")
		return nil, err
	}

	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (c *AuthUseCase) Login(request *model.LoginRequest) (*model.LoginResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		return nil, err
	}

	user, err := c.UserRepository.FindByEmail(c.DB, request.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	expirationStr := c.Config.GetString("JWT_ACCESS_EXPIRATION")
	if expirationStr == "" {
		expirationStr = "5h"
	}
	expirationTime, err := time.ParseDuration(expirationStr)
	if err != nil {
		expirationTime = 5 * time.Hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(expirationTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte(c.Config.GetString("JWT_ACCESS_SECRET")))
	if err != nil {
		c.Log.Error().Err(err).Msg("failed to generate token")
		return nil, err
	}

	return &model.LoginResponse{
		Token:     tokenString,
		ExpiresAt: time.Now().Add(expirationTime).Unix(),
		User: model.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}
