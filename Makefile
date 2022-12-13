.PHONY: test
test:
	mkdir -p artifacts/
	go test -race -cover -coverprofile=artifacts/cover.out ./...

.PHONY: coverage
coverage: test
	go tool cover -html=artifacts/cover.out

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: ci
ci: fmt lint test

.PHONY: build
build: generate ci build-unchecked

.PHONY: build-unchecked
build-unchecked:
	go build -o ./artifacts/kiora ./cmd/kiora

PROTO_TARGETS = $(wildcard internal/dto/kioraproto/schema/*.capnp)
PROTO_OUTPUTS = $(patsubst internal/dto/kioraproto/schema/%.capnp,internal/dto/kioraproto/%.capnp.go,$(PROTO_TARGETS))
$(PROTO_OUTPUTS): $(PROTO_TARGETS)
	test -n "$(GOPATH)" || (echo 'missing $$GOPATH' ; exit 1)
	capnp compile -I$(GOPATH)/src/capnproto.org/go/capnp/std -ogo:internal/dto/kioraproto --src-prefix internal/dto/kioraproto/schema $^

.PHONY: generate
generate: $(PROTO_OUTPUTS)
