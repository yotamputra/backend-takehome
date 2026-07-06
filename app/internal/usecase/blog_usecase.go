package usecase

import (
	"app/internal/entity"
	"app/internal/model"
	"app/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BlogUseCase struct {
	DB             *gorm.DB
	Log            *zerolog.Logger
	Validate       *validator.Validate
	BlogRepository *repository.BlogRepository
}

func NewBlogUseCase(db *gorm.DB, log *zerolog.Logger, validate *validator.Validate, blogRepository *repository.BlogRepository) *BlogUseCase {
	return &BlogUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		BlogRepository: blogRepository,
	}
}

func (c *BlogUseCase) Create(request *model.CreateBlogRequest, authorID string) (*model.BlogResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		return nil, err
	}

	authorUuid, err := uuid.Parse(authorID)
	if err != nil {
		return nil, errors.New("invalid author id")
	}

	blog := &entity.Blog{
		Title:    request.Title,
		Content:  request.Content,
		AuthorID: authorUuid,
	}

	if err := c.BlogRepository.Create(c.DB, blog); err != nil {
		c.Log.Error().Err(err).Msg("failed to create blog post")
		return nil, err
	}

	// Fetch again to get the preloaded author
	createdBlog, err := c.BlogRepository.FindById(c.DB, blog.ID.String())
	if err != nil {
		return nil, err
	}

	return c.toBlogResponse(createdBlog), nil
}

func (c *BlogUseCase) GetById(id string) (*model.BlogResponse, error) {
	blog, err := c.BlogRepository.FindById(c.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blog post not found")
		}
		return nil, err
	}

	return c.toBlogResponse(blog), nil
}

func (c *BlogUseCase) GetAll(page, size int) ([]model.BlogResponse, int64, error) {
	blogs, total, err := c.BlogRepository.FindAll(c.DB, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]model.BlogResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = *c.toBlogResponse(&blog)
	}

	return responses, total, nil
}

func (c *BlogUseCase) Update(request *model.UpdateBlogRequest, currentUserID string) (*model.BlogResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		return nil, err
	}

	blog, err := c.BlogRepository.FindById(c.DB, request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blog post not found")
		}
		return nil, err
	}

	// Authorization check
	if blog.AuthorID.String() != currentUserID {
		return nil, errors.New("forbidden: you do not have permission to modify this post")
	}

	blog.Title = request.Title
	blog.Content = request.Content

	if err := c.BlogRepository.Update(c.DB, blog); err != nil {
		c.Log.Error().Err(err).Msg("failed to update blog post")
		return nil, err
	}

	return c.toBlogResponse(blog), nil
}

func (c *BlogUseCase) Delete(id string, currentUserID string) error {
	blog, err := c.BlogRepository.FindById(c.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("blog post not found")
		}
		return err
	}

	// Authorization check
	if blog.AuthorID.String() != currentUserID {
		return errors.New("forbidden: you do not have permission to delete this post")
	}

	if err := c.BlogRepository.Delete(c.DB, blog); err != nil {
		c.Log.Error().Err(err).Msg("failed to delete blog post")
		return err
	}

	return nil
}

func (c *BlogUseCase) toBlogResponse(blog *entity.Blog) *model.BlogResponse {
	comments := make([]model.CommentResponse, len(blog.Comments))
	for i, comment := range blog.Comments {
		comments[i] = model.CommentResponse{
			ID:         comment.ID.String(),
			PostID:     comment.PostID.String(),
			AuthorName: comment.AuthorName,
			Content:    comment.Content,
			CreatedAt:  comment.CreatedAt,
		}
	}

	return &model.BlogResponse{
		ID:      blog.ID.String(),
		Title:   blog.Title,
		Content: blog.Content,
		Author: model.UserResponse{
			ID:        blog.Author.ID.String(),
			Name:      blog.Author.Name,
			Email:     blog.Author.Email,
			CreatedAt: blog.Author.CreatedAt,
			UpdatedAt: blog.Author.UpdatedAt,
		},
		CreatedAt: blog.CreatedAt,
		UpdatedAt: blog.UpdatedAt,
		Comments:  comments,
	}
}
