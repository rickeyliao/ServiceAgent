SHELL := /bin/bash

# The name of the executable (default is current directory name)
BASENAME := $(shell echo $${PWD\#\#*/})
TARGET := nbssa
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
#VERSION := 1.0.1
#BUILD := `git rev-parse HEAD`
#BUILDTIME := `date "+%Y-%m-%d/%H:%M:%S/%Z"`

include .Makefile

# Use linker flags to provide version/build settings to the target
#LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.BuildTime=$(BUILDTIME) -linkmode=external -v"
LDFLAGS=-x -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.BuildTime=$(BUILDTIME) -w"
#LDFLAGS=-race -x -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.BuildTime=$(BUILDTIME)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./test/*")

.PHONY: all build clean install uninstall fmt simplify check run Test preinst

all: check build

Test:
	@echo $(LDFLAGS)
	@echo $(eval srcname:= $(shell which ${BASENAME}))
	@echo $(BUILDTIME)
	@echo $$(which ${BASENAME}) $(subst ${BASENAME},$(TARGET),$(srcname))

rpcservice := app/pb
dhtrpc := dht/pb
assetdir := ui/asset
wifiapdir := wifiap/control

proto:
	protoc -I=$(rpcservice)  --go_out=plugins=grpc:${rpcservice}   ${rpcservice}/*.proto
	protoc -I=$(dhtrpc)  --go_out=plugins=grpc:${dhtrpc}   ${dhtrpc}/*.proto

staticfile2bin:
	go-bindata -o $(assetdir)/res.go -pkg=asset ui/xadmin/...
	go-bindata -o $(wifiapdir)/res.go -pkg=control wifiap/staticfile/...

$(TARGET): staticfile2bin proto $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)
	@rm -f $(rpcservice)/*.pb.go $(assetdir)/res.go
#install: preinst
#	#@go install $(LDFLAGS)
#	@mv $$(which ${BASENAME})  $(subst ${BASENAME},$(TARGET),$$(which ${BASENAME}))

preinst: $(TARGET)
	@go install $(LDFLAGS)

install: preinst
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
