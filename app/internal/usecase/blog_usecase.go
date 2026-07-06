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

func (c *BlogUseCase) Create(request *model.CreateBlogRequest, authorID int) (*model.BlogResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		return nil, err
	}

	blog := &entity.Blog{
		Title:    request.Title,
		Content:  request.Content,
		AuthorID: authorID,
	}

	if err := c.BlogRepository.Create(c.DB, blog); err != nil {
		c.Log.Error().Err(err).Msg("failed to create blog post")
		return nil, err
	}

	// Fetch again to get the preloaded author
	createdBlog, err := c.BlogRepository.FindById(c.DB, blog.ID)
	if err != nil {
		return nil, err
	}

	return c.toBlogResponse(createdBlog), nil
}

func (c *BlogUseCase) GetById(id int) (*model.BlogResponse, error) {
	blog, err := c.BlogRepository.FindById(c.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blog post not found")
		}
		return nil, err
	}

	return c.toBlogResponse(blog), nil
}

func (c *BlogUseCase) GetAll() ([]model.BlogResponse, error) {
	blogs, err := c.BlogRepository.FindAll(c.DB)
	if err != nil {
		return nil, err
	}

	responses := make([]model.BlogResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = *c.toBlogResponse(&blog)
	}

	return responses, nil
}

func (c *BlogUseCase) Update(request *model.UpdateBlogRequest, currentUserID int) (*model.BlogResponse, error) {
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
	if blog.AuthorID != currentUserID {
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

func (c *BlogUseCase) Delete(id int, currentUserID int) error {
	blog, err := c.BlogRepository.FindById(c.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("blog post not found")
		}
		return err
	}

	// Authorization check
	if blog.AuthorID != currentUserID {
		return errors.New("forbidden: you do not have permission to delete this post")
	}

	if err := c.BlogRepository.Delete(c.DB, blog); err != nil {
		c.Log.Error().Err(err).Msg("failed to delete blog post")
		return err
	}

	return nil
}

func (c *BlogUseCase) toBlogResponse(blog *entity.Blog) *model.BlogResponse {
	return &model.BlogResponse{
		ID:      blog.ID,
		Title:   blog.Title,
		Content: blog.Content,
		Author: model.UserResponse{
			ID:        blog.Author.ID,
			Name:      blog.Author.Name,
			Email:     blog.Author.Email,
			CreatedAt: blog.Author.CreatedAt,
			UpdatedAt: blog.Author.UpdatedAt,
		},
		CreatedAt: blog.CreatedAt,
		UpdatedAt: blog.UpdatedAt,
	}
}
