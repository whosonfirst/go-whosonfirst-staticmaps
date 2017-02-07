# go-whosonfirst-staticmap

Go package for rendering static maps of Who's On First places.

## Important

Too soon. Move along.

## Install

You will need to have both `Go` and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Example

_Note that all error handle has been removed for the sake of brevity._

```
import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-staticmap"
	"image/png"		
	"log"
)

func main() {

	wofid := 85688637
	png := fmt.Sprintf("%d.png", wofid)
	
	sm, _ := staticmap.NewStaticMap(wofid)
	im, _ := sm.Render()

	fh, _ := os.Create(png)
	defer fh.Close()

	png.Encode(file, im)
}
```

## Tools

### wof-staticmap

```
./bin/wof-staticmap -h
Usage of ./bin/wof-staticmap:
  -data-root string
    	Where to look for Who's On First source data. (default "https://whosonfirst.mapzen.com/data")
  -id int
    	A valid Who's On First to render.
  -image-height int
    	... (default 480)
  -image-width int
    	... (default 640)
  -save-as string
    	Save the map to this path. If empty then the map will saved as {WOFID}.png.
```

### wof-staticmapd

_Please write me_

## See also

* github.com/flopp/go-staticmaps
