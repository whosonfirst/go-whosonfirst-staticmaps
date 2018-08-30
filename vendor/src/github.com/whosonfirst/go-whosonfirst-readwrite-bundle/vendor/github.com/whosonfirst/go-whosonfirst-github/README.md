# go-whosonfirst-github

Go package for working with Who's On First GitHub repositories.

## Tools

### Webhooks

_All of the webhook tools need some documentation loving..._

#### wof-create-hook

_Please write me_

```
./bin/wof-create-hook -token {TOKEN} -hook-url {URL} -hook-secret {SECRET} -org whosonfirst-data -repo whosonfirst-data-venue-us-il
```

You can also create webhooks for all of the repositories in an organization by passing the `-repo '*'` flag. You can still filter the list of repos by setting the `-prefix` flag.

```
./bin/wof-create-hook -token {TOKEN} -hook-url {URL} -hook-secret {SECRET} -org whosonfirst-data -repo '*' -prefix whosonfirst-data
fetching repo list...🕓 
2017/04/05 15:42:24 webhook already configured for whosonfirst-data, skipping
2017/04/05 15:42:24 created webhook for whosonfirst-data-venue-us-wv
2017/04/05 15:42:25 created webhook for whosonfirst-data-venue-us-ne
2017/04/05 15:42:25 created webhook for whosonfirst-data-venue-us-wi
2017/04/05 15:42:25 created webhook for whosonfirst-data-venue-us-nv
2017/04/05 15:42:25 created webhook for whosonfirst-data-venue-us-ar
2017/04/05 15:42:25 created webhook for whosonfirst-data-venue-us-ms

...and so on
```

#### wof-list-hooks

_Please write me_

#### wof-update-hook

_Please write me_

```
./bin/wof-update-hook -token {TOKEN} -hook-url {URL} -hook-secret {NEW_SECRET} -org whosonfirst-data -repo whosonfirst-data-venue-us-il
```

You can also update webhooks for all of the repositories in an organization by passing the `-repo '*'` flag. You can still filter the list of repos by setting the `-prefix` flag.

```
./bin/wof-update-hook -token {TOKEN} -hook-url {URL} -hook-secret {NEW_SECRET} -org whosonfirst-data -repo '*' -prefix whosonfirst-data
fetching repo list...🕓 
2017/04/05 15:42:24 edited webhook for whosonfirst-data
2017/04/05 15:42:24 edited webhook for whosonfirst-data-venue-us-wv
2017/04/05 15:42:25 edited webhook for whosonfirst-data-venue-us-ne
2017/04/05 15:42:25 edited webhook for whosonfirst-data-venue-us-wi

...and so on
```

### Repos

#### wof-clone-repos

Clone (or update from `master`) Who's On First data repositories in parallel.

```
./bin/wof-clone-repos -h
Usage of ./bin/wof-clone-repos:
  -destination string
    	Where to clone repositories to (default "/usr/local/data")
  -dryrun
    	Go through the motions but don't actually clone (or update) anything
  -giturl
    	Clone using Git URL (rather than default HTTPS)
  -org string
    	The name of the organization to clone repositories from (default "whosonfirst-data")
  -prefix string
    	Limit repositories to only those with this prefix (default "whosonfirst-data")
  -procs int
    	The number of concurrent processes to clone with (default 20)
```

#### wof-list-repos

Print (to STDOUT) the list of repository names for an organization.

```
./bin/wof-list-repos -h
Usage of ./bin/wof-list-repos:
  -exclude string
    	Exclude repositories with this prefix
  -forked
    	Only include repositories that have been forked
  -not-forked
    	Only include repositories that have not been forked
  -org string
    	The name of the organization to clone repositories from (default "whosonfirst-data")
  -prefix string
    	Limit repositories to only those with this prefix (default "whosonfirst-data")
  -token string
    	A valid GitHub API access token
  -updated-since string
    	A valid ISO8601 duration string (months are currently not supported)
```

For example:

```
./bin/wof-list-repos -org whosonfirst -prefix '' -forked | sort
Clustr
emoji-search
flamework
flamework-api
flamework-artisanal-integers
flamework-aws
flamework-geo
flamework-invitecodes
flamework-multifactor-auth
flamework-storage
flamework-tools
go-pubsocketd
go-ucd
now
privatesquare
py-flamework-api
py-machinetag
py-slack-api
python-edtf
reachable
redis-tools
slackcat
suncalc-go
walk
watchman
whereonearth-metropolitan-area
youarehere-www
```

## Caveats

### Things this package doesn't deal with (yet)

* Anything that requires a GitHub API access token
* Anything other than the `master` branch of a repository
* The ability to exclude specific repositories

## See also

* https://github.com/whosonfirst-data/
* https://github.com/google/go-github
* https://en.wikipedia.org/wiki/ISO_8601#Durations