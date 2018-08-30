package provider

import (
	"errors"
	"flag"
	"github.com/whosonfirst/go-staticmaps"
	"sort"
)

func ValidProviders() []string {

	provider_names := []string{"rasterzen"}
	providers := sm.GetTileProviders()

	for name, _ := range providers {
		provider_names = append(provider_names, name)
	}

	sort.Strings(provider_names)
	return provider_names
}

func NewTileProviderFromFlags() (sm.TileProvider, error) {

	fl_map := make(map[string]*flag.Flag)

	flag.VisitAll(func(fl *flag.Flag) {
		fl_map[fl.Name] = fl
	})

	pr, ok := fl_map["provider"]

	if !ok {
		return nil, errors.New("Missing -provider flag")
	}

	name := pr.Value.(flag.Getter).Get().(string) // U SO WEIRD GO...
	args := make([]interface{}, 0)

	if name == "rasterzen" {

		key, ok := fl_map["nextzen-api-key"]

		if !ok {
			return nil, errors.New("Missing -nextzen-api-key flag")
		}

		args = append(args, key.Value.(flag.Getter).Get()) // SO SO WEIRD...
	}

	return NewTileProviderFromString(name, args...)
}

func NewTileProviderFromString(name string, args ...interface{}) (sm.TileProvider, error) {

	var tp sm.TileProvider
	var err error

	if name == "rasterzen" {

		if len(args) == 0 {
			return nil, errors.New("Missing API key")
		}

		api_key := args[0].(string)

		tp, err = NewRasterzenTileProvider(api_key)

	} else {

		providers := sm.GetTileProviders()
		p, ok := providers[name]

		if ok {
			tp = p
		} else {
			err = errors.New("Invalid tile provider")
		}

	}

	return tp, err
}
