package domain

import "time"

type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description *string   `json:"description"`
	Completed   bool      `json:"completed" gorm:"default:false;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Todo) TableName() string {
	return "todos"
}

type TodoFilter struct {
	Completed *bool
	Search    string
	Limit     int
	Offset    int
}

func (f *TodoFilter) Validate() {
	if f.Limit <= 0 {
		f.Limit = 5
	}

	if f.Limit > 100 {
		f.Limit = 100
	}

	if f.Offset < 0 {
		f.Offset = 0
	}
}
