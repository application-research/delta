SHELL=/usr/bin/env bash
GO_BUILD_IMAGE?=golang:1.19
VERSION=$(shell git describe --always --tag --dirty)
DOCKER_COMPOSE_FILE=docker-compose.yml

.PHONY: all
all: build

.PHONY: build
build:
	git submodule update --init --recursive
	make -C extern/filecoin-ffi
	go generate
	go build -tags netgo -ldflags="-s -w -X main.Commit=$(shell git rev-parse HEAD) -X main.Version=$(VERSION)" -o delta

.PHONY: clean
clean:
	rm -f delta
	git submodule deinit --all -f

install:
	install -C -m 0755 delta /usr/local/bin

.PHONY: docker-compose-build
docker-compose-build:
	BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
	COMMIT=$(shell git rev-parse HEAD) \
	VERSION=$(shell git describe --always --tag --dirty) \
	WALLET_DIR=$(WALLET_DIR) \
	DESCRIPTION=$(DESCRIPTION)
	docker-compose -f $(DOCKER_COMPOSE_FILE) build --build-arg WALLET_DIR=$(WALLET_DIR)

.PHONY: docker-compose-build
docker-compose-build:
	BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
	COMMIT=$(shell git rev-parse HEAD) \
	VERSION=$(shell git describe --always --tag --dirty) \
	WALLET_DIR=$(WALLET_DIR) \
	DESCRIPTION=$(DESCRIPTION) \
	TAG=$(TAG) \
	docker-compose -f $(DOCKER_COMPOSE_FILE) build --build-arg WALLET_DIR=$(WALLET_DIR) --build-arg REPO=$(REPO)

.PHONY: docker-compose-up
docker-compose-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up

.PHONY: docker-compose-run
docker-compose-run: docker-compose-build docker-compose-up

.PHONY: docker-compose-down
docker-compose-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: prepare-spec docker-push
docker-push:
	docker build -t delta:$(VERSION) .
	docker tag delta:$(VERSION) 0utercore/delta:$(VERSION)
	docker push 0utercore/delta:$(VERSION)
