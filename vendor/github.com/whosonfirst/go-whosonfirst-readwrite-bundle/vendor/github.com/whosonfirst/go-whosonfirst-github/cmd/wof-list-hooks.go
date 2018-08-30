package main

// https://godoc.org/github.com/google/go-github/github#Hook
// https://developer.github.com/v3/orgs/hooks/#list-hooks

// curl -s -i -H "Authorization: token {TOKEN}" https://api.github.com/orgs/whosonfirst/hooks
// curl -s -i -H "Authorization: token {TOKEN}" https://api.github.com/repos/whosonfirst-data/whosonfirst-data/hooks

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	"log"
	"os"
	"strings"
)

func main() {

	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")
	token := flag.String("oauth2-token", "", "...")

	org_hooks := flag.Bool("org-hooks", false, "Show hooks for the organization itself, not child repositories")

	prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")
	forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")

	flag.Parse()

	client, ctx, err := util.NewClientAndContext(*token)

	if err != nil {
		log.Fatal(err)
	}

	if *org_hooks {

		opts := github.ListOptions{PerPage: 100}

		for {

			hooks, rsp, err := client.Organizations.ListHooks(ctx, *org, &opts)

			if err != nil {
				log.Fatal(err)
			}

			for _, h := range hooks {

				log.Println(fmt.Sprintf("%s has webhook %s (active: %t)", *org, *h.URL, *h.Active))

				// please add a flag to toggle display of the actual webhook URL...
				// log.Println(fmt.Sprintf("%s has webhook %s (active: %t)", *org, h.Config["url"], *h.Active))

				log.Println(h)
			}

			if rsp.NextPage == 0 {
				break
			}

			opts.Page = rsp.NextPage
		}

		os.Exit(0)
	}

	// Get all the webhooks for all the repositories for an organization

	repos_opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, repos_rsp, err := client.Repositories.ListByOrg(ctx, *org, repos_opts)

		if err != nil {
			log.Fatal(err)
		}

		for _, r := range repos {

			if *prefix != "" && !strings.HasPrefix(*r.Name, *prefix) {
				continue
			}

			if *forked && !*r.Fork {
				continue
			}

			if *not_forked && *r.Fork {
				continue
			}

			hooks_opts := github.ListOptions{PerPage: 100}

			hooks, _, err := client.Repositories.ListHooks(ctx, *org, *r.Name, &hooks_opts)

			if err != nil {
				log.Fatal(err)
			}

			for _, h := range hooks {

				log.Println(fmt.Sprintf("%s has webhook %s (active: %t)", *r.Name, *h.URL, *h.Active))

				// please add a flag to toggle display of the actual webhook URL...
				// log.Println(fmt.Sprintf("%s has webhook %s (active: %t)", *r.Name, h.Config["url"], *h.Active))
			}

		}

		if repos_rsp.NextPage == 0 {
			break
		}

		repos_opts.ListOptions.Page = repos_rsp.NextPage
	}

	os.Exit(0)
}
