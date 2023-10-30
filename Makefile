PACKAGE := github.com/Driver-C/tryssh/cmd/version
GO_VERSION := $(shell go version | awk '{print $$3}')
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS :=

GIT_TAG = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)

OS_ARCH_LIST=darwin:amd64 darwin:arm64 freebsd:amd64 linux:amd64 linux:arm linux:arm64 windows:amd64 windows:arm64

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= $(GIT_TAG)

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X '$(PACKAGE).TrysshVersion=$(BINARY_VERSION)'
else
	# If cannot find any information that can be used as a version number, change it to debug
	BINARY_VERSION = "debug"
endif

LDFLAGS += -X '$(PACKAGE).BuildGoVersion=$(GO_VERSION)'
LDFLAGS += -X '$(PACKAGE).BuildTime=$(BUILD_TIME) UTC'

.PHONY: default
default: build

.PHONY: build
build: tidy
	@go build -v -trimpath -ldflags "$(LDFLAGS)" ./

.PHONY: tidy
tidy: clean
	@go mod tidy

.PHONY: clean
clean:
	@go clean
	@rm -rf ./release

.PHONY: multi
multi: tidy
	@$(foreach n, $(OS_ARCH_LIST),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		target_suffix=$(BINARY_VERSION)-$${os}-$${arch};\
		bin_name="tryssh";\
		if [ $${os} = "windows" ]; then bin_name="tryssh.exe"; fi;\
		echo "[==> Build $${os}-$${arch} start... <==]";\
		mkdir -p ./release/$${os}-$${arch};\
		cp ./LICENSE ./release/$${os}-$${arch}/;\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} go build -v -trimpath -ldflags "$(LDFLAGS)" \
		-o ./release/$${os}-$${arch}/$${bin_name};\
		cd ./release;\
		zip -rq tryssh-$${target_suffix}.zip $${os}-$${arch};\
		rm -rf $${os}-$${arch};\
		cd ..;\
		echo "[==> Build $${os}-$${arch} done <==]";\
	)
