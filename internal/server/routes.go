package server

import (
	"html/template"
	"io"
	"net/http"

	"github.com/huntermotko/site/cmd/api/about"
	"github.com/huntermotko/site/cmd/api/resources"
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
  e.GET("/resources", s.Resources)
	e.GET("/tools", s.Tools)
	e.Logger.Fatal(e.Start(":8080"))
	return e
}

type Name struct {
  First string
  Last string
}

func (s *Server) Home(c echo.Context) error {
  name := Name{
    First: `
 ___  ___  ___  ___  ________   _________  _______   ________     
|\  \|\  \|\  \|\  \|\   ___  \|\___   ___\\  ___ \ |\   __  \    
\ \  \\\  \ \  \\\  \ \  \\ \  \|___ \  \_\ \   __/|\ \  \|\  \   
 \ \   __  \ \  \\\  \ \  \\ \  \   \ \  \ \ \  \_|/_\ \   _  _\  
  \ \  \ \  \ \  \\\  \ \  \\ \  \   \ \  \ \ \  \_|\ \ \  \\  \| 
   \ \__\ \__\ \_______\ \__\\ \__\   \ \__\ \ \_______\ \__\\ _\ 
    \|__|\|__|\|_______|\|__| \|__|    \|__|  \|_______|\|__|\|__|
`,
    Last: `
 _____ ______   ________  _________  ___  __    ________     
|\   _ \  _   \|\   __  \|\___   ___\\  \|\  \ |\   __  \    
\ \  \\\__\ \  \ \  \|\  \|___ \  \_\ \  \/  /|\ \  \|\  \   
 \ \  \\|__| \  \ \  \\\  \   \ \  \ \ \   ___  \ \  \\\  \  
  \ \  \    \ \  \ \  \\\  \   \ \  \ \ \  \\ \  \ \  \\\  \ 
   \ \__\    \ \__\ \_______\   \ \__\ \ \__\\ \__\ \_______\
    \|__|     \|__|\|_______|    \|__|  \|__| \|__|\|_______|
`,
  }
	return c.Render(200, "index", name)
}
func (s *Server) About(c echo.Context) error {
	a := about.GetAbout()
	return c.Render(200, "about", a)
}
func (s *Server) Resources(c echo.Context) error {
  res := resources.GetResources()
	return c.Render(200, "resources", res)
}
func (s *Server) Tools(c echo.Context) error {
  tls := tools.GetTools()
	return c.Render(200, "tools", tls)
}
func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
