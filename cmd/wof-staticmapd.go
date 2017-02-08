package main

import (
	"bytes"
	"errors"
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

func main() {

	whoami, err := user.Current()
	default_creds := ""

	if err == nil {
		default_creds = fmt.Sprintf("shared:%s/.aws/credentials:default", whoami.HomeDir)
	}

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var cache_s3 = flag.Bool("cache-s3", false, "...")

	var s3_credentials = flag.String("s3-credentials", default_creds, "...")
	var s3_bucket = flag.String("s3-bucket", "whosonfirst.mapzen.com", "...")
	var s3_prefix = flag.String("s3-prefix", "static", "...")
	var s3_region = flag.String("s3-region", "us-east-1", "...")

	var height = flag.Int("image-height", 480, "The height in pixels for rendered maps.")
	var width = flag.Int("image-width", 640, "The width in pixels for rendered maps.")
	var root = flag.String("data-root", "https://whosonfirst.mapzen.com/data", "Where to look for Who's On First source data.")

	flag.Parse()

	var sm storagemaster.Provider

	if *cache_s3 {

		cfg := provider.S3Config{
			Bucket:      *s3_bucket,
			Prefix:      *s3_prefix,
			Region:      *s3_region,
			Credentials: *s3_credentials,
		}

		sm, err = provider.NewS3Provider(cfg)

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
			http.Error(rsp, "Invalid ID parameter", http.StatusBadRequest)
			return
		}

		// log.Println("rendering", wofid)

		sm, err := staticmap.NewStaticMap(int64(wofid))

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		sm.DataRoot = *root
		sm.Width = *width
		sm.Height = *height

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

				fname := fmt.Sprintf("%d-%d-%d.png", wofid, *width, *height)

				rel_path := filepath.Join(root, fname)

				extras, err := storagemaster.NewStorageMasterExtras()

				if err != nil {
					msg := fmt.Sprintf("failed to PUT %s because %s\n", rel_path, err)
					log.Println(msg)
					return
				}

				extras.Set("acl", "public-read")
				extras.Set("content-type", "image/png")
				
				err = sm.Put(rel_path, buffer.Bytes(), extras)

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
