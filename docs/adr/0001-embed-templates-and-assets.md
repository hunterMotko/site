# Templates and static assets are embedded in the binary

Templates were parsed from the relative path `internal/app/views/*/*.html` and
static files served from `internal/app/public`, so the compiled binary only ran
when the working directory happened to be the repo root. It worked under `air`
and under a Dockerfile that copied the whole tree, and would have failed
anywhere else. Both are now embedded with `embed.FS`, making the binary the
entire deployment artifact — nothing beside it to forget to copy.

## Consequences

- Editing a template or stylesheet requires a rebuild; there is no reloading a
  file on disk. A stale process must also be killed, because it serves the old
  copy out of its own binary.
- The runtime container image needs only the binary and CA certificates.
