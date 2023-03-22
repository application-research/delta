SHELL=/usr/bin/env bash
GO_BUILD_IMAGE?=golang:1.19
VERSION=$(shell git describe --always --tag --dirty)

.PHONY: all
all: build

.PHONY: build
build:
	git submodule update --init --recursive
	make -C extern/filecoin-ffi
	go generate
	go build -tags netgo -ldflags="-s -w -X main.Commit=$(shell git rev-parse HEAD) -X main.Version=$(VERSION)" -o delta

.PHONE: clean
clean:
	rm -f delta
	git submodule deinit --all -f

install:
	install -C -m 0755 delta /usr/local/bin

