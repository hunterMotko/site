BINARY_NAME=app
SOURCES=$(wildcard ./cmd/main.go)

build: $(SOURCES)
	go build -o $(BINARY_NAME) $(SOURCES)

## dev: run air dev 
.PHONY: dev
dev: 
	@air -c .air.toml 


