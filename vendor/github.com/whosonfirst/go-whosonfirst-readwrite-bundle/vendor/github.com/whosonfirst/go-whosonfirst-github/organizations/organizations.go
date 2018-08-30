package organizations

import (
	"github.com/google/go-github/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	_ "log"
	"strings"
	"time"
)

type ListOptions struct {
	Prefix      string
	Exclude     string
	Forked      bool
	NotForked   bool
	AccessToken string
	PushedSince *time.Time
}

func NewDefaultListOptions() *ListOptions {

	opts := ListOptions{
		Prefix:      "",
		Exclude:     "",
		Forked:      false,
		NotForked:   false,
		AccessToken: "",
		PushedSince: nil,
	}

	return &opts
}

func ListRepos(org string, opts *ListOptions) ([]string, error) {

	repos := make([]string, 0)

	cb := func(r *github.Repository) error {
		repos = append(repos, *r.Name)
		return nil
	}

	err := ListReposWithCallback(org, opts, cb)

	return repos, err
}

func ListReposWithCallback(org string, opts *ListOptions, cb func(repo *github.Repository) error) error {

	client, ctx, err := util.NewClientAndContext(opts.AccessToken)

	if err != nil {
		return err
	}

	gh_opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {

		possible, resp, err := client.Repositories.ListByOrg(ctx, org, gh_opts)

		if err != nil {
			return err
		}

		for _, r := range possible {

			if opts.Prefix != "" && !strings.HasPrefix(*r.Name, opts.Prefix) {
				continue
			}

			if opts.Exclude != "" && strings.HasPrefix(*r.Name, opts.Exclude) {
				continue
			}

			if opts.Forked && !*r.Fork {
				continue
			}

			if opts.NotForked && *r.Fork {
				continue
			}

			if opts.PushedSince != nil {

				if r.PushedAt.Before(*opts.PushedSince) {
					continue
				}
			}

			err := cb(r)

			if err != nil {
				return err
			}

		}

		if resp.NextPage == 0 {
			break
		}

		gh_opts.ListOptions.Page = resp.NextPage
	}

	return nil
}
