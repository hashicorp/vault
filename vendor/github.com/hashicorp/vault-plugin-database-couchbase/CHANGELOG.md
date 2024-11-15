## Unreleased

## v0.12.0
### Sept 4, 2024

IMPROVEMENTS:
* Updated dependencies: (https://github.com/hashicorp/vault-plugin-database-couchbase/pull/80)

BUG FIXES:
* allow custom username templates to use the lowercase function (https://github.com/hashicorp/vault-plugin-database-couchbase/pull/81)

## v0.11.0
IMPROVEMENTS:
* Updated dependencies:
  * `github.com/jackc/pgx/v4` v4.18.1 -> v4.18.2
  * `google.golang.org/protobuf` v1.32.0 -> v1.33.0
  * `github.com/hashicorp/go-plugin` v1.5.2 to -> v1.6.0

## v0.10.1
* Revert dependency update causing build failures on 32-bit systems
  * github.com/couchbase/gocb/v2 v2.7.1 -> v2.6.5

## v0.10.0
* Updated dependencies:
  * github.com/couchbase/gocb/v2 v2.6.3 -> v2.6.5
  * github.com/hashicorp/go-hclog v1.5.0 -> v1.6.2
  * github.com/hashicorp/vault/sdk v0.10.0 -> v0.10.2
  * golang.org/x/mod v0.12.0 -> v0.15.0
  * github.com/opencontainers/runc v1.1.6 -> v1.1.12
  * github.com/docker/docker v24.0.5+incompatible -> v24.0.9+incompatible

## v0.9.4

IMPROVEMENTS:
* Updated indirect dependency `golang.org/x/net` v0.9.0 -> v0.15.0 due to vulnerability GO-2023-1988 v0.9.0

## v0.9.3

IMPROVEMENTS:

* Updated dependencies:
  * `github.com/hashicorp/vault/sdk` v0.9.0 -> v0.10.0
  * `github.com/stretchr/testify` v1.8.3 -> v1.8.4
  * `golang.org/x/mod` v0.9.0 -> v0.12.0

## v0.9.2

CHANGES:
* Renaming  `cmd/couchbase-database-plugin/main.go` to `cmd/vault-plugin-database-couchbase/main.go` [[GH-50](https://github.com/hashicorp/vault-plugin-database-couchbase/pull/50)]

## v0.9.1

IMPROVEMENTS:
* Updated dependencies:
   * `github.com/couchbase/gocb/v2` v2.3.3 -> v2.6.3
   * `github.com/hashicorp/go-hclog` v1.0.0 -> v1.5.0
   * `github.com/hashicorp/go-version` v1.3.0 -> v1.6.0
   * `github.com/hashicorp/vault/sdk` v0.5.3 -> v0.9.0
   * `github.com/ory/dockertest/v3` v3.8.0 -> v3.10.0
   * `github.com/stretchr/testify` v1.7.0 -> v1.8.3
   * `golang.org/x/mod` v0.9.0 added
