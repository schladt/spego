BASE_NAME=spego-stub
GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
OS=$(shell go env GOOS)
# ARCH=$(shell go env GOARCH)
# EXT=$(shell go env GOEXE)
# BINARY_PATH=bin/$(BASE_NAME)-$(OS)-$(ARCH)$(EXT)

all: build

build:
	go env -w GOOS=windows
	go env -w GOARCH=386
	go env -w CC=i686-w64-mingw32-gcc
	go env -w CGO_ENABLED=1
	$(GOBUILD) -v -ldflags="-s -w" -o bin/$(BASE_NAME)-windows-386.exe .

	go env -w GOARCH=amd64
	go env -w CC=x86_64-w64-mingw32-gcc
	go env -w CGO_ENABLED=1
	$(GOBUILD) -v -ldflags="-s -w" -o bin/$(BASE_NAME)-windows-amd64.exe .

run:
	build
	./$(BINARY_PATH)

deps:
	$(GOGET) github.com/brimstone/go-shellcode