# Go parameters

SHELL=powershell
GOCMD=go
GOBUILD=$(GOCMD) build -ldflags '-w -s -H windowsgui'
BINARY_NAME=LocalizeEpub
BINARY_UNIX=$(BINARY_NAME)_unix
Version=1.0.1

all: get_module system build

get_module:
	$(GOCMD) mod download

system:
ifeq ($(OS),Windows_NT)
    SYSTEM := .exe
else
	UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        SYSTEM := 
    endif
    ifeq ($(UNAME_S),Darwin)
        SYSTEM = .app
    endif
endif

build:
	$(GOBUILD) -o $(BINARY_NAME)_$(Version)$(SYSTEM) -v

build-windows_amd64:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_LDFLAGS="-static" $(GOBUILD) -o $(BINARY_NAME)_$(Version)_windows_amd64.exe -v

build-linux_amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc CXX=g++ CGO_LDFLAGS="-static" $(GOBUILD) -o $(BINARY_NAME)_$(Version)_linux_amd64 -v

build-darwin_amd64:
    CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ CGO_LDFLAGS="-static" $(GOBUILD) -o $(BINARY_NAME)_$(Version)_darwin_amd64.app -v

build-darwin_arm64:
    CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC=oa64-clang CXX=oa64-clang++ CGO_LDFLAGS="-static" $(GOBUILD) -o $(BINARY_NAME)_$(Version)_darwin_amd64.app -v

