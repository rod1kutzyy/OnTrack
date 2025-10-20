package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rod1kutzyy/OnTrack/internal/entity"
	"github.com/rod1kutzyy/OnTrack/internal/usecase"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type TodoHandler struct {
	usecase *usecase.TodoUsecase
}

func NewTodoHandler(uc *usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{usecase: uc}
}

func mapErrorToResponse(err error) (int, ErrorResponse) {
	switch {
	case errors.Is(err, entity.ErrTodoNotFound):
		return http.StatusNotFound, ErrorResponse{Message: entity.ErrTodoNotFound.Error()}
	case errors.Is(err, entity.ErrTodoTitleRequired):
		return http.StatusBadRequest, ErrorResponse{Message: entity.ErrTodoTitleRequired.Error()}
	case errors.Is(err, entity.ErrInvalidTodo):
		return http.StatusBadRequest, ErrorResponse{Message: entity.ErrInvalidTodo.Error()}
	case errors.Is(err, entity.ErrTodoDatabase):
		return http.StatusInternalServerError, ErrorResponse{Message: entity.ErrTodoDatabase.Error()}
	default:
		return http.StatusInternalServerError, ErrorResponse{Message: "internal server error"}
	}
}

func (h *TodoHandler) Create(c *gin.Context) {
	var todo entity.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.usecase.Create(ctx, &todo); err != nil {
		code, response := mapErrorToResponse(err)
		c.JSON(code, response)
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	todos, err := h.usecase.GetAll(ctx)
	if err != nil {
		code, response := mapErrorToResponse(err)
		c.JSON(code, response)
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid todo ID",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	todo, err := h.usecase.GetByID(ctx, uint(id))
	if err != nil {
		code, response := mapErrorToResponse(err)
		c.JSON(code, response)
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid todo ID",
		})
		return
	}

	var todo entity.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	todo.ID = uint(id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.usecase.Update(ctx, &todo); err != nil {
		code, response := mapErrorToResponse(err)
		c.JSON(code, response)
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid todo ID",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.usecase.Delete(ctx, uint(id)); err != nil {
		code, response := mapErrorToResponse(err)
		c.JSON(code, response)
		return
	}

	c.Status(http.StatusNoContent)
}
