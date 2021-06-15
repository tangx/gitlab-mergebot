
PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = v$(shell cat .version)
COMMIT_SHA ?= $(shell git describe --always)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD=CGO_ENABLED=0 go build -a -ldflags "-X ${PKG}/version.Version=${VERSION}+sha.${COMMIT_SHA}"

APP ?= gitlab-mergebot
WORKSPACE ?= ./cmd/$(APP)

up:
	cd $(WORKSPACE) && go run .

clean:
	rm -rf out

tidy:
	go mod tidy

build:
	make build.$(APP) GOOS=darwin GOARCH=amd64
	make build.$(APP) GOOS=linux GOARCH=amd64
	make build.$(APP) GOOS=linux GOARCH=arm64

build.$(APP): tidy
	@echo "Building $(APP) for $(GOOS)/$(GOARCH)"
	cd $(WORKSPACE) && $(GOBUILD) -o ../../out/$(APP)-$(GOOS)-$(GOARCH)

install: build.$(APP)
	mv out/$(APP)-$(GOOS)-$(GOARCH) ${GOPATH}/bin/$(APP)

