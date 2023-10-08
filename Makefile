PACKAGE := tryssh/cmd/version
VERSION := v0.3.3
GO_VERSION := $(shell go version | awk '{print $$3}')
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS :=

LDFLAGS += -X '$(PACKAGE).TrysshVersion=$(VERSION)'
LDFLAGS += -X '$(PACKAGE).BuildGoVersion=$(GO_VERSION)'
LDFLAGS += -X '$(PACKAGE).BuildTime=$(BUILD_TIME) UTC'

.PHONY: default
default: build

.PHONY: build
build: tidy
	@go build -ldflags "$(LDFLAGS)" ./

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: clean
clean:
	@go clean
