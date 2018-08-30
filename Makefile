CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test ! -d src; then mkdir src; fi
	if test ! -d src/github.com/whosonfirst/go-whosonfirst-staticmap; then mkdir -p src/github.com/whosonfirst/go-whosonfirst-staticmap; fi
	cp staticmap.go src/github.com/whosonfirst/go-whosonfirst-staticmap/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite-bundle"
	# @GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-staticmaps"
	@GOPATH=$(GOPATH) go get -u "github.com/Wessie/appdirs"
	@GOPATH=$(GOPATH) go get -u "github.com/flopp/go-coordsparser"
	@GOPATH=$(GOPATH) go get -u "github.com/fogleman/gg"
	@GOPATH=$(GOPATH) go get -u "github.com/golang/geo/..."
	@GOPATH=$(GOPATH) go get -u "github.com/tkrajina/gpxgo/gpx"
	cp -r /usr/local/whosonfirst/go-staticmaps src/github.com/whosonfirst/
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	mv src/github.com/whosonfirst/go-whosonfirst-readwrite-bundle/vendor/github.com/whosonfirst/go-whosonfirst-cli src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-readwrite-bundle/vendor/github.com/whosonfirst/go-whosonfirst-readwrite src/github.com/whosonfirst/

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-staticmap cmd/wof-staticmap.go
	# @GOPATH=$(GOPATH) go build -o bin/wof-staticmapd cmd/wof-staticmapd.go

fmt:
	go fmt cmd/*.go
	go fmt staticmap.go
