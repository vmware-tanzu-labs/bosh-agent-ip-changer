NAME ?= opsmanutils
OUTPUT = ./bin/$(NAME)
GO_SOURCES = $(shell find . -type f -name '*.go')
GOBIN ?= $(shell go env GOPATH)/bin
VERSION ?= 0.0.0
GITSHA = $(shell git rev-parse HEAD)
GITDIRTY = $(shell git diff --quiet HEAD || echo "dirty")
GOLANGCI_LINT_VERSION := $(shell golangci-lint --version 2>/dev/null)
LDFLAGS_VERSION = -X github.com/vmware-tanzu-labs/opsman-utils/command.VERSION=$(VERSION) \
				  -X github.com/vmware-tanzu-labs/opsman-utils/command.COMMIT=$(GITSHA)

.PHONY: all
all: build lint test ## Runs build, lint and test

.PHONY: clean
clean: ## Clean testcache and delete build output
	go clean -testcache
	@rm -rf bin/
	@rm -rf dist/
	@rm -f coverage.html

$(OUTPUT): $(GO_SOURCES)
	@echo "Building $(VERSION)"
	go build -o $(OUTPUT) -ldflags "$(LDFLAGS_VERSION)" ./cmd

.PHONY: build
build: $(OUTPUT) ## Build the main binary

.PHONY: test
test: ## Run the unit tests
	go test -short ./...

test-bench: ## Run all the benchmark tests
	go test -bench=. -benchmem ./...

test-all: test test-bench ## Run all of the tests

.PHONY: install
install: build ## Copy build to GOPATH/bin
	@cp $(OUTPUT) $(GOBIN)
	@echo "[OK] CLI binary installed under $(GOBIN)"

.PHONY: coverage
coverage: ## Run the tests with coverage and race detection
	go test -v --race -coverprofile=c.out -covermode=atomic ./...

.PHONY: report
report: ## Show coverage in an html report
	go tool cover -html=c.out -o coverage.html

.PHONY: generate
generate: ## Generate fakes
	go generate ./...

.PHONY: release
release: $(GO_SOURCES) ## Cross-compile binary for various operating systems
	@rm -rf dist
	@mkdir -p dist
	GOOS=darwin   GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-w $(LDFLAGS_VERSION)" -o $(OUTPUT)     ./cmd/$(NAME) && tar -czf dist/$(NAME)-darwin-amd64.tgz  -C bin $(NAME) && rm -f $(OUTPUT)
	GOOS=linux    GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-w $(LDFLAGS_VERSION)" -o $(OUTPUT)     ./cmd/$(NAME) && tar -czf dist/$(NAME)-linux-amd64.tgz   -C bin $(NAME) && rm -f $(OUTPUT)
	GOOS=windows  GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS_VERSION)"    -o $(OUTPUT).exe ./cmd/$(NAME) && zip -j dist/$(NAME)-windows-amd64.zip $(OUTPUT).exe && rm -f $(OUTPUT).exe

.PHONY: lint
lint: ## Validate style and syntax
ifdef GOLANGCI_LINT_VERSION
	golangci-lint run
else
	@echo "Installing latest golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
	@echo "[OK] golangci-lint installed"
	./bin/golangci-lint run
endif

.PHONY: tidy
tidy: ## Remove unused dependencies
	go mod tidy

.PHONY: list
list: ## Print the current module's dependencies.
	go list -m all

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'