.PHONY: start
start:
	go run cmd/apiserver/main.go
.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -timeout 30s ./...

.DEFAULT_GOAL := build