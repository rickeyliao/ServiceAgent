SHELL := /bin/bash

# The name of the executable (default is current directory name)
BASENAME := $(shell echo $${PWD\#\#*/})
TARGET := sa
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
#LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -linkmode=external -v"
LDFLAGS=-race -x -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)" 

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./test/*")

.PHONY: all build clean install uninstall fmt simplify check run Test

all: check install

Test:
	@echo $(LDFLAGS)
	@echo $(GOBIN)

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)

install:
	@go install $(LDFLAGS)
	@mv $$(which ${BASENAME}) $(subst $(BASENAME),$(TARGET),$$(which ${BASENAME}))

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /test/); do golint $${d}; done
	@go tool vet ${SRC}

run: install
	@$(TARGET)
