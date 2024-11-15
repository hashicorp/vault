## Unreleased


## 0.17.0
### Sept 5, 2024


### Build:
* Build with go 1.22.6


### Dependency updates:
* `github.com/hashicorp/vault/api` v1.12.2 -> v1.14.0
* `github.com/hashicorp/vault/sdk` v0.11.1 -> v0.13.0


## 0.16.0
### May 20, 2024

IMPROVEMENTS:
* Updated dependencies [PR-52](https://github.com/hashicorp/vault-plugin-auth-oci/pull/52)
* Updated dependencies:
  * `github.com/hashicorp/go-plugin` v1.5.2 -> v1.6.0 to enable running the plugin in containers

## 0.15.1
### February 6, 2024

CHANGES:
* Downgrades github.com/oracle/oci-go-sdk to v59.0.0 due to an incompatibility with Vault
 
## 0.15.0
### February 6, 2024

IMPROVEMENTS:
* Update go.mod version to 1.21
* Update dependencies:
  * github.com/oracle/oci-go-sdk v24.3.0 -> v65.57.0
  * github.com/hashicorp/go-hclog v1.5.0 -> v1.6.2
  * github.com/hashicorp/vault/api v1.10.0 -> v1.11.0
  * github.com/hashicorp/vault/sdk v0.10.0 -> v0.10.2

## 0.14.2
### September 5, 2023

IMPROVEMENTS:
* Update dependencies:
  * golang.org/x/net v0.9.0 -> v0.15.0

## 0.14.1
### September 5, 2023

IMPROVEMENTS:
* Update dependencies:
  * github.com/hashicorp/vault/api v1.9.1 -> v1.9.2
  * github.com/hashicorp/vault/sdk v0.9.0 -> v0.9.2

## 0.14.0

* Add display attributes for OpenAPI OperationID's [GH-29](https://github.com/hashicorp/vault-plugin-auth-oci/pull/29)
* update dependencies
  * `github.com/hashicorp/vault/api` v1.9.1 [GH-31](https://github.com/hashicorp/vault-plugin-auth-oci/pull/31)

## 0.13.1

CHANGES:
* Repond with a 400 instead of 401 to login errors. [GH-27](https://github.com/hashicorp/vault-plugin-auth-oci/pull/27)

IMPROVEMENTS:

* Return success message when writing role [GH-27](https://github.com/hashicorp/vault-plugin-auth-oci/pull/27)
* Return error messages when failing to login [GH-27](https://github.com/hashicorp/vault-plugin-auth-oci/pull/27)
* enable plugin multiplexing [GH-25](https://github.com/hashicorp/vault-plugin-auth-oci/pull/25)
* update dependencies
  * `github.com/hashicorp/vault/api` v1.9.0
  * `github.com/hashicorp/vault/sdk` v0.8.1
