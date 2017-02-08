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

_work in progress..._

```
./bin/wof-staticmapd -h
Usage of ./bin/wof-staticmapd:
  -cache
    	...
  -cache-provider string
    	... (default "s3")
  -data-root string
    	Where to look for Who's On First source data. (default "https://whosonfirst.mapzen.com/data")
  -gracehttp.log
    	Enable logging. (default true)
  -height int
    	The default height in pixels for rendered maps. (default 480)
  -host string
    	The hostname to listen for requests on (default "localhost")
  -port int
    	The port number to listen for requests on (default 8080)
  -s3-bucket string
    	... (default "whosonfirst.mapzen.com")
  -s3-credentials string
    	... (default "shared:/Users/asc/.aws/credentials:default")
  -s3-prefix string
    	... (default "static")
  -s3-region string
    	... (default "us-east-1")
  -size value
    	Zero or more custom {LABEL}={WIDTH}x{HEIGHT} parameters.
  -width int
    	The default width in pixels for rendered maps. (default 640)
```

## Examples

### Default sizes

```
./bin/wof-staticmapd 
curl http://127.0.0.1:8080/?id=1108794405&size=sq
```

![](images/1108794405.png)

### Custom sizes

```
./bin/wof-staticmapd -size sq=100x100
curl http://127.0.0.1:8080/?id=1108794405&size=sq
```

![](images/1108794405-sq.png)

## See also

* github.com/flopp/go-staticmaps
