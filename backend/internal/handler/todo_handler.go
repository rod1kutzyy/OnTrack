package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rod1kutzyy/OnTrack/internal/domain"
	"github.com/rod1kutzyy/OnTrack/internal/dto"
	"github.com/rod1kutzyy/OnTrack/internal/logger"
	"github.com/rod1kutzyy/OnTrack/internal/usecase"
	"github.com/rod1kutzyy/OnTrack/internal/validator"
)

type TodoHandler struct {
	todoUseCase usecase.TodoUseCase
	validator   *validator.TodoValidator
}

func NewTodoHandler(todoUseCase usecase.TodoUseCase, validator *validator.TodoValidator) *TodoHandler {
	return &TodoHandler{
		todoUseCase: todoUseCase,
		validator:   validator,
	}
}

// @Summary Create a new Todo
// @Description Creates a new todo item and stores it in the database
// @Tags todos
// @Accept json
// @Produce json
// @Param input body dto.CreateTodoRequest true "Todo creation data"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req dto.CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.WithError(err).Warn("Failed to parse request body")
		response := dto.NewErrorResponseWithCode(
			"Bad Request",
			"Invalid request data format",
			"INVALID_JSON",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if validationErrors := h.validator.ValidateCreateTodo(req); len(validationErrors) > 0 {
		logger.Logger.WithField("errors", validationErrors).Warn("Validation failed")
		response := dto.NewValidationErrorResponse(validationErrors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	todo, err := h.todoUseCase.CreateTodo(c.Request.Context(), req)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to create todo")
		response := dto.NewErrorResponse("Internal Server Error", "Failed to create todo")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	responseDTO := h.mapTodoToDTO(todo)
	response := dto.NewSuccessResponse(responseDTO, "Todo created successfully")
	c.JSON(http.StatusCreated, response)
}

// @Summary Get Todo by ID
// @Description Returns a single todo item by its ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		response := dto.NewErrorResponseWithCode(
			"Bad Request",
			"Invalid todo ID format",
			"INVALID_ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	todo, err := h.todoUseCase.GetTodoByID(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response := dto.NewErrorResponseWithCode(
				"Not Found",
				fmt.Sprintf("Todo with ID %d not found", id),
				"TODO_NOT_FOUND",
			)
			c.JSON(http.StatusNotFound, response)
		}

		logger.Logger.WithError(err).Error("Failed to get todo")
		response := dto.NewErrorResponse("Internal Server Error", "Failed to retrieve todo")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	responseDTO := h.mapTodoToDTO(todo)
	response := dto.NewSuccessResponse(responseDTO, "")
	c.JSON(http.StatusOK, response)
}

// @Summary Get all Todos
// @Description Returns a paginated list of todos, optionally filtered by completion status or search term
// @Tags todos
// @Produce json
// @Param completed query bool false "Filter by completion status"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos [get]
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	var filter dto.TodoFilterRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.Logger.WithError(err).Warn("Invalid query parameters")
		response := dto.NewErrorResponse("Bad Request", "Invalid query parameters")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if validationErrors := h.validator.ValidateFilter(filter); len(validationErrors) > 0 {
		logger.Logger.WithField("errors", validationErrors).Warn("Filter validation failed")
		response := dto.NewValidationErrorResponse(validationErrors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	filter.Validate()

	domainFilter := domain.TodoFilter{
		Completed: filter.Completed,
		Search:    filter.Search,
		Limit:     filter.Limit,
		Offset:    filter.GetOffset(),
	}

	todos, total, err := h.todoUseCase.GetAllTodos(c.Request.Context(), domainFilter)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to get todos")
		response := dto.NewErrorResponse("Internal Server Error", "Failed to retrieve todos")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	todoDTOs := make([]dto.TodoResponse, len(todos))
	for i, todo := range todos {
		todoDTOs[i] = h.mapTodoToDTO(&todo)
	}

	listResponse := dto.TodoListResponse{
		Items:      todoDTOs,
		Pagination: dto.NewPaginationResponse(total, filter.Page, filter.Limit),
	}

	response := dto.NewSuccessResponse(listResponse, "")
	c.JSON(http.StatusOK, response)
}

// @Summary Update Todo
// @Description Updates an existing todo item by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param input body dto.UpdateTodoRequest true "Updated todo data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		response := dto.NewErrorResponseWithCode(
			"Bad Request",
			"Invalid todo ID",
			"INVALID_ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.WithError(err).Warn("Failed to parse request body")
		response := dto.NewErrorResponseWithCode(
			"Bad Request",
			"Invalid request data",
			"INVALID_JSON",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if validationErrors := h.validator.ValidateUpdateTodo(req); len(validationErrors) > 0 {
		logger.Logger.WithField("errors", validationErrors).Warn("Update validation failed")
		response := dto.NewValidationErrorResponse(validationErrors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	todo, err := h.todoUseCase.UpdateTodo(c.Request.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response := dto.NewErrorResponseWithCode(
				"Not Found",
				fmt.Sprintf("Todo with ID %d not found", id),
				"TODO_NOT_FOUND",
			)
			c.JSON(http.StatusNotFound, response)
			return
		}

		logger.Logger.WithError(err).Error("Failed to update todo")
		response := dto.NewErrorResponse("Internal Server Error", "Failed to update todo")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	responseDTO := h.mapTodoToDTO(todo)
	response := dto.NewSuccessResponse(responseDTO, "Todo updated successfully")
	c.JSON(http.StatusOK, response)
}

// @Summary Delete Todo
// @Description Deletes a todo item by its ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 204 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		response := dto.NewErrorResponseWithCode(
			"Bad Request",
			"Invalid todo ID",
			"INVALID_ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.todoUseCase.DeleteTodo(c.Request.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			response := dto.NewErrorResponseWithCode(
				"Not Found",
				fmt.Sprintf("Todo with ID %d not found", id),
				"TODO_NOT_FOUND",
			)
			c.JSON(http.StatusNotFound, response)
			return
		}

		logger.Logger.WithError(err).Error("Failed to delete todo")
		response := dto.NewErrorResponse("Internal Server Error", "Failed to delete todo")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(nil, "Todo deleted successfully")
	c.JSON(http.StatusNoContent, response)
}

// @Summary Toggle Todo completion
// @Description Toggles the completion status (done/undone) of a todo item by its ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos/{id}/toggle [patch]
func (h *TodoHandler) ToggleTodoComplete(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		response := dto.NewErrorResponseWithCode(
			"Bad Request",
			"Invalid todo ID",
			"INVALID_ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	todo, err := h.todoUseCase.ToggleTodoComplete(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response := dto.NewErrorResponseWithCode(
				"Not Found",
				fmt.Sprintf("Todo with ID %d not found", id),
				"TODO_NOT_FOUND",
			)
			c.JSON(http.StatusNotFound, response)
			return
		}

		logger.Logger.WithError(err).Error("Failed to toggle todo completion")
		response := dto.NewErrorResponse("Internal Server Error", "Failed to toggle completion status")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	responseDTO := h.mapTodoToDTO(todo)
	response := dto.NewSuccessResponse(responseDTO, "Completion status toggled successfully")
	c.JSON(http.StatusOK, response)
}

func (h *TodoHandler) mapTodoToDTO(todo *domain.Todo) dto.TodoResponse {
	description := ""
	if todo.Description != nil {
		description = *todo.Description
	}

	return dto.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

func (h *TodoHandler) parseIDParam(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}

	return uint(id), nil
}
