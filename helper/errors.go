package helper

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	CodeValidation           = "VALIDATION_ERROR"
	CodeBadRequest           = "BAD_REQUEST"
	CodeUnauthorized         = "UNAUTHORIZED"
	CodeForbidden            = "FORBIDDEN"
	CodeNotFound             = "NOT_FOUND"
	CodeConflict             = "CONFLICT"
	CodeUnsupportedMediaType = "UNSUPPORTED_MEDIA_TYPE"
	CodeNotAcceptable        = "NOT_ACCEPTABLE"
	CodeTooManyRequests      = "TOO_MANY_REQUESTS"
	CodeInternal             = "INTERNAL_ERROR"
)

type AppError struct {
	Status  int
	Code    string
	Message string
	Fields  map[string]string
	cause   error
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
func (e *AppError) Unwrap() error                 { return e.cause }
func (e *AppError) WithCause(err error) *AppError { e.cause = err; return e }
func BadRequest(message string) *AppError {
	return &AppError{fiber.StatusBadRequest, CodeBadRequest, message, nil, nil}
}
func Unauthorized(message string) *AppError {
	return &AppError{fiber.StatusUnauthorized, CodeUnauthorized, message, nil, nil}
}
func Forbidden(message string) *AppError {
	return &AppError{fiber.StatusForbidden, CodeForbidden, message, nil, nil}
}
func NotFound(message string) *AppError {
	return &AppError{fiber.StatusNotFound, CodeNotFound, message, nil, nil}
}
func Conflict(message string) *AppError {
	return &AppError{fiber.StatusConflict, CodeConflict, message, nil, nil}
}
func Validation(fields map[string]string) *AppError {
	return &AppError{fiber.StatusUnprocessableEntity, CodeValidation, "validasi gagal", fields, nil}
}
func NotAcceptable(message string) *AppError {
	return &AppError{fiber.StatusNotAcceptable, CodeNotAcceptable, message, nil, nil}
}
func UnsupportedMediaType(message string) *AppError {
	return &AppError{fiber.StatusUnsupportedMediaType, CodeUnsupportedMediaType, message, nil, nil}
}
func TooManyRequests(message string) *AppError {
	return &AppError{fiber.StatusTooManyRequests, CodeTooManyRequests, message, nil, nil}
}
func Internal(cause error) *AppError {
	return &AppError{fiber.StatusInternalServerError, CodeInternal, "terjadi kesalahan pada server", nil, cause}
}

func Fail(_ *fiber.Ctx, status int, message string) error {
	switch status {
	case fiber.StatusBadRequest:
		return BadRequest(message)
	case fiber.StatusUnauthorized:
		return Unauthorized(message)
	case fiber.StatusForbidden:
		return Forbidden(message)
	case fiber.StatusNotFound:
		return NotFound(message)
	case fiber.StatusConflict:
		return Conflict(message)
	case fiber.StatusUnsupportedMediaType:
		return UnsupportedMediaType(message)
	case fiber.StatusUnprocessableEntity:
		return Validation(map[string]string{"_": message})
	default:
		return Internal(nil)
	}
}
