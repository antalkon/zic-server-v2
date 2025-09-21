package middleware

import (
	"backend/internal/repository"
	tokenjwt "backend/pkg/token_jwt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	authRepo *repository.AuthRepository
}

func NewAuthMiddleware(authRepo *repository.AuthRepository) *AuthMiddleware {
	return &AuthMiddleware{
		authRepo: authRepo,
	}
}

func (m *AuthMiddleware) AuthRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			// 1) Никогда не перехватываем сам refresh-эндпоинт, иначе петля
			if path == "/api/v1/auth/refresh-token" {
				return next(c)
			}

			// 2) Достаём access_token
			cookie, err := c.Cookie("access_token")
			if err != nil || cookie == nil || cookie.Value == "" {
				return m.unauthorized(c)
			}

			claims, err := tokenjwt.DecodeJWT(cookie.Value)
			if err != nil {
				return m.unauthorized(c)
			}

			// 3) Проверяем срок действия токена
			if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
				return m.unauthorized(c)
			}

			// 4) Достаём пользователя и роль
			user, err := m.authRepo.GetUserById(claims.UserID)
			if err != nil {
				return m.unauthorized(c)
			}
			role, err := m.authRepo.GetUserRole(user.RoleID)
			if err != nil {
				return m.unauthorized(c)
			}

			// 5) Кладём в контекст
			c.Set("user_id", user.ID)
			c.Set("user_role", role.Name)

			return next(c)
		}
	}
}

// Поведение при неавторизованности:
// - для API: 401 JSON
// - для страниц: редирект на /login
func (m *AuthMiddleware) unauthorized(c echo.Context) error {
	path := c.Request().URL.Path

	// API → 401 JSON
	if strings.HasPrefix(path, "/api/") {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error":   "unauthorized",
			"message": "Access token is missing or invalid",
		})
	}

	// Страницы → редирект на refresh-token с next
	next := url.QueryEscape(c.Request().URL.RequestURI())
	return c.Redirect(http.StatusTemporaryRedirect, "/api/v1/auth/refresh-token?next="+next)
}
