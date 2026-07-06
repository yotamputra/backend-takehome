package repository

import (
	"app/internal/entity"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BlogRepository struct {
	Log *zerolog.Logger
}

func NewBlogRepository(log *zerolog.Logger) *BlogRepository {
	return &BlogRepository{
		Log: log,
	}
}

func (r *BlogRepository) Create(db *gorm.DB, blog *entity.Blog) error {
	return db.Create(blog).Error
}

func (r *BlogRepository) FindById(db *gorm.DB, id int) (*entity.Blog, error) {
	var blog entity.Blog
	if err := db.Preload("Author").Preload("Comments").First(&blog, id).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) FindAll(db *gorm.DB) ([]entity.Blog, error) {
	var blogs []entity.Blog
	if err := db.Preload("Author").Preload("Comments").Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *BlogRepository) Update(db *gorm.DB, blog *entity.Blog) error {
	return db.Save(blog).Error
}

func (r *BlogRepository) Delete(db *gorm.DB, blog *entity.Blog) error {
	return db.Delete(blog).Error
}
