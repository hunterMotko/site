# huntermotko.dev

Personal site. One Go module, server-rendered with Echo and `html/template`,
templates and static assets compiled in via `embed.FS`, optional SQLite for
request metrics. Deployed on a DigitalOcean droplet behind nginx.

Run `make ci` before pushing — it reproduces the entire GitHub Actions workflow
(`fmt-check`, `vet`, `test`, `govulncheck`) locally.

## Agent skills

### Issue tracker

Issues live as markdown files under `.scratch/<feature-slug>/` in this repo.
See `docs/agents/issue-tracker.md`.

### Domain docs

Single-context: `CONTEXT.md` at the repo root, ADRs in `docs/adr/`.
See `docs/agents/domain.md`.

## Repo conventions worth knowing

- **`local/` is gitignored** and holds private planning notes. Never read it for
  project decisions, write skill output into it, or link to it from a tracked
  file.
- **Templates and assets are embedded**, so a rebuild is required to see changes
  to `internal/app/views/` or `internal/app/public/` — and a stale process must
  be killed first, since it serves the old copy from its own binary.
- **`.env` is tracked and must never hold secrets.** Precedence is the real
  environment, then `.env.local` (gitignored), then `.env`.
