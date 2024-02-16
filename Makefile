PACKAGE     = ceyes
DATE       ?= $(shell date +%FT%T%z)
VERSION    ?= $(shell git describe --tags --always)

PKG_LIST    = $(shell go list ./... | grep -v /vendor/ | grep -v /scripts/)

GO          = go
GOLINT      = golangci-lint
GODOC       = godoc
GOFMT       = gofmt

V           = 0
Q           = $(if $(filter 1,$V),,@)
M           = $(shell printf "\033[0;35m▶\033[0m")


.PHONY: all
all: clean build

# Help
go-version: ## Print go version
	$Q $(GO) version

# Vendor
.PHONY: vendor
vendor: ## Create vendor directory from go.sum
	$(info $(M) running mod vendor…) @
	$Q $(GO) mod vendor

# Tidy
.PHONY: tidy
tidy: ## Update go.sum with go.mod
	$(info $(M) running mod tidy…) @
	$Q $(GO) mod tidy

# Build
.PHONY: build
build: ## Build the project
	$(info $(M) running build…) @
	$Q $(GO) build -o ./$(PACKAGE) cmd/main.go

# Lint
.PHONY: lint
lint: ## Run linter check on project
	$(info $(M) running $(GOLINT)…)
	$Q $(GOLINT) run

# Test
.PHONY: test
test: ## Run all unit tests
	$(info $(M) running all tests...)
	$Q $(GO) test -v ./...

# Clean
.PHONY: clean
clean: ## Clean the project
	$(info $(M) cleaning…) @
	$Q rm -f $(PACKAGE)

.PHONY: cover
cover: ## Run test cover and open it into a html browser
	$(info $(M) running coverprofile…) @
	$Q $(GO) test -v ./... -coverprofile=/tmp/c.out && go tool cover -html=/tmp/c.out

# Check
.PHONY: check
check: vendor lint test