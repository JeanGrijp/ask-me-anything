package responses

import (
	"strconv"

	"github.com/JeanGrijp/ask-me-anything/internal/validators"
	"github.com/go-playground/validator/v10"
)

func NewSuccessResponse(data any) Response {
	return Response{
		Status: "success",
		Data:   data,
	}
}

func NewSuccessMessage(message string, data any) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, err ...error) Response {
	resp := Response{
		Status:  "error",
		Message: message,
	}

	if len(err) > 0 && err[0] != nil {
		resp.Message += ": " + err[0].Error()
	}

	return resp
}

func NewValidationError(message string, fields []FieldError) Response {
	return Response{
		Status:  "error",
		Message: message,
		Fields:  fields,
	}
}

func ConvertValidationErrors(err error) []FieldError {
	if err == nil {
		return nil
	}

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	var errors []FieldError
	for _, e := range ve {
		translated := e.Translate(validators.Translator)
		errors = append(errors, FieldError{
			Field:   e.Field(),
			Message: translated,
		})
	}

	return errors
}

func ParseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}

func ParseIntPointer(value string) *int {
	if value == "" {
		return nil
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &i
}

func ParseBoolPointer(value string) *bool {
	if value == "" {
		return nil
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}
	return &b
}

func ToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
