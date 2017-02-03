CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test ! -d src; then mkdir src; fi
	if test ! -d src/github.com/whosonfirst/go-whosonfirst-staticmap; then mkdir -p src/github.com/whosonfirst/go-whosonfirst-staticmap; fi
	cp staticmap.go src/github.com/whosonfirst/go-whosonfirst-staticmap/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/flopp/go-staticmaps"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"
	@GOPATH=$(GOPATH) go get -u "github.com/aws/aws-sdk-go"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"

vendor-deps: deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-render-staticmap cmd/wof-render-staticmap.go

fmt:
	go fmt cmd/*.go
	go fmt staticmap.go
