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
ci: generate fmt lint test

.PHONY: build
build: generate ci build-unchecked

.PHONY: build-unchecked
build-unchecked:
	go build -o ./artifacts/kiora ./cmd/kiora

.PHONY: run
run: build
	./artifacts/kiora -c ./testdata/kiora.dot --raft.data-dir artifacts/kiora-raft-data

.PHONY: run-cluster
run-cluster:
	./testdata/run-cluster.sh

PROTO_TARGETS = $(wildcard internal/dto/kioraproto/schema/*.proto)
PROTO_OUTPUTS = $(patsubst internal/dto/kioraproto/schema/%.proto,internal/dto/kioraproto/%.pb.go,$(PROTO_TARGETS))
$(PROTO_OUTPUTS): $(PROTO_TARGETS)
	cd internal/dto/kioraproto/schema && protoc --go_opt=paths=source_relative --go_out=../ $(patsubst internal/dto/kioraproto/schema/%,%,$^)

.PHONY: generate
generate: $(PROTO_OUTPUTS)

.PHONY: generate-clean
generate-clean:
	rm $(PROTO_OUTPUTS)