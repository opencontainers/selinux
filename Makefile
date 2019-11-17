export GO111MODULE=off
GO ?= go
BUILDTAGS := selinux

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
