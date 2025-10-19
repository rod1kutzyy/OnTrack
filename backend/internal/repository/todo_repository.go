package repository

import (
	"context"

	"github.com/rod1kutzyy/OnTrack/internal/entity"
	"github.com/rod1kutzyy/OnTrack/internal/usecase"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) usecase.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) GetAll(ctx context.Context) ([]entity.Todo, error) {
	var todos []entity.Todo
	err := r.db.WithContext(ctx).Find(&todos).Error
	return todos, err
}

func (r *todoRepository) GetByID(ctx context.Context, id uint) (*entity.Todo, error) {
	var todo entity.Todo
	err := r.db.WithContext(ctx).First(&todo).Error
	return &todo, err
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Todo{}, id).Error
}
