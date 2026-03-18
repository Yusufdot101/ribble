.PHONY:build
build:
	@go build -o tmp/main ./cmd/api

.PHONY:run
run:
	@./tmp/main

.PHONY:run-dev
run-dev:
	@air

.PHONY:test
test:
	@go test -v ./...

.PHONY:test-coverage
test-coverage:
	@go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
