package repository

import (
	"context"

	"github.com/rod1kutzyy/OnTrack/internal/domain"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *domain.Todo) error
	GetByID(ctx context.Context, id uint) (*domain.Todo, error)
	GetAll(ctx context.Context, filter domain.TodoFilter) ([]domain.Todo, error)
	Update(ctx context.Context, todo *domain.Todo) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context, filter domain.TodoFilter)
}
