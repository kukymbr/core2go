package errors

import (
	"fmt"
	"net/http"
)

// NewExposable creates new ExposableError instance.
func NewExposable(code int, messages ...any) error {
	if len(messages) == 0 && code >= 100 && code <= 599 {
		messages = []any{
			http.StatusText(code),
		}
	}

	return &ExposableError{
		Code:    code,
		Message: fmt.Sprint(messages...),
	}
}

// ExposableError is an error to expose in the user's interface.
// If Localized field is not empty, interface should use it to show to the user.
type ExposableError struct {
	Code      int    `json:"code" example:"404"`
	Message   string `json:"message" example:"err.notfound"`
	Localized string `json:"localized,omitempty" example:"Not Found"`
}

// Error returns error as a string.
func (e *ExposableError) Error() string {
	return e.Message
}
