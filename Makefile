VERSION ?= $(shell git describe --tags)

IMAGE = codeclimate/codeclimate-staticcheck
PKG = github.com/alexkappa/codeclimate-staticcheck
PKGS = $(shell go list ./... | grep -v /vendor/)

BFLAGS = -a -tags netgo
LDFLAGS = "-s -w"

OS ?= darwin
ARCH ?= amd64

build:
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o bin/codeclimate-staticcheck-$(OS)-$(ARCH) $(BFLAGS) -ldflags $(LDFLAGS)

image:
	@docker build -t $(IMAGE) .
	@docker tag $(IMAGE):latest $(IMAGE):$(VERSION)

analyze:
	@codeclimate analyze --dev
