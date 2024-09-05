BINARY_NAME=app
SOURCES=$(wildcard ./cmd/main.go)

build: $(SOURCES)
	go build -o $(BINARY_NAME) $(SOURCES)

## dev: run air dev 
.PHONY: dev
dev: 
	@air -c .air.toml 

## css: build tailwindcss
.PHONY: css
css:
	npx tailwindcss -i ./internal/app/public/css/index.css -o ./internal/app/public/css/out.css --minify

## css-watch: watch build tailwindcss
.PHONY: css-watch
css-watch:
	npx tailwindcss -i ./internal/app/public/css/index.css -o ./internal/app/public/css/out.css --watch

