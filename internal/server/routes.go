package server

import (
	"html/template"
	"io"
	"net/http"

	"github.com/huntermotko/site/cmd/api/about"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(
	w io.Writer,
	name string,
	data any,
	c echo.Context,
) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("internal/app/views/*/*.html")),
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "internal/app/public")

	e.Renderer = newTemplate()
	e.GET("/", s.Home)
	e.GET("/about", s.About)

	// Return the handler rather than calling e.Start(). Starting Echo's own
	// server here would block forever, so the *http.Server configured in
	// NewServer() — and its read/write/idle timeouts — would never be used.
	return e
}

func (s *Server) Home(c echo.Context) error {
	a := about.GetAbout()
	return c.Render(200, "index", a)
}

func (s *Server) About(c echo.Context) error {
	a := about.GetAbout()
	return c.Render(200, "about", a)
}
