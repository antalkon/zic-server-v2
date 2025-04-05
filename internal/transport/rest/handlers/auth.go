package handlers

import (
	"net/http"
	"os"

	"backend/internal/models"
	"backend/internal/transport/rest/req"
	"backend/internal/transport/rest/res"
	"backend/internal/transport/service"
	"backend/internal/utils"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	auth     *service.AuthService
	validate *validator.Validate
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{
		auth:     auth,
		validate: validator.New(),
	}
}

func (h *AuthHandler) SignUpUser(c echo.Context) error {
	var client models.User

	if err := c.Bind(&client); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	if err := h.validate.Struct(client); err != nil {
		code, msg := utils.ValidationError()
		return c.JSON(code, msg)
	}

	_, err := h.auth.SignUpClient(&client, c)
	if err != nil {
		code, msg := utils.ConflictError()
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusCreated, res.SignUpRes{
		Message: "User created successfully",
	})
}

func (h *AuthHandler) SignInUser(c echo.Context) error {
	var req req.SignInReq

	if err := c.Bind(&req); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	if err := h.validate.Struct(req); err != nil {
		code, msg := utils.ValidationError()
		return c.JSON(code, msg)
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	tokens, err := h.auth.SignInClient(&user, c)
	if err != nil {
		code, msg := utils.UnauthorizedError()
		return c.JSON(code, msg)
	}

	setAuthCookies(c, tokens)

	return c.JSON(http.StatusOK, res.SignInRes{
		Message: "User signed in successfully",
	})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil || cookie == nil || cookie.Value == "" {
		return c.JSON(utils.MissingTokenError())
	}

	tokens, err := h.auth.RefreshToken(cookie.Value)
	if err != nil {
		return c.JSON(utils.InternalServerError("failed to refresh token: " + err.Error()))
	}

	setAuthCookies(c, tokens)

	return c.JSON(http.StatusOK, res.SignInRes{
		Message: "Tokens updated successfully",
	})
}

func (h *AuthHandler) SignOutUser(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		return c.JSON(utils.MissingTokenError())
	}

	// Удаляем refresh токен из БД
	if err := h.auth.SignOut(cookie.Value); err != nil {
		return c.JSON(utils.InternalServerError("failed to sign out: " + err.Error()))
	}

	// Чистим refresh_token cookie
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Чистим access_token cookie
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return c.JSON(http.StatusOK, res.SignOutRes{
		Message: "User signed out successfully",
	})
}

// setAuthCookies устанавливает access/refresh токены в куки
func setAuthCookies(c echo.Context, tokens models.AuthTokens) {
	isDev := os.Getenv("APP_ENV") == "development"

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   !isDev,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   !isDev,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 30, // 30 дней
	})
}
