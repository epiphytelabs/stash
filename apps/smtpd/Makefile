.PHONY: all binaries compress lint test vendor

commands = smtpd

binaries = $(addprefix dist/, $(commands))

all: binaries

binaries: $(binaries)

compress: $(binaries)
	upx-ucl -1 $^

lint:
	go run ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run

test:
	env TEST=true go test -covermode atomic -coverprofile coverage.txt -mod=vendor ./...

vendor:
	go mod tidy
	go mod vendor
	go run vendor/github.com/goware/modvendor/main.go -copy="**/*.c **/*.h"

$(binaries): dist/%: $(sources)
	go build -o $@ -mod=vendor --ldflags="-s -w" ./cmd/$*
