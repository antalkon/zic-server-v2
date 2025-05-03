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

// 400 Bad Request
func BadRequestError() (int, *ErrorResponse) {
	return NewError(http.StatusBadRequest, "Bad Request")
}

// 404 Not Found
func NotFoundError() (int, *ErrorResponse) {
	return NewError(http.StatusNotFound, "Resource not found")
}

// 409 Conflict
func ConflictError() (int, *ErrorResponse) {
	return NewError(http.StatusConflict, "Conflict: data already exists or violates constraints")
}

// 401 Unauthorized
func UnauthorizedError() (int, *ErrorResponse) {
	return NewError(http.StatusUnauthorized, "Unauthorized access")
}

// 403 Forbidden
func ForbiddenError() (int, *ErrorResponse) {
	return NewError(http.StatusForbidden, "Access forbidden")
}

// 500 Internal Server Error
func InternalServerError(msg string) (int, *ErrorResponse) {
	return NewError(http.StatusInternalServerError, msg)
}

// ----------------------------------------
// Дополнительные кастомные ошибки

// 400 Bad Request — ошибка валидации
func ValidationError() (int, *ErrorResponse) {
	return NewError(http.StatusBadRequest, "Validation failed")
}

// 401 Unauthorized — токен истёк или недействителен
func TokenExpiredError() (int, *ErrorResponse) {
	return NewError(http.StatusUnauthorized, "Token expired or invalid")
}

// 401 Unauthorized — отсутствует refresh токен
func MissingTokenError() (int, *ErrorResponse) {
	return NewError(http.StatusUnauthorized, "Missing refresh token")
}

// 409 Conflict — логин уже существует
func LoginAlreadyExistsError() (int, *ErrorResponse) {
	return NewError(http.StatusConflict, "User with this login already exists")
}

// 429 Too Many Requests
func RateLimitExceededError() (int, *ErrorResponse) {
	return NewError(http.StatusTooManyRequests, "Too many requests, please try again later")
}

// 405 Method Not Allowed
func MethodNotAllowedError() (int, *ErrorResponse) {
	return NewError(http.StatusMethodNotAllowed, "Method not allowed on this endpoint")
}

// 409 Conflict — кастомное сообщение
func ConflictCustomError(msg string) (int, *ErrorResponse) {
	return NewError(http.StatusConflict, msg)
}
