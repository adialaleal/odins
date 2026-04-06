.PHONY: build install tidy lint test clean

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

clean:
	rm -f $(BINARY)

run:
	go run . $(ARGS)
