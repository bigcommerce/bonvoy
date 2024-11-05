GOBIN := ${GOPATH}/bin
GOCMD := go
GOBUILD := CGO_ENABLED=0 $(GOCMD) build
GOCLEAN := $(GOCMD) clean
PATH := ${GOBIN}:${PATH}

export PATH
OBJECT=bonvoy

default: deps build

clean:
	$(GOCLEAN)
	rm -f $(OBJECT)

build:
	$(GOBUILD) -v -o ${OBJECT}

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
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o ${OBJECT}-linux-amd64
