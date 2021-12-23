all: generate fmt build vet test

fmt:
	goimports -w -local "github.com/wedaly" .

generate:
	go generate ./...

build:
	go build -o gospelunk main.go

install:
	go install

test:
	go test ./...

vet:
	go vet ./...

bench:
	go test ./... -bench=.

clean:
	rm -rf gospelunk
	go clean ./...
