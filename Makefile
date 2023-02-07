SHELL=/usr/bin/env bash
GO_BUILD_IMAGE?=golang:1.19

.PHONY: all
all: build

.PHONY: build
build:
	go build -tags netgo -ldflags '-s -w' -o stg-dealer

.PHONE: clean
clean:
	rm -f stg-dealer