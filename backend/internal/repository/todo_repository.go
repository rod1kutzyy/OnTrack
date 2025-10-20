package repository

import (
	"context"
	"errors"

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
	if err := r.db.WithContext(ctx).Create(todo).Error; err != nil {
		return entity.ErrTodoDatabase
	}
	return nil
}

func (r *todoRepository) GetAll(ctx context.Context) ([]entity.Todo, error) {
	var todos []entity.Todo
	err := r.db.WithContext(ctx).Find(&todos).Error
	if err != nil {
		return nil, entity.ErrTodoDatabase
	}
	return todos, nil
}

func (r *todoRepository) GetByID(ctx context.Context, id uint) (*entity.Todo, error) {
	var todo entity.Todo
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrTodoNotFound
		}
		return nil, entity.ErrTodoDatabase
	}
	return &todo, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	err := r.db.WithContext(ctx).Save(todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ErrTodoNotFound
		}
		return entity.ErrTodoDatabase
	}
	return nil
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Todo{}, id)
	if result.Error != nil {
		return entity.ErrTodoDatabase
	}
	if result.RowsAffected == 0 {
		return entity.ErrTodoNotFound
	}
	return nil
}
