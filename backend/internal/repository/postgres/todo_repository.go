package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/rod1kutzyy/OnTrack/internal/domain"
	"github.com/rod1kutzyy/OnTrack/internal/repository"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &todoRepository{
		db: db,
	}
}

func (r *todoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	if err := r.db.WithContext(ctx).Create(todo).Error; err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

func (r *todoRepository) GetByID(ctx context.Context, id uint) (*domain.Todo, error) {
	var todo domain.Todo

	if err := r.db.WithContext(ctx).First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("todo with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

func (r *todoRepository) GetAll(ctx context.Context, filter domain.TodoFilter) ([]domain.Todo, error) {
	var todos []domain.Todo

	query := r.db.WithContext(ctx).Model(&domain.Todo{})

	if filter.Completed != nil {
		query = query.Where("completed = ?", *filter.Completed)
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where(
			"title ILIKE ? OR description ILIKE ?",
			searchPattern, searchPattern,
		)
	}

	query = query.Order("created_at DESC")

	query = query.Limit(filter.Limit).Offset(filter.Offset)

	if err := query.Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	result := r.db.WithContext(ctx).Save(todo)

	if result.Error != nil {
		return fmt.Errorf("failed to update todo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", todo.ID)
	}

	return nil
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Todo{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete todo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return nil
}

func (r *todoRepository) Count(ctx context.Context, filter domain.TodoFilter) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&domain.Todo{})

	if filter.Completed != nil {
		query = query.Where("completed = ?", *filter.Completed)
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where(
			"title ILIKE ? OR description ILIKE ?",
			searchPattern, searchPattern,
		)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count todos: %w", err)
	}

	return count, nil
}
