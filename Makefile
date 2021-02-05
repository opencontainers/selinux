GO ?= go

all: build build-cross

define go-build
	GOOS=$(1) GOARCH=$(2) $(GO) build ./...
endef

.PHONY: build
build:
	$(call go-build,linux,amd64)

.PHONY: build-cross
build-cross:
	$(call go-build,linux,386)
	$(call go-build,linux,arm)
	$(call go-build,linux,arm64)
	$(call go-build,linux,ppc64le)
	$(call go-build,linux,s390x)
	$(call go-build,linux,mips64le)
	$(call go-build,windows,amd64)
	$(call go-build,windows,386)

BUILD_PATH := $(shell pwd)/build
BUILD_BIN_PATH := ${BUILD_PATH}/bin
GOLANGCI_LINT := ${BUILD_BIN_PATH}/golangci-lint
GOLANGCI_INSTALL := https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
GOLANGCI_VERSION := v1.31.0

.PHONY: check-gopath
check-gopath:
ifndef GOPATH
	$(error GOPATH is not set)
endif

.PHONY: test
test: check-gopath
	go test -timeout 3m ${TESTFLAGS} -v ./...

${GOLANGCI_LINT}:
	curl -sSfL ${GOLANGCI_INSTALL} | sh -s -- -b ${BUILD_BIN_PATH} ${GOLANGCI_VERSION}

.PHONY: lint
lint: ${GOLANGCI_LINT}
	${GOLANGCI_LINT} version
	${GOLANGCI_LINT} linters
	${GOLANGCI_LINT} run

.PHONY: vendor
vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
