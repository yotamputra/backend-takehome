package repository

import (
	"app/internal/entity"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CommentRepository struct {
	Log *zerolog.Logger
}

func NewCommentRepository(log *zerolog.Logger) *CommentRepository {
	return &CommentRepository{
		Log: log,
	}
}

func (r *CommentRepository) Create(db *gorm.DB, comment *entity.Comment) error {
	return db.Create(comment).Error
}

func (r *CommentRepository) FindByPostId(db *gorm.DB, postId string) ([]entity.Comment, error) {
	var comments []entity.Comment
	err := db.Where("post_id = ?", postId).Order("created_at desc").Find(&comments).Error
	return comments, err
}
