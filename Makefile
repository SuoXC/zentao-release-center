.PHONY: build run clean tidy

build:
	go build -o bin/release-center .

run:
	go run .

clean:
	rm -rf bin/

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

test:
	go test -v ./...
