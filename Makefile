PROG_NAME = janna
USERNAME = vterdunov
IMAGE_NAME = $(USERNAME)/$(PROG_NAME)

TAG ?= $(COMMIT)

COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%S')
PROJECT ?= github.com/$(USERNAME)/${PROG_NAME}

GO_VARS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_LDFLAGS := -ldflags '-extldflags "-fno-PIC -static" \
	-X ${PROJECT}/internal/version.Commit=${COMMIT} \
	-X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}' \
	-tags 'osusergo netgo static_build'

GOLANGCI_LINTER_IMAGE = golangci/golangci-lint:v1.17.1

all: lint docker

.PHONY: docker
docker: ## Build Docker container
	docker build --tag=$(IMAGE_NAME):$(COMMIT) --tag=$(IMAGE_NAME):latest --build-arg=GITHUB_TOKEN=${GITHUB_TOKEN} --file build/Dockerfile .

.PHONY: compile
compile: ## Build server binary
	$(GO_VARS) go build -v $(GO_LDFLAGS) -o $(PROG_NAME) ./cmd/api/main.go
	$(GO_VARS) go build -v $(GO_LDFLAGS) -o worker ./cmd/worker/main.go

.PHONY: compile-worker-debug
compile-worker-debug: ## Build worker without compiler optomizations
	$(GO_VARS) go build -v -gcflags "all=-N -l" -o worker ./cmd/worker/main.go

.PHONY: test
test: ## Run tests. With -race flag
	go test -race -count=100 -v ./...

.PHONY: push
push: ## Push docker container to registry
	docker tag $(IMAGE_NAME):$(COMMIT) $(IMAGE_NAME):$(TAG)
	docker push $(IMAGE_NAME):$(TAG)

.PHONY: run
run: ## Extract env variables from .env and run server with race detector
	@env `cat .env | grep -v ^# | xargs` go run -race ./cmd/server/main.go

compile-and-run: compile ## Extract env variables from .env. Compile and run server
	@env `cat .env | grep -v ^# | xargs` ./$(PROG_NAME)

.PHONY: lint
lint: ## Run linters
	@echo Linting...
	@docker run --tty --rm -v $(CURDIR):/lint -v $$HOME/go/pkg/mod:/go/pkg/mod -w /lint $(GOLANGCI_LINTER_IMAGE) golangci-lint run

.PHONY: tools
tools:
	GO111MODULE=off go get github.com/vektra/mockery/.../

.PHONY: generate
generate:
	go generate ./...

.PHONY: compose-run
compose-run: compile ## Start whole stack of Janna services for development
	docker-compose -f deploy/docker-compose.dev.yml up --build --scale worker=3

.PHONY: compose-run-debug
compose-run-debug: compile compile-worker-debug
	docker-compose -f deploy/docker-compose.dev.yml up --build

.PHONY: compose-clean
compose-clean: ## Clean containers and its data
	docker-compose -f deploy/docker-compose.dev.yml down --volumes

.PHONY: help
help: ## Display this help message
	@echo "Please use \`make <target>\` where <target> is one of:"
	@cat $(MAKEFILE_LIST) | grep -e "^[-a-zA-Z_\.]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
