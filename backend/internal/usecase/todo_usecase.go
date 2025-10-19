package usecase

import (
	"context"
	"errors"

	"github.com/rod1kutzyy/OnTrack/internal/entity"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	GetAll(ctx context.Context) ([]entity.Todo, error)
	GetByID(ctx context.Context, id uint) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id uint) error
}

type TodoUsecase struct {
	repo TodoRepository
}

func NewTodoUsecase(repo TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

// TODO: Add DTO and validation and Update usecase for DTO

func (uc *TodoUsecase) Create(ctx context.Context, todo *entity.Todo) error {
	if todo.Title == "" {
		return errors.New("title is required")
	}

	return uc.repo.Create(ctx, todo)
}

func (uc *TodoUsecase) GetAll(ctx context.Context) ([]entity.Todo, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *TodoUsecase) GetByID(ctx context.Context, id uint) (*entity.Todo, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *TodoUsecase) Update(ctx context.Context, todo *entity.Todo) error {
	return uc.repo.Update(ctx, todo)
}

func (uc *TodoUsecase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
