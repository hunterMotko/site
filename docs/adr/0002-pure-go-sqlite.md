# Pure-Go SQLite over mattn/go-sqlite3

Request metrics are stored in SQLite. `mattn/go-sqlite3` is the better-known
driver and is faster, but it requires cgo — which produces a dynamically linked
binary that must match the runtime image's libc, forcing a fat runtime image and
defeating the single-static-binary deployment in
[ADR-0001](0001-embed-templates-and-assets.md). `modernc.org/sqlite` is a
pure-Go translation, so `CGO_ENABLED=0` holds.

At this site's traffic the speed difference is irrelevant and the simpler build
is not, so the trade goes to the pure-Go driver.

## Consequences

- The driver name is `sqlite`, not `sqlite3`. This is the usual thing people get
  wrong when switching.
- Do not "upgrade" to `mattn/go-sqlite3` for performance. It would reintroduce
  cgo and break the static build.
