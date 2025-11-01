package dto

import "time"

type TodoResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TodoListResponse struct {
	Items      []TodoResponse     `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

func NewPaginationResponse(total int64, page, limit int) PaginationResponse {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse{
		Total:       total,
		TotalPages:  totalPages,
		CurrentPage: page,
		PerPage:     limit,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Tag     string      `json:"tag,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

func NewSuccessResponse(data interface{}, message string) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

func NewErrorResponse(err, message string) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error:   err,
		Message: message,
	}
}

func NewErrorResponseWithCode(err, message, code string) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error:   err,
		Message: message,
		Code:    code,
	}
}

func NewValidationErrorResponse(errors []ValidationError) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error:   "Validation Error",
		Message: "Request validation failed",
		Code:    "VALIDATION_ERROR",
		Details: errors,
	}
}
