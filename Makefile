.PHONY: test-backend lint-frontend lint-backend fmt-backend fmt-frontend lint fmt integration coverage ci build build-unchecked run run-cluster generate clean
test:
	mkdir -p artifacts/
	go test -short -race -cover -coverprofile=artifacts/cover.out ./...

lint-backend:
	golangci-lint run ./...
fmt-backend:
	go fmt ./...

lint-frontend:
	cd frontend && npm run lint
fmt-frontend:
	cd frontend && npm run prettier --write ./src

lint: lint-backend lint-frontend
fmt: fmt-backend fmt-frontend

integration:
	go test -count=1 ./integration

coverage: test
	go tool cover -html=artifacts/cover.out

ci: fmt lint test

build: build-frontend build-backend

build-backend:
	go build -o ./artifacts/tuku ./cmd/tuku
	go build -o ./artifacts/kiora ./cmd/kiora

build-frontend:
	cd frontend && npm run build
	rm -r ./internal/server/frontend/assets
	cp -r ./frontend/build ./internal/server/frontend/assets

run:
	./artifacts/kiora -c ./testdata/kiora.dot

run-cluster:
	./testdata/run-cluster.sh

generate:
	mockgen -source ./lib/kiora/kioradb/db.go > mocks/mock_kioradb/db.go
	mockgen -source ./lib/kiora/config/provider.go > mocks/mock_config/provider.go
	mockgen -source ./internal/clustering/broadcaster.go > mocks/mock_clustering/broadcaster.go
	mockgen -source ./internal/services/bus.go > mocks/mock_services/bus.go
	oapi-codegen -generate gorilla,spec,types -package apiv1 ./internal/server/api/apiv1/api.yaml > ./internal/server/api/apiv1/apiv1.gen.go
	cd frontend && npm exec openapi -- --useOptions -i ../internal/server/api/apiv1/api.yaml -o src/api

clean:
	rm -rf ./artifacts
	rm -rf ./frontend/build
	rm -rf ./internal/server/frontend/assets
