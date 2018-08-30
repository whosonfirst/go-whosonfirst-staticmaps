package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-readwrite-github/reader"
	"github.com/whosonfirst/go-whosonfirst-readwrite-github/writer"
	"log"
	"strings"
)

func main() {

	var source = flag.String("source", "", "...")
	var target = flag.String("target", "", "...")

	flag.Parse()

	parts := strings.Split(*source, "#")

	if len(parts) != 2 {
		log.Fatal("Invalid source")
	}

	repo := parts[0]
	branch := parts[1]

	r, err := reader.NewGitHubReader(repo, branch)

	if err != nil {
		log.Fatal(err)
	}

	w, err := writer.NewGitHubWriter(*target)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		fh, err := r.Read(path)

		if err != nil {
			log.Fatal(err)
		}

		err = w.Write(path, fh)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("copied %s to %s\n", r.URI(path), w.URI(path))
	}

}
