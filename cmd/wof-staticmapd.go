package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/whosonfirst/go-storagemaster"
	"github.com/whosonfirst/go-storagemaster/provider"
	"github.com/whosonfirst/go-whosonfirst-staticmap"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
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

	whoami, err := user.Current()
	default_creds := ""

	if err == nil {
		default_creds = fmt.Sprintf("shared:%s/.aws/credentials:default", whoami.HomeDir)
	}

	var sizes NamedSizes

	flag.Var(&sizes, "size", "Zero or more custom {LABEL}={WIDTH}x{HEIGHT} parameters.")

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var cache_s3 = flag.Bool("cache-s3", false, "...")

	var s3_credentials = flag.String("s3-credentials", default_creds, "...")
	var s3_bucket = flag.String("s3-bucket", "whosonfirst.mapzen.com", "...")
	var s3_prefix = flag.String("s3-prefix", "static", "...")
	var s3_region = flag.String("s3-region", "us-east-1", "...")

	var height = flag.Int("height", 480, "The default height in pixels for rendered maps.")
	var width = flag.Int("width", 640, "The default width in pixels for rendered maps.")
	var root = flag.String("data-root", "https://whosonfirst.mapzen.com/data", "Where to look for Who's On First source data.")

	flag.Parse()

	sz_map, err := sizes.ToMap()

	if err != nil {
		log.Fatal(err)
	}

	sz_map["default"] = CustomSize{Width: *width, Height: *height}

	var storage storagemaster.Provider

	if *cache_s3 {

		cfg := provider.S3Config{
			Bucket:      *s3_bucket,
			Prefix:      *s3_prefix,
			Region:      *s3_region,
			Credentials: *s3_credentials,
		}

		storage, err = provider.NewS3Provider(cfg)

		if err != nil {
			log.Fatal(err)
		}
	}

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

		if *cache_s3 {

			go func() {

				root, err := uri.Id2Path(wofid)

				if err != nil {
					log.Println(err)
					return
				}

				fname := fmt.Sprintf("%d.png", wofid)

				if sz_label != "default" {
					fname = fmt.Sprintf("%d-%s.png", wofid, sz_label)
				}

				rel_path := filepath.Join(root, fname)

				extras, err := storagemaster.NewStoragemasterExtras()

				if err != nil {
					msg := fmt.Sprintf("failed to PUT %s because %s\n", rel_path, err)
					log.Println(msg)
					return
				}

				extras.Set("acl", "public-read")
				extras.Set("content-type", "image/png")

				err = storage.Put(rel_path, buffer.Bytes(), extras)

				if err != nil {
					msg := fmt.Sprintf("failed to PUT %s because %s\n", rel_path, err)
					log.Println(msg)
					return
				}
			}()
		}

		rsp.Header().Set("Content-Type", "image/png")
		rsp.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		rsp.Write(buffer.Bytes())
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err = gracehttp.Serve(&http.Server{Addr: endpoint, Handler: mux})

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)

}
