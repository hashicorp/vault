# Microsoft Graph Core SDK for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/microsoftgraph/msgraph-sdk-go-core/)](https://pkg.go.dev/github.com/microsoftgraph/msgraph-sdk-go-core/) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=microsoftgraph_msgraph-sdk-go-core&metric=coverage)](https://sonarcloud.io/dashboard?id=microsoftgraph_msgraph-sdk-go-core) [![Sonarcloud Status](https://sonarcloud.io/api/project_badges/measure?project=microsoftgraph_msgraph-sdk-go-core&metric=alert_status)](https://sonarcloud.io/dashboard?id=microsoftgraph_msgraph-sdk-go-core)

Get started with the Microsoft Graph Core SDK for Go by integrating the [Microsoft Graph API](https://docs.microsoft.com/graph/overview) into your Go application! You can also have a look at the [Go documentation](https://pkg.go.dev/github.com/microsoftgraph/msgraph-sdk-go-core/)

> **Note:** Although you can use this library directly, we recommend you use the [v1](https://github.com/microsoftgraph/msgraph-sdk-go) or [beta](https://github.com/microsoftgraph/msgraph-beta-sdk-go) library which rely on this library and additionally provide a fluent style Go API and models.
>
> **Note:** The Microsoft Graph Go SDK is currently in Release Candidate (RC) version starting from version 0.34.1. The SDK is still undergoing testing but minimum breaking changes should be expected. Checkout the [known limitations](https://github.com/microsoftgraph/msgraph-sdk-go-core/issues/1).

## Samples and usage guide

- [Middleware usage](https://github.com/microsoftgraph/msgraph-sdk-design/)

## 1. Installation

```Shell
go get github.com/microsoftgraph/msgraph-sdk-go-core
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

## 2. Getting started

### 2.1 Register your application

Register your application by following the steps at [Register your app with the Microsoft Identity Platform](https://docs.microsoft.com/graph/auth-register-app-v2).

### 2.2 Create an AuthenticationProvider object

An instance of the **GraphRequestAdapterBase** class handles building client. To create a new instance of this class, you need to provide an instance of **AuthenticationProvider**, which can authenticate requests to Microsoft Graph.

For an example of how to get an authentication provider, see [choose a Microsoft Graph authentication provider](https://docs.microsoft.com/graph/sdks/choose-authentication-providers?tabs=Go).

> Note: we are working to add the getting started information for Go to our public documentation, in the meantime the following sample should help you getting started.

```Golang
import (
    azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    a          "github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
    "context"
)

cred, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
    TenantID: "<the tenant id from your app registration>",
    ClientID: "<the client id from your app registration>",
    UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
        fmt.Println(message.Message)
        return nil
    },
})

if err != nil {
    fmt.Printf("Error creating credentials: %v\n", err)
}

auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"Mail.Read", "Mail.Send"})
if err != nil {
    fmt.Printf("Error authentication provider: %v\n", err)
    return
}

```

### 2.3 Get a Request Adapter object

You must get a **GraphRequestAdapterBase** object to make requests against the service.

```Golang
import core "github.com/microsoftgraph/msgraph-sdk-go-core"

adapter, err := core.NewGraphRequestAdapterBase(auth)
if err != nil {
    fmt.Printf("Error creating adapter: %v\n", err)
    return
}
```

## 3. Make requests against the service

After you have a **GraphRequestAdapterBase** that is authenticated, you can begin making calls against the service. The requests against the service look like our [REST API](https://docs.microsoft.com/graph/api/overview?view=graph-rest-1.0).

### 3.1 Get the user's details

To retrieve the user's details

```Golang
import abs "github.com/microsoft/kiota-abstractions-go"

requestInf := abs.NewRequestInformation()
targetUrl, err := url.Parse("https://graph.microsoft.com/v1.0/me")
if err != nil {
    fmt.Printf("Error parsing URL: %v\n", err)
}
requestInf.SetUri(*targetUrl)

// User is your own type that implements Parsable or comes from the service library
user, err := adapter.SendAsync(*requestInf, func() { return &User }, nil)

if err != nil {
    fmt.Printf("Error getting the user: %v\n", err)
}
```

## 4. Issues

For known issues, see [issues](https://github.com/MicrosoftGraph/msgraph-sdk-go-core/issues).

## 5. Contributions

The Microsoft Graph SDK is open for contribution. To contribute to this project, see [Contributing](https://github.com/microsoftgraph/msgraph-sdk-go-core/blob/main/CONTRIBUTING.md).

## 6. License

Copyright (c) Microsoft Corporation. All Rights Reserved. Licensed under the [MIT license](LICENSE).

## 7. Third-party notices

[Third-party notices](THIRD%20PARTY%20NOTICES)
