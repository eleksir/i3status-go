#!/usr/bin/env gmake -f

BUILDOPTS=-ldflags="-s -w" -a -gcflags=all=-l -trimpath

BINARY=i3status-go

## Use calssic targets where first one is deafult target
all: clean build

## This target compiles binary
build:
	go build ${BUILDOPTS} -o ${BINARY} ./cmd/${BINARY}


## Remove binary with golang compiler' means
clean:
	rm -rf ${BINARY}


## Misc target, for development purposes. Updates vendored libs, brutal way.
upgrade:
	$(RM) -r vendor
	go get -d -u -t ./...
	go mod tidy
	go mod vendor

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
