PACKAGE := tryssh/cmd/version
GO_VERSION := $(shell go version | awk '{print $$3}')
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS :=

GIT_TAG = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X '$(PACKAGE).TrysshVersion=$(VERSION)'
endif

LDFLAGS += -X '$(PACKAGE).BuildGoVersion=$(GO_VERSION)'
LDFLAGS += -X '$(PACKAGE).BuildTime=$(BUILD_TIME) UTC'

.PHONY: default
default: build

.PHONY: build
build: tidy
	@go build -v -ldflags "$(LDFLAGS)" ./

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: clean
clean:
	@go clean
