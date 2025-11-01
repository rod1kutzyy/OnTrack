package dto

type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	Completed   *bool   `json:"completed" binding:"omitempty"`
}

type TodoFilterRequest struct {
	Completed *bool  `form:"completed"`
	Search    string `form:"search" binding:"omitempty,max=100"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (f *TodoFilterRequest) GetOffset() int {
	if f.Page <= 1 {
		return 0
	}

	return (f.Page - 1) * f.Limit
}

func (f *TodoFilterRequest) Validate() {
	if f.Page < 1 {
		f.Page = 1
	}

	if f.Limit <= 0 {
		f.Limit = 10
	}

	if f.Limit > 100 {
		f.Limit = 100
	}
}

func (f *TodoFilterRequest) HasSearchField() bool {
	return f.Search != ""
}

func (f *TodoFilterRequest) HasCompletedField() bool {
	return f.Completed != nil
}
