package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (cv *CustomValidator) Validate(i interface{}) []ErrorResponse {
	if err := cv.validator.Struct(i); err != nil {
		var errors []ErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("failed on the '%s' validation", err.Tag())
			if err.Param() != "" {
				msg = fmt.Sprintf("failed on the '%s' validation with parameter '%s'", err.Tag(), err.Param())
			}
			errors = append(errors, ErrorResponse{
				Field:   strings.ToLower(err.Field()),
				Message: msg,
			})
		}
		return errors
	}
	return nil
}
