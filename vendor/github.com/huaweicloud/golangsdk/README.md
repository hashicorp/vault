# Golangsdk: a Huawei clouds SDK for Golang
[![Go Report Card](https://goreportcard.com/badge/github.com/huaweicloud/golangsdk?branch=master)](https://goreportcard.com/badge/github.com/huaweicloud/golangsdk)
[![Build Status](https://travis-ci.org/huaweicloud/golangsdk.svg?branch=master)](https://travis-ci.org/huaweicloud/golangsdk)
[![Coverage Status](https://coveralls.io/repos/github/huaweicloud/golangsdk/badge.svg?branch=master)](https://coveralls.io/github/huaweicloud/golangsdk?branch=master)
[![LICENSE](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://github.com/huaweicloud/golangsdk/blob/master/LICENSE)

Golangsdk is a Huawei clouds Go SDK.
Golangsdk is based on [Gophercloud](https://github.com/gophercloud/gophercloud)
which is an OpenStack Go SDK and has a great design.
Golangsdk has added and removed some features to support Huawei clouds.

## Added features

- **autoscaling**: auto scaling service
- **kms**: key management service
- **rds**: relational database service
- **smn**: simple message notification service
- **drs**: disaster recovery service

## Removed features

- **blockstorage**: block storage service
- **cdn**: content delivery network service
- **compute**: compute service
- **db**: database service
- **imageservice**: image service
- **objectstorage**: object storage service
- **orchestration**: orchestration service
- **sharedfilesystems**: share file system service

## Useful links

* [Reference documentation](http://godoc.org/github.com/huaweicloud/golangsdk)
* [Effective Go](https://golang.org/doc/effective_go.html)

## How to install

Before installing, you need to ensure that your [GOPATH environment variable](https://golang.org/doc/code.html#GOPATH)
is pointing to an appropriate directory where you want to install Golangsdk:

```bash
mkdir $HOME/go
export GOPATH=$HOME/go
```

To protect yourself against changes in your dependencies, we highly recommend choosing a
[dependency management solution](https://github.com/golang/go/wiki/PackageManagementTools) for
your projects, such as [godep](https://github.com/tools/godep). Once this is set up, you can install
golangsdk as a dependency like so:

```bash
go get github.com/huaweicloud/golangsdk

# Edit your code to import relevant packages from "github.com/huaweicloud/golangsdk"

godep save ./...
```

This will install all the source files you need into a `Godeps/_workspace` directory, which is
referenceable from your own source files when you use the `godep go` command.

## Getting started

### Credentials

Because you'll be hitting an API, you will need to retrieve your Huawei clouds
credentials and either store them as environment variables or in your local Go
files. The first method is recommended because it decouples credential
information from source code, allowing you to push the latter to your version
control system without any security risk.

You will need to retrieve the following:

* username
* password
* a valid IAM identity URL

### Authentication

Once you have access to your credentials, you can begin plugging them into
Golangsdk. The next step is authentication, and this is handled by a base
"Provider" struct. To get one, you can either pass in your credentials
explicitly, or tell Golangsdk to use environment variables:

```go
import (
  "github.com/huaweicloud/golangsdk"
  "github.com/huaweicloud/golangsdk/openstack"
  "github.com/huaweicloud/golangsdk/openstack/utils"
)

// Option 1: Pass in the values yourself
opts := golangsdk.AuthOptions{
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

The `ProviderClient` is the top-level client that all of your Huawei clouds services
derive from. The provider contains all of the authentication details that allow
your Go code to access the API - such as the base URL and token ID.

### Provision a rds instance

Once we have a base Provider, we inject it as a dependency into each Huawei clouds
service. In order to work with the rds API, we need a rds service
client; which can be created like so:

```go
client, err := openstack.NewRdsServiceV1(provider, golangsdk.EndpointOpts{
  Region: os.Getenv("OS_REGION_NAME"),
})
```

We then use this `client` for any rds API operation we want. In our case,
we want to provision a rds instance - so we invoke the `Create` method and pass
in the name and the flavor ID (database specification) we're
interested in:

```go
import "github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"

instance, err := instances.Create(client, instances.CreateOpts{
  Name:      "My new rds instance!",
  FlavorRef: "flavor_id",
}).Extract()
```

The above code sample creates a new rds instance with the parameters, and embodies the
new resource in the `instance` variable (a
[`instances.Instance`](http://godoc.org/github.com/huaweicloud/golangsdk) struct).

## Advanced Usage

Have a look at the [FAQ](./FAQ.md) for some tips on customizing the way Golangsdk works.

## Backwards-Compatibility Guarantees

None. Vendor it and write tests covering the parts you use.

## Contributing

See the [contributing guide](./.github/CONTRIBUTING.md).

## Help and feedback

If you're struggling with something or have spotted a potential bug, feel free
to submit an issue to our [bug tracker](https://github.com/huaweicloud/golangsdk/issues).

## Thank You

We'd like to extend special thanks and appreciation to the following:

### OpenLab

<a href="http://openlabtesting.org/"><img src="assets/openlab.png" width="600px"></a>

OpenLab is providing a full CI environment to test each PR and merge for a variety of OpenStack releases.

### VEXXHOST

<a href="https://vexxhost.com/"><img src="assets/vexxhost.png" width="600px"></a>

VEXXHOST is providing their services to assist with the development and testing of Golangsdk.

## License

Golangsdk is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.

