SHELL = /bin/bash

NAME := jx-v2-tekton-converter
ORG := jstrachan
ORG_REPO := $(ORG)/$(NAME)
RELEASE_ORG_REPO := $(ORG_REPO)
REV := $(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
ROOT_PACKAGE := github.com/$(ORG_REPO)
BRANCH     := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y%m%d-%H:%M:%S)
#GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
GO_VERSION := 1.12

GO := GO111MODULE=on go
BUILD_TARGET = build
CGO_ENABLED = 0

# set dev version unless VERSION is explicitly set via environment
#VERSION ?= $(shell echo "$$(git for-each-ref refs/tags/ --count=1 --sort=-version:refname --format='%(refname:short)' 2>/dev/null)-dev+$(REV)" | sed 's/^v//')

MAIN_SRC_FILE=./main.go
BUILDFLAGS :=  -ldflags \
  " -X main.buildTime=$(BUILD_DATE) \
		-X main.gitCommit=$(REV) \
		-X main.version=$(VERSION)"

.PHONY: build
build:
	go build -o bin/$(NAME) main.go

release: build test

fmt:
	go fmt ./...

test:
	go test ./...

linux: ## Build for Linux
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o build/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x build/linux/$(NAME)


bindata:
	go-bindata -o pkg/assets/assets.go -pkg assets resources/git/git-clone.yaml

.PHONY: goreleaser
goreleaser:
	step-go-releaser --organisation=$(ORG) --revision=$(REV) --branch=$(BRANCH) --build-date=$(BUILD_DATE) --root-package=$(ROOT_PACKAGE) --go-version=$(GO_VERSION) --version=$(VERSION)

.PHONY: clean
clean: ## Clean the generated artifacts
	rm -rf bin release dist


