package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-staticmap"
	_ "github.com/whosonfirst/go-whosonfirst-uri"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type CustomSize struct {
	Width  int
	Height int
}

type NamedSizes []string

func (p *NamedSizes) String() string {
	return strings.Join(*p, "\n")
}

func (p *NamedSizes) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func (p *NamedSizes) ToMap() (map[string]CustomSize, error) {

	m := make(map[string]CustomSize)

	for _, str_pair := range *p {

		pair := strings.Split(str_pair, "=")

		name := pair[0]
		dims := strings.Split(pair[1], "x")

		w, err := strconv.Atoi(dims[0])

		if err != nil {
			return nil, err
		}

		h, err := strconv.Atoi(dims[1])

		if err != nil {
			return nil, err
		}

		m[name] = CustomSize{Width: w, Height: h}
	}

	return m, nil
}

func main() {

	var sizes NamedSizes

	flag.Var(&sizes, "size", "Zero or more custom {LABEL}={WIDTH}x{HEIGHT} parameters.")

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	/*
	var cache = flag.Bool("cache", false, "Cache rendered maps")
	var cache_provider = flag.String("cache-provider", "s3", "A valid cache provider. Valid options are: s3")

	var s3_credentials = flag.String("s3-credentials", default_creds, "A string descriptor for your AWS credentials. Valid options are: env:;shared:PATH_TO_SHARED_CREDENTIALS_FILE:SHARED_CREDENTIALS_PROFILE; iam:")
	var s3_bucket = flag.String("s3-bucket", "whosonfirst.mapzen.com", "A valid S3 bucket where cached files are stored.")
	var s3_prefix = flag.String("s3-prefix", "static", "An optional subdirectory (prefix) where cached files are stored in S3.")
	var s3_region = flag.String("s3-region", "us-east-1", "A valid AWS S3 region")
	*/
	
	var height = flag.Int("height", 480, "The default height in pixels for rendered maps.")
	var width = flag.Int("width", 640, "The default width in pixels for rendered maps.")
	var root = flag.String("data-root", "https://whosonfirst.mapzen.com/data", "Where to look for Who's On First source data.")

	flag.Parse()

	sz_map, err := sizes.ToMap()

	if err != nil {
		log.Fatal(err)
	}

	sz_map["default"] = CustomSize{Width: *width, Height: *height}

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		query := req.URL.Query()

		str_wofid := query.Get("id")

		if str_wofid == "" {
			http.Error(rsp, "Missing ID parameter", http.StatusBadRequest)
			return
		}

		wofid, err := strconv.Atoi(str_wofid)

		if err != nil {
			http.Error(rsp, "Invalid 'id' parameter", http.StatusBadRequest)
			return
		}

		sz_label := query.Get("size")

		if sz_label == "" {
			sz_label = "default"
		}

		sz, ok := sz_map[sz_label]

		if !ok {
			http.Error(rsp, "Invalid 'size' parameter", http.StatusBadRequest)
			return
		}

		sm, err := staticmap.NewStaticMap(int64(wofid))

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		sm.DataRoot = *root
		sm.Width = sz.Width
		sm.Height = sz.Height

		im, err := sm.Render()

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		buffer := new(bytes.Buffer)

		err = png.Encode(buffer, im)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "image/png")
		rsp.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		rsp.Write(buffer.Bytes())
	}

	ping := func(rsp http.ResponseWriter, req *http.Request) {

		rsp.Header().Set("Content-Type", "text/plain")
		rsp.Write([]byte("PONG"))
	}
	
	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/ping", ping)	

	err = http.ListenAndServe(endpoint, mux)	

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)

}
