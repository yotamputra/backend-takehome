package usecase

import (
	"app/internal/entity"
	"app/internal/model"
	"app/internal/repository"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CommentUseCase struct {
	DB                *gorm.DB
	Log               *zerolog.Logger
	Validate          *validator.Validate
	CommentRepository *repository.CommentRepository
	BlogRepository    *repository.BlogRepository
}

func NewCommentUseCase(db *gorm.DB, log *zerolog.Logger, validate *validator.Validate, commentRepository *repository.CommentRepository, blogRepository *repository.BlogRepository) *CommentUseCase {
	return &CommentUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		CommentRepository: commentRepository,
		BlogRepository:    blogRepository,
	}
}

func (c *CommentUseCase) toCommentResponse(comment *entity.Comment) *model.CommentResponse {
	return &model.CommentResponse{
		ID:         comment.ID,
		PostID:     comment.PostID,
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt.Unix(),
	}
}

func (c *CommentUseCase) Create(request *model.CreateCommentRequest, postId int, authorName string) (*model.CommentResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		return nil, err
	}

	// Verify post exists
	_, err := c.BlogRepository.FindById(c.DB, postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		c.Log.Error().Err(err).Msg("failed to find post")
		return nil, errors.New("internal server error")
	}

	comment := &entity.Comment{
		PostID:     postId,
		AuthorName: authorName,
		Content:    request.Content,
	}

	if err := c.CommentRepository.Create(c.DB, comment); err != nil {
		c.Log.Error().Err(err).Msg("failed to create comment")
		return nil, errors.New("failed to create comment")
	}

	return c.toCommentResponse(comment), nil
}

func (c *CommentUseCase) GetByPostId(postId int) ([]model.CommentResponse, error) {
	// Verify post exists
	_, err := c.BlogRepository.FindById(c.DB, postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		c.Log.Error().Err(err).Msg("failed to find post")
		return nil, errors.New("internal server error")
	}

	comments, err := c.CommentRepository.FindByPostId(c.DB, postId)
	if err != nil {
		c.Log.Error().Err(err).Msg("failed to find comments")
		return nil, errors.New("internal server error")
	}

	responses := make([]model.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = *c.toCommentResponse(&comment)
	}

	return responses, nil
}
