package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-readwrite-bundle"
	"github.com/whosonfirst/go-whosonfirst-staticmaps"
	"github.com/whosonfirst/go-whosonfirst-staticmaps/provider"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	valid_providers := provider.ValidProviders()
	str_providers := strings.Join(valid_providers, ",")

	desc_providers := fmt.Sprintf("A valid go-staticmaps provider. Valid providers are: %s", str_providers)

	str_valid := bundle.ValidReadersString()

	desc := fmt.Sprintf("One or more valid Who's On First reader DSN strings. DSN strings MUST contain a 'reader=SOURCE' pair followed by any additional pairs required by that reader. Supported reader sources are: %s.", str_valid)

	var dsn_flags flags.MultiDSNString
	flag.Var(&dsn_flags, "dsn", desc)

	var height = flag.Int("height", 480, "The height in pixels of your new map.")
	var width = flag.Int("width", 640, "The width in pixels of your new map.")
	var saveas = flag.String("save-as", "map.png", "Save the map to this path.")

	var api_key = flag.String("nextzen-api-key", "", "A valid Nextzen API key. Required if -provider is 'rasterzen'")

	var tile_provider = flag.String("provider", "stamen-toner", desc_providers)

	if *tile_provider == "rasterzen" && *api_key == "" {
		log.Fatal("Missing Nextzen API key for rasterzen provider")
	}

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

	tp, err := provider.NewTileProviderFromFlags()

	if err != nil {
		log.Fatal(err)
	}

	sm, err := staticmaps.NewStaticMap(tp, r)

	if err != nil {
		log.Fatal(err)
	}

	sm.Width = *width
	sm.Height = *height

	fh, err := os.OpenFile(*saveas, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()

	err = sm.RenderAsPNG(fh, wof_ids...)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
