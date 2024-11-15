# NOTE: This module will go out of support by March 31, 2023.  For authenticating with Azure AD, use module [azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity) instead.  For help migrating from `auth` to `azidentiy` please consult the [migration guide](https://aka.ms/azsdk/go/identity/migration).  General information about the retirement of this and other legacy modules can be found [here](https://azure.microsoft.com/updates/support-for-azure-sdk-libraries-that-do-not-conform-to-our-current-azure-sdk-guidelines-will-be-retired-as-of-31-march-2023/).

## Authentication

Typical SDK operations must be authenticated and authorized. The `autorest.Authorizer`
interface allows use of any auth style in requests, such as inserting an OAuth2
Authorization header and bearer token received from Azure AD.

The SDK itself provides a simple way to get an authorizer which first checks
for OAuth client credentials in environment variables and then falls back to
Azure's [Managed Service Identity]() when available, e.g. when on an Azure
VM. The following snippet from [the previous section](#use) demonstrates
this helper.

```go
import "github.com/Azure/go-autorest/autorest/azure/auth"

// create a VirtualNetworks client
vnetClient := network.NewVirtualNetworksClient("<subscriptionID>")

// create an authorizer from env vars or Azure Managed Service Idenity
authorizer, err := auth.NewAuthorizerFromEnvironment()
if err != nil {
    handle(err)
}

vnetClient.Authorizer = authorizer

// call the VirtualNetworks CreateOrUpdate API
vnetClient.CreateOrUpdate(context.Background(),
// ...
```

The following environment variables help determine authentication configuration:

- `AZURE_ENVIRONMENT`: Specifies the Azure Environment to use. If not set, it
  defaults to `AzurePublicCloud`. Not applicable to authentication with Managed
  Service Identity (MSI).
- `AZURE_AD_RESOURCE`: Specifies the AAD resource ID to use. If not set, it
  defaults to `ResourceManagerEndpoint` for operations with Azure Resource
  Manager. You can also choose an alternate resource programmatically with
  `auth.NewAuthorizerFromEnvironmentWithResource(resource string)`.

### More Authentication Details

The previous is the first and most recommended of several authentication
options offered by the SDK because it allows seamless use of both service
principals and [Azure Managed Service Identity][]. Other options are listed
below.

> Note: If you need to create a new service principal, run `az ad sp create-for-rbac -n "<app_name>"` in the
> [azure-cli](https://github.com/Azure/azure-cli). See [these
> docs](https://docs.microsoft.com/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest)
> for more info. Copy the new principal's ID, secret, and tenant ID for use in
> your app, or consider the `--sdk-auth` parameter for serialized output.

[azure managed service identity]: https://docs.microsoft.com/azure/active-directory/msi-overview

- The `auth.NewAuthorizerFromEnvironment()` described above creates an authorizer
  from the first available of the following configuration:

      1. **Client Credentials**: Azure AD Application ID and Secret.

          - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
          - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
          - `AZURE_CLIENT_SECRET`: Specifies the app secret to use.

      2. **Client Certificate**: Azure AD Application ID and X.509 Certificate.

          - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
          - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
          - `AZURE_CERTIFICATE_PATH`: Specifies the certificate Path to use.
          - `AZURE_CERTIFICATE_PASSWORD`: Specifies the certificate password to use.

      3. **Resource Owner Password**: Azure AD User and Password. This grant type is *not
         recommended*, use device login instead if you need interactive login.

          - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
          - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
          - `AZURE_USERNAME`: Specifies the username to use.
          - `AZURE_PASSWORD`: Specifies the password to use.

      4. **Azure Managed Service Identity**: Delegate credential management to the
         platform. Requires that code is running in Azure, e.g. on a VM. All
         configuration is handled by Azure. See [Azure Managed Service
         Identity](https://docs.microsoft.com/azure/active-directory/msi-overview)
         for more details.

- The `auth.NewAuthorizerFromFile()` method creates an authorizer using
  credentials from an auth file created by the [Azure CLI][]. Follow these
  steps to utilize:

  1. Create a service principal and output an auth file using `az ad sp create-for-rbac --sdk-auth > client_credentials.json`.
  2. Set environment variable `AZURE_AUTH_LOCATION` to the path of the saved
     output file.
  3. Use the authorizer returned by `auth.NewAuthorizerFromFile()` in your
     client as described above.

- The `auth.NewAuthorizerFromCLI()` method creates an authorizer which
  uses [Azure CLI][] to obtain its credentials.
  
  The default audience being requested is `https://management.azure.com` (Azure ARM API).
  To specify your own audience, export `AZURE_AD_RESOURCE` as an evironment variable.
  This is read by `auth.NewAuthorizerFromCLI()` and passed to Azure CLI to acquire the access token.
  
  For example, to request an access token for Azure Key Vault, export
  ```
  AZURE_AD_RESOURCE="https://vault.azure.net"
  ```
  
- `auth.NewAuthorizerFromCLIWithResource(AUDIENCE_URL_OR_APPLICATION_ID)` - this method is self contained and does
  not require exporting environment variables. For example, to request an access token for Azure Key Vault:
  ```
  auth.NewAuthorizerFromCLIWithResource("https://vault.azure.net")
  ```

  To use `NewAuthorizerFromCLI()` or `NewAuthorizerFromCLIWithResource()`, follow these steps:

  1. Install [Azure CLI v2.0.12](https://docs.microsoft.com/cli/azure/install-azure-cli) or later. Upgrade earlier versions.
  2. Use `az login` to sign in to Azure.

  If you receive an error, use `az account get-access-token` to verify access.

  If Azure CLI is not installed to the default directory, you may receive an error
  reporting that `az` cannot be found.  
  Use the `AzureCLIPath` environment variable to define the Azure CLI installation folder.

  If you are signed in to Azure CLI using multiple accounts or your account has
  access to multiple subscriptions, you need to specify the specific subscription
  to be used. To do so, use:

  ```
  az account set --subscription <subscription-id>
  ```

  To verify the current account settings, use:

  ```
  az account list
  ```

[azure cli]: https://github.com/Azure/azure-cli

- Finally, you can use OAuth's [Device Flow][] by calling
  `auth.NewDeviceFlowConfig()` and extracting the Authorizer as follows:

  ```go
  config := auth.NewDeviceFlowConfig(clientID, tenantID)
  a, err := config.Authorizer()
  ```

[device flow]: https://oauth.net/2/device-flow/
