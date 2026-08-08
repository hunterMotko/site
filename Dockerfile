# The builder image is the source of truth for the Go version. go.mod's `go`
# directive is a minimum, not a pin, and CI installs whatever this line names —
# so the toolchain that builds here is the toolchain CI tests with.
FROM golang:1.25-alpine AS build

WORKDIR /src

# Dependencies are copied and downloaded before the source, so editing a handler
# does not invalidate the module cache layer and re-download the world.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 is what makes the output a static binary. It works because the
# SQLite driver is modernc.org/sqlite, a pure-Go implementation — with
# mattn/go-sqlite3 this would need cgo and a libc-matched runtime image.
#
# -trimpath keeps local filesystem paths out of the binary; -s -w drop the
# symbol table and DWARF data, which are of no use in production.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/app ./cmd/main.go

# Templates and static assets are compiled into the binary via embed.FS, so the
# runtime image needs the binary and CA certificates and nothing else. The
# previous version shipped the entire golang image — toolchain, module cache,
# and source — as the runtime, around a gigabyte to serve three pages.
FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 10001 app

# /data is the mount point for the metrics volume. Creating it here, owned by
# the unprivileged user, is what makes the volume writable: Docker initialises a
# fresh named volume from the image's directory at that path, ownership
# included. Without this the volume mounts root:root, the app cannot open the
# database, and metrics silently degrade to off — the failure looks like a
# configuration mistake rather than a permissions one.
RUN mkdir -p /data && chown app:app /data

COPY --from=build /out/app /app/app

# Runs unprivileged. Nothing here needs root, and a container process that does
# not need root should not have it — note this rules out binding a port below
# 1024, which is correct anyway with nginx terminating TLS in front.
USER app
WORKDIR /app

EXPOSE 8080
ENTRYPOINT ["/app/app"]
