.PHONY: build install tidy lint test clean ai-packs ai-check

BINARY = odins
VERSION ?= dev

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY) .

install: build
	mv $(BINARY) /usr/local/bin/$(BINARY)

tidy:
	go mod tidy

lint:
	go vet ./...

test:
	go test -v -race ./...

ai-packs:
	go run ./tools/ai-pack-gen

ai-check:
	go run ./tools/ai-pack-gen --check

clean:
	rm -f $(BINARY)

run:
	go run . $(ARGS)
