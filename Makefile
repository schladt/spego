BASE_NAME=spego
GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)
EXT=$(shell go env GOEXE)
BINARY_PATH=bin/$(BASE_NAME)-$(OS)-$(ARCH)$(EXT)

all: build final

build:
	$(GOBUILD) -v -ldflags="-s -w" -o $(BINARY_PATH) .

run:
	build
	./$(BINARY_PATH)

final:
ifeq ($(OS),linux)
	$(eval OUTFILE := $(shell sha256sum $(BINARY_PATH)))
else ifeq ($(OS),darwin)
	$(eval OUTFILE := $(shell shasum -a 256 $(BINARY_PATH)))
else ifeq ($(OS),windows)
	$(eval OUTFILE := $(shell PowerShell.exe Get-FileHash $(BINARY_PATH)))
else
	$(error Unsupported: $(NATIVE))
endif
	@echo "Binary created : $(OUTFILE)"

deps:
	$(GOGET) github.com/jteeuwen/go-bindata/...