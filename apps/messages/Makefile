.PHONY: all binaries compress lint vendor

binaries = dist/messages
sources  = $(shell find . -name '*.go')

all: binaries

binaries: $(binaries)

compress: $(binaries)
	upx-ucl -1 $^

lint:
	go run ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run

vendor:
	go mod tidy
	go mod vendor
	go run vendor/github.com/goware/modvendor/main.go -copy="**/*.c **/*.h"

dist/messages: $(sources)
	go build -o $@ -mod=vendor --ldflags="-s -w" .

