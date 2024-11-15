## Unreleased

## v0.16.0
Improvements:
* Updated dependencies:
  * `github.com/hashicorp/vault/sdk` v0.12.0 -> v0.13.0
  * `github.com/hashicorp/go-retryablehttp` v0.7.6 -> v0.7.7
  * `golang.org/x/net` v0.23.0 -> v0.25.0
  * `golang.org/x/crypto` v0.21.0 -> v0.23.0
  * `golang.org/x/text` v0.16.0 -> v0.17.0
  * `golang.org/x/sys` v0.21.0 -> v0.24.0
  * `golang.org/x/sync` v0.7.0 -> v0.8.0
  * `github.com/docker/docker` v24.0.9 -> v25.0.6
  * `google.golang.org/grpc` v1.60.1 -> v1.65.0
  * `google.golang.org/protobuf` v1.33.0 -> v1.34.1
  * `github.com/opencontainers/runc` v1.1.12 -> v1.1.14

## v0.15.0
Improvements:
* Updated dependencies:
  * `github.com/hashicorp/vault/sdk` v0.11.0 -> v0.12.0
  * `github.com/hashicorp/go-retryablehttp` v0.7.5 -> v0.7.6
  * `github.com/stretchr/testify` v1.8.4 -> v1.9.0
  * `golang.org/x/net` v0.17.0 -> v0.23.0
  * `golang.org/x/crypto` v0.17.0 -> v0.21.0
  * `github.com/docker/docker` v24.0.7 -> v24.0.9

Improvements:
* Updated dependencies:
  * `github.com/hashicorp/vault/sdk` v0.10.2 -> v0.11.0
  * `google.golang.org/protobuf` v1.31.0 -> 1.33.0

## v0.14.0
Changes:
* Building with go 1.21.7

Improvements:
* Updated dependencies:
  * github.com/hashicorp/go-retryablehttp v0.7.4 -> v0.7.5
  * github.com/hashicorp/go-secure-stdlib/tlsutil v0.1.2 -> v0.1.3
  * github.com/hashicorp/vault/sdk v0.9.2 -> v0.10.2
  * golang.org/x/net v0.8.0 -> v0.17.0
  * golang.org/x/crypto v0.6.0 -> v0.17.0
  * github.com/opencontainers/runc v1.1.6 -> v1.1.12
  * github.com/docker/docker v24.0.5 -> v24.0.7
  * google.golang.org/grpc v1.57.0 -> v1.57.1

## v0.13.3
Improvements:
* Updated dependencies:
  * github.com/hashicorp/go-retryablehttp v0.7.2 -> v0.7.4
  * github.com/hashicorp/vault/sdk v0.9.0 -> v0.9.2
  * github.com/stretchr/testify v1.8.2 -> v1.8.4

## v0.13.2
* Dependency upgrades 

## v0.13.1
* No new features

## v0.13.0
* No new features

## v0.12.0
* No new features

## 0.11.1 (Aug 1, 2022)

Bug Fixes:
* Fix bug in boolean parsing for initialize [[GH-38](https://github.com/hashicorp/vault-plugin-database-elasticsearch/pull/38)]

## 0.11.0 (May 25th, 2022)

Improvements:
* Use the new `/_security` base API path instead of `/_xpack/security` when managing elasticsearch [[GH-37](https://github.com/hashicorp/vault-plugin-database-elasticsearch/pull/37)]
