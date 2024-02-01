FROM golang:1.22-rc-bookworm AS builder
WORKDIR /app 
COPY go.mod go.sum ./ 
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /site

FROM builder AS test
RUN go test -v ./...

FROM ubuntu:18.04 AS release
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /site /site
ENTRYPOINT ["/site"]
