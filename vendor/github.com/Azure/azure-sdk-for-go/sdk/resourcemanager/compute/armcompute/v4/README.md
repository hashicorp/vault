# Azure Compute Module for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4)](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4)

The `armcompute` module provides operations for working with Azure Compute.

[Source code](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/resourcemanager/compute/armcompute)

# Getting started

## Prerequisites

- an [Azure subscription](https://azure.microsoft.com/free/)
- Go 1.18 or above (You could download and install the latest version of Go from [here](https://go.dev/doc/install). It will replace the existing Go on your machine. If you want to install multiple Go versions on the same machine, you could refer this [doc](https://go.dev/doc/manage-install).)

## Install the package

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

Install the Azure Compute module:

```sh
go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4
```

## Authorization

When creating a client, you will need to provide a credential for authenticating with Azure Compute.  The `azidentity` module provides facilities for various ways of authenticating with Azure including client/secret, certificate, managed identity, and more.

```go
cred, err := azidentity.NewDefaultAzureCredential(nil)
```

For more information on authentication, please see the documentation for `azidentity` at [pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity).

## Client Factory

Azure Compute module consists of one or more clients. We provide a client factory which could be used to create any client in this module.

```go
clientFactory, err := armcompute.NewClientFactory(<subscription ID>, cred, nil)
```

You can use `ClientOptions` in package `github.com/Azure/azure-sdk-for-go/sdk/azcore/arm` to set endpoint to connect with public and sovereign clouds as well as Azure Stack. For more information, please see the documentation for `azcore` at [pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore).

```go
options := arm.ClientOptions {
    ClientOptions: azcore.ClientOptions {
        Cloud: cloud.AzureChina,
    },
}
clientFactory, err := armcompute.NewClientFactory(<subscription ID>, cred, &options)
```

## Clients

A client groups a set of related APIs, providing access to its functionality.  Create one or more clients to access the APIs you require using client factory.

```go
client := clientFactory.NewLogAnalyticsClient()
```

## More sample code

- [Availability Set](https://aka.ms/azsdk/go/mgmt/samples?path=sdk/resourcemanager/compute/availabilityset)
- [Virtual Machine](https://aka.ms/azsdk/go/mgmt/samples?path=sdk/resourcemanager/compute/create_vm)
- [Dedicated Host](https://aka.ms/azsdk/go/mgmt/samples?path=sdk/resourcemanager/compute/dedicated_host)
- [Disk](https://aka.ms/azsdk/go/mgmt/samples?path=sdk/resourcemanager/compute/disk)
- [Gallery](https://aka.ms/azsdk/go/mgmt/samples?path=sdk/resourcemanager/compute/gallery)
- [Proximity Placement Group](https://aka.ms/azsdk/go/mgmt/samples?path=sdk/resourcemanager/compute/proximity)
- [Snapshot](https://aka.ms/azsdk/go/mgmt/samples?path=sdk/resourcemanager/compute/snapshot)
- [Virtual Machine Scale Set](https://aka.ms/azsdk/go/mgmt/samples?path=sdk/resourcemanager/compute/vmscaleset)

## Major Version Upgrade

Go uses [semantic import versioning](https://github.com/golang/go/wiki/Modules#semantic-import-versioning) to ensure a good backward compatibility for modules. For Azure Go management SDK, we usually upgrade module version according to cooresponding service's API version. Regarding it could be a complicated experience for major version upgrade, we will try our best to keep the SDK API stable and release new version in backward compatible way. However, if any unavoidable breaking changes and a new major version releases for SDK modules, you could use these commands under your module folder to upgrade:

```sh
go install github.com/icholy/gomajor@latest
gomajor get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute@latest
```

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Compute` label.

# Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit [https://cla.microsoft.com](https://cla.microsoft.com).

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the
[Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information, see the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/)
or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any
additional questions or comments.
