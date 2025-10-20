package entity

import (
	"errors"
	"time"
)

var (
	ErrTodoNotFound      = errors.New("todo not found")
	ErrInvalidTodo       = errors.New("invalid todo data")
	ErrTodoTitleRequired = errors.New("todo's title is required")
	ErrTodoDatabase      = errors.New("todo database error")
)

type Todo struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Done        bool      `gorm:"default:false" json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
