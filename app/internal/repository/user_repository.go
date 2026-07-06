package repository

import (
	"app/internal/entity"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserRepository struct {
	Log *zerolog.Logger
}

func NewUserRepository(log *zerolog.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) Create(db *gorm.DB, user *entity.User) error {
	return db.Create(user).Error
}

func (r *UserRepository) FindByEmail(db *gorm.DB, email string) (*entity.User, error) {
	var user entity.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindById(db *gorm.DB, id string) (*entity.User, error) {
	var user entity.User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
