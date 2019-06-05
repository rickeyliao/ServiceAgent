SHELL := /bin/bash

# The name of the executable (default is current directory name)
BASENAME := $(shell echo $${PWD\#\#*/})
TARGET := nbssa
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 1.0.1
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
#LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -linkmode=external -v"
LDFLAGS=-x -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)" 
#LDFLAGS=-race -x -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)" 

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./test/*")

.PHONY: all build clean install uninstall fmt simplify check run Test preinstmac installmac

all: check build

Test:
	@echo $(LDFLAGS)
	@echo $(eval srcname:= $(shell which ${BASENAME}))
	@echo $(srcname)
	@echo $$(which ${BASENAME}) $(subst ${BASENAME},$(TARGET),$(srcname))

rpcservice := app/pb


proto:
	protoc -I=$(rpcservice)  --go_out=plugins=grpc:${rpcservice}   ${rpcservice}/*.proto

$(TARGET): proto $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)
	@rm -f $(rpcservice)/*.pb.go
install:
	@go install $(LDFLAGS)
	@mv $$(which ${BASENAME})  $(subst ${BASENAME},$(TARGET),$$(which ${BASENAME}))

preinstmac: $(TARGET)
	@go install $(LDFLAGS)

installmac: preinstmac
	@$(eval srcname:= $(shell which ${BASENAME}))
	@mv $(srcname) $(subst $(BASENAME),$(TARGET),$(srcname))

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /test/); do golint $${d}; done
	#@go vet ${SRC}

run: install
	@$(TARGET)
