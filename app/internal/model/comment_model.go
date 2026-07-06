package model

import "time"

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type CommentResponse struct {
	ID         string `json:"id"`
	PostID     string `json:"post_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
