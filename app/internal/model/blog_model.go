package model

import "time"

type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required"`
}

type UpdateBlogRequest struct {
	ID      string `json:"-" validate:"required"`
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required"`
}

type BlogResponse struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Author    UserResponse      `json:"author"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	
	Comments  []CommentResponse `json:"comments"`
}
