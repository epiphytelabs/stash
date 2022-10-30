## development #################################################################

FROM ddollar/go:1.19 AS development

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /src

EXPOSE 2525

CMD [ "go", "run", "./vendor/github.com/cespare/reflex", "-c", ".reflex" ]

## package #####################################################################

FROM ddollar/go:1.19 AS package

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /src

COPY . .

RUN make binaries compress

## production ##################################################################

FROM ddollar/ubuntu:lts AS production

ENV GOPATH=/go
ENV PATH=$PATH:/opt/bin

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /

COPY --from=package /src/dist/smtpd /opt/bin/

EXPOSE 2525

CMD [ "/opt/bin/smtpd" ]
