package responses

type Status string

const (
	Success Status = "success"
	Error   Status = "error"
)

type Response struct {
	Status  Status       `json:"status" example:"success"`
	Message string       `json:"message,omitempty" example:"Operação realizada com sucesso"`
	Fields  []FieldError `json:"fields,omitempty" example:"[{\"field\":\"email\",\"message\":\"Email inválido\"}]"`
	Data    interface{}  `json:"data,omitempty"`
}

// swagger:model Pagination
// @name Pagination
// @Description: Estruturas de resposta para paginação
type Pagination struct {
	Limit   int         `json:"limit" example:"10"`
	Page    int         `json:"page" example:"1"`
	Total   int         `json:"total" example:"100"`
	Content interface{} `json:"content"`
}

// swagger:model FieldError
// @name FieldError
// @Description: Estruturas de erro de validação
type FieldError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email inválido"`
}

// swagger:model ValidationErrorResponse
// @name ValidationErrorResponse
// @Description: Estruturas de resposta de erro de validação
type ValidationErrorResponse struct {
	Status  Status       `json:"status" example:"error"`
	Message string       `json:"message" example:"Erro de validação"`
	Fields  []FieldError `json:"fields" example:"[{\"field\":\"email\",\"message\":\"Email inválido\"}]"`
}
