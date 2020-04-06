export GO111MODULE=off
GO ?= go
BUILDTAGS := selinux

all: build build-cross

define go-build
	GOOS=$(1) GOARCH=$(2) $(GO) build -tags $(BUILDTAGS) ./...
endef

define go-build-noselinux
	GOOS=$(1) GOARCH=$(2) $(GO) build ./...
endef

.PHONY:
build:
	$(call go-build,linux,amd64)

.PHONY:
build-cross:
	$(call go-build,linux,386)
	$(call go-build,linux,arm)
	$(call go-build,linux,arm64)
	$(call go-build,linux,ppc64le)
	$(call go-build,linux,s390x)
	$(call go-build,windows,amd64)
	$(call go-build,windows,386)
	$(call go-build-noselinux,linux,amd64)
	$(call go-build-noselinux,linux,arm)
	$(call go-build-noselinux,linux,arm64)
	$(call go-build-noselinux,linux,ppc64le)
	$(call go-build-noselinux,linux,s390x)
	$(call go-build-noselinux,windows,amd64)
	$(call go-build-noselinux,windows,386)

BUILD_PATH := $(shell pwd)/build
BUILD_BIN_PATH := ${BUILD_PATH}/bin
GOLANGCI_LINT := ${BUILD_BIN_PATH}/golangci-lint

check-gopath:
ifndef GOPATH
	$(error GOPATH is not set)
endif

.PHONY: test
test: check-gopath
	bash -c "diff  <(grep '^func [A-Z]' go-selinux/selinux_stub.go) <(grep '^func [A-Z]' go-selinux/selinux_linux.go)"
	go test -timeout 3m -tags "${BUILDTAGS}" ${TESTFLAGS} -v ./...
	go test -timeout 3m ${TESTFLAGS} -v ./...

${GOLANGCI_LINT}:
	export \
		VERSION=v1.23.7 \
		URL=https://raw.githubusercontent.com/golangci/golangci-lint \
		BINDIR=${BUILD_BIN_PATH} && \
	curl -sfL $$URL/$$VERSION/install.sh | sh -s $$VERSION

.PHONY:
lint: ${GOLANGCI_LINT}
	${GOLANGCI_LINT} version
	${GOLANGCI_LINT} linters
	${GOLANGCI_LINT} run

vendor:
	export GO111MODULE=on \
		$(GO) mod tidy && \
		$(GO) mod vendor && \
		$(GO) mod verify
