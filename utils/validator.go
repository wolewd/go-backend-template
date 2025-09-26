package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) []string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	return formatValidationErrors(validationErrors)
}

func formatValidationErrors(validationErrors validator.ValidationErrors) []string {
	var errs []string

	for _, fieldErr := range validationErrors {
		field := strings.ToLower(fieldErr.Field())
		tag := fieldErr.Tag()
		param := fieldErr.Param()

		errs = append(errs, buildErrorMessage(field, tag, param))
	}

	return errs
}

func buildErrorMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, param)
	default:
		return fmt.Sprintf("%s is not valid", field)
	}
}
