package responses

import (
	"net/http"
)

// SendError sends an error response
func SendError(w http.ResponseWriter, status int, message string) {
	response := NewErrorResponse(message)
	JSON(w, status, response)
}

// SendValidationError sends a validation error response
func SendValidationError(w http.ResponseWriter, err error) {
	fields := ConvertValidationErrors(err)
	response := NewValidationError("Validation failed", fields)
	JSON(w, http.StatusBadRequest, response)
}

// SendSuccess sends a success response
func SendSuccess(w http.ResponseWriter, data interface{}) {
	response := NewSuccessResponse(data)
	JSON(w, http.StatusOK, response)
}

// SendCreated sends a created response
func SendCreated(w http.ResponseWriter, data interface{}) {
	response := NewSuccessResponse(data)
	JSON(w, http.StatusCreated, response)
}
