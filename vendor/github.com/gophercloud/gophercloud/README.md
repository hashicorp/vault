# Gophercloud: an OpenStack SDK for Go
[![Coverage Status](https://coveralls.io/repos/github/gophercloud/gophercloud/badge.svg?branch=v1)](https://coveralls.io/github/gophercloud/gophercloud?branch=v1)

Gophercloud is an OpenStack Go SDK.

## Useful links

* [Reference documentation](http://godoc.org/github.com/gophercloud/gophercloud)
* [Effective Go](https://golang.org/doc/effective_go.html)

## How to install

Reference a Gophercloud package in your code:

```go
import "github.com/gophercloud/gophercloud"
```

Then update your `go.mod`:

```shell
go mod tidy
```

## Getting started

### Credentials

Because you'll be hitting an API, you will need to retrieve your OpenStack
credentials and either store them in a `clouds.yaml` file, as environment
variables, or in your local Go files. The first method is recommended because
it decouples credential information from source code, allowing you to push the
latter to your version control system without any security risk.

You will need to retrieve the following:

* A valid Keystone identity URL
* Credentials. These can be a username/password combo, a set of Application
  Credentials, a pre-generated token, or any other supported authentication
  mechanism.

For users who have the OpenStack dashboard installed, there's a shortcut. If
you visit the `project/api_access` path in Horizon and click on the
"Download OpenStack RC File" button at the top right hand corner, you can
download either a `clouds.yaml` file or an `openrc` bash file that exports all
of your access details to environment variables. To use the `clouds.yaml` file,
place it at `~/.config/openstack/clouds.yaml`. To use the `openrc` file, run
`source openrc` and you will be prompted for your password.

### Authentication

Once you have access to your credentials, you can begin plugging them into
Gophercloud. The next step is authentication, which is handled by a base
"Provider" struct. There are number of ways to construct such a struct.

**With `gophercloud/utils`**

The [github.com/gophercloud/utils](https://github.com/gophercloud/utils)
library provides the `clientconfig` package to simplify authentication. It
provides additional functionality, such as the ability to read `clouds.yaml`
files. To generate a "Provider" struct using the `clientconfig` package:

```go
import (
	"github.com/gophercloud/utils/openstack/clientconfig"
)

// You can also skip configuring this and instead set 'OS_CLOUD' in your
// environment
opts := new(clientconfig.ClientOpts)
opts.Cloud = "devstack-admin"

provider, err := clientconfig.AuthenticatedClient(opts)
```

A provider client is a top-level client that all of your OpenStack service
clients derive from. The provider contains all of the authentication details
that allow your Go code to access the API - such as the base URL and token ID.

Once we have a base Provider, we inject it as a dependency into each OpenStack
service. For example, in order to work with the Compute API, we need a Compute
service client. This can be created like so:

```go
client, err := clientconfig.NewServiceClient("compute", opts)
```

**Without `gophercloud/utils`**

> *Note*
> gophercloud doesn't provide support for `clouds.yaml` file so you need to
> implement this functionality yourself if you don't wish to use
> `gophercloud/utils`.

You can also generate a "Provider" struct without using the `clientconfig`
package from `gophercloud/utils`. To do this, you can either pass in your
credentials explicitly or tell Gophercloud to use environment variables:

```go
import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

// Option 1: Pass in the values yourself
opts := gophercloud.AuthOptions{
  IdentityEndpoint: "https://openstack.example.com:5000/v2.0",
  Username: "{username}",
  Password: "{password}",
}

// Option 2: Use a utility function to retrieve all your environment variables
opts, err := openstack.AuthOptionsFromEnv()
```

Once you have the `opts` variable, you can pass it in and get back a
`ProviderClient` struct:

```go
provider, err := openstack.AuthenticatedClient(opts)
```

As above, you can then use this provider client to generate a service client
for a particular OpenStack service:

```go
client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
	Region: os.Getenv("OS_REGION_NAME"),
})
```

### Provision a server

We can use the Compute service client generated above for any Compute API
operation we want. In our case, we want to provision a new server. To do this,
we invoke the `Create` method and pass in the flavor ID (hardware
specification) and image ID (operating system) we're interested in:

```go
import "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

server, err := servers.Create(client, servers.CreateOpts{
	Name:      "My new server!",
	FlavorRef: "flavor_id",
	ImageRef:  "image_id",
}).Extract()
```

The above code sample creates a new server with the parameters, and embodies the
new resource in the `server` variable (a
[`servers.Server`](http://godoc.org/github.com/gophercloud/gophercloud) struct).

## Advanced Usage

Have a look at the [FAQ](./docs/FAQ.md) for some tips on customizing the way Gophercloud works.

## Backwards-Compatibility Guarantees

Gophercloud versioning follows [semver](https://semver.org/spec/v2.0.0.html).

Before `v1.0.0`, there were no guarantees. Starting with v1, there will be no breaking changes within a major release.

See the [Release instructions](./RELEASE.md).

## Contributing

See the [contributing guide](./.github/CONTRIBUTING.md).

## Help and feedback

If you're struggling with something or have spotted a potential bug, feel free
to submit an issue to our [bug tracker](https://github.com/gophercloud/gophercloud/issues).
