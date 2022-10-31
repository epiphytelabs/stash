.PHONY: lint mail

all:

lint:
	-make -C api lint
	-make -C apps/messages lint
	-make -C apps/smtpd lint

mail:
	cat example/mail | sed -e "s|%%DATE%%|$(shell date '+%Y-%m-%d %H:%M:%S')|" | nc localhost 25
