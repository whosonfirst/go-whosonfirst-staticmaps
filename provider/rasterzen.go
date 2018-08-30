package provider

import (
	"bytes"
	"github.com/whosonfirst/go-rasterzen/nextzen"
	"github.com/whosonfirst/go-rasterzen/tile"
	"github.com/whosonfirst/go-staticmaps"
)

type RasterzenTileProvider struct {
	sm.TileProvider
	api_key     string
	name        string
	attribution string
	tileSize    int
	urlPattern  string
	shards      []string
}

func NewRasterzenTileProvider(apikey string) (sm.TileProvider, error) {

	t := &DefaultTileProvider{
		name:        "rasterzen",
		attribution: "OSM, Nextzen",
		tileSize:    256,
		urlPattern:  "",
		shards:      []string{},
	}

	return t
}

func (t *RasterzenTileProvider) Name() string {
	return t.name
}

func (t *RasterzenTileProvider) Attribution() string {
	return t.attribution
}

func (t *RasterzenTileProvider) TileSize() int {
	return t.tileSize
}

func (t *RasterzenTileProvider) URLPattern() string {
	return t.urlPattern
}

func (t *RasterzenTileProvider) Shards() []string {
	return t.shards
}

func (t *RasterzenTileProvider) TileURL(zoom int, x int, y int) string {

	return fmt.Sprintf(t.URLPattern(), shard, zoom, x, y)
}

func (t *RasterzenTileProvider) FetchTile(z int, x int, y int) ([]byte, error) {

	raw, err := nextzen.FetchTile(z, x, y, t.api_key)

	if err != nil {
		return nil, err
	}

	cropped, err := nextzen.CropTile(z, x, y, raw)

	if err != nil {
		return nil, err
	}

	wr := new(bytes.Buffer)
	im, err := tile.ToPNG(cropped, wr)

	if err != nil {
		return nil, err
	}

	return wr.Bytes(), nil
}
