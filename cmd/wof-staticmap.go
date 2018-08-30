package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-readwrite-bundle"	
	"github.com/whosonfirst/go-whosonfirst-staticmap"
	"image/png"
	"log"
	"os"
	"strconv"
)

func main() {

	str_valid := bundle.ValidReadersString()

	desc := fmt.Sprintf("DSN strings MUST contain a 'reader=SOURCE' pair followed by any additional pairs required by that reader. Supported reader sources are: %s.", str_valid)

	var dsn_flags flags.MultiDSNString
	flag.Var(&dsn_flags, "dsn", desc)

	var height = flag.Int("height", 480, "The height in pixels of your new map.")
	var width = flag.Int("width", 640, "The width in pixels of your new map.")
	var saveas = flag.String("save-as", "", "Save the map to this path. If empty then the map will saved as {WOFID}.png.")

	// deprecated
	// var wofid = flag.Int64("id", 0, "A valid Who's On First to render.")
	// var root = flag.String("data-root", "https://whosonfirst.mapzen.com/data", "Where to look for Who's On First source data.")

	flag.Parse()

	wof_ids := make([]int64, 0)
	
	for _, str_id := range flag.Args() {

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		wof_ids = append(wof_ids, id)
	}

	if len(wof_ids) == 0 {
		log.Fatal("No IDs to render...")
	}

	r, err := bundle.NewMultiReaderFromFlags(dsn_flags)

	if err != nil {
		log.Fatal(err)
	}
	
	sm, err := staticmap.NewStaticMap(r)

	if err != nil {
		log.Fatal(err)
	}

	sm.Width = *width
	sm.Height = *height

	im, err := sm.Render(wof_ids...)

	if err != nil {
		log.Fatal(err)
	}

	if *saveas == "" {
		*saveas = fmt.Sprintf("%s.png", "debug") // FIX ME
	}
	
	fh, err := os.Create(*saveas)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()

	err = png.Encode(fh, im)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
