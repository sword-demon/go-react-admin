// Package errors provides unified error types for the application
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents business error codes
type ErrorCode int

const (
	// Common errors (1000-1999)
	ErrInvalidParams ErrorCode = 1000 + iota
	ErrInternalServer
	ErrNotFound
	ErrUnauthorized
	ErrForbidden
	ErrConflict
	ErrBadRequest
	ErrTooManyRequests
)

const (
	// User errors (2000-2999)
	ErrUserNotFound ErrorCode = 2000 + iota
	ErrUserAlreadyExists
	ErrUserInvalidCredentials
	ErrUserInvalidPassword
	ErrUserDisabled
	ErrUserInvalidStatus
)

const (
	// Role errors (3000-3999)
	ErrRoleNotFound ErrorCode = 3000 + iota
	ErrRoleAlreadyExists
	ErrRoleInUse
	ErrRoleInvalidKey
)

const (
	// Dept errors (4000-4999)
	ErrDeptNotFound ErrorCode = 4000 + iota
	ErrDeptAlreadyExists
	ErrDeptHasChildren
	ErrDeptHasUsers
)

const (
	// Menu errors (5000-5999)
	ErrMenuNotFound ErrorCode = 5000 + iota
	ErrMenuAlreadyExists
	ErrMenuHasChildren
)

const (
	// Permission errors (6000-6999)
	ErrPermissionDenied ErrorCode = 6000 + iota
	ErrPermissionNotFound
	ErrPermissionInvalidPattern
)

// Error represents a business error with code and message
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"` // Internal error (not exposed to client)
}

// Error implements error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap returns the internal error
func (e *Error) Unwrap() error {
	return e.Err
}

// HTTPStatus returns the HTTP status code for the error
func (e *Error) HTTPStatus() int {
	switch e.Code {
	case ErrNotFound, ErrUserNotFound, ErrRoleNotFound, ErrDeptNotFound, ErrMenuNotFound, ErrPermissionNotFound:
		return http.StatusNotFound
	case ErrUnauthorized, ErrUserInvalidCredentials:
		return http.StatusUnauthorized
	case ErrForbidden, ErrPermissionDenied:
		return http.StatusForbidden
	case ErrConflict, ErrUserAlreadyExists, ErrRoleAlreadyExists, ErrDeptAlreadyExists, ErrMenuAlreadyExists:
		return http.StatusConflict
	case ErrBadRequest, ErrInvalidParams:
		return http.StatusBadRequest
	case ErrTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// New creates a new Error with code and message
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with code and message
func Wrap(code ErrorCode, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is checks if the error is of the specified code
func Is(err error, code ErrorCode) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// GetCode extracts the error code from error
func GetCode(err error) ErrorCode {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ErrInternalServer
}

// Predefined errors for common cases
var (
	// Common
	ErrRecordNotFound      = New(ErrNotFound, "record not found")
	ErrRecordAlreadyExists = New(ErrConflict, "record already exists")
	ErrInvalidInput        = New(ErrInvalidParams, "invalid input parameters")

	// User
	ErrUsernameExists       = New(ErrUserAlreadyExists, "username already exists")
	ErrInvalidCredentials   = New(ErrUserInvalidCredentials, "invalid username or password")
	ErrInvalidOldPassword   = New(ErrUserInvalidPassword, "old password is incorrect")
	ErrUserNotFoundError    = New(ErrUserNotFound, "user not found")
	ErrUserDisabledError    = New(ErrUserDisabled, "user is disabled")

	// Role
	ErrRoleKeyExists        = New(ErrRoleAlreadyExists, "role key already exists")
	ErrRoleNotFoundError    = New(ErrRoleNotFound, "role not found")
	ErrRoleInUseError       = New(ErrRoleInUse, "role is in use by users")

	// Dept
	ErrDeptNotFoundError    = New(ErrDeptNotFound, "department not found")
	ErrDeptHasChildrenError = New(ErrDeptHasChildren, "department has children")
	ErrDeptHasUsersError    = New(ErrDeptHasUsers, "department has users")

	// Menu
	ErrMenuNotFoundError    = New(ErrMenuNotFound, "menu not found")
	ErrMenuHasChildrenError = New(ErrMenuHasChildren, "menu has children")

	// Permission
	ErrPermissionDeniedError = New(ErrPermissionDenied, "permission denied")
)