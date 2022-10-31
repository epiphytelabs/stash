.PHONY: all build lint mail vendor

all: build

build:
	-make -C api
	-make -C apps/messages
	-make -C apps/smtpd

lint:
	-make -C api lint
	-make -C apps/messages lint
	-make -C apps/smtpd lint

mail:
	cat example/mail | sed -e "s|%%DATE%%|$(shell date '+%Y-%m-%d %H:%M:%S')|" | nc localhost 25

vendor:
	-make -C api vendor
	-make -C apps/messages vendor
	-make -C apps/smtpd vendor
