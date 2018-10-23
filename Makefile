
VERSION ?= v3.1
PREFIX ?= 192.168.1.52/tenx_containers

BINARY_BUILD_IMAGE = golang:1.9.2-alpine3.7
DOCKER_BUILD_IMAGE = $(PREFIX)/httpmq-proxy:$(VERSION)

all: build-binary

build-binary: clean
	go build -v -o bin/httpmqProxy main.go
	@#docker run --rm -v $(shell pwd):/go/src/$(shell basename $(shell pwd)) \
	#-w /go/src/$(shell basename $(shell pwd))/cmd/ ${BINARY_BUILD_IMAGE} \
	#/bin/sh -c "go build -v -o ../deploy/httpmqProxy"

build-image: build-binary
	docker build -t ${DOCKER_BUILD_IMAGE} deploy/

clean:
	rm -rf bin/httpmqProxy

.PHONY: all build-binary build-image clean
