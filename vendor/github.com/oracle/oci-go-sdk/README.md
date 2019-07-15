# Oracle Cloud Infrastructure Golang SDK
[![wercker status](https://app.wercker.com/status/09bc4818e7b1d70b04285331a9bdbc41/s/master "wercker status")](https://app.wercker.com/project/byKey/09bc4818e7b1d70b04285331a9bdbc41)

This is the Go SDK for Oracle Cloud Infrastructure. This project is open source and maintained by Oracle Corp. 
The home page for the project is [here](https://godoc.org/github.com/oracle/oci-go-sdk/).
>***WARNING:***: To avoid automatically consuming breaking changes if we have to rev the major version of the Go SDK, 
please consider using the  [Go dependency management tool](https://github.com/golang/dep), or vendoring the SDK. 
This will allow you to pin to a specific version of the Go SDK in your project, letting you control how and when you move to the next major version.

## Dependencies
- Install [Go programming language](https://golang.org/dl/).
- Install [GNU Make](https://www.gnu.org/software/make/), using the package manager or binary distribution tool appropriate for your platform.
 


## Installing
Use the following command to install this SDK:

```
go get -u github.com/oracle/oci-go-sdk
```
Alternatively you can git clone this repo.

## Working with the Go SDK
To start working with the Go SDK, you import the service package, create a client, and then use that client to make calls.

### Configuring 
Before using the SDK, set up a config file with the required credentials. See [SDK and Tool Configuration](https://docs.us-phoenix-1.oraclecloud.com/Content/API/Concepts/sdkconfig.htm) for instructions.

Note that the Go SDK does not support profile inheritance or defining custom values in the configuration file.

Once a config file has been setup, call `common.DefaultConfigProvider()` function as follows:

 ```go
 // Import necessary packages
 import (
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/identity" // Identity or any other service you wish to make requests to
)
 
 //...
 
configProvider := common.DefaultConfigProvider()
```

 Or, to configure the SDK programmatically instead, implement the `ConfigurationProvider` interface shown below:
 ```go
// ConfigurationProvider wraps information about the account owner
type ConfigurationProvider interface {
	KeyProvider
	TenancyOCID() (string, error)
	UserOCID() (string, error)
	KeyFingerprint() (string, error)
	Region() (string, error)
}
```
Or simply use one of  structs exposed by the `oci-go-sdk` that already implement the above [interface](https://godoc.org/github.com/oracle/oci-go-sdk/common#ConfigurationProvider)

### Making a Request
To make a request to an Oracle Cloud Infrastructure service, create a client for the service and then use the client to call a function from the service.

- *Creating a client*: All packages provide a function to create clients, using the naming convention `New<ServiceName>ClientWithConfigurationProvider`,
such as `NewVirtualNetworkClientWithConfigurationProvider` or `NewIdentityClientWithConfigurationProvider`. To create a new client, 
pass a struct that conforms to the `ConfigurationProvider` interface, or use the `DefaultConfigProvider()` function in the common package.

For example: 
```go
config := common.DefaultConfigProvider()
client, err := identity.NewIdentityClientWithConfigurationProvider(config)
if err != nil { 
     panic(err)
}
```

- *Making calls*: After successfully creating a client, requests can now be made to the service. Generally all functions associated with an operation
accept [`context.Context`](https://golang.org/pkg/context/) and a struct that wraps all input parameters. The functions then return a response struct
that contains the desired data, and an error struct that describes the error if an error occurs.

For example:
```go
id := "your_group_id"
response, err := client.GetGroup(context.Background(), identity.GetGroupRequest{GroupId:&id})
if err != nil {
	//Something happened
	panic(err)
}
//Process the data in response struct
fmt.Println("Group's name is:", response.Name)
```

## Organization of the SDK
The `oci-go-sdk` contains the following:
- **Service packages**: All packages except `common` and any other package found inside `cmd`. These packages represent 
the Oracle Cloud Infrastructure services supported by the Go SDK. Each package represents a service. 
These packages include methods to interact with the service, structs that model 
input and output parameters, and a client struct that acts as receiver for the above methods.

- **Common package**: Found in the `common` directory. The common package provides supporting functions and structs used by service packages.
Includes HTTP request/response (de)serialization, request signing, JSON parsing, pointer to reference and other helper functions. Most of the functions
in this package are meant to be used by the service packages.

- **cmd**: Internal tools used by the `oci-go-sdk`.

## Examples
Examples can be found [here](https://github.com/oracle/oci-go-sdk/tree/master/example)

## Documentation
Full documentation can be found [on the godocs site](https://godoc.org/github.com/oracle/oci-go-sdk/).

## Help
* The [Issues](https://github.com/oracle/oci-go-sdk/issues) page of this GitHub repository.
* [Stack Overflow](https://stackoverflow.com/), use the [oracle-cloud-infrastructure](https://stackoverflow.com/questions/tagged/oracle-cloud-infrastructure) and [oci-go-sdk](https://stackoverflow.com/questions/tagged/oci-go-sdk) tags in your post.
* [Developer Tools](https://community.oracle.com/community/cloud_computing/bare-metal/content?filterID=contentstatus%5Bpublished%5D~category%5Bdeveloper-tools%5D&filterID=contentstatus%5Bpublished%5D~objecttype~objecttype%5Bthread%5D) of the Oracle Cloud forums.
* [My Oracle Support](https://support.oracle.com).


## Contributing
`oci-go-sdk` is an open source project. See [CONTRIBUTING](/CONTRIBUTING.md) for details.

Oracle gratefully acknowledges the contributions to oci-go-sdk that have been made by the community.


## License
Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

This SDK and sample is dual licensed under the Universal Permissive License 1.0 and the Apache License 2.0.

See [LICENSE](/LICENSE.txt) for more details.

## Changes
See [CHANGELOG](/CHANGELOG.md).

## Known Issues
You can find information on any known issues with the SDK here and under the [Issues](https://github.com/oracle/oci-go-sdk/issues) tab of this project's GitHub repository.

## Building and Testing
### Dev Dependencies
- Install [Testify](https://github.com/stretchr/testify) with the command:
```sh
go get github.com/stretchr/testify
```
- Install [go lint](https://github.com/golang/lint) with the command:
```
go get -u github.com/golang/lint/golint
```
### Build
Building is provided by the make file at the root of the project. To build the project execute.

```
make build
```

To run the tests:
```
make test
```
