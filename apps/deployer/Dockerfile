FROM node:alpine

RUN apk update && \
	apk add --no-cache \
	bash \
	python3 \
	py3-pip \
	make \
	go \
	&& pip3 install --upgrade pip \
	&& pip3 install \
	awscli \
	&& rm -rf /var/cache/apk/*
