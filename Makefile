
VERSION ?= v3.1
PREFIX ?= 192.168.1.52/tenx_containers

BINARY_BUILD_IMAGE = golang:1.9.2-alpine3.7
DOCKER_BUILD_IMAGE = $(PREFIX)/httpmq-proxy:$(VERSION)

all: build

build: clean
	CGO_ENABLED=0;GOOS="linux" go build -v -o httpmqProxy

mac: clean
	CGO_ENABLED=0;GOOS="darwin" go build -v -o httpmqProxy

clean:
	rm -rf httpmqProxy

.PHONY: all build-binary build-image clean
