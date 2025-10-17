package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	BadRequest          = "Invalid request"
	Unauthorized        = "Unauthorized"
	Forbidden           = "Access forbidden"
	NotFound            = "Resource not found"
	Conflict            = "Data conflict"
	UnprocessableEntity = "Unprocessable entity"
	TooManyRequests     = "Too many requests, please try again later"
	InternalServerError = "Internal server error"
	ServiceUnavailable  = "Service temporarily unavailable"
)

var ErrorMessages = map[int]string{

	http.StatusBadRequest:          BadRequest,
	http.StatusUnauthorized:        Unauthorized,
	http.StatusForbidden:           Forbidden,
	http.StatusNotFound:            NotFound,
	http.StatusConflict:            Conflict,
	http.StatusUnprocessableEntity: UnprocessableEntity,
	http.StatusTooManyRequests:     TooManyRequests,
	http.StatusInternalServerError: InternalServerError,
	http.StatusServiceUnavailable:  ServiceUnavailable,
}

var (
	ErrBadRequest         = errors.New(strings.ToLower(BadRequest))
	ErrUnauthorized       = errors.New(strings.ToLower(Unauthorized))
	ErrForbidden          = errors.New(strings.ToLower(Forbidden))
	ErrNotFound           = errors.New(strings.ToLower(NotFound))
	ErrConflict           = errors.New(strings.ToLower(Conflict))
	ErrUnprocessable      = errors.New(strings.ToLower(UnprocessableEntity))
	ErrTooManyRequests    = errors.New(strings.ToLower(TooManyRequests))
	ErrInternalServer     = errors.New(strings.ToLower(InternalServerError))
	ErrServiceUnavailable = errors.New(strings.ToLower(ServiceUnavailable))
	ErrEmptyField         = errors.New("empty field")

	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenExpired            = errors.New("token expired")
	ErrInvalidToken            = errors.New("invalid token")
	ErrUserNotFoundInCtx       = errors.New("user not found in context")
)

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"details,omitempty"`
}

type AppError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	return e.Err.Error()
}
func NewAppError(statusCode int, err error) *AppError {
	return &AppError{
		Message:    ErrorMessages[statusCode],
		StatusCode: statusCode,
		Err:        err,
	}
}

func FailedToFind(entity string, err error) error {
	return fmt.Errorf("failed to find %s: %w", entity, err)
}

func FailedToSave(entity string, err error) error {
	return fmt.Errorf("failed to save %s: %w", entity, err)
}

func FailedToUpdate(entity string, err error) error {
	return fmt.Errorf("failed to update %s: %w", entity, err)
}

func FailedToDelete(entity string, err error) error {
	return fmt.Errorf("failed to delete %s: %w", entity, err)
}

func EmptyField(field string) error {
	return fmt.Errorf("%s cannot be empty", field)
}

func InvalidParam(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
