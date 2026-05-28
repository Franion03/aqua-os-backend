.PHONY: build run test clean

build:
	go build -o bin/aquaos-backend ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

clean:
	rm -rf bin/ aquaos.db

fmt:
	go fmt ./...

tidy:
	go mod tidy
