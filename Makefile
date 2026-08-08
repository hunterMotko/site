BINARY_NAME=app
MAIN=./cmd/main.go

.PHONY: build dev test vet fmt fmt-check vuln ci print-go-version clean

## build: compile the binary
build:
	go build -o $(BINARY_NAME) $(MAIN)

## dev: run with live reload
dev:
	@air -c .air.toml

## test: run the suite with the race detector
test:
	go test -race ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

## fmt-check: fail rather than rewrite, so CI reports formatting instead of
## silently fixing it on a branch nobody looks at again
fmt-check:
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then \
		echo "gofmt needed:"; echo "$$out"; exit 1; \
	fi

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## ci: everything the GitHub workflow runs. The workflow's job bodies are these
## targets and nothing else — when CI and the laptop run different commands they
## drift, and "passes locally, red in CI" becomes routine.
ci: fmt-check vet test vuln

## print-go-version: the Go minor the Dockerfile builds with. CI reads this
## instead of go.mod, because setup-go treats the go directive as an exact pin
## rather than the minimum it is.
print-go-version:
	@grep -m1 '^FROM golang:' Dockerfile | sed 's/.*golang:\([0-9.]*\).*/\1/'

## deploy: run ON THE DROPLET, from the repo checkout. Recorded here rather than
## remembered — see docs/adr/0003-droplet-pulls-and-builds.md for why the box
## builds instead of pulling a published image.
##
## Requires, on the droplet and not in git:
##   .env.production        config; same keys as .env
##   /srv/site/resume.pdf   bind-mounted read-only by compose
deploy:
	@test -f .env.production || { \
		echo "missing .env.production — see .env for the keys"; exit 1; }
	git pull
	docker compose up -d --build
	docker compose ps

## deploy-check: confirm the running site is actually healthy after a deploy.
## `docker compose ps` reports the container is up, which is not the same thing.
deploy-check:
	@curl -fsS -o /dev/null -w 'GET  / -> %{http_code}\n' http://127.0.0.1:8080/
	@curl -fsS -o /dev/null -w 'HEAD / -> %{http_code}\n' -I http://127.0.0.1:8080/
	@curl -fsS http://127.0.0.1:8080/robots.txt | grep -q 'Sitemap: https://' \
		&& echo 'robots.txt carries an absolute sitemap URL (SITE_URL is set)' \
		|| { echo 'SITE_URL is NOT set — canonical and og: tags are missing'; exit 1; }

clean:
	rm -f $(BINARY_NAME)
	rm -rf tmp/
