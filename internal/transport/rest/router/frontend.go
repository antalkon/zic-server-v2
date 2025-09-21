// router/frontend.go
package router

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct{ templates *template.Template }

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if data == nil {
		data = echo.Map{}
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

// Принимаем готовый echo.MiddlewareFunc
func SetupFrontendRoutes(e *echo.Echo, authMW echo.MiddlewareFunc) {
	// Шаблоны
	t := template.Must(template.New("").ParseGlob("web/public/**/*.html"))
	e.Renderer = &TemplateRenderer{templates: t}

	// Статика
	e.Static("/static", "web/static")
	e.Static("/public", "web/public")

	// Публичные
	e.GET("/login", func(c echo.Context) error { return c.File("web/public/pages/login.html") })
	e.GET("/", func(c echo.Context) error { return c.File("web/public/pages/index.html") })

	// Защищённые
	protected := e.Group("", authMW) // <-- тут БЕЗ .AuthRequired()

	protected.GET("/dashboard", func(c echo.Context) error {
		data := echo.Map{"Title": "Dashboard"}
		return c.Render(http.StatusOK, "dashboard", data)
	})

	protected.GET("/rooms", func(c echo.Context) error {
		data := echo.Map{"Title": "Кабинеты"}
		return c.Render(http.StatusOK, "rooms", data)
	})
}
