package staticmap

import (
	_ "errors"
	"fmt"
	"github.com/golang/geo/s2"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-staticmaps"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"	
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-uri"	
	"image"
	_ "log"
	"strings"
)

type StaticMap struct {
	TileProvider string
	Fill         string
	Width        int
	Height       int
	reader       reader.Reader
}

func NewStaticMap(r reader.Reader) (*StaticMap, error) {

	sm := StaticMap{
		TileProvider: "stamen-toner",
		Width:        800,
		Height:       640,
		Fill:         "0xFF00967F",
	reader: r,
	}

	return &sm, nil
}

func (s *StaticMap) Render(ids ...int64) (image.Image, error) {

	// this probably deserves to be a utility function
	// somewhere... (20180830/thisisaaronland)
	
	features := make([]geojson.Feature, 0)

	for _, id := range ids {
		
		uri, err := uri.Id2RelPath(id)

		if err != nil {
			return nil, err
		}

		fh, err := s.reader.Read(uri)

		if err != nil {
			return nil, err
		}

		f, err := feature.LoadFeatureFromReader(fh)

		if err != nil {
			return nil, err
		}

		features = append(features, f)
	}
	
	ctx := sm.NewContext()

	tileProviders := sm.GetTileProviders()
	tp := tileProviders[s.TileProvider]

	if tp != nil {
		ctx.SetTileProvider(tp)
	}

	ctx.SetSize(s.Width, s.Height)

	err := s.setExtent(ctx, features...)

	if err != nil {
		return nil, err
	}

	err = s.addMarkers(ctx, features...)

	if err != nil {
		return nil, err
	}

	return ctx.Render()
}

func (s *StaticMap) setExtent(ctx *sm.Context, features ...geojson.Feature) error {

	swlat := 0.0
	swlon := 0.0
	nelat := 0.0
	nelon := 0.0

	areas := make([]*sm.Area, 0)

	for _, f := range features {

		b := f.Bytes()

		geom_type := gjson.GetBytes(b, "geometry.type").String()
		coords := gjson.GetBytes(b, "geometry.coordinates")

		if geom_type == "Polygon" || geom_type == "MultiPolygon" {

			bbox := gjson.GetBytes(b, "bbox").Array()

			f_swlat := bbox[1].Float()
			f_swlon := bbox[0].Float()
			f_nelat := bbox[3].Float()
			f_nelon := bbox[2].Float()

			if f_swlat < swlat {
				swlat = f_swlat
			}

			if f_swlon < swlon {
				swlon = f_swlon
			}

			if f_nelat > nelat {
				nelat = f_nelat
			}

			if f_nelon < swlon {
				swlon = f_nelon
			}

			if geom_type == "Polygon" {

				for _, poly := range coords.Array() {

					area, err := s.poly2area(poly)

					if err != nil {
						return err
					}

					areas = append(areas, area)
				}
			} else {

				for _, multi := range coords.Array() {

					for _, poly := range multi.Array() {

						area, err := s.poly2area(poly)

						if err != nil {
							return err
						}

						areas = append(areas, area)
					}

				}
			}

		} else {

			latlon := coords.Array()

			lat := latlon[1].Float()
			lon := latlon[0].Float()

			if lat < swlat {
				swlat = lat
			}

			if lon < swlon {
				swlon = lon
			}

			if lat > nelat {
				nelat = lat
			}

			if lon > nelon {
				nelon = lon
			}

		}

	}

	if swlat == nelat && swlon == nelon {

		ctx.SetCenter(s2.LatLngFromDegrees(swlat, swlon))

	} else {

		s2_bbox, err := sm.CreateBBox(nelat, swlon, swlat, nelon)

		if err != nil {
			return err
		}

		ctx.SetBoundingBox(*s2_bbox)

		for _, a := range areas {
			ctx.AddArea(a)
		}

	}

	return nil
}

func (s *StaticMap) addMarkers(ctx *sm.Context, features ...geojson.Feature) error {

	for _, f := range features {

		b := f.Bytes()
		
		geom_type := gjson.GetBytes(b, "geometry.type").String()
		
		if geom_type == "Point" {

			label_lat := gjson.GetBytes(b, "properties.lbl:latitude")
			label_lon := gjson.GetBytes(b, "properties.lbl:longitude")

			if label_lat.Exists() && label_lon.Exists() {

				label_marker := fmt.Sprintf("color:%s|%0.6f,%0.6f", s.Fill, label_lat.Float(), label_lon.Float())

				markers, err := sm.ParseMarkerString(label_marker)

				if err != nil {
					return err
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
					return err
				}

				for _, marker := range markers {
					ctx.AddMarker(marker)
				}

			}
		}

	}

	return nil
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
