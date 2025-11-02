package validator

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/rod1kutzyy/OnTrack/internal/dto"
)

type TodoValidator struct {
	validate *validator.Validate
}

func NewTodoValidator() *TodoValidator {
	v := validator.New()

	v.RegisterValidation("notblank", motBlankValidator)

	return &TodoValidator{
		validate: v,
	}
}

func (tv *TodoValidator) ValidateCreateTodo(req dto.CreateTodoRequest) []dto.ValidationError {
	err := tv.validate.Struct(req)
	if err == nil {
		description := ""

		if req.Description != nil {
			description = *req.Description
		}

		return tv.validateTodoBusinessRules(req.Title, description)
	}

	return tv.transformValidationErrors(err)
}

func (tv *TodoValidator) ValidateUpdateTodo(req dto.UpdateTodoRequest) []dto.ValidationError {
	var errors []dto.ValidationError

	if req.Title != nil {
		trimmedTitle := strings.TrimSpace(*req.Title)
		if trimmedTitle == "" {
			errors = append(errors, dto.ValidationError{
				Field:   "title",
				Message: "Title cannot be empty or contain only spaces",
				Tag:     "notblank",
			})
		}

		if utf8.RuneCountInString(trimmedTitle) > 255 {
			errors = append(errors, dto.ValidationError{
				Field:   "title",
				Message: "Title must not exceed 255 characters",
				Tag:     "max",
			})
		}

		businessErrors := tv.validateTodoBusinessRules(trimmedTitle, "")
		errors = append(errors, businessErrors...)
	}

	if req.Description != nil {
		if utf8.RuneCountInString(*req.Description) > 1000 {
			errors = append(errors, dto.ValidationError{
				Field:   "description",
				Message: "Description must not exceed 1000 characters",
				Tag:     "max",
			})
		}
	}

	return errors
}

func (tv *TodoValidator) ValidateFilter(filter dto.TodoFilterRequest) []dto.ValidationError {
	var errors []dto.ValidationError

	if len(filter.Search) > 100 {
		errors = append(errors, dto.ValidationError{
			Field:   "search",
			Message: "Search query is too long (max 100 characters)",
			Tag:     "max",
		})
	}

	if filter.Page < 1 {
		errors = append(errors, dto.ValidationError{
			Field:   "page",
			Message: "Page number must be positive",
			Tag:     "min",
		})
	}

	if filter.Limit <= 0 || filter.Limit > 100 {
		errors = append(errors, dto.ValidationError{
			Field:   "limit",
			Message: "Limit must be between 1 and 100",
			Tag:     "range",
		})
	}

	return errors
}

func (tv *TodoValidator) validateTodoBusinessRules(title, description string) []dto.ValidationError {
	var errors []dto.ValidationError

	if strings.TrimSpace(title) == "" {
		errors = append(errors, dto.ValidationError{
			Field:   "title",
			Message: "Title cannot be empty or contain only whitespace",
			Tag:     "notblank",
		})
	}

	if isOnlyDigits(title) {
		errors = append(errors, dto.ValidationError{
			Field:   "title",
			Message: "Title cannot contain only digits",
			Tag:     "content_quality",
		})
	}

	if hasExcessiveSpecialChars(title) {
		errors = append(errors, dto.ValidationError{
			Field:   "title",
			Message: "Title contains too many special characters",
			Tag:     "content_quality",
		})
	}

	if description != "" && strings.TrimSpace(title) == strings.TrimSpace(description) {
		errors = append(errors, dto.ValidationError{
			Field:   "description",
			Message: "Description should not be identical to title",
			Tag:     "unique_content",
		})
	}

	return errors
}

func (tv *TodoValidator) transformValidationErrors(err error) []dto.ValidationError {
	var errors []dto.ValidationError

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		errors = append(errors, dto.ValidationError{
			Field:   "unknown",
			Message: err.Error(),
			Tag:     "validation",
		})
		return errors
	}

	for _, fieldError := range validationErrors {
		errors = append(errors, dto.ValidationError{
			Field:   strings.ToLower(fieldError.Field()),
			Message: getErrorMessage(fieldError),
			Tag:     fieldError.Tag(),
			Value:   fieldError.Value(),
		})
	}

	return errors
}

func getErrorMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "notblank":
		return fmt.Sprintf("%s cannot be empty or contain only whitespace", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func motBlankValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return strings.TrimSpace(value) != ""
}

func isOnlyDigits(s string) bool {
	trimmed := strings.TrimSpace(s)

	if trimmed == "" {
		return false
	}

	for _, r := range trimmed {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func hasExcessiveSpecialChars(s string) bool {
	if len(s) == 0 {
		return false
	}

	specialCount := 0
	totalCount := utf8.RuneCountInString(s)

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == ' ' || r >= 0x0400) {
			specialCount++
		}
	}

	return float64(specialCount)/float64(totalCount) > 0.3
}
