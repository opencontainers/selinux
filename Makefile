BUILDTAGS := selinux
RUNC_LINK := $(CURDIR)/Godeps/_workspace/src/github.com/opencontainers/selinux
export GOPATH := $(CURDIR)/Godeps/_workspace

$(RUNC_LINK):
	ln -sfn $(CURDIR) $(RUNC_LINK)

.PHONY: test
test: $(RUNC_LINK) | $(RUNC_LINK)
		go test -timeout 3m -tags "$(BUILDTAGS)" ${TESTFLAGS} -v ./...

clean:
	rm -f $(RUNC_LINK)
