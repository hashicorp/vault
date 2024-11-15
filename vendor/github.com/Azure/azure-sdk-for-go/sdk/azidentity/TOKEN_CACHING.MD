## Token caching in the Azure Identity client module

Token caching helps apps:

- Improve their resilience and performance.
- Reduce the number of requests sent to Microsoft Entra ID to obtain access tokens.
- Reduce the number of times users are prompted to authenticate.

When an app needs to access a protected Azure resource, it typically needs to obtain an access token from Entra ID by sending an HTTP request and sometimes prompting a user to authenticate interactively. Credentials with caches (see [the below table](#credentials-supporting-token-caching) for a list) store access tokens either [in memory](#in-memory-token-caching) or, optionally, [on disk](#persistent-token-caching). These credentials return cached tokens whenever possible, to avoid unnecessary token requests or user interaction. Both cache implementations are safe for concurrent use.

#### Caching can't be disabled

Whether a credential caches tokens isn't configurable. If a credential has a cache of either kind, it requests a new token only when it can't provide one from its cache. Azure SDK service clients have an additional, independent layer of in-memory token caching, to prevent redundant token requests. This cache works with any credential type, even a custom implementation defined outside the Azure SDK, and can't be disabled. Disabling token caching is therefore impossible when using Azure SDK clients or most `azidentity` credential types. However, in-memory caches can be cleared by constructing new credential and client instances.

### In-memory token caching

Credential types that support caching store tokens in memory by default and require no configuration to do so. Each instance of these types has its own cache, and two credential instances never share an in-memory cache.

### Persistent token caching

Some credential types support opt-in persistent token caching (see [the below table](#credentials-supporting-token-caching) for a list). This feature enables credentials to store and retrieve tokens across process executions, so an application doesn't need to authenticate every time it runs.

Persistent caches are encrypted at rest using a mechanism that depends on the operating system:

| Operating system | Encryption facility                   |
|------------------|---------------------------------------|
| Linux            | kernel key retention service (keyctl) |
| macOS            | Keychain                              |
| Windows          | Data Protection API (DPAPI)           |

Persistent caching requires encryption. When the required encryption facility is unuseable, or the application is running on an unsupported OS, the persistent cache constructor returns an error. This doesn't mean that authentication is impossible, only that credentials can't persist authentication data and the application will need to reauthenticate the next time it runs. See the [package documentation][example] for example code showing how to configure persistent caching and access cached data.

### Credentials supporting token caching

The following table indicates the state of in-memory and persistent caching in each credential type.

**Note:** in-memory caching is enabled by default for every type supporting it. Persistent token caching must be enabled explicitly. See the [package documentation][user_example] for an example showing how to do this for credential types authenticating users. For types that authenticate service principals, set the `Cache` field on the constructor's options as shown in [this example][sp_example].

| Credential                     | In-memory token caching                                             | Persistent token caching |
|--------------------------------|---------------------------------------------------------------------|--------------------------|
| `AzureCLICredential`           | Not Supported                                                       | Not Supported            |
| `AzureDeveloperCLICredential`  | Not Supported                                                       | Not Supported            |
| `AzurePipelinesCredential`     | Supported                                                           | Supported                |
| `ClientAssertionCredential`    | Supported                                                           | Supported                |
| `ClientCertificateCredential`  | Supported                                                           | Supported                |
| `ClientSecretCredential`       | Supported                                                           | Supported                |
| `DefaultAzureCredential`       | Supported if the target credential in the default chain supports it | Not Supported            |
| `DeviceCodeCredential`         | Supported                                                           | Supported                |
| `EnvironmentCredential`        | Supported                                                           | Not Supported            |
| `InteractiveBrowserCredential` | Supported                                                           | Supported                |
| `ManagedIdentityCredential`    | Supported                                                           | Not Supported            |
| `OnBehalfOfCredential`         | Supported                                                           | Not Supported            |
| `UsernamePasswordCredential`   | Supported                                                           | Supported                |
| `WorkloadIdentityCredential`   | Supported                                                           | Supported                |

[sp_example]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#example-package-PersistentServicePrincipalAuthentication
[user_example]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#example-package-PersistentUserAuthentication
