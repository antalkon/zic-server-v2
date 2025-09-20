package router

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer реализует echo.Renderer для html/template
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if data == nil {
		data = echo.Map{}
	}
	// сюда можно подмешивать глобальные данные, если нужно:
	// data.(echo.Map)["User"] = c.Get("user")
	return t.templates.ExecuteTemplate(w, name, data)
}

// SetupFrontendRoutes — подключаем статику и шаблоны
func SetupFrontendRoutes(e *echo.Echo) {
	// 1) Парсим все html из public/pages и public/components
	// важно: в файлах должны быть {{ define "имя" }} ... {{ end }}
	t := template.Must(
		template.New("").
			// при желании можно добавить FuncMap: .Funcs(template.FuncMap{"upper": strings.ToUpper})
			ParseGlob("web/public/**/*.html"),
	)
	e.Renderer = &TemplateRenderer{templates: t}

	// 2) Статика
	e.Static("/static", "web/static")
	e.Static("/public", "web/public")

	// 3) Страницы
	// login/index у тебя как обычные файлы (без шаблонов) — оставим так
	e.GET("/login", func(c echo.Context) error {
		return c.File("web/public/pages/login.html")
	})
	e.GET("/", func(c echo.Context) error {
		return c.File("web/public/pages/index.html")
	})

	// 4) dashboard рендерим как шаблон (используется define "dashboard")
	e.GET("/dashboard", func(c echo.Context) error {
		// при необходимости передавай данные в шаблон
		data := echo.Map{
			"Title": "Dashboard",
		}
		return c.Render(http.StatusOK, "dashboard", data)
	})
}
