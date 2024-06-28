GOBIN    	:= $(GOPATH)/bin
PATH     	:= $(GOBIN):$(PATH)
GOLANGCI_VERSION=v1.59.1
CGO_ENABLED=0

all: tools lint test

test:
	go test -v -cover -race ./...

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOBIN} ${GOLANGCI_VERSION}

lint:
	$(info Running Go code checkers and linters)
	golangci-lint --version
	golangci-lint run

build:
	go build ./cmd/wait-for


.PHONY: all build test
