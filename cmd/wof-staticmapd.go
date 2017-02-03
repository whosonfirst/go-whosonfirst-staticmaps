package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/whosonfirst/go-whosonfirst-staticmap"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var cache = flag.Bool("cache", false, "...")

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

		sm, err := staticmap.NewStaticMap(int64(wofid))

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

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

		if *cache {
			// cache image here...
		}

		rsp.Header().Set("Content-Type", "image/png")
		rsp.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		rsp.Write(buffer.Bytes())
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := gracehttp.Serve(&http.Server{Addr: endpoint, Handler: mux})

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)

}
