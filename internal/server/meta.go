package server

import (
	"encoding/json"
	"html/template"
	"log"
)

// Meta is everything the shared <head> and navbar need that is not page
// content: the title bar, the description crawlers and link unfurls read, and
// the path, which drives both the canonical URL and which nav item is marked
// current.
type Meta struct {
	Title       string
	Description string
	Path        string
	SiteURL     string
	Dev         bool

	// HasResume mirrors whether the /resume.pdf route was registered. The
	// contact block omits the download link when it is false, so there is no
	// state in which the site advertises a file it will 404 on.
	HasResume bool

	// JSONLD is the structured-data block, marshalled in Go rather than
	// hand-written in the template. html/template will not escape inside a
	// <script> element the way it does in HTML body context, so building the
	// JSON by string concatenation in a template is how apostrophes in prose
	// end up producing invalid structured data — or worse, breaking out of the
	// script tag. Marshalling makes correct escaping the default.
	JSONLD template.JS
}

// CanonicalURL is the absolute URL for this page, or "" when SITE_URL is unset.
// Templates omit the canonical and og:url tags entirely when it is empty, which
// is the correct behaviour: no canonical tag is strictly better than a wrong
// one, since a wrong one actively redirects ranking signal somewhere else.
func (m Meta) CanonicalURL() string {
	if m.SiteURL == "" {
		return ""
	}
	return m.SiteURL + m.Path
}

// OGImage is the absolute URL of the link-preview image. Relative paths are
// ignored by every unfurler — Slack, LinkedIn, iMessage — so without SITE_URL
// there is nothing useful to emit.
func (m Meta) OGImage() string {
	if m.SiteURL == "" {
		return ""
	}
	return m.SiteURL + "/images/headshot.jpeg"
}

// IsCurrent reports whether a nav href is the page being rendered, so the
// navbar can set aria-current correctly. It was previously hardcoded on the
// Home link, which told every assistive technology that /about and /work were
// the home page.
func (m Meta) IsCurrent(href string) bool {
	return m.Path == href
}

// newMeta builds the per-page metadata, attaching the site-wide structured
// data block.
func (s *Server) newMeta(title, description, path string) Meta {
	return Meta{
		Title:       title,
		Description: description,
		Path:        path,
		SiteURL:     s.cfg.SiteURL,
		Dev:         s.cfg.Dev,
		HasResume:   s.hasResume,
		JSONLD:      s.personJSONLD(),
	}
}

// personJSONLD emits schema.org Person data. This is what lets a search engine
// connect the name to the profiles and the credentials rather than treating the
// page as anonymous prose — and alumniOf is where the training programs get
// stated as machine-readable fact instead of four logos.
func (s *Server) personJSONLD() template.JS {
	type org struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		URL  string `json:"url,omitempty"`
	}
	doc := map[string]any{
		"@context": "https://schema.org",
		"@type":    "Person",
		"name":     "Hunter Motko",
		"jobTitle": "Platform & Backend Engineer",
		"email":    "mailto:huntermotko.dev@gmail.com",
		"sameAs": []string{
			"https://github.com/hunterMotko",
			"https://linkedin.com/in/hunter-motko",
		},
		"knowsAbout": []string{
			"Go", "Python", "PostgreSQL", "SQLite", "Docker",
			"CI/CD", "Data pipelines", "Linux server administration",
		},
		"alumniOf": []org{
			{Type: "EducationalOrganization", Name: "Hack Reactor", URL: "https://www.hackreactor.com/"},
			{Type: "EducationalOrganization", Name: "Qwasar Silicon Valley", URL: "https://www.qwasar.io/"},
			{Type: "EducationalOrganization", Name: "The Last Mile", URL: "https://www.thelastmile.org/"},
			{Type: "EducationalOrganization", Name: "The Next Chapter", URL: "https://www.nextchapterbk.com/"},
		},
	}
	if s.cfg.SiteURL != "" {
		doc["url"] = s.cfg.SiteURL
	}

	b, err := json.Marshal(doc)
	if err != nil {
		// The document is a literal built above with no dynamic input, so this
		// cannot fail in practice. Log and emit nothing rather than panicking:
		// missing structured data must not take a page down.
		log.Printf("marshal person json-ld: %v", err)
		return ""
	}
	return template.JS(b)
}
