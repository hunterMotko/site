package server

import (
	"html/template"
	"io"
	"net/http"

	"github.com/huntermotko/site/cmd/api/about"
	"github.com/huntermotko/site/cmd/api/tools"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
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

	e.GET("/", s.healthHandler)
	e.GET("/", s.Home)
	e.GET("/about", s.About)
	e.GET("/tools", s.Tools)

	e.Logger.Fatal(e.Start(":8080"))
	return e
}

func (s *Server) Home(c echo.Context) error {
	return c.Render(200, "index", nil)
}

func (s *Server) About(c echo.Context) error {
  a := about.GetAbout()
	return c.Render(200, "about", a)
}

func (s *Server) Tools(c echo.Context) error {
	return c.Render(200, "tools", tools.GetTools())
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
