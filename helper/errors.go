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
	Errors  map[string][]string
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
	return &AppError{Status: fiber.StatusBadRequest, Code: CodeBadRequest, Message: message}
}

func Unauthorized(message string) *AppError {
	return &AppError{Status: fiber.StatusUnauthorized, Code: CodeUnauthorized, Message: message}
}

func Forbidden(message string) *AppError {
	return &AppError{Status: fiber.StatusForbidden, Code: CodeForbidden, Message: message}
}

func NotFound(message string) *AppError {
	return &AppError{Status: fiber.StatusNotFound, Code: CodeNotFound, Message: message}
}

func Conflict(message string) *AppError {
	return &AppError{Status: fiber.StatusConflict, Code: CodeConflict, Message: message}
}

func Validation(fields map[string]string) *AppError {
	errMap := make(map[string][]string)
	for k, v := range fields {
		errMap[k] = []string{v}
	}
	return &AppError{
		Status:  fiber.StatusUnprocessableEntity,
		Code:    CodeValidation,
		Message: "Validasi gagal",
		Fields:  fields,
		Errors:  errMap,
	}
}

func ValidationErrors(errors map[string][]string) *AppError {
	return &AppError{
		Status:  fiber.StatusUnprocessableEntity,
		Code:    CodeValidation,
		Message: "Validasi gagal",
		Errors:  errors,
	}
}

func ValidationField(field string, message string) *AppError {
	return &AppError{
		Status:  fiber.StatusUnprocessableEntity,
		Code:    CodeValidation,
		Message: "Validasi gagal",
		Errors:  map[string][]string{field: {message}},
	}
}

func NotAcceptable(message string) *AppError {
	return &AppError{Status: fiber.StatusNotAcceptable, Code: CodeNotAcceptable, Message: message}
}

func UnsupportedMediaType(message string) *AppError {
	return &AppError{Status: fiber.StatusUnsupportedMediaType, Code: CodeUnsupportedMediaType, Message: message}
}

func TooManyRequests(message string) *AppError {
	return &AppError{Status: fiber.StatusTooManyRequests, Code: CodeTooManyRequests, Message: message}
}

func Internal(cause error) *AppError {
	return &AppError{
		Status:  fiber.StatusInternalServerError,
		Code:    CodeInternal,
		Message: "Terjadi kesalahan internal pada server",
		cause:   cause,
	}
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
	case fiber.StatusTooManyRequests:
		return TooManyRequests(message)
	case fiber.StatusUnprocessableEntity:
		return ValidationField("general", message)
	default:
		return Internal(nil)
	}
}
