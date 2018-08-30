package staticmap

import (
	_ "errors"
	"fmt"
	"github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"image"
	"io/ioutil"
	_ "log"
	"net/http"
	"strings"
)

type StaticMap struct {
	DataRoot     string
	TileProvider string
	Fill         string
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
		Fill:         "0xFF00967F",
		wofid:        wofid,
	}

	return &sm, nil
}

func (s *StaticMap) Render() (image.Image, error) {

	b, err := s.Fetch()

	if err != nil {
		return nil, err
	}

	ctx := sm.NewContext()

	tileProviders := sm.GetTileProviders()
	tp := tileProviders[s.TileProvider]

	if tp != nil {
		ctx.SetTileProvider(tp)
	}

	ctx.SetSize(s.Width, s.Height)

	geom_type := gjson.GetBytes(b, "geometry.type").String()
	coords := gjson.GetBytes(b, "geometry.coordinates")

	if geom_type == "Polygon" || geom_type == "MultiPolygon" {

		bbox := gjson.GetBytes(b, "bbox").Array()

		swlat := bbox[1].Float()
		swlon := bbox[0].Float()
		nelat := bbox[3].Float()
		nelon := bbox[2].Float()

		s2_bbox, err := sm.CreateBBox(nelat, swlon, swlat, nelon)

		if err != nil {
			return nil, err
		}

		ctx.SetBoundingBox(*s2_bbox)

		areas := make([]*sm.Area, 0)

		if geom_type == "Polygon" {

			for _, poly := range coords.Array() {

				area, err := s.poly2area(poly)

				if err != nil {
					return nil, err
				}

				areas = append(areas, area)
			}
		} else {

			for _, multi := range coords.Array() {

				for _, poly := range multi.Array() {

					area, err := s.poly2area(poly)

					if err != nil {
						return nil, err
					}

					areas = append(areas, area)
				}

			}
		}

		for _, a := range areas {
			ctx.AddArea(a)
		}

	} else {

		latlon := coords.Array()

		lat := latlon[1].Float()
		lon := latlon[0].Float()

		ctx.SetCenter(s2.LatLngFromDegrees(lat, lon))
	}

	if geom_type == "Point" {

		label_lat := gjson.GetBytes(b, "properties.lbl:latitude")
		label_lon := gjson.GetBytes(b, "properties.lbl:longitude")

		if label_lat.Exists() && label_lon.Exists() {

			label_marker := fmt.Sprintf("color:%s|%0.6f,%0.6f", s.Fill, label_lat.Float(), label_lon.Float())

			markers, err := sm.ParseMarkerString(label_marker)

			if err != nil {
				return nil, err
			}

			for _, marker := range markers {
				ctx.AddMarker(marker)
			}

		} else {

			geom_lat := gjson.GetBytes(b, "properties.geom:latitude").Float()
			geom_lon := gjson.GetBytes(b, "properties.geom:longitude").Float()

			geom_marker := fmt.Sprintf("color:%s|%0.6f,%0.6f", s.Fill, geom_lat, geom_lon)

			markers, err := sm.ParseMarkerString(geom_marker)

			if err != nil {
				return nil, err
			}

			for _, marker := range markers {
				ctx.AddMarker(marker)
			}

		}
	}

	return ctx.Render()
}

// please put me in a utility function or something... (20170203/thisisaaronland)

func (s *StaticMap) Fetch() ([]byte, error) {

	url, err := uri.Id2AbsPath(s.DataRoot, s.wofid)

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

func (s *StaticMap) poly2area(poly gjson.Result) (*sm.Area, error) {

	fill := fmt.Sprintf("fill:%s", s.Fill)

	args := []string{
		"color:0x000000",
		fill,
		"weight:2",
	}

	for _, ring := range poly.Array() {
		pt := ring.Array()
		lat := pt[1].Float()
		lon := pt[0].Float()

		str_pt := fmt.Sprintf("%0.6f,%0.6f", lat, lon)
		args = append(args, str_pt)
	}

	str_args := strings.Join(args, "|")

	area, err := sm.ParseAreaString(str_args)

	if err != nil {
		return nil, err
	}

	return area, nil
}
