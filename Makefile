ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/_output

GOBIN ?= $$(go env GOPATH)/bin

z := $(shell mkdir -p $(OUTPUT_DIR))

.PHONY: all
all: fmt lint vet check test build

#----------------------------------------------------------------------------------
# Repo setup
#----------------------------------------------------------------------------------
GOLANGCI_LINT ?= go tool golangci-lint
GOLANGCI_LINT_FMT ?= $(GOLANGCI_LINT) fmt

.PHONY: fmt
fmt:
	$(GOLANGCI_LINT_FMT) ./...

.PHONY: fmt-changed
fmt-changed: ## Format only the changed code with golangci-lint (skip deleted files)
	git status -s -uno | awk '{print $$2}' | grep '.*.go$$' | xargs -r -I{} bash -lc '[ -f "{}" ] && $(GOLANGCI_LINT_FMT) "{}" || true'

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: mod-download
mod-download:
	go mod download all

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)

#----------------------------------------------------------------------------------
# Analyze
#----------------------------------------------------------------------------------

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run ./...

.PHONY: check
check:
	go tool govulncheck ./...
	go tool staticcheck ./...

#----------------------------------------------------------------------------------
# Test & Coverage
#----------------------------------------------------------------------------------
GO_TEST_ARGS ?=
GO_TEST_REPORT_ARGS ?= --jsonfile $(OUTPUT_DIR)/test-report.json
TEST_ARGS ?= -race -outputdir=$(OUTPUT_DIR)
TEST_COVERAGE_ARGS ?= --cover --covermode=atomic --coverprofile=cover.out

.PHONY: test
test:
	go tool gotestsum --format=standard-verbose --rerun-fails-abort-on-data-race --packages="./..." $(GO_TEST_ARGS) -- $(TEST_ARGS)

.PHONY: test-with-coverage
test-with-coverage: TEST_ARGS += $(TEST_COVERAGE_ARGS)
test-with-coverage: test

.PHONY: test-with-report
test-with-report: GO_TEST_ARGS += $(GO_TEST_REPORT_ARGS)
test-with-report: test

.PHONY: test-with-coverage-and-report
test-with-coverage-and-report: TEST_ARGS += $(TEST_COVERAGE_ARGS)
test-with-coverage-and-report: GO_TEST_ARGS += $(GO_TEST_REPORT_ARGS)
test-with-coverage-and-report: test

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: validate-test-coverage
validate-test-coverage: install-go-test-coverage 
	${GOBIN}/go-test-coverage --config=./.testcoverage.yaml

.PHONY: view-test-coverage
view-test-coverage:
	go tool cover -html $(OUTPUT_DIR)/cover.out

#----------------------------------------------------------------------------------
# Run Application
#----------------------------------------------------------------------------------

.PHONY: start
start:
	go run main.go run