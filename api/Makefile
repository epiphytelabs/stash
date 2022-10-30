.PHONY: all binaries ca compress deploy lint mail test vendor

commands = api
image    = epiphytelabs/stash

binaries = $(addprefix dist/, $(commands))

all: binaries

binaries: $(binaries)

ca:
	$(if $(DOMAIN),,$(error DOMAIN is not set))
	mkdir -p db/certs/public
	mkdir -p db/certs/private
	$(eval tmp := $(shell mktemp -d))
	sed -e "s/%%DOMAIN%%/$(DOMAIN)/g" config/ca.cnf > $(tmp)/ca.cnf
	openssl req -new -nodes -x509 -days 3650 -subj "/CN=ca.stash" -extensions v3_ca \
	  -config $(tmp)/ca.cnf -keyout db/certs/private/ca.pem -out db/certs/public/ca.pem
	rm -rf $(tmp)

compress: $(binaries)
	upx-ucl -1 $^

deploy:
	git diff-index --quiet HEAD
	git diff --quiet @{upstream}
	ssh ubuntu@stash.dollar.io "cd stash && git pull && docker compose up -d"

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
