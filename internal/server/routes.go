package server

import (
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/huntermotko/site/cmd/api/about"
	"github.com/huntermotko/site/cmd/api/work"
	"github.com/huntermotko/site/internal/app"
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

// newTemplate parses the views out of the embedded filesystem rather than off
// disk. The glob is relative to the embed root, so it is "views/*/*.html" and
// not the old repo-relative "internal/app/views/*/*.html" — which only resolved
// when the process happened to be started from the repo root.
func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseFS(app.Views, "views/*/*.html")),
	}
}

// homeVM and the other view models embed their content struct rather than
// nesting it under a field. Embedding promotes the fields, so templates keep
// referring to {{ .Positioning }} directly while {{ .Meta }} becomes available
// alongside — adding page metadata without rewriting every reference in the
// templates.
type homeVM struct {
	about.About
	Meta Meta
}

type aboutVM struct {
	about.About
	Meta Meta
}

type workVM struct {
	work.Work
	Meta Meta
}

// page registers a handler for GET and HEAD.
//
// Echo's e.GET binds GET alone, so a HEAD request got 405 Method Not Allowed.
// That is not a curiosity: uptime monitors, link checkers, and several link
// unfurlers probe with HEAD by default, and a 405 reads to them as the site
// being broken. Go's net/http suppresses the body for HEAD responses on its
// own, so the same handler serves both correctly.
func page(e *echo.Echo, path string, h echo.HandlerFunc) {
	e.Match([]string{http.MethodGet, http.MethodHead}, path, h)
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())

	if s.db != nil {
		e.Use(s.recordPageview)
	}

	// Serving the embedded public/ subtree at the root, so asset URLs stay
	// "/css/index.css" rather than gaining a "/public" prefix the templates
	// do not use.
	pub, err := fs.Sub(app.Public, "public")
	if err != nil {
		panic("embedded public assets missing: " + err.Error())
	}
	e.StaticFS("/", pub)

	e.Renderer = newTemplate()

	page(e, "/", s.Home)
	page(e, "/work", s.Work)
	page(e, "/about", s.About)
	page(e, "/robots.txt", s.Robots)
	page(e, "/sitemap.xml", s.Sitemap)

	// The resume route exists only when the file does, so the contact block
	// never advertises a download that 404s. Checked once at startup rather
	// than per request: the file does not appear while the process is running,
	// and a stat on every request to decide whether to render a link is a
	// syscall for nothing.
	if s.cfg.ResumePath != "" {
		if _, err := os.Stat(s.cfg.ResumePath); err == nil {
			s.hasResume = true
			e.File("/resume.pdf", s.cfg.ResumePath)
		}
	}

	if s.db != nil && s.cfg.StatsToken != "" {
		page(e, "/stats", s.Stats)
	}

	// Return the handler rather than calling e.Start(). Starting Echo's own
	// server here would block forever, so the *http.Server configured in
	// NewServer() — and its read/write/idle timeouts — would never be used.
	return e
}

func (s *Server) Home(c echo.Context) error {
	return c.Render(http.StatusOK, "index", homeVM{
		About: about.GetAbout(),
		Meta: s.newMeta(
			"Hunter Motko — Platform & Backend Engineer",
			"I build production systems and run the machines they live on. Go, TypeScript, Postgres, Docker, Linux. Remote, US teams.",
			"/",
		),
	})
}

func (s *Server) About(c echo.Context) error {
	return c.Render(http.StatusOK, "about", aboutVM{
		About: about.GetAbout(),
		Meta: s.newMeta(
			"About — Hunter Motko",
			"How I got here, what I have shipped, and the route through The Last Mile and The Next Chapter that got me into engineering.",
			"/about",
		),
	})
}

func (s *Server) Work(c echo.Context) error {
	return c.Render(http.StatusOK, "work", workVM{
		Work: work.GetWork(),
		Meta: s.newMeta(
			"Work — Hunter Motko",
			"Two production systems: an inventory platform for a shed builder, and statewide public-defense data commissioned by the Illinois Supreme Court.",
			"/work",
		),
	})
}
