package model

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type CommentResponse struct {
	ID         int    `json:"id"`
	PostID     int    `json:"post_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
}
