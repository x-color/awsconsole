.DEFAULT_GOAL := help

help: ## Display help messages
	@grep -E '^[a-zA-Z_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Setup the project
	@go mod tidy

fmt: ## Format code
	@go fmt ./...

lint: ## Run linter
	@go vet ./...

build: ## Build the binary
	@go build -o awsconsole ./cmd/awsconsole/main.go

install: ## Install the binary
	@go install ./cmd/awsconsole/main.go
