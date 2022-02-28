

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

all: help

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

##@ Build

build: ## Build the translation binary to the ./bin/ directory
	echo "Compiling for every OS and Platform"
	echo "FreeBSD"
	GOOS=freebsd GOARCH=386 go build -o bin/translate-freebsd-386
	GOOS=freebsd GOARCH=arm go build -o bin/translate-freebsd-arm
	GOOS=freebsd GOARCH=amd64 go build -o bin/translate-freebsd-amd64
	GOOS=freebsd GOARCH=arm64 go build -o bin/translate-freebsd-arm64
	echo "MacOS X"
	GOOS=darwin GOARCH=amd64 go build -o bin/translate-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/translate-darwin-arm64
	echo "Linux"
	GOOS=linux GOARCH=386 go build -o bin/translate-linux-386
	GOOS=linux GOARCH=amd64 go build -o bin/translate-linux-amd64
	GOOS=linux GOARCH=arm go build -o bin/translate-linux-arm
	GOOS=linux GOARCH=arm64 go build -o bin/translate-linux-arm64
	echo "Windows"
	GOOS=windows GOARCH=386 go build -o bin/translate-windows-386
	GOOS=windows GOARCH=amd64 go build -o bin/translate-windows-amd64
	GOOS=windows GOARCH=arm go build -o bin/translate-windows-arm
	GOOS=windows GOARCH=arm64 go build -o bin/translate-windows-arm64

run: ## Run the translation program without building
	go run .