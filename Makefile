.PHONY: test
test:
	go test -race -cover ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
build:
	go build ./cmd/kiora

.PHONY: ci
ci: fmt lint test