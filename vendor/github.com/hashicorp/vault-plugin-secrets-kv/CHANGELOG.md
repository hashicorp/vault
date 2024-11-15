## Unreleased

## v0.20.0

BUG FIXES:

* Fixes panic from occurring when renewing KVv1 secrets.

CHANGES:

* Updated dependencies:
  * `github.com/hashicorp/vault/api` v1.13.0 -> v1.14.0
  * `github.com/hashicorp/vault/sdk` v0.12.0 -> v0.13.0
  * `github.com/go-test/deep` v1.1.0 -> v1.1.1
  * `github.com/docker/docker` v25.0.5 -> v25.0.6
  * `github.com/hashicorp/go-retryablehttp` v0.7.1 -> v0.7.7
  * `google.golang.org/grpc` v1.64.0 -> v1.64.1

## v0.19.0

CHANGES:

* Updated dependencies:
  * `github.com/golang/protobuf` v1.5.3 -> v1.5.4
  * `github.com/hashicorp/go-hclog` v1.6.2 -> v1.6.3
  * `github.com/hashicorp/vault/api` v1.11.0 -> v1.13.0
  * `github.com/hashicorp/vault/sdk` v0.10.2 -> v0.12.0
  * `google.golang.org/protobuf` v1.33.0 -> v1.34.1

## 0.18.0

CHANGES:

* Updated dependencies:
  * `github.com/hashicorp/go-plugin` v1.5.2 -> v1.6.0 to enable running the plugin in containers
  * Bump golang.org/x/net from 0.18.0 to 0.23.0 (#150)
  * Bump github.com/docker/docker (#148)
  * Bump google.golang.org/protobuf from 1.32.0 to 1.33.0 (#147)
  * Bump github.com/go-jose/go-jose/v3 from 3.0.1 to 3.0.3 (#145)

## 0.17.0

CHANGES:

* Updated dependencies:
  * github.com/hashicorp/go-hclog v1.5.0 -> v1.6.2
  * github.com/hashicorp/go-secure-stdlib/parseutil v0.1.7 -> v0.1.8
  * github.com/hashicorp/vault/api v1.9.2 -> v1.11.0
  * github.com/hashicorp/vault/sdk v0.10.0 -> v0.10.2
  * google.golang.org/protobuf v1.31.0 -> v1.32.0
  * github.com/go-jose/go-jose/v3@v3.0.0 -> v3.0.1
  * golang.org/x/crypto@v0.12.0 -> v0.17.0
  * golang.org/x/net@v0.14.0 -> v0.18.0
  * google.golang.org/grpc@v1.57.0 -> v1.61.0

## 0.16.2

CHANGES:

* Updated dependencies:
  * `github.com/hashicorp/vault/sdk` v0.9.3-0.20230831152851-56ce89544e64 -> v0.10.0
  * `github.com/hashicorp/go-plugin` v1.4.10 -> v1.5.2

## 0.16.1

CHANGES:

* Updated dependencies:
  * `golang.org/x/crypto` v0.6.0 -> v0.12.0
  * `golang.org/x/net` v0.8.0 -> v0.14.0
  * `golang.org/x/text` v0.8.0 -> v0.12.0
## 0.16.0

CHANGES:

* Events: now include `data_path`, `operation`, and `modified` [GH-124](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/124)
* Updated dependencies:
   * `github.com/hashicorp/vault/api` v1.9.0 -> v1.9.2
   * `github.com/hashicorp/vault/sdk` v0.9.0 -> v0.9.3-0.20230831152851-56ce89544e64
   * `google.golang.org/protobuf` v1.30.0 ->  v1.31.0

## 0.15.0

IMPROVEMENTS:

* Add display attributes for OpenAPI OperationID's [GH-104](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/104)
* Add versions to delete and undelete events [GH-122](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/122)
