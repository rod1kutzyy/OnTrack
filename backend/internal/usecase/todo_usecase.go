package usecase

import (
	"context"

	"github.com/rod1kutzyy/OnTrack/internal/domain"
	"github.com/rod1kutzyy/OnTrack/internal/dto"
)

type TodoUseCase interface {
	CreateTodo(ctx context.Context, req dto.CreateTodoRequest) (*domain.Todo, error)
	GetTodoByID(ctx context.Context, id uint) (*domain.Todo, error)
	GetAllTodos(ctx context.Context, filter domain.TodoFilter) ([]domain.Todo, int64, error)
	UpdateTodo(ctx context.Context, id uint, req dto.UpdateTodoRequest) (*domain.Todo, error)
	DeleteTodo(ctx context.Context, id uint) error
	ToggleTodoComplete(ctx context.Context, id uint) (*domain.Todo, error)
}
