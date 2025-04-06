package middleware

import (
	"backend/internal/repository"
	tokenjwt "backend/pkg/token_jwt"
	"net/http"
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
			cookie, err := c.Cookie("access_token")
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/api/v1/auth/refresh-token")
			}

			if cookie.Value == "" {
				return c.Redirect(http.StatusTemporaryRedirect, "/api/v1/auth/refresh-token")
			}

			claims, err := tokenjwt.DecodeJWT(cookie.Value)
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/api/v1/auth/refresh-token")
			}

			// Проверяем срок действия токена
			if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
				return c.Redirect(http.StatusTemporaryRedirect, "/api/v1/auth/refresh-token")
			}

			user, err := m.authRepo.GetUserById(claims.UserID)
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/api/v1/auth/refresh-token")
			}

			// Получаем роль пользователя
			role, err := m.authRepo.GetUserRole(user.RoleID)
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/api/v1/auth/refresh-token")
			}

			// Устанавливаем ID пользователя и роль в контекст
			c.Set("user_id", user.ID)
			c.Set("user_role", role.Name)

			return next(c)
		}
	}
}
