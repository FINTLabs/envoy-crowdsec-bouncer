.PHONY: all
all: fmt lint vet check test build

.PHONY: test
test:
	go test ./...

.PHONY: start
start:
	go run main.go run

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: check
check:
	govulncheck ./...
	staticcheck ./...