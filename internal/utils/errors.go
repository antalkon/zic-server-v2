package utils

import "net/http"

type ErrorResponse struct {
	Message string `json:"message"`
}

// NewError возвращает объект ошибки с кодом
func NewError(code int, msg string) (int, *ErrorResponse) {
	return code, &ErrorResponse{Message: msg}
}

// ----------------------------------------
// Базовые статусы

func BadRequestError() (int, *ErrorResponse) {
	return NewError(http.StatusBadRequest, "Bad Request")
}

func NotFoundError() (int, *ErrorResponse) {
	return NewError(http.StatusNotFound, "Resource not found")
}

func ConflictError() (int, *ErrorResponse) {
	return NewError(http.StatusConflict, "Conflict: data already exists or violates constraints")
}

func UnauthorizedError() (int, *ErrorResponse) {
	return NewError(http.StatusUnauthorized, "Unauthorized access")
}

func ForbiddenError() (int, *ErrorResponse) {
	return NewError(http.StatusForbidden, "Access forbidden")
}

func InternalServerError(msg string) (int, *ErrorResponse) {
	return NewError(http.StatusInternalServerError, msg)
}

// ----------------------------------------
// Дополнительные кастомные ошибки

func ValidationError() (int, *ErrorResponse) {
	return NewError(http.StatusBadRequest, "Validation failed")
}

func TokenExpiredError() (int, *ErrorResponse) {
	return NewError(http.StatusUnauthorized, "Token expired or invalid")
}

func MissingTokenError() (int, *ErrorResponse) {
	return NewError(http.StatusUnauthorized, "Missing refresh token")
}

func LoginAlreadyExistsError() (int, *ErrorResponse) {
	return NewError(http.StatusConflict, "User with this login already exists")
}

func RateLimitExceededError() (int, *ErrorResponse) {
	return NewError(http.StatusTooManyRequests, "Too many requests, please try again later")
}

func MethodNotAllowedError() (int, *ErrorResponse) {
	return NewError(http.StatusMethodNotAllowed, "Method not allowed on this endpoint")
}
