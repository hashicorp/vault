## Unreleased

## v0.19.2

BUGS:

* fix a bug that prevented logins when validating vm names, vmss names, and resource groups  (https://github.com/hashicorp/vault-plugin-auth-azure/pull/172)

## v0.19.1

BUGS:

* fix an endless loop of warning spamming the logs on login error (https://github.com/hashicorp/vault-plugin-auth-azure/pull/170)

## v0.19.0

IMPROVEMENTS:
* Add login field validation for subscription id, resource group name, vmss name, and vm name
* Bump Go version to 1.22.6
* Updated dependencies:
  * `github.com/docker/docker` v25.0.5+incompatible -> v25.0.6+incompatible
  * `github.com/hashicorp/go-retryablehttp` v0.7.1 -> v0.7.7
  * `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.11.1 -> v1.14.0
  * `github.com/Azure/azure-sdk-for-go/sdk/azidentity` v1.5.2 -> v1.7.0
  * `github.com/coreos/go-oidc/v3` v3.10.0 -> v3.11.0
  * `github.com/hashicorp/vault/api` v1.13.0 -> v1.14.0
  * `github.com/hashicorp/vault/sdk` v0.12.0 -> v0.13.0
  * `github.com/microsoftgraph/msgraph-sdk-go` v1.42.0 -> v1.47.0
  * `github.com/microsoftgraph/msgraph-sdk-go-core` v1.1.0 -> v1.2.1
  * `golang.org/x/oauth2` v0.20.0 -> v0.23.0

## v0.18.1

BUGS:

* fix an endless loop of warning spamming the logs on login error

## v0.18.0

FEATURES:
* Add support for Workload Identify Federation (https://github.com/hashicorp/vault-plugin-auth-azure/pull/151)

IMPROVEMENTS:
* Bump github.com/coreos/go-oidc to v3 (https://github.com/hashicorp/vault-plugin-auth-azure/pull/157)
* Updated dependencies: (https://github.com/hashicorp/vault-plugin-auth-azure/pull/154)
* Updated dependencies:
  * `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.9.1 -> v1.10.0
  * `github.com/go-jose/go-jose/v3` v3.0.1 -> v3.0.3
  * `github.com/hashicorp/vault/api` v1.11.0 -> v1.12.0
  * `github.com/hashicorp/vault/sdk` v0.10.2 -> v0.11.0
  * `github.com/microsoftgraph/msgraph-sdk-go` v1.32.0 -> v1.35.0
  * `github.com/microsoftgraph/msgraph-sdk-go-core` v1.0.1 -> v1.1.0
  * `golang.org/x/oauth2` v0.16.0 -> v0.17.0

## v0.17.1

BUGS:

* fix an endless loop of warning spamming the logs on login error

## v0.17.0
IMPROVEMENTS:
* Make `framework.WALPrefix`` a local path [[GH-137](https://github.com/hashicorp/vault-plugin-auth-azure/pull/137)]
* Updated dependencies:
   * `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.7.1 -> v1.9.1
   * `github.com/Azure/azure-sdk-for-go/sdk/azidentity` v1.3.1 -> v1.5.1
   * `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi` v1.1.0 -> v1.2.0
   * `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources` v1.1.1 -> v1.2.0
   * `github.com/google/uuid` v1.3.1 -> v1.6.0
   * `github.com/hashicorp/go-hclog` v1.5.0 -> v1.6.2
   * `github.com/hashicorp/vault/api` v1.9.2 -> v1.11.0
   * `github.com/hashicorp/vault/sdk` v0.9.2 -> v0.10.2
   * `github.com/microsoftgraph/msgraph-sdk-go` v1.13.0 -> v1.32.0
   * `github.com/microsoftgraph/msgraph-sdk-go-core` v1.0.0 -> v1.0.1
   * `golang.org/x/oauth2` v0.11.0 -> v0.16.0

## v0.16.2
IMPROVEMENTS:
    * Prevent write-ahead-log data from being replicated to performance secondaries [GH-137](https://github.com/hashicorp/vault-plugin-auth-azure/pull/137)
    * Added Azure API configurable retry options [GH-133](https://github.com/hashicorp/vault-plugin-auth-azure/pull/133)

## v0.16.1
IMPROVEMENTS:
* Updated dependencies:
  * github.com/Azure/azure-sdk-for-go/sdk/azcore v1.7.1
  * github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.3.1
  * golang.org/x/oauth2 v0.11.0

## v0.16.0

FEATURES:
* Allow Workload Identity Federation based Azure resources to authenticate with Vault via appID

IMPROVEMENTS:
* Replaces the deprecated [go-autorest](https://github.com/Azure/go-autorest) client with [msgraph-sdk-go](https://github.com/microsoftgraph/msgraph-sdk-go) [[GH-121]](https://github.com/hashicorp/vault-plugin-auth-azure/pull/121)
* Updated dependencies:
  * `github.com/hashicorp/vault/sdk` v0.9.1 -> v0.9.2
  * `github.com/microsoftgraph/msgraph-sdk-go` v1.12.0 -> v1.13.0
  * `golang.org/x/oauth2` v0.9.0 -> v0.10.0

## v0.15.1

BUG FIXES:

* Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials [[GH-118]](https://github.com/hashicorp/vault-plugin-auth-azure/pull/118)

## v0.15.0

IMPROVEMENTS:
* Add display attributes for OpenAPI OperationID's [[GH-106](https://github.com/hashicorp/vault-plugin-auth-azure/pull/106)]
* Updated dependencies:
   * `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.4.0 -> v1.6.0
   * `github.com/Azure/azure-sdk-for-go/sdk/azidentity` v1.2.2 -> v1.3.0
   * `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4` v4.1.0 -> v4.2.1
   * `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi` v1.0.0 -> v1.1.0
   * `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources` v1.0.0 -> v1.1.1
   * `github.com/Azure/go-autorest/autorest` v0.11.28 -> v0.11.29
   * `github.com/hashicorp/vault/api` v1.9.0 -> v1.9.1
   * `github.com/hashicorp/vault/sdk` v0.8.1 -> v0.9.0
   * `golang.org/x/oauth2` v0.6.0 -> v0.8.0
   * `mvdan.cc/gofumpt` v0.3.1 -> v0.5.0
 * Downgraded dependencies:
   * `github.com/Azure/azure-sdk-for-go` v68.0.0+incompatible -> v67.2.0+incompatible

## v0.14.0

IMPROVEMENTS:

* Enable multiplexing [[GH-96](https://github.com/hashicorp/vault-plugin-auth-azure/pull/96)]
* Upgrade to Go 1.20.2
* Updated dependencies:
   * `github.com/Azure/azure-sdk-for-go` v67.2.0+incompatible -> v68.0.0+incompatible
   * `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.3.1 -> v1.4.0
   * `github.com/Azure/azure-sdk-for-go/sdk/azidentity` v1.2.1 -> v1.2.2
   * `github.com/hashicorp/go-hclog` v1.4.0 -> v1.5.0
   * `github.com/hashicorp/go-uuid` v1.0.2 -> v1.0.3
   * `golang.org/x/oauth2` v0.4.0 -> v0.6.0

## v0.13.1

BUG FIXES:

* Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials [[GH-118]](https://github.com/hashicorp/vault-plugin-auth-azure/pull/118)

## v0.13.0

FEATURES:

* Add rotate root support to Azure Auth [GH-88](https://github.com/hashicorp/vault-plugin-auth-azure/pull/88)
* Allow any Azure resource that supports managed identities to authenticate with Vault [GH-71](https://github.com/hashicorp/vault-plugin-auth-azure/pull/71)
* Adds support for Virtual Machine Scale Set Flex Authentication [GH-63](https://github.com/hashicorp/vault-plugin-auth-azure/pull/63)

IMPROVEMENTS:
* Updates dependencies
  * `github.com/Azure/azure-sdk-for-go v67.2.0+incompatible`[[GH-88](https://github.com/hashicorp/vault-plugin-auth-azure/pull/88)]
  * `github.com/Azure/go-autorest/autorest v0.11.28` [[GH-65](https://github.com/hashicorp/vault-plugin-auth-azure/pull/65)]
  * `github.com/hashicorp/vault/api v1.8.3` [[GH-82](https://github.com/hashicorp/vault-plugin-auth-azure/pull/82)]
  * `github.com/hashicorp/vault/sdk v0.7.0` [[GH-82](https://github.com/hashicorp/vault-plugin-auth-azure/pull/82)]
  * `github.com/hashicorp/go-hclog v1.4.0`[[GH-78](https://github.com/hashicorp/vault-plugin-auth-azure/pull/78)]
  * `golang.org/x/oauth2 v0.4.0`[[GH-84](https://github.com/hashicorp/vault-plugin-auth-azure/pull/84)]

* Upgrades to Go 1.19 [[GH-65](https://github.com/hashicorp/vault-plugin-auth-azure/pull/65)]
