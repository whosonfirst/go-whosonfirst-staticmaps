package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"github.com/whosonfirst/iso8601duration"
	"log"
	"os"
	"time"
)

func main() {

	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")
	prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")
	exclude := flag.String("exclude", "", "Exclude repositories with this prefix")
	updated_since := flag.String("updated-since", "", "A valid ISO8601 duration string (months are currently not supported)")	
	forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")
	token := flag.String("token", "", "A valid GitHub API access token")

	flag.Parse()

	opts := organizations.NewDefaultListOptions()

	opts.Prefix = *prefix
	opts.Exclude = *exclude
	opts.Forked = *forked
	opts.NotForked = *not_forked
	opts.AccessToken = *token

	if *updated_since != "" {

		// maybe also this https://github.com/araddon/dateparse ?
		
		d, err := duration.FromString(*updated_since)

		if err != nil {
			log.Fatal(err)
		}

		now := time.Now()
		since := now.Add(- d.ToDuration())

		opts.PushedSince = &since
	}

	repos, err := organizations.ListRepos(*org, opts)

	if err != nil {
		log.Fatal(err)
	}

	for _, name := range repos {
		fmt.Println(name)
	}

	os.Exit(0)
}
