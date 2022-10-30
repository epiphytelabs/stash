## development #################################################################

FROM ddollar/go:1.19 AS development

ARG DOMAIN

ENV VITE_API_ENDPOINT=https://${DOMAIN}/graph
ENV VITE_CLIENT_PORT=443
ENV VITE_DOMAIN=${DOMAIN}

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /src
COPY package.json package-lock.json ./
RUN npm install

COPY . .

EXPOSE 3000

CMD ["bash", "-c", "npm rebuild && npx vite --host"]

## package ####################################################################

FROM ddollar/go:1.19 AS package

ARG DOMAIN

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /src

COPY --from=development /src /src

ENV VITE_API_ENDPOINT=https://${DOMAIN}/graph
ENV VITE_DOMAIN=${DOMAIN}

RUN npx vite build

## production ##################################################################

FROM nginx

COPY --from=package /src/nginx.conf /etc/nginx/nginx.conf
COPY --from=package /src/dist /usr/share/nginx/html
