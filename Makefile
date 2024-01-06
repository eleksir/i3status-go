#!/usr/bin/env gmake -f

BUILDOPTS=-ldflags="-s -w" -a -gcflags=all=-l -trimpath
FILELIST=batteries.go clock.go collection.go config.go cputemp.go cron.go i3.go ifmon.go la.go main.go memory.go \
         pulseaudio.go signal.go spawn.go stdin.go vpn.go

BINARY=i3status-go

## Use calssic targets where first one is deafult target
all: clean build

## This target compiles binary
build:
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY} ${FILELIST}


## Remove binary with golang compiler' means
clean:
	go clean


## Misc target, for development purposes. Updates vendored libs, brutal way.
upgrade:
	$(RM) -r vendor
	go get -d -u -t ./...
	go mod tidy
	go mod vendor

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
