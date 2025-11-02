package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/rod1kutzyy/OnTrack/internal/domain"
	"github.com/rod1kutzyy/OnTrack/internal/dto"
	"github.com/rod1kutzyy/OnTrack/internal/logger"
	"github.com/rod1kutzyy/OnTrack/internal/repository"
)

type todoUseCase struct {
	todoRepo repository.TodoRepository
}

func NewTodoUseCase(todoRepo repository.TodoRepository) TodoUseCase {
	return &todoUseCase{
		todoRepo: todoRepo,
	}
}

func (uc *todoUseCase) CreateTodo(ctx context.Context, req dto.CreateTodoRequest) (*domain.Todo, error) {
	logger.Logger.WithField("title", req.Title).Info("Creating new todo")

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty or contain only spaces")
	}

	var description *string
	if req.Description != nil {
		trimmedDescription := strings.TrimSpace(*req.Description)
		if trimmedDescription != "" {
			description = &trimmedDescription
		}
	}

	todo := &domain.Todo{
		Title:       title,
		Description: description,
		Completed:   false,
	}

	if err := uc.todoRepo.Create(ctx, todo); err != nil {
		logger.Logger.WithError(err).Error("Failed to create todo")
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	logger.Logger.WithField("id", todo.ID).Info("Todo created successfully")
	return todo, nil
}

func (uc *todoUseCase) GetTodoByID(ctx context.Context, id uint) (*domain.Todo, error) {
	logger.Logger.WithField("id", id).Debug("Fetching todo by ID")

	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		logger.Logger.WithError(err).WithField("id", id).Warn("Todo not found")
		return nil, err
	}

	return todo, nil
}

func (uc *todoUseCase) GetAllTodos(ctx context.Context, filter domain.TodoFilter) ([]domain.Todo, int64, error) {
	logger.Logger.WithField("filter", filter).Debug("Fetching todos with filter")

	filter.Validate()

	todos, err := uc.todoRepo.GetAll(ctx, filter)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to fetch todos")
		return nil, 0, fmt.Errorf("failed to fetch todos: %w", err)
	}

	count, err := uc.todoRepo.Count(ctx, filter)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to count todos")
		return nil, 0, fmt.Errorf("failed to count todos: %w", err)
	}

	logger.Logger.WithField("count", len(todos)).Debug("Todos fetched successfully")
	return todos, count, nil
}

func (uc *todoUseCase) UpdateTodo(ctx context.Context, id uint, req dto.UpdateTodoRequest) (*domain.Todo, error) {
	logger.Logger.WithField("id", id).Info("Updating todo")

	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, fmt.Errorf("title cannot be empty")
		}
		todo.Title = title
	}

	if req.Description != nil {
		trimmedDescription := strings.TrimSpace(*req.Description)
		if trimmedDescription == "" {
			todo.Description = nil
		} else {
			todo.Description = &trimmedDescription
		}
	}

	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	if err := uc.todoRepo.Update(ctx, todo); err != nil {
		logger.Logger.WithError(err).Error("Failed to update todo")
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	logger.Logger.WithField("id", id).Info("Todo updated successfully")
	return todo, nil
}

func (uc *todoUseCase) DeleteTodo(ctx context.Context, id uint) error {
	logger.Logger.WithField("id", id).Info("Deleting todo")

	if err := uc.todoRepo.Delete(ctx, id); err != nil {
		logger.Logger.WithError(err).Error("Failed to delete todo")
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	logger.Logger.WithField("id", id).Info("Todo deleted successfully")
	return nil
}

func (uc *todoUseCase) ToggleTodoComplete(ctx context.Context, id uint) (*domain.Todo, error) {
	logger.Logger.WithField("id", id).Info("Toggling todo completion status")

	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	todo.Completed = !todo.Completed

	if err := uc.todoRepo.Update(ctx, todo); err != nil {
		logger.Logger.WithError(err).Error("Failed to toggle todo completion")
		return nil, fmt.Errorf("failed to toggle todo completion: %w", err)
	}

	logger.Logger.WithField("id", id).WithField("completed", todo.Completed).Info("Todo completion status toggled")
	return todo, nil
}
