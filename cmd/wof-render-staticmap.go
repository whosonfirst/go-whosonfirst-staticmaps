package main

import (
	_ "flag"
	"github.com/fogleman/gg"
	"github.com/whosonfirst/go-whosonfirst-staticmap"
	"log"
)

func main() {

	wofid := int64(85668077)
	sm, err := staticmap.NewStaticMap(wofid)

	if err != nil {
		log.Fatal(err)
	}

	im, err := sm.Render()

	if err != nil {
		log.Fatal(err)
	}

	gg.SavePNG("test.png", im)
}
