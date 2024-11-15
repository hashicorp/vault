# Kiota Http Library for Go

![Go](https://github.com/microsoft/kiota-http-go/actions/workflows/go.yml/badge.svg)

The Kiota HTTP Library for Go is the Go HTTP library implementation with [net/http](https://pkg.go.dev/net/http).

A [Kiota](https://github.com/microsoft/kiota) generated project will need a reference to a HTTP package to make HTTP requests to an API endpoint.

Read more about Kiota [here](https://github.com/microsoft/kiota/blob/main/README.md).

## Using the Kiota Http Library for Go

```Shell
go get github.com/microsoft/kiota-http-go
```

```Golang
httpAdapter, err := kiotahttp.NewNetHttpRequestAdapter(authProvider)
```

## Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft 
trademarks or logos is subject to and must follow 
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.
