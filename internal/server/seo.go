package server

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// origin resolves the absolute base URL robots.txt and sitemap.xml must
// contain; neither format has a relative form.
//
// SITE_URL wins when set; the Host fallback keeps dev and CI working with no
// configuration. Safe here in a way it would NOT be for canonical/og:url, which
// stay gated on SITE_URL — a forged Host on a canonical tag misdirects other
// people's traffic, because the tag is embedded in a page real visitors fetch.
// A sitemap is fetched by the crawler, which sends the true Host, so forging it
// only returns a bogus document to whoever forged it.
func origin(c echo.Context, siteURL string) string {
	if siteURL != "" {
		return siteURL
	}
	return c.Scheme() + "://" + c.Request().Host
}

// Robots serves /robots.txt, generated rather than kept as a static file so the
// Sitemap line always carries the origin the request actually arrived on.
//
// It deliberately Disallows nothing except /stats. Listing /stats is a
// judgement call in the other direction from a private admin path: it is
// already unguessable-free (the token is the guard, not the URL), and keeping
// it out of the index avoids a page of traffic counts showing up in search
// results for your own name.
func (s *Server) Robots(c echo.Context) error {
	var b strings.Builder
	b.WriteString("User-agent: *\n")
	b.WriteString("Allow: /\n")
	b.WriteString("Disallow: /stats\n")
	b.WriteString("\n")
	fmt.Fprintf(&b, "Sitemap: %s/sitemap.xml\n", origin(c, s.cfg.SiteURL))
	return c.String(http.StatusOK, b.String())
}

type sitemapURL struct {
	Loc        string `xml:"loc"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	NS      string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// Sitemap serves /sitemap.xml, built from the pages actually served rather than
// kept as a static file. A static sitemap advertising a 404 is a Search Console
// error, and one omitting a live page may never get that page crawled; both
// happen the first time a route changes and the file does not.
//
// /stats and /resume.pdf are excluded: the first is private, and the second is
// a document, not a page — listing it competes with the pages that should rank.
func (s *Server) Sitemap(c echo.Context) error {
	base := origin(c, s.cfg.SiteURL)
	set := sitemapURLSet{
		NS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs: []sitemapURL{
			{Loc: base + "/", ChangeFreq: "monthly", Priority: "1.0"},
			{Loc: base + "/work", ChangeFreq: "monthly", Priority: "0.9"},
			{Loc: base + "/about", ChangeFreq: "monthly", Priority: "0.8"},
		},
	}
	return c.XMLPretty(http.StatusOK, set, "  ")
}
