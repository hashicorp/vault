# packngo
Packet Go Api Client

![](https://www.packet.net/media/images/xeiw-packettwitterprofilew.png)


Installation
------------

`go get github.com/packethost/packngo`

Usage
-----

To authenticate to the Packet API, you must have your API token exported in env var `PACKET_API_TOKEN`.

This code snippet initializes Packet API client, and lists your Projects:

```go
package main

import (
	"log"

	"github.com/packethost/packngo"
)

func main() {
	c, err := packngo.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	ps, _, err := c.Projects.List(nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range ps {
		log.Println(p.ID, p.Name)
	}
}

```

This lib is used by the official [terraform-provider-packet](https://github.com/terraform-providers/terraform-provider-packet).

You can also learn a lot from the `*_test.go` sources. Almost all out tests touch the Packet API, so you can see how auth, querying and POSTing works. For example [devices_test.go](devices_test.go).



Acceptance Tests
----------------

If you want to run tests against the actual Packet API, you must set envvar `PACKET_TEST_ACTUAL_API` to non-empty string for the `go test`. The device tests wait for the device creation, so it's best to run a few in parallel.

To run a particular test, you can do

```
$ PACKNGO_TEST_ACTUAL_API=1 go test -v -run=TestAccDeviceBasic
```

If you want to see HTTP requests, set the `PACKNGO_DEBUG` env var to non-empty string, for example:

```
$ PACKNGO_DEBUG=1 PACKNGO_TEST_ACTUAL_API=1 go test -v -run=TestAccVolumeUpdate
```


Committing
----------

Before committing, it's a good idea to run `gofmt -w *.go`. ([gofmt](https://golang.org/cmd/gofmt/))
