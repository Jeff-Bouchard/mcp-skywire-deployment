.DEFAULT_GOAL := help

PROJECT_ROOT:=github.com/SkycoinPro/skycoin-service-discovery

SHELL := /usr/bin/env bash

VERSION := $(shell git describe --always)

RFC_3339 := "+%Y-%m-%dT%H:%M:%SZ"
DATE := $(shell date -u $(RFC_3339))
COMMIT := $(shell git rev-list -1 HEAD)

GOLANGCI_LINT_VERSION:=latest

OPTS?=GO111MODULE=on
TEST_OPTS:=TEST_SERVICEDISC_REDIS_ADDR=redis://127.0.0.1:6379
TEST_ARGS:=-tags no_ci -cover -timeout=5m

BIN := ${PWD}/bin

SKYWIRE_UTILITIES_REPO := github.com/skycoin/skywire-utilities
BUILDINFO_PATH := $(SKYWIRE_UTILITIES_REPO)/pkg/buildinfo

BUILDINFO_VERSION := -X $(BUILDINFO_PATH).version=$(VERSION)
BUILDINFO_DATE := -X $(BUILDINFO_PATH).date=$(DATE)
BUILDINFO_COMMIT := -X $(BUILDINFO_PATH).commit=$(COMMIT)

BUILDINFO?=$(BUILDINFO_VERSION) $(BUILDINFO_DATE) $(BUILDINFO_COMMIT)
BUILD_OPTS?="-ldflags=$(BUILDINFO)"
BUILD_OPTS_DEPLOY?="-ldflags=$(BUILDINFO) -w -s"

install-tools: ## Install linters
	${OPTS} go install golang.org/x/tools/cmd/goimports@latest
	${OPTS} go get -u github.com/incu6us/goimports-reviser/v2

generate:
	go generate ./...

dep: ## Tidies and vendors dependencies.
	$(OPTS) go mod tidy -v
	$(OPTS) go mod vendor -v

format: dep ## Formats the code. Must have goimports and goimports-reviser installed (use make install-tools).
	goimports -w -local $(PROJECT_ROOT) ./pkg
	goimports -w -local $(PROJECT_ROOT) ./cmd
	find . -type f -name '*.go' -not -path "./vendor/*"  -exec goimports-reviser -project-name ${PROJECT_ROOT} -file-path {} \;

lint: format ## Formats and lints golang code.
	${OPTS} golangci-lint run -c .golangci.yml ./...
	$(OPTS) go vet -all ./...

test: ## Run tests (TODO: Have better redis setup logic. E.g. custom redis config).
	docker run --name integration-test-db -e POSTGRES_USERNAME=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DATABASE=postgres -p 8383:5432 -d postgres 
	docker run -d -p 6379:6379 --name integration-test-redis redis
	sleep 1
	$(TEST_OPTS) go clean -testcache
	$(TEST_OPTS) go test $(TEST_ARGS) ./pkg/...
	docker stop integration-test-redis && docker rm integration-test-redis
	docker stop integration-test-db && docker rm integration-test-db

check: lint test ## Run lint and check

build: ## Builds service-discovery into ./bin
	mkdir -p ${BIN}; go build ${BUILD_OPTS} -o ${BIN} ./cmd/*

build-deploy: ## Build for deployment Docker image
	${OPTS} go build ${BUILD_OPTS_DEPLOY} -o /release/service-discovery ./cmd/service-discovery

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$|^##.*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
