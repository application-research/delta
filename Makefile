SHELL=/usr/bin/env bash
GO_BUILD_IMAGE?=golang:1.19

.PHONY: all
all: build

.PHONY: build
build:
	git submodule update --init --recursive
	make -C extern/filecoin-ffi
	go build -tags netgo -ldflags '-s -w' -o delta

.PHONE: clean
clean:
	rm -f delta
	git submodule deinit --all -f

install:
	install -C -m 0755 delta /usr/local/bin

.PHONY: generate-swagger
generate-swagger:
	scripts/swagger/swag.sh
