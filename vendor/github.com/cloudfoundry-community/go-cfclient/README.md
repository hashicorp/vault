# go-cfclient

[![build workflow](https://github.com/cloudfoundry-community/go-cfclient/actions/workflows/build.yml/badge.svg?branch=master)](https://github.com/cloudfoundry-community/go-cfclient/actions/workflows/build.yml)
[![GoDoc](https://godoc.org/github.com/cloudfoundry-community/go-cfclient?status.svg)](http://godoc.org/github.com/cloudfoundry-community/go-cfclient)
[![Report card](https://goreportcard.com/badge/github.com/cloudfoundry-community/go-cfclient)](https://goreportcard.com/report/github.com/cloudfoundry-community/go-cfclient)

## Overview

`cfclient` is a package to assist you in writing apps that need to interact with [Cloud Foundry](http://cloudfoundry.org).
It provides functions and structures to retrieve and update


## Usage

```
go get github.com/cloudfoundry-community/go-cfclient
```

NOTE: Currently this project is not versioning its releases and so breaking changes might be introduced.
Whilst hopefully notifications of breaking changes are made via commit messages, ideally your project will use a local
vendoring system to lock in a version of `go-cfclient` that is known to work for you.
This will allow you to control the timing and maintenance of upgrades to newer versions of this library.

Some example code:

```go
package main

import (
	"fmt"

	"github.com/cloudfoundry-community/go-cfclient"
)

func main() {
	c := &cfclient.Config{
		ApiAddress: "https://api.10.244.0.34.xip.io",
		Username:   "admin",
		Password:   "secret",
	}
	client, _ := cfclient.NewClient(c)
	apps, _ := client.ListApps()
	fmt.Println(apps)
}
```

### Paging Results

The API supports paging results via query string parameters. All of the v3 ListV3*ByQuery functions support paging. Only a subset of v2 function calls support paging the results:

- ListSpacesByQuery
- ListOrgsByQuery
- ListAppsByQuery
- ListServiceInstancesByQuery
- ListUsersByQuery

You can iterate over the results page-by-page using a function similar to this one:

```go
func processSpacesOnePageAtATime(client *cfclient.Client) error {
	page := 1
	pageSize := 50

	q := url.Values{}
	q.Add("results-per-page", strconv.Itoa(pageSize))

	for {
		// get the current page of spaces
		q.Set("page", strconv.Itoa(page))
		spaces, err := client.ListSpacesByQuery(q)
		if err != nil {
			fmt.Printf("Error getting spaces by query: %s", err)
			return err
		}

		// do something with each space
		fmt.Printf("Page %d:\n", page)
		for _, s := range spaces {
			fmt.Println("  " + s.Name)
		}

		// if we hit an empty page or partial page, that means we're done
		if len(spaces) < pageSize {
			break
		}

		// next page
		page++
	}
	return nil
}
```

## Development

```shell
make all
```

### Errors

If the Cloud Foundry error definitions change at <https://github.com/cloudfoundry/cloud_controller_ng/blob/master/vendor/errors/v2.yml>
then the error predicate functions in this package need to be regenerated.

To do this, simply use Go to regenerate the code:

```shell
make generate
```

## Contributing

Pull requests welcome. Please ensure you run all the unit tests, go fmt the code, and golangci-lint via `make all`
