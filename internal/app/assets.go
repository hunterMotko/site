// Package app carries the site's templates and static assets, embedded into
// the binary at compile time.
//
// Before this, templates were parsed from the relative path
// "internal/app/views/*/*.html" and static files served from
// "internal/app/public" — so the compiled binary only ran when the working
// directory happened to be the repo root. It worked under `air` and under a
// Dockerfile that COPYs the whole tree, and would have failed anywhere else.
// Embedding removes the working-directory dependency entirely: the binary is
// the deployment artifact, with nothing beside it to forget to copy.
package app

import "embed"

// Views holds the html/template sources, parsed by the server at startup.
// Paths inside are rooted at "views/", not at the repo root.
//
//go:embed views
var Views embed.FS

// Public holds everything served over HTTP as-is: CSS, JS, images. Callers
// want it rooted at "/", so serve it through fs.Sub(Public, "public") rather
// than mounting it directly — otherwise every asset URL carries a /public
// prefix that the templates don't use.
//
//go:embed public
var Public embed.FS
