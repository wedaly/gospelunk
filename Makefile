.PHONY: all fmt build test install vet bench clean

all: fmt build vet test

fmt:
	gofmt -s -w .
	goimports -w -local "github.com/wedaly/gospelunk" .

build:
	go build -o gospelunk github.com/wedaly/gospelunk

test:
	go test ./...

install:
	go install

vet:
	go vet ./...

clean:
	rm -rf gospelunk
	go clean ./...
