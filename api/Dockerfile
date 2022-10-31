## development #################################################################

FROM ddollar/go:1.19 AS development

ENV MODE=development

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /src

EXPOSE 4000

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

COPY --from=package /src/dist/api /opt/bin/

EXPOSE 4000

CMD [ "/opt/bin/api" ]
