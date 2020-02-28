export GO111MODULE=off
GO ?= go
BUILDTAGS := selinux

all: build build-cross

define go-build
	GOOS=$(1) GOARCH=$(2) $(GO) build -tags $(BUILDTAGS) ./...
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

check-gopath:
ifndef GOPATH
	$(error GOPATH is not set)
endif

.PHONY: test
test: check-gopath
	go test -timeout 3m -tags "${BUILDTAGS}" ${TESTFLAGS} -v ./...
	go test -timeout 3m ${TESTFLAGS} -v ./...

.PHONY:
lint:
	@out="$$(golint go-selinux)"; \
	if [ -n "$$out" ]; then \
		echo "$$out"; \
		exit 1; \
	fi

vendor:
	export GO111MODULE=on \
		$(GO) mod tidy && \
		$(GO) mod vendor && \
		$(GO) mod verify
