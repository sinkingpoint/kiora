.PHONY: test
test:
	mkdir -p artifacts/
	go test -short -race -cover -coverprofile=artifacts/cover.out ./...

.PHONY: integration
integration:
	go test -timeout=1m -count=1 ./integration

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
	./artifacts/kiora -c ./testdata/kiora.dot

.PHONY: run-cluster
run-cluster:
	./testdata/run-cluster.sh

.PHONY: generate
generate:
	mockgen -source ./lib/kiora/kioradb/db.go > mocks/mock_kioradb/db.go
	mockgen -source ./lib/kiora/config/provider.go > mocks/mock_config/provider.go
	mockgen -source ./internal/clustering/broadcaster.go > mocks/mock_clustering/broadcaster.go
	mockgen -source ./internal/services/bus.go > mocks/mock_services/bus.go

.PHONY: generate-clean
generate-clean:
	rm $(PROTO_OUTPUTS)