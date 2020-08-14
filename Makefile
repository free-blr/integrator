CGO_ENABLED ?= 0
PREFIX ?= bin/
PWD := $(shell pwd)
DB_MIGRATIONS_PATH ?= $(PWD)/migrations/
GOLANG_CI_LINT := $(shell command -v golangci-lint 2> /dev/null)

.PHONY: all
all: lint test build

.PHONY: lint
lint:
ifndef GOLANG_CI_LINT
	wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b /usr/local/bin v1.18.0
endif
	CGO_ENABLED=0 golangci-lint --deadline=5m run ./...

.PHONY: build
build: build-bot # todo: build in docker, don't be loh

.PHONY: build-stats
build-bot:
	CGO_ENABLED=$(CGO_ENABLED) GOOS="$(GOOS)" go build -mod=vendor -a -installsuffix cgo -o "$(PREFIX)bot" cmd/bot/main.go

.PHONY: test
test:
	DB_MIGRATIONS_PATH=$(DB_MIGRATIONS_PATH) go test --count=1 -p=1 -mod=vendor -v $(shell go list ./... | grep -v internal/repository)
	DB_MIGRATIONS_PATH=$(DB_MIGRATIONS_PATH) go test -tags --count=1 -p=1 -mod=vendor -v $(shell go list ./... | grep internal/repository)

.PHONY: clean
clean:
	rm -rf bin
	rm -rf output