package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-staticmap"
	"image/png"
	"log"
	"os"
)

func main() {

	var wofid = flag.Int64("id", 0, "A valid Who's On First to render.")
	var height = flag.Int("image-height", 480, "The height in pixels of your new map.")
	var width = flag.Int("image-width", 640, "The width in pixels of your new map.")
	var root = flag.String("data-root", "https://whosonfirst.mapzen.com/data", "Where to look for Who's On First source data.")
	var saveas = flag.String("save-as", "", "Save the map to this path. If empty then the map will saved as {WOFID}.png.")

	flag.Parse()

	if *wofid == 0 {
		log.Fatal("Missing WOF ID")
	}

	if *saveas == "" {
		*saveas = fmt.Sprintf("%d.png", *wofid)
	}

	sm, err := staticmap.NewStaticMap(*wofid)

	if err != nil {
		log.Fatal(err)
	}

	sm.DataRoot = *root
	sm.Width = *width
	sm.Height = *height

	im, err := sm.Render()

	if err != nil {
		log.Fatal(err)
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
