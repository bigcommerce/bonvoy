GOBIN := ${GOPATH}/bin
GOCMD := go
GOBUILD := CGO_ENABLED=0 $(GOCMD) build
GOCLEAN := $(GOCMD) clean
PATH := ${GOBIN}:${PATH}

export PATH
BINARY_NAME=./bin/bonvoy

default: deps build

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux-amd64

build:
	mkdir -p bin
	$(GOBUILD) -v -o $(BINARY_NAME) .

deps:
	$(GOCMD) mod tidy
	$(GOCMD) install go.uber.org/mock/mockgen@latest
	$(GOCMD) install github.com/mfridman/tparse@latest

generate: deps
	$(GOCMD) generate -v ./...

lint:
	$(GOLANGCI_LINT) run

test:
	$(GOCMD) test -v -tags=integration $$(go list ./... | grep -v vendor/)

test-unit:
	$(GOCMD) test -v -coverprofile=c.out $$(go list ./... | grep -v vendor/)

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o $(BINARY_NAME)-linux-amd64
