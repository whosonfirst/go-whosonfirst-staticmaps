package staticmap

import (
	"fmt"
	"github.com/flopp/go-staticmaps"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"image"
	"io/ioutil"
	_ "log"
	"net/http"
)

type StaticMap struct {
	DataRoot     string
	TileProvider string
	Width        int
	Height       int
	wofid        int64
}

func NewStaticMap(wofid int64) (*StaticMap, error) {

	sm := StaticMap{
		DataRoot:     "https://whosonfirst.mapzen.com/data",
		TileProvider: "stamen-toner",
		Width:        800,
		Height:       640,
		wofid:        wofid,
	}

	return &sm, nil
}

func (s *StaticMap) Render() (image.Image, error) {

	b, err := s.Fetch()

	if err != nil {
		return nil, err
	}

	// id := gjson.GetBytes(b, "properties.wof:id")

	bbox := gjson.GetBytes(b, "bbox").Array()

	swlat := bbox[1].Float()
	swlon := bbox[0].Float()
	nelat := bbox[3].Float()
	nelon := bbox[2].Float()

	ctx := sm.NewContext()

	tileProviders := sm.GetTileProviders()
	tp := tileProviders[s.TileProvider]

	if tp != nil {
		ctx.SetTileProvider(tp)
	}

	ctx.SetSize(s.Width, s.Height)

	geom_lat := gjson.GetBytes(b, "properties.geom:latitude").Float()
	geom_lon := gjson.GetBytes(b, "properties.geom:longitude").Float()

	geom_marker := fmt.Sprintf("color:red|%0.6f,%0.6f", geom_lat, geom_lon)

	markers, err := sm.ParseMarkerString(geom_marker)

	if err != nil {
		return nil, err
	}

	for _, marker := range markers {
		ctx.AddMarker(marker)
	}

	/*
		area_string := "color:0x00FF00|fill:0x00FF007F|weight:2|"

		area, err := sm.ParseAreaString(area_string)

		if err != nil {
		   return nil, err
		   }


		ctx.AddArea(area)
	*/

	s2_bbox, err := sm.CreateBBox(nelat, swlon, swlat, nelon)

	if err != nil {
		return nil, err
	}

	ctx.SetBoundingBox(*s2_bbox)

	return ctx.Render()
}

// please put me in a utility function or something... (20170203/thisisaaronland)

func (s *StaticMap) Fetch() ([]byte, error) {

	url, err := uri.Id2AbsPath(s.DataRoot, int(s.wofid))

	if err != nil {
		return nil, err
	}

	rsp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	return b, nil
}
