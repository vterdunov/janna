PROG_NAME = janna
USERNAME = vterdunov
IMAGE_NAME = $(USERNAME)/$(PROG_NAME)

COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%S')
PROJECT ?= github.com/$(USERNAME)/${PROG_NAME}

GO_VARS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_LDFLAGS := -ldflags '-extldflags "-fno-PIC -static" \
	-X ${PROJECT}/internal/version.Commit=${COMMIT} \
	-X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}' \
	-tags 'osusergo netgo static_build'

GOLANGCI_LINTER_IMAGE = golangci/golangci-lint:v1.17.1

.PHONY: help
help: ## Display this help message
	@echo "Please use \`make <target>\` where <target> is one of:"
	@cat $(MAKEFILE_LIST) | grep -e "^[-a-zA-Z_\.]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: compile
compile: ## Build binary
	$(GO_VARS) go build -v $(GO_LDFLAGS) -o $(PROG_NAME) ./cmd/server/server.go

.PHONY: run
run: ## Extract env variables from .env and run server with race detector
	@env `cat .env | grep -v ^# | xargs` go run -race ./cmd/server/server.go

compile-and-run: compile ## Extract env variables from .env. Compile and run server
	@env `cat .env | grep -v ^# | xargs` ./$(PROG_NAME)

.PHONY: lint
lint: ## Run linters
	@echo Linting...
	@docker run --tty --rm -v $(CURDIR):/lint -w /lint $(GOLANGCI_LINTER_IMAGE) /bin/sh -c 'if ! [ -d 'vendor' ]; then go mod download; fi && golangci-lint run'

.PHONY: mock
mock:
	@mockery -dir internal/usecase -output internal/usecase/ -outpkg usecase_test -case snake -all -testonly
