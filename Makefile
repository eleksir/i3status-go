#!/usr/bin/env gmake -f

BUILDOPTS=-ldflags="-s -w" -a -gcflags=all=-l -trimpath -buildvcs=false

BINARY=i3status-go
TEST1=battery-test
TEST2=cmdrun-test

## Use calssic targets where first one is deafult target
all: clean build

## This target compiles binary
build:
	go build ${BUILDOPTS} -o ${BINARY} ./cmd/${BINARY}

battery-test:
	rm -rf ${TEST1}
	go build ${BUILDOPTS} -o ${TEST1} ./cmd/${TEST1}

cmdrun-test:
	rm -rf ${TEST1}
	go build ${BUILDOPTS} -o ${TEST2} ./cmd/${TEST2}

## Remove binary with golang compiler' means
clean:
	rm -rf ${BINARY} ${TEST1} ${TEST2}

## Misc target, for development purposes. Updates vendored libs, brutal way.
upgrade:
	$(RM) -r vendor
	go get -u -t ./...
	go mod tidy
	go mod vendor

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
