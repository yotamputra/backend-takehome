package http

import (
	"app/internal/delivery/middleware"
	"app/internal/model"
	"app/internal/usecase"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type BlogController struct {
	BlogUseCase *usecase.BlogUseCase
	Log         *zerolog.Logger
}

func NewBlogController(blogUseCase *usecase.BlogUseCase, log *zerolog.Logger) *BlogController {
	return &BlogController{
		BlogUseCase: blogUseCase,
		Log:         log,
	}
}

// Create godoc
// @Summary      Create a new blog post
// @Description  Create a new blog post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        request body model.CreateBlogRequest true "Create Blog Request"
// @Security     BearerAuth
// @Success      200 {object} model.WebResponse[model.BlogResponse]
// @Failure      400 {object} model.WebResponse[any]
// @Failure      401 {object} model.WebResponse[any]
// @Router       /api/posts [post]
func (c *BlogController) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil {
		c.jsonResponse(w, http.StatusUnauthorized, model.WebResponse[any]{Errors: "Unauthorized"})
		return
	}

	var request model.CreateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid request body"})
		return
	}

	response, err := c.BlogUseCase.Create(&request, user.ID)
	if err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: err.Error()})
		return
	}

	c.jsonResponse(w, http.StatusOK, model.WebResponse[*model.BlogResponse]{Data: response})
}

// GetById godoc
// @Summary      Get a blog post by ID
// @Description  Get a blog post by ID
// @Tags         posts
// @Produce      json
// @Param        id path string true "Blog ID"
// @Success      200 {object} model.WebResponse[model.BlogResponse]
// @Failure      404 {object} model.WebResponse[any]
// @Router       /api/posts/{id} [get]
func (c *BlogController) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid ID format"})
		return
	}

	response, err := c.BlogUseCase.GetById(idStr)
	if err != nil {
		c.jsonResponse(w, http.StatusNotFound, model.WebResponse[any]{Errors: err.Error()})
		return
	}

	c.jsonResponse(w, http.StatusOK, model.WebResponse[*model.BlogResponse]{Data: response})
}

// GetAll godoc
// @Summary      List all blog posts
// @Description  List all blog posts
// @Tags         posts
// @Produce      json
// @Param        page    query     int  false  "Page number"  default(1)
// @Param        size    query     int  false  "Page size"    default(10)
// @Success      200 {object} model.WebResponse[[]model.BlogResponse]
// @Router       /api/posts [get]
func (c *BlogController) GetAll(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	responses, total, err := c.BlogUseCase.GetAll(page, size)
	if err != nil {
		c.jsonResponse(w, http.StatusInternalServerError, model.WebResponse[any]{Errors: err.Error()})
		return
	}

	totalPage := int(math.Ceil(float64(total) / float64(size)))

	c.jsonResponse(w, http.StatusOK, model.WebResponse[[]model.BlogResponse]{
		Data: responses,
		Paging: &model.PageMetadata{
			Page:      page,
			Size:      size,
			TotalItem: total,
			TotalPage: totalPage,
		},
	})
}

// Update godoc
// @Summary      Update a blog post
// @Description  Update a blog post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id path string true "Blog ID"
// @Param        request body model.UpdateBlogRequest true "Update Blog Request"
// @Security     BearerAuth
// @Success      200 {object} model.WebResponse[model.BlogResponse]
// @Failure      400 {object} model.WebResponse[any]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      403 {object} model.WebResponse[any]
// @Router       /api/posts/{id} [put]
func (c *BlogController) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil {
		c.jsonResponse(w, http.StatusUnauthorized, model.WebResponse[any]{Errors: "Unauthorized"})
		return
	}

	idStr := r.PathValue("id")
	if idStr == "" {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid ID format"})
		return
	}

	var request model.UpdateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid request body"})
		return
	}
	request.ID = idStr

	response, err := c.BlogUseCase.Update(&request, user.ID)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "forbidden") {
			status = http.StatusForbidden
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.jsonResponse(w, status, model.WebResponse[any]{Errors: err.Error()})
		return
	}

	c.jsonResponse(w, http.StatusOK, model.WebResponse[*model.BlogResponse]{Data: response})
}

// Delete godoc
// @Summary      Delete a blog post
// @Description  Delete a blog post
// @Tags         posts
// @Produce      json
// @Param        id path string true "Blog ID"
// @Security     BearerAuth
// @Success      200 {object} model.WebResponse[any]
// @Failure      400 {object} model.WebResponse[any]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      403 {object} model.WebResponse[any]
// @Router       /api/posts/{id} [delete]
func (c *BlogController) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil {
		c.jsonResponse(w, http.StatusUnauthorized, model.WebResponse[any]{Errors: "Unauthorized"})
		return
	}

	idStr := r.PathValue("id")
	if idStr == "" {
		c.jsonResponse(w, http.StatusBadRequest, model.WebResponse[any]{Errors: "invalid ID format"})
		return
	}

	if err := c.BlogUseCase.Delete(idStr, user.ID); err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "forbidden") {
			status = http.StatusForbidden
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.jsonResponse(w, status, model.WebResponse[any]{Errors: err.Error()})
		return
	}

	msg := "Successfully delete the post"
	c.jsonResponse(w, http.StatusOK, model.WebResponse[any]{Message: &msg})
}

func (c *BlogController) jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
