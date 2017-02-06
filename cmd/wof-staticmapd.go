package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/whosonfirst/go-whosonfirst-staticmap"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

/* oh god put me in a library */

// all of this S3 stuff is cloned from https://github.com/thisisaaronland/go-iiif/blob/master/aws/s3.go
// and probably deserves to be moved in to a bespoke package some day... (20170131/thisisaaronland)

type S3Connection struct {
	service *s3.S3
	bucket  string
	prefix  string
}

type S3Config struct {
	Bucket      string
	Prefix      string
	Region      string
	Credentials string // see notes below
}

func NewS3Connection(s3cfg S3Config) (*S3Connection, error) {

	// https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/

	cfg := aws.NewConfig()
	cfg.WithRegion(s3cfg.Region)

	if strings.HasPrefix(s3cfg.Credentials, "env:") {

		creds := credentials.NewEnvCredentials()
		cfg.WithCredentials(creds)

	} else if strings.HasPrefix(s3cfg.Credentials, "shared:") {

		details := strings.Split(s3cfg.Credentials, ":")

		if len(details) != 3 {
			return nil, errors.New("Shared credentials need to be defined as 'shared:CREDENTIALS_FILE:PROFILE_NAME'")
		}

		creds := credentials.NewSharedCredentials(details[1], details[2])
		cfg.WithCredentials(creds)

	} else if strings.HasPrefix(s3cfg.Credentials, "iam:") {

		// assume an IAM role suffient for doing whatever

	} else {

		return nil, errors.New("Unknown S3 config")
	}

	sess := session.New(cfg)

	if s3cfg.Credentials != "" {

		_, err := sess.Config.Credentials.Get()

		if err != nil {
			return nil, err
		}
	}

	service := s3.New(sess)

	c := S3Connection{
		service: service,
		bucket:  s3cfg.Bucket,
		prefix:  s3cfg.Prefix,
	}

	return &c, nil
}

func (conn *S3Connection) Put(key string, body []byte, content_type string) error {

	key = conn.prepareKey(key)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(conn.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(content_type),
	}

	_, err := conn.service.PutObject(params)

	if err != nil {
		return err
	}

	return nil
}

func (conn *S3Connection) prepareKey(key string) string {

	if conn.prefix == "" {
		return key
	}

	return filepath.Join(conn.prefix, key)
}

/* end of oh god put me in a library */

func main() {

	whoami, err := user.Current()
	default_creds := ""

	if err == nil {
		default_creds = fmt.Sprintf("shared:%s/.aws/credentials:default", whoami.HomeDir)
	}

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var s3_credentials = flag.String("s3-credentials", default_creds, "...")
	var s3_bucket = flag.String("s3-bucket", "whosonfirst.mapzen.com", "...")
	var s3_prefix = flag.String("s3-prefix", "static", "...")
	var s3_region = flag.String("s3-region", "us-east-1", "...")

	var height = flag.Int("image-height", 480, "...")
	var width = flag.Int("image-width", 640, "...")

	var cache = flag.Bool("cache", false, "...")

	flag.Parse()

	cfg := S3Config{
		Bucket:      *s3_bucket,
		Prefix:      *s3_prefix,
		Region:      *s3_region,
		Credentials: *s3_credentials,
	}

	conn, err := NewS3Connection(cfg)

	if err != nil {
		log.Fatal(err)
	}

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		query := req.URL.Query()

		str_wofid := query.Get("id")

		if str_wofid == "" {
			http.Error(rsp, "Missing ID parameter", http.StatusBadRequest)
			return
		}

		wofid, err := strconv.Atoi(str_wofid)

		if err != nil {
			http.Error(rsp, "Invalid ID parameter", http.StatusBadRequest)
			return
		}

		sm, err := staticmap.NewStaticMap(int64(wofid))

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		sm.Width = *width
		sm.Height = *height

		im, err := sm.Render()

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		buffer := new(bytes.Buffer)

		err = png.Encode(buffer, im)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if *cache {

			go func() {

				root, err := uri.Id2Path(wofid)

				if err != nil {
					log.Println(err)
					return
				}

				fname := fmt.Sprintf("%d-%d-%d.png", wofid, *width, *height)

				rel_path := filepath.Join(root, fname)

				err = conn.Put(rel_path, buffer.Bytes(), "image/png")

				if err != nil {
					msg := fmt.Sprintf("failed to PUT %s because %s\n", rel_path, err)
					log.Println(msg)
					return
				}
			}()
		}

		rsp.Header().Set("Content-Type", "image/png")
		rsp.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		rsp.Write(buffer.Bytes())
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err = gracehttp.Serve(&http.Server{Addr: endpoint, Handler: mux})

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)

}
