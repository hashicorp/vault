# linodego

![Tests](https://img.shields.io/github/actions/workflow/status/linode/linodego/ci.yml?branch=main)
[![Release](https://img.shields.io/github/v/release/linode/linodego)](https://github.com/linode/linodego/releases/latest)
[![GoDoc](https://godoc.org/github.com/linode/linodego?status.svg)](https://godoc.org/github.com/linode/linodego)
[![Go Report Card](https://goreportcard.com/badge/github.com/linode/linodego)](https://goreportcard.com/report/github.com/linode/linodego)

Go client for [Linode REST v4 API](https://techdocs.akamai.com/linode-api/reference/api)

## Installation

```sh
go get -u github.com/linode/linodego
```

## Documentation

See [godoc](https://godoc.org/github.com/linode/linodego) for a complete reference.

The API generally follows the naming patterns prescribed in the [OpenAPIv3 document for Linode APIv4](https://techdocs.akamai.com/linode-api/reference/api).

Deviations in naming have been made to avoid using "Linode" and "Instance" redundantly or inconsistently.

A brief summary of the features offered in this API client are shown here.

## Examples

### General Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

func main() {
	apiKey, ok := os.LookupEnv("LINODE_TOKEN")
	if !ok {
		log.Fatal("Could not find LINODE_TOKEN, please assert it is set.")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	linodeClient := linodego.NewClient(oauth2Client)
	linodeClient.SetDebug(true)

	res, err := linodeClient.GetInstance(context.Background(), 4090913)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", res)
}
```

### Pagination

#### Auto-Pagination Requests

```go
kernels, err := linodego.ListKernels(context.Background(), nil)
// len(kernels) == 218
```

Or, use a page value of "0":

```go
opts := linodego.NewListOptions(0,"")
kernels, err := linodego.ListKernels(context.Background(), opts)
// len(kernels) == 218
```

#### Single Page

```go
opts := linodego.NewListOptions(2,"")
// or opts := linodego.ListOptions{PageOptions: &linodego.PageOptions{Page: 2}, PageSize: 500}
kernels, err := linodego.ListKernels(context.Background(), opts)
// len(kernels) == 100
```

ListOptions are supplied as a pointer because the Pages and Results
values are set in the supplied ListOptions.

```go
// opts.Results == 218
```

> **_NOTES:_**  
>	- The ListOptions will be mutated by list endpoint functions.
>	- Instances of ListOptions should NOT be shared across multiple list endpoint functions.
>	- The resulting number of results and pages can be accessed through the user-supplied ListOptions instance.

#### Filtering

```go
f := linodego.Filter{}
f.AddField(linodego.Eq, "mine", true)
fStr, err := f.MarshalJSON()
if err != nil {
    log.Fatal(err)
}
opts := linodego.NewListOptions(0, string(fStr))
stackscripts, err := linodego.ListStackscripts(context.Background(), opts)
```

### Error Handling

#### Getting Single Entities

```go
linode, err := linodego.GetInstance(context.Background(), 555) // any Linode ID that does not exist or is not yours
// linode == nil: true
// err.Error() == "[404] Not Found"
// err.Code == "404"
// err.Message == "Not Found"
```

#### Lists

For lists, the list is still returned as `[]`, but `err` works the same way as on the `Get` request.

```go
linodes, err := linodego.ListInstances(context.Background(), linodego.NewListOptions(0, "{\"foo\":bar}"))
// linodes == []
// err.Error() == "[400] [X-Filter] Cannot filter on foo"
```

Otherwise sane requests beyond the last page do not trigger an error, just an empty result:

```go
linodes, err := linodego.ListInstances(context.Background(), linodego.NewListOptions(9999, ""))
// linodes == []
// err = nil
```

### Response Caching

By default, certain endpoints with static responses will be cached into memory. 
Endpoints with cached responses are identified in their [accompanying documentation](https://pkg.go.dev/github.com/linode/linodego?utm_source=godoc).

The default cache entry expiry time is `15` minutes. Certain endpoints may override this value to allow for more frequent refreshes (e.g. `client.GetRegion(...)`).
The global cache expiry time can be customized using the `client.SetGlobalCacheExpiration(...)` method.

Response caching can be globally disabled or enabled for a client using the `client.UseCache(...)` method.

The global cache can be cleared and refreshed using the `client.InvalidateCache()` method.

### Writes

When performing a `POST` or `PUT` request, multiple field related errors will be returned as a single error, currently like:

```go
// err.Error() == "[400] [field1] foo problem; [field2] bar problem; [field3] baz problem"
```

## Tests

Run `make testunit` to run the unit tests. 

Run `make testint` to run the integration tests. The integration tests use fixtures.

To update the test fixtures, run `make fixtures`.  This will record the API responses into the `fixtures/` directory.
Be careful about committing any sensitive account details.  An attempt has been made to sanitize IP addresses and
dates, but no automated sanitization will be performed against `fixtures/*Account*.yaml`, for example.

To prevent disrupting unaffected fixtures, target fixture generation like so: `make ARGS="-run TestListVolumes" fixtures`.

## Discussion / Help

Join us at [#linodego](https://gophers.slack.com/messages/CAG93EB2S) on the [gophers slack](https://gophers.slack.com)

## License

[MIT License](LICENSE)
