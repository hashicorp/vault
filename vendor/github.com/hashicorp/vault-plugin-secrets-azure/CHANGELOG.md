## Unreleased

## v0.20.1
### October 14, 2024

IMPROVEMENTS:

* Prevent noisy logs for non-existent or deleted out-of-band errors (https://github.com/hashicorp/vault-plugin-secrets-azure/pull/220)

## v0.20.0
IMPROVEMENTS:
* Bump Go version to 1.22.6
* Updated dependencies [[GH-208]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/208):
  * `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.11.1 -> v1.14.0
  * `github.com/Azure/azure-sdk-for-go/sdk/azidentity` v1.6.0 -> v1.7.0
  * `github.com/go-test/deep` v1.1.0 -> v1.1.1
  * `github.com/hashicorp/vault/api` v1.13.0 -> v1.14.0
  * `github.com/hashicorp/vault/sdk` v0.12.0 -> v0.13.0
  * `github.com/microsoftgraph/msgraph-sdk-go` v1.42.0 -> v1.47.0
  * `github.com/microsoftgraph/msgraph-sdk-go-core` v1.1.0 -> v1.2.1
  * `github.com/docker/docker` v25.0.5+incompatible -> v25.0.6+incompatible [[GH-217]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/217)

FEATURES:
* Adds ability to limit the lifetime of service principal secrets in Azure through `explicit_max_ttl` on roles ([GH-199](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/199))

## v0.19.2
IMPROVEMENTS:
* Updated dependencies [[GH-215]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/215)
  * `github.com/Azure/azure-sdk-for-go/sdk/azidentity` v1.5.2 ->  v1.6.0
  * `github.com/hashicorp/go-retryablehttp` v0.7.1 -> v0.7.7
  * `golang.org/x/crypto` v0.21.0 -> v0.24.0
  * `golang.org/x/net` v0.23.0 -> v0.26.0
  * `golang.org/x/sys` v0.18.0 -> v0.21.0

## v0.19.1

BUG FIXES:
* Fix segmentation fault when unassigning role assignments [[GH-213]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/213)

## v0.19.0

IMPROVEMENTS:
* Updated dependencies:
  * `github.com/microsoftgraph/msgraph-sdk-go` v1.40.0 -> v1.42.0

FEATURES:
* Adds secret-less configuration of Azure secret engine using plugin Workload Identity Federation (https://github.com/hashicorp/vault-plugin-secrets-azure/pull/188)

## v0.18.1

BUG FIXES:
* Use applicationObjectID instead of clientID in GetApplication filter [[GH-200]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/200)

## v0.18.0

CHANGES:

* `/config` endpoint no longer supports a `password_policy` parameter [[GH-181]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/181)

BUGS:

* Prevent panic when unassigning roles [[GH-191]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/191)

IMPROVEMENTS:

* Updated dependencies [[GH-182]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/182) [[GH-197]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/197)
   * `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.9.1 -> v1.11.1
   * `github.com/Azure/azure-sdk-for-go/sdk/azidentity` v1.5.1 -> v1.5.2
   * `github.com/hashicorp/go-hclog` v1.6.2 -> v1.6.3
   * `github.com/hashicorp/vault/api` v1.11.0 -> v1.13.0
   * `github.com/hashicorp/vault/sdk` v0.10.2 -> v0.12.0
   * `github.com/microsoftgraph/msgraph-sdk-go` v1.32.0 -> v1.40.0
   * `github.com/microsoftgraph/msgraph-sdk-go-core` v1.0.1 -> v1.1.0
* `google.golang.org/protobuf` v1.32.0 -> v1.33.0 [[GH-184]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/184)
* `github.com/docker/docker` v25.0.2+incompatible -> v25.0.5+incompatible [[GH-185]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/185)
* `golang.org/x/net` 0.21.0 -> 0.23.0 [[GH-195]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/195)

## v0.17.3

BUG FIXES:
* Fix segmentation fault when unassigning role assignments [[GH-213]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/213)

IMPROVEMENTS:
* Update dependencies [[GH-214]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/214)
  * `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.9.1 -> v1.11.1
  * `github.com/Azure/azure-sdk-for-go/sdk/azidentity` v1.5.1 -> v1.6.0
  * `github.com/hashicorp/go-hclog` v1.6.2 -> v1.6.3
  * `github.com/docker/docker` v25.0.2+incompatible -> v25.0.5+incompatible
  * `github.com/go-jose/go-jose/v3` v3.0.1 -> v3.0.3
  * `github.com/hashicorp/go-retryablehttp` v0.7.1 -> v0.7.7
  * `golang.org/x/crypto` v0.17.0 -> v0.24.0
  * `golang.org/x/net` v0.19.0 -> v0.26.0

## v0.17.2

BUG FIXES:
* Use applicationObjectID instead of clientID in GetApplication filter [[GH-200]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/200)

## v0.17.1

BUG FIXES:
* Add nil check for response when unassigning roles [[GH-191]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/191)

## v0.17.0

IMPROVEMENTS:

* Update dependencies [[GH-176]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/176)
  * github.com/Azure/azure-sdk-for-go/sdk/azcore v1.9.0 -> v1.9.1
  * github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.4.0 -> v1.5.1
  * github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2 v2.1.1 -> v2.2.0
  * github.com/google/uuid v1.3.1 -> v1.6.0
  * github.com/hashicorp/go-hclog v1.5.0 -> v1.6.2
  * github.com/hashicorp/vault/api v1.10.0 -> v1.11.0
  * github.com/hashicorp/vault/sdk v0.10.0 -> v0.10.2
  * github.com/microsoftgraph/msgraph-sdk-go v1.22.0 -> v1.32.0
  * github.com/microsoftgraph/msgraph-sdk-go-core v1.0.0 -> v1.0.1

## v0.16.3

IMPROVEMENTS:

* Add sign_in_audience and tags fields to application registration [GH-174](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/174)
* Prevent write-ahead-log data from being replicated to performance secondaries [GH-164](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/164)
* Update dependencies [[GH-161]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/161)
  * github.com/Azure/azure-sdk-for-go v68.0.0
* Update dependencies [[GH-162]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/162)
  * golang.org/x/crypto v0.13.0
  * golang.org/x/net v0.15.0
  * golang.org/x/sys v0.12.0
  * golang.org/x/text v0.13.0

## v0.16.2

IMPROVEMENTS:

* Update dependencies [[GH-160]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/160)
  * github.com/hashicorp/vault/api v1.9.1 -> v1.10.0
  * github.com/hashicorp/vault/sdk v0.9.0 -> v0.10.0

## v0.16.1

BUG FIXES:

* Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials [[GH-150]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/150)

## v0.16.0

IMPROVEMENTS:

* permanently delete app during WAL rollback [GH-138](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/138)
* enable plugin multiplexing [GH-134](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/134)
* add display attributes for OpenAPI OperationID's [GH-141](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/141)
* update dependencies
  * `github.com/hashicorp/vault/api` v1.9.1 [GH-145](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/145)
  * `github.com/hashicorp/vault/sdk` v0.9.0 [GH-141](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/141)
  * `github.com/hashicorp/go-hclog` v1.5.0 [GH-140](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/140)
  * `github.com/Azure/go-autorest/autorest` v0.11.29 [GH-144](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/144)

## v0.15.1

BUG FIXES:

* Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials [[GH-150]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/150)

## v0.15.0

CHANGES:

* Changes user-agent header value to use correct Vault version information and include
  the plugin type and name in the comment section [[GH-123]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/123)

FEATURES:

* Adds ability to persist an application for the lifetime of a role [[GH-98]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/98)

IMPROVEMENTS:

* Updated dependencies [[GH-109](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/109)]
    * `github.com/Azure/azure-sdk-for-go v67.0.0+incompatible`
    * `github.com/Azure/go-autorest/autorest v0.11.28`
    * `github.com/Azure/go-autorest/autorest/azure/auth v0.5.11`
    * `github.com/hashicorp/go-hclog v1.3.1`
    * `github.com/hashicorp/go-uuid v1.0.3`
    * `github.com/hashicorp/vault/api v1.8.2`
    * `github.com/hashicorp/vault/sdk v0.6.1`
    * `github.com/mitchellh/mapstructure v1.5.0`
* Upgraded to go 1.19 [[GH-109](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/109)]

## v0.14.2

BUG FIXES:

* Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials [[GH-150]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/150)

## v0.14.1

BUG FIXES:

* Adds WAL rollback mechanism to clean up Role Assignments during partial failure [[GH-110]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/110)

## v0.14.0

IMPROVEMENTS:

* Add option to permanently delete AzureAD objects [[GH-104](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/104)]

CHANGES:

* Remove deprecated AAD graph code [[GH-101](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/101)]
* Remove partner ID from user agent string [[GH-95](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/95)]

## v0.11.4

CHANGES:

* Sets `use_microsoft_graph_api` to true by default [[GH-90](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/90)]

BUG FIXES:

* Fixes environment not being used when using MS Graph [[GH-87](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/87)]
