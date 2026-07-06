package model

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int   `json:"total_page"`
}

type WebResponse[T any] struct {
	Data    T             `json:"data"`
	Paging  *PageMetadata `json:"paging,omitempty"`
	Message *string       `json:"message,omitempty"`
	Errors  any           `json:"errors,omitempty"`
}

func Ok[T any](data T) WebResponse[T] {
	return WebResponse[T]{
		Data: data,
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
