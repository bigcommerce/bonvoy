GOBIN := ${GOPATH}/bin
PATH := ${GOBIN}:${PATH}

export PATH
OBJECT=bonvoy

default: deps build

clean:
	@rm ${OBJECT}

build:
	@go build -v -o ${OBJECT}

deps:
	@go mod tidy
	@go get -u github.com/rakyll/gotest

test:
	@gotest -v $$(go list ./... | grep -v vendor/) -tags=integration

test-unit:
	@gotest $$(go list ./... | grep -v vendor/)

build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${OBJECT}-linux-amd64
