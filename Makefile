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
check: staticcheck govulncheck
	$(STATICCHECK) ./...
	$(GOVULNCHECK) ./...

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

STATICCHECK ?= $(LOCALBIN)/staticcheck
GOVULNCHECK ?= $(LOCALBIN)/govulncheck

.PHONY: govulncheck
govulncheck: $(GOVULNCHECK) ## Download govulncheck locally if necessary.
$(GOVULNCHECK): $(LOCALBIN)
	test -s $(LOCALBIN)/govulncheck || GOBIN=$(LOCALBIN) go install golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: staticcheck
staticcheck: $(STATICCHECK) ## Download staticcheck locally if necessary.
$(STATICCHECK): $(LOCALBIN)
	test -s $(LOCALBIN)/staticcheck || GOBIN=$(LOCALBIN) go install honnef.co/go/tools/cmd/staticcheck@latest