package entity

import "errors"

// ErrNotFound é retornado quando um registro não é encontrado no banco.
var ErrNotFound = errors.New("not found")

// ValidationError representa um erro de validação de domínio (400 Bad Request).
// Permite que handlers distingam erros de validação de erros de infraestrutura
// sem depender de heurísticas frágeis como errors.Unwrap.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// NewValidationError cria um erro de validação com a mensagem fornecida.
func NewValidationError(msg string) error {
	return &ValidationError{Message: msg}
}
