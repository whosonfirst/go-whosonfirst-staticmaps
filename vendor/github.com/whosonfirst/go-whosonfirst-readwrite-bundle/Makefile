CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-readwrite-bundle
	cp *.go src/github.com/whosonfirst/go-whosonfirst-readwrite-bundle/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-cli"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-github/organizations"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite-github/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite-http/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite-mysql/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite-s3/..."
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite-sqlite/..."

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt cmd/*.go

bin: 	self
	rm -rf bin/*
	@GOPATH=$(GOPATH) go build -o bin/wof-bundle-reader cmd/wof-bundle-reader.go
