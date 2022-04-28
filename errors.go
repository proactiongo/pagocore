package pagocore

import (
	"fmt"
	"net/http"
)

// Errors
var (
	ErrUserBlocked      = NewError(http.StatusForbidden, "requested user is blocked")
	ErrPassFailed       = NewError(http.StatusUnauthorized, "password is invalid")
	ErrRoleNotAllowed   = NewError(http.StatusForbidden, "unexpected user role")
	ErrTokenInvalid     = NewError(http.StatusUnauthorized, "token is invalid")
	ErrTokenExpired     = NewError(http.StatusUnauthorized, "token is expired")
	ErrTokenUnsupported = NewError(http.StatusUnprocessableEntity, "unsupported sign method")
	ErrNotFound         = NewError(http.StatusNotFound, "not found")
)

// NewError creates a new Error instance
func NewError(code int, messages ...interface{}) *Error {
	if len(messages) == 0 && code >= 100 && code <= 599 {
		messages = []interface{}{
			http.StatusText(code),
		}
	}
	return &Error{
		Code:    code,
		Message: fmt.Sprint(messages...),
	}
}

// Error defines the response error
type Error struct {
	Code      int    `json:"code" example:"403"`
	Message   string `json:"message" example:"Access denied"`
	Localized string `json:"localized,omitempty" example:"Доступ запрещен"`
}

// Error as a string
func (e *Error) Error() string {
	return e.Message
}
