# huntermotko.dev

My personal site — [huntermotko.dev](https://huntermotko.dev).

Server-rendered Go. Echo and `html/template`, three direct dependencies, no
JavaScript toolchain. Templates and static assets are compiled into the binary
with `embed.FS`, so what ships is one static file that runs from anywhere.

Page content lives in Go rather than a CMS or a pile of Markdown — see
`cmd/api/`. It is a small enough site that a struct is the simplest thing that
works, and it means a typo is a compile error.

## Running it

```sh
make dev                # live reload (air)
make build && ./app     # or just build it
make ci                 # fmt, vet, race tests, govulncheck — the whole CI workflow
```

Configuration comes from the environment. `.env` is tracked and documents every
key with empty values; put real ones in `.env.local`, which is gitignored. The
site runs with none of them set.

To run the container instead:

```sh
docker compose up -d --build
```

## Layout

```
cmd/api/            page content, as Go data
internal/app/       templates and static assets (embedded)
internal/server/    routes, config, rendering, SEO
internal/database/  optional SQLite request counts
```

## Docs

- [`docs/deploy.md`](docs/deploy.md) — deploying, and what fails quietly if it is
  misconfigured
- [`docs/adr/`](docs/adr/) — the decisions worth explaining, and what they cost
- [`CONTEXT.md`](CONTEXT.md) — the vocabulary the site and the code share
