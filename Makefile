.PHONY: test-backend
test:
	mkdir -p artifacts/
	go test -short -race -cover -coverprofile=artifacts/cover.out ./...

.PHONY: lint-backend
lint:
	golangci-lint run ./...

.PHONY: fmt-backend
fmt:
	go fmt ./...

.PHONY: lint-frontend
lint-frontend:
	cd frontend && npm run lint

.PHONY: fmt-frontend
fmt-frontend:
	cd frontend && npm run prettier --write ./src

.PHONY: lint
lint: lint-backend lint-frontend

.PHONY: fmt
fmt: fmt-backend fmt-frontend

.PHONY: integration
integration:
	go test -timeout=1m -count=1 ./integration

.PHONY: coverage
coverage: test
	go tool cover -html=artifacts/cover.out

.PHONY: ci
ci: generate fmt lint test

.PHONY: build
build: generate ci build-unchecked

.PHONY: build-unchecked
build-unchecked:
	cd frontend && npm run build
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

.PHONY: clean
clean:
	rm -rf ./artifacts
	rm -rf ./frontend/build
