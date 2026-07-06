package http

import (
	"app/internal/delivery/middleware"
	"app/internal/model"
	"app/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CommentController struct {
	CommentUseCase *usecase.CommentUseCase
}

func NewCommentController(commentUseCase *usecase.CommentUseCase) *CommentController {
	return &CommentController{
		CommentUseCase: commentUseCase,
	}
}

func (c *CommentController) jsonResponse(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Create godoc
// @Summary      Add a comment to a blog post
// @Description  Add a comment to a blog post
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id path int true "Blog ID"
// @Param        request body model.CreateCommentRequest true "Create Comment Request"
// @Security     BearerAuth
// @Success      200 {object} model.WebResponse[model.CommentResponse]
// @Failure      400 {object} model.WebResponse[any]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      404 {object} model.WebResponse[any]
// @Router       /api/posts/{id}/comments [post]
func (c *CommentController) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil {
		c.jsonResponse(w, http.StatusUnauthorized, model.WebResponse[any]{Errors: "Unauthorized"})
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid ID format"})
		return
	}

	var request model.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid request body"})
		return
	}

	response, err := c.CommentUseCase.Create(&request, id, user.Name)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.jsonResponse(w, status, model.WebResponse[any]{Errors: err.Error()})
		return
	}

	c.jsonResponse(w, http.StatusOK, model.WebResponse[*model.CommentResponse]{Data: response})
}

// GetByPostId godoc
// @Summary      List all comments for a blog post
// @Description  List all comments for a blog post
// @Tags         comments
// @Produce      json
// @Param        id path int true "Blog ID"
// @Success      200 {object} model.WebResponse[[]model.CommentResponse]
// @Failure      400 {object} model.WebResponse[any]
// @Failure      404 {object} model.WebResponse[any]
// @Router       /api/posts/{id}/comments [get]
func (c *CommentController) GetByPostId(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid ID format"})
		return
	}

	responses, err := c.CommentUseCase.GetByPostId(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.jsonResponse(w, status, model.WebResponse[any]{Errors: err.Error()})
		return
	}

	c.jsonResponse(w, http.StatusOK, model.WebResponse[[]model.CommentResponse]{Data: responses})
}
