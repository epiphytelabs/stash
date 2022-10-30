.PHONY: lint

all:

lint:
	make -C api lint
	make -C apps/smtpd lint
