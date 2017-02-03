package main

import (
	"flag"
	"github.com/fogleman/gg"
	"github.com/whosonfirst/go-whosonfirst-staticmap"
	"log"
)

func main() {

	var wofid = flag.Int64("id", 0, "...")

	flag.Parse()

	if *wofid == 0 {
		log.Fatal("Missing WOF ID")
	}

	sm, err := staticmap.NewStaticMap(*wofid)

	if err != nil {
		log.Fatal(err)
	}

	im, err := sm.Render()

	if err != nil {
		log.Fatal(err)
	}

	gg.SavePNG("test.png", im)
}
