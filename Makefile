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

test:
	go test -v $$(go list ./... | grep -v vendor/) -tags=integration

test-unit:
	go test $$(go list ./... | grep -v vendor/)

