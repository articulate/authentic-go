GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

test:
	ginkgo

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

deps:
	curl -s https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	${GOPATH}/bin/dep ensure

.PHONY: fmt fmtcheck
