package model

type WebResponse[T any] struct {
	Data    T       `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	Errors  any     `json:"errors,omitempty"`
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
