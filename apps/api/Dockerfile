FROM golang:alpine as buildAPI

RUN apk update && \
	apk add \
	build-base

WORKDIR /api
COPY apps/api .
RUN go mod download
RUN go build -o inari cmd/web/main.go

FROM nginx:alpine

RUN apk update && \
	apk add \
	bash

WORKDIR /
COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=buildAPI /api/inari /inari
COPY apps/api/start.sh /docker-entrypoint.d/start-inari-web.sh
