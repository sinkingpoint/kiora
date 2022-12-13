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

PROTO_TARGETS = $(wildcard internal/dto/kioraproto/schema/*.capnp)
PROTO_OUTPUTS = $(patsubst internal/dto/kioraproto/schema/%.capnp,internal/dto/kioraproto/%.capnp.go,$(PROTO_TARGETS))
$(PROTO_OUTPUTS): $(PROTO_TARGETS)
	capnp compile -I$(GOPATH)/src/capnproto.org/go/capnp/std -ogo:internal/dto/kioraproto --src-prefix internal/dto/kioraproto/schema $^

.PHONY: generate
generate: $(PROTO_OUTPUTS)
