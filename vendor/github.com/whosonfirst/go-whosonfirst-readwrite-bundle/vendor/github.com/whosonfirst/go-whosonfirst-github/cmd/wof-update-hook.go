package main

// https://developer.github.com/v3/repos/hooks/#edit-a-hook
// https://godoc.org/github.com/google/go-github/github#RepositoriesService.EditHook
// https://godoc.org/github.com/google/go-github/github#OrganizationsService.EditHook
// https://godoc.org/github.com/google/go-github/github#Hook

import (
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/google/go-github/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	"log"
	"strings"
	"time"
)

func main() {

	org := flag.String("org", "", "")
	repo := flag.String("repo", "", "")
	token := flag.String("token", "", "...")
	prefix := flag.String("prefix", "", "Limit repositories to only those with this prefix")

	name := flag.String("hook-name", "web", "")
	url := flag.String("hook-url", "", "")
	content_type := flag.String("hook-content-type", "json", "")
	secret := flag.String("hook-secret", "", "")

	flag.Parse()

	if *token == "" {
		log.Fatal("Missing OAuth2 token")
	}

	client, ctx, err := util.NewClientAndContext(*token)

	if err != nil {
		log.Fatal(err)
	}

	lookup := make(map[string]*github.Hook)

	if *repo == "" {

		log.Fatal("Organizations are not supported yet")

	} else {

		repos := make([]string, 0)

		if *repo == "*" {

			done := make(chan bool)

			go func() {

				sp := spinner.New(spinner.CharSets[38], 200*time.Millisecond)
				sp.Prefix = "fetching repo list..."
				sp.Start()

				for {

					select {
					case <-done:
						sp.Stop()
						return
					}
				}
			}()

			repos_opts := &github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{PerPage: 100},
			}

			for {

				repos_list, repos_rsp, err := client.Repositories.ListByOrg(ctx, *org, repos_opts)

				if err != nil {
					log.Fatal(err)
				}

				for _, r := range repos_list {

					if *prefix != "" && !strings.HasPrefix(*r.Name, *prefix) {
						continue
					}

					repos = append(repos, *r.Name)

					hooks_opts := github.ListOptions{PerPage: 100}

					hooks, _, err := client.Repositories.ListHooks(ctx, *org, *r.Name, &hooks_opts)

					if err != nil {
						log.Fatal(err)
					}

					for _, h := range hooks {

						if h.Config["url"] == *url {
							lookup[*r.Name] = h
							break
						}
					}

				}

				if repos_rsp.NextPage == 0 {
					break
				}

				repos_opts.ListOptions.Page = repos_rsp.NextPage
			}

			done <- true

		} else {

			log.Fatal("Please get hooks for a single repo")
			repos = append(repos, *repo)
		}

		for _, r := range repos {

			hook, ok := lookup[r]

			if !ok {
				log.Println(fmt.Sprintf("webhook not configured for %s, skipping", r))
				continue
			}

			if *secret != "" {
				hook.Config["secret"] = *secret
			}

			if *content_type != "" {
				hook.Config["content_type"] = *content_type
			}

			if *name != "" {
				hook.Name = name
			}

			_, _, err = client.Repositories.EditHook(ctx, *org, r, *hook.ID, hook)

			if err != nil {
				log.Fatal(fmt.Sprintf("failed to edit webhook for %s, because %s", r, err.Error()))
			}

			log.Println(fmt.Sprintf("edited webhook for %s", r))
		}
	}

}
