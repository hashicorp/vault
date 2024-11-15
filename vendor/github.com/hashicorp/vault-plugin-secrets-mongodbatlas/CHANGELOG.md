## Unreleased

## v0.13.0
### Sept 10, 2024
IMPROVEMENTS:
* Updated dependencies [GH-79](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/79)

## v0.12.0
### May 20, 2024
IMPROVEMENTS:
* Updated dependencies:
  * `github.com/hashicorp/go-hclog` v1.5.0 -> v1.6.3
  * `github.com/hashicorp/vault/api` v1.11.0 -> v1.13.0
  * `github.com/hashicorp/vault/sdk` v0.10.2 -> v0.12.0
  * `github.com/mongodb-forks/digest` v1.0.5 -> v1.1.0
  * `go.mongodb.org/atlas` v0.33.0 -> v0.36.0
  * `golang.org/x/net` v0.17.0 -> v0.23.0
  * `github.com/docker/docker` v24.0.7+incompatible -> v24.0.9+incompatible
* Upgrade `github.com/go-jose/go-jose/v3` to `github.com/go-jose/go-jose/v4`

## v0.11.0
### February 7, 2024
IMPROVEMENTS:
* Updated dependencies [GH-65](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/65):
   * `github.com/hashicorp/vault/api` v1.9.2 -> v1.11.0
   * `github.com/hashicorp/vault/sdk` v0.9.2 -> v0.10.2
* Bump golang.org/x/crypto from 0.6.0 to 0.17.0 [GH-64](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/64)
* Bump github.com/go-jose/go-jose/v3 from 3.0.0 to 3.0.1 [GH-63](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/63)
* Bump google.golang.org/grpc from 1.53.0 to 1.56.3 [GH-62](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/62)
* Bump golang.org/x/net from 0.8.0 to 0.17.0 [GH-58](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/58)
* Bump [github.com/docker/docker](https://github.com/docker/docker) from 24.0.5+incompatible to 24.0.7+incompatible [GH-67](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/67)
* Bump google.golang.org/grpc from 1.57.0 to 1.57.1 [GH-66](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/66)

## 0.10.2
### October 25, 2023

IMPROVEMENTS:
* Improve handing of fetching API keys for revocation of leases [GH-59](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/59)

## 0.10.1
### September 1, 2023

IMPROVEMENTS:
* Update dependencies:
  * github.com/hashicorp/vault/api v1.9.1 -> v1.9.2
  * github.com/hashicorp/vault/sdk v0.9.0 -> v0.9.2
  * github.com/mongodb-forks/digest v1.0.4 -> v1.0.5
  * go.mongodb.org/atlas v0.25.0 -> v0.33.0

## 0.10.0
### May 23, 2023

IMPROVEMENTS:
* Update Go version to 1.20.2
* Add display attributes for OpenAPI OperationID
* enable plugin multiplexing [GH-35](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/35)
* update dependencies
  * `github.com/hashicorp/vault/api` v1.9.1
  * `github.com/hashicorp/vault/sdk` v0.8.1
  * `go.mongodb.org/atlas` v0.25.0

## 0.9.1
### February 9, 2023

Bug Fixes:
* Fix a bug that did not allow WAL rollback to handle partial failures when
  creating API keys [GH-32](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/32)

Improvements:
* Update dependencies [GH-33](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/pull/33)
  * github.com/hashicorp/vault/api v1.8.3
  * github.com/hashicorp/vault/sdk v0.7.0

## 0.9.0
### February 6, 2023

Improvements:
* Change how the Vault version is acquired for the user agent string. This
  change is transparent to users.
