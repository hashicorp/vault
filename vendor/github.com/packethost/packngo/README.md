# packngo

[![](https://img.shields.io/badge/stability-maintained-green.svg)](https://github.com/packethost/standards/blob/master/maintained-statement.md)
[![Release](https://img.shields.io/github/v/release/packethost/packngo)](https://github.com/packethost/packngo/releases/latest)
[![GoDoc](https://godoc.org/github.com/packethost/packngo?status.svg)](https://godoc.org/github.com/packethost/packngo)
[![Go Report Card](https://goreportcard.com/badge/github.com/packethost/packngo)](https://goreportcard.com/report/github.com/packethost/packngo)
[![Slack](https://slack.equinixmetal.com/badge.svg)](https://slack.equinixmetal.com/)
[![Twitter Follow](https://img.shields.io/twitter/follow/equinixmetal.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=equinixmetal)

A Golang client for the Equinix Metal API. ([Packet is now Equinix Metal](https://blog.equinix.com/blog/2020/10/06/equinix-metal-metal-and-more/))

## Installation

To import this library into your Go project:

```go
import "github.com/packethost/packngo"
```

Reference a particular version with:

```sh
go get github.com/packethost/packngo@v0.2.0
```

## Stability and Compatibility

This repository is [Maintained](https://github.com/packethost/standards/blob/master/maintained-statement.md) meaning that this software is supported by Equinix Metal and its community - available to use in production environments.

Packngo is currently provided with a major version of [v0](https://blog.golang.org/v2-go-modules). We'll try to avoid breaking changes to this library, but they will certainly happen as we work towards a stable v1 library. See [CHANGELOG.md](CHANGELOG.md) for details on the latest additions, removals, fixes, and breaking changes.

While packngo provides an interface to most of the [Equinix Metal API](https://metal.equinix.com/developers/api/), the API is regularly adding new features. To request or contribute support for more API end-points or added fields, [create an issue](https://github.com/packethost/packngo/issues/new).

See [SUPPORT.md](SUPPORT.md) for any other issues.

## Usage

To authenticate to the Equinix Metal API, you must have your API token exported in env var `PACKET_AUTH_TOKEN`.

This code snippet initializes Equinix Metal API client, and lists your Projects:

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

This library is used by the official [terraform-provider-packet](https://github.com/packethost/terraform-provider-packet).

You can also learn a lot from the `*_test.go` sources. Almost all out tests touch the Equinix Metal API, so you can see how auth, querying and POSTing works. For example [devices_test.go](devices_test.go).

<details>
<summary>Linked Resources</summary>

### Linked resources in Get\* and List\* functions

The Equinix Metal API includes references to related entities for a wide selection of resource types, indicated by `href` fields. The Equinix Metal API allows for these entities to be included in the API response, saving the user from making more round-trip API requests. This is useful for linked resources, e.g members of a project, devices in a project. Similarly, by excluding entities that are included by default, you can reduce the API response time and payload size.

Control of this behavior is provided through [common attributes](https://metal.equinix.com/developers/api/common-parameters/) that can be used to toggle, by field name, which referenced resources will be included as values in API responses. The API exposes this feature through `?include=` and `?exclude=` query parameters which accept a comma-separated list of field names. These field names can be dotted to reference nested entities.

Most of the packngo `Get` functions take references to `GetOptions` parameters (or `ListOptions` for `List` functions). These types include an `Include` and `Exclude` slice that will be converted to query parameters upon request.

For example, if you want to list users in a project, you can fetch the project via `Projects.Get(pid, nil)` call. The result of this call will be a `Project` which has a `Users []User` attribute. The items in the `[]User` slice only have a non-zero URL attribute, the rest of the fields will be type defaults. You can then parse the ID of the User resources and fetch them consequently.

Optionally, you can use the ListOptions struct in the project fetch call to include the Users (`members` JSON tag).  Then, every item in the `[]User` slice will have all (not only the `Href`) attributes populated.

```go
Projects.Get(pid, &packngo.ListOptions{Includes: []{'members'}})
```

The following is a more comprehensive illustration of Includes and Excludes.

```go
import (
	"log"

	"github.com/packethost/packngo"
)

func listProjectsAndUsers(lo *packngo.ListOptions) {
	c, err := packngo.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	ps, _, err := c.Projects.List(lo)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listing for listOptions %+v\n", lo)
	for _, p := range ps {
		log.Printf("project resource %s has %d users", p.Name, len(p.Users))
		for _, u := range p.Users {
			if u.Email != "" && u.FullName != "" {
				log.Printf("  user %s has email %s\n", u.FullName, u.Email)
			} else {
				log.Printf("  only got user link %s\n", u.URL)
			}
		}
	}
}

func main() {
	loMembers := &packngo.ListOptions{Includes: []string{"members"}}
	loMembersOut := &packngo.ListOptions{Excludes: []string{"members"}}
	listProjectsAndUsers(loMembers)
	listProjectsAndUsers(nil)
	listProjectsAndUsers(loMembersOut)
}
```

</details>

## Contributing

See [CONTIBUTING.md](CONTRIBUTING.md).
