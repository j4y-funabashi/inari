#!/usr/bin/env bash

export GO111MODULE=on

build_func() {
    env GOOS=linux \
	go build -ldflags '-d -s -w' -a -tags netgo \
	-installsuffix netgo \
	-o "bin/${1}" "${1}/main.go"
}

build_func timeline
