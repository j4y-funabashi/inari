FROM golang:1.25-alpine AS builder

RUN apk update && \
	apk add \
	build-base

WORKDIR /api
COPY . .
RUN go mod download
RUN go build -o inari cmd/web/main.go

FROM alpine:3.23

RUN apk update && \
	apk add \
	bash

WORKDIR /
COPY --from=builder /api/inari /inari
EXPOSE 8080

CMD ["/inari"]
