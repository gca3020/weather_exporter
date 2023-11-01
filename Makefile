TAG?=latest
OWNER?=gca3020
SERVER?=docker.io
IMG_NAME?=weather_exporter
PLATFORMS?="linux/amd64,linux/arm/v7,linux/arm64"

MAIN_PACKAGE_PATH := ./cmd/weather_exporter
BINARY_NAME := weather_exporter

Version := $(shell git describe --always --tags --dirty)
GitCommit := $(shell git rev-parse --short HEAD)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	CGO_ENABLED=1 go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	CGO_ENABLED=1 go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	go build -ldflags "-X main.Release=${Version} -X main.SHA=${GitCommit}" -o build/weather_exporter cmd/weather_exporter/main.go

## build-local: build a local container for testing
.PHONY: build-local
build-local:
	@docker buildx create --use --name=multiarch --node multiarch && \
	docker buildx build \
		--progress=plain \
		--build-arg Version=$(Version) --build-arg GitCommit=$(GitCommit) \
		--platform linux/amd64 \
		--output "type=docker,push=false" \
		--tag $(SERVER)/$(OWNER)/$(IMG_NAME):$(Version) .
