# Shush - Comment Removal CLI Tool
.PHONY: help build install clean test lint fmt vet run dev release tag push check deps update-deps

# Default target
help: ## Show this help message
	@echo "Shush - Comment Removal CLI Tool"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build commands
build: ## Build the binary
	go build -ldflags="-w -s" -o shush ./cmd/shush

build-dev: ## Build without optimizations (faster, for development)
	go build -o shush ./cmd/shush

install: build ## Build and install to GOPATH/bin
	go install ./cmd/shush

# Cross-platform builds
build-all: ## Build for all platforms
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o dist/shush-linux-amd64 ./cmd/shush
	GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o dist/shush-linux-arm64 ./cmd/shush
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o dist/shush-darwin-amd64 ./cmd/shush
	GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o dist/shush-darwin-arm64 ./cmd/shush

# Development commands
run: build-dev ## Build and run with example file
	./shush --help

dev: build-dev ## Quick development cycle (build + show help)
	./shush --llm | head -20

test: ## Run tests
	go test -v ./...

test-race: ## Run tests with race detection
	go test -race -v ./...

test-cover: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Code quality - NOTE: We don't use golangci-lint in this project
lint: ## Run basic linting (go vet only)
	@echo "Running basic linting with go vet..."
	go vet ./...

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

check: fmt vet test ## Run all checks (format, vet, test)

# Dependency management
deps: ## Download dependencies
	go mod download

tidy: ## Tidy dependencies
	go mod tidy

update-deps: ## Update dependencies
	go get -u ./...
	go mod tidy

# Release commands
version: ## Show current version
	@grep 'var version' cmd/shush/main.go | cut -d'"' -f2

release: clean check build ## Prepare for release (clean, check, build)
	@echo "Ready for release. Current version: $$(make version)"
	@echo "To create release:"
	@echo "  1. Update version in cmd/shush/main.go"
	@echo "  2. Run: make tag"
	@echo "  3. Run: make push"

tag: ## Create and push git tag for current version
	$(eval VERSION := $(shell make version))
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@echo "Created tag v$(VERSION)"
	@echo "Run 'make push' to push tag and trigger release"

push: ## Push current branch and tags
	git push upstream
	git push upstream --tags

# Utility commands
clean: ## Clean build artifacts
	rm -f shush
	rm -rf dist/
	rm -f coverage.out coverage.html

demo: build-dev ## Run demo commands
	@echo "=== Shush Demo ==="
	@echo "1. Version info:"
	./shush --version
	@echo ""
	@echo "2. Help output:"
	./shush --help
	@echo ""
	@echo "3. LLM guide preview:"
	./shush --llm | head -10

# Development helpers
watch: ## Watch for changes and rebuild (requires entr)
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | entr -r make dev; \
	else \
		echo "entr not installed. Install with your package manager"; \
		echo "  macOS: brew install entr"; \
		echo "  Ubuntu: apt install entr"; \
	fi

size: build ## Show binary size
	@ls -lh shush | awk '{print "Binary size: " $$5}'

# Create dist directory
dist:
	mkdir -p dist