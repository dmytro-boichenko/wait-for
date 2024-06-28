GOBIN    	:= $(GOPATH)/bin
PATH     	:= $(GOBIN):$(PATH)
GOLANGCI_VERSION=v1.55.2
CGO_ENABLED=0

all: tools lint test

PACKAGES := go list ./...

## run the unit-tests
test:
	go test -v -cover -race ./...

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOBIN} ${GOLANGCI_VERSION}

lint:
	$(info Running Go code checkers and linters)
	golangci-lint --version
	golangci-lint $(V) run --timeout 5m

build:
	go build ./cmd/wait-for


.PHONY: all build test
