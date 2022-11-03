.PHONY: all build init lint mail vendor

projects = api apps/messages apps/smtpd filters/email

all: build

build:
	$(call make-sub)

init:
	$(if $(wildcard .env),$(error .env already exists))
	echo "POSTGRES_PASSWORD=$(shell pwgen -Bs1 32)" >> .env

lint:
	$(call make-sub,lint)

mail:
	cat example/mail | sed -e "s|%%DATE%%|$(shell date '+%Y-%m-%d %H:%M:%S')|" | nc localhost 25

vendor:
	$(call make-sub,vendor)

make-sub = @$(foreach project,$(projects),echo $(project):; make -C $(project) $(1) | sed 's/^/  /'; )
