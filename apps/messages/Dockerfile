## development #################################################################

FROM ddollar/go:1.19 AS development

ENV MODE=development
ENV VITE_CLIENT_PORT=443

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /src
COPY package.json package-lock.json ./
RUN npm install

COPY . .

EXPOSE 4000

CMD [ "go", "run", "./vendor/github.com/cespare/reflex", "-c", ".reflex" ]

## package ####################################################################

FROM ddollar/go:1.19 AS package

ARG DOMAIN

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /src

COPY --from=development /src /src

ENV VITE_API_ENDPOINT=https://${DOMAIN}/graph
ENV VITE_DOMAIN=${DOMAIN}

RUN npx vite build
RUN make binaries compress

## production ##################################################################

FROM ddollar/ubuntu:lts AS production

ENV GOPATH=/go
ENV PATH=$PATH:/opt/bin

WORKDIR /

COPY --from=package /src/dist/messages /opt/bin/

EXPOSE 4000

CMD ["/opt/bin/messages"]
