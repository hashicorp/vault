## Unreleased

## v0.14.3

BUG FIXES:

* fix an edge case where add an LDAP user or service account can be added to more than one role or set (https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/123)

## v0.14.2

BUG FIXES:

* fix a panic on static role creation when the config is unset (https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/119)

* fix case sensitivity issues in the role rotation process (https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/118)

## v0.14.1

BUG FIXES:
* fix a panic on init when static roles have names defined as hierarchical paths (https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/115)

## v0.14.0

### IMPROVEMENTS:

* update dependencies [GH-113](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/113)
  * `github.com/go-ldap/ldap/v3` v3.4.6 -> v3.4.8
  * `github.com/hashicorp/go-hclog` v1.6.2 -> v1.6.3
  * `github.com/hashicorp/go-secure-stdlib/parseutil` v0.1.7 -> v0.1.8
  * `github.com/hashicorp/vault/api` v1.13.0 -> v1.14.0
  * `github.com/hashicorp/vault/sdk` v0.12.0 -> v0.13.0
  * `golang.org/x/text` v0.14.0 -> v0.18.0
  * `github.com/hashicorp/go-retryablehttp` v0.7.1 -> v0.7.7
* bump .go-version to 1.22.6

## v0.13.1

BUG FIXES:
* fix a panic on init when static roles have names defined as hierarchical paths (https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/116)

## v0.13.0

FEATURES:
* Enable role and set names with hierarchical paths
  * https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/102
  * https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/104
  * https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/105

IMPROVEMENTS:
* Updated dependencies (https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/101):
   * `github.com/go-ldap/ldap/v3` v3.4.4 -> v3.4.6
   * `github.com/hashicorp/go-hclog` v1.5.0 -> v1.6.2
   * `github.com/hashicorp/go-secure-stdlib/parseutil` v0.1.7 -> v0.1.8
   * `github.com/hashicorp/vault/api` v1.9.2 -> v1.13.0
   * `github.com/hashicorp/vault/sdk` v0.11.1-0.20240325190132-c20eae3e84c5 -> v0.12.0
   * `github.com/stretchr/testify` v1.8.4 -> v1.9.0

## v0.12.1

### BUG FIXES:
* Fix inability to rotate-root when using `userattr=userPrincipalName` and `upndomain` is not set [GH-91](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/91)

## v0.12.0

* update dependencies [GH-90](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/90)
  * Bump golang.org/x/crypto from 0.7.0 to 0.17.0 (#87)
  * Bump github.com/go-jose/go-jose/v3 from 3.0.0 to 3.0.1 (#86)
  * Bump google.golang.org/grpc from 1.53.0 to 1.56.3 (#84)
  * Bump golang.org/x/net from 0.8.0 to 0.17.0 (#81)

## v0.11.3

### FEATURES:
* add `skip_static_role_import_rotation` and `skip_import_rotation` to allow users to retain the existing role password
on import (note: Vault will not know the role password until it is rotated) [GH-83](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/83)

### BUG FIXES:
* Revert back to armon/go-metrics [GH-88](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/88)

### IMPROVEMENTS:
* add rotate-root support when using userattr=userPrincipalName

## v0.11.2

### IMPROVEMENTS:

* update dependencies [GH-XXX](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/XXX)
  * github.com/hashicorp/go-metrics v0.5.1
  * github.com/hashicorp/vault/api v1.9.2
  * github.com/hashicorp/vault/sdk v0.9.2
  * github.com/stretchr/testify v1.8.4
  * golang.org/x/text v0.12.0

## v0.11.1

### IMPROVEMENTS:
* prevent overwriting of schema and password_policy values on update of config [GH-75](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/75)

## v0.11.0

### IMPROVEMENTS:

* enable plugin multiplexing [GH-55](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/55)
* update dependencies
  * `github.com/hashicorp/vault/api` v1.9.1
  * `github.com/hashicorp/vault/sdk` v0.9.0

## v0.10.0

CHANGES:

* CreateOperation should only be implemented alongside ExistenceCheck [[GH-50]](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/50)

IMPROVEMENTS:

* Update golang.org/x/text to v0.3.8 [[GH-48]](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/48)

## v0.9.0

FEATURES:

- Adds service account check-out functionality for `ad`, `openldap`, and `racf` schemas.

IMPROVEMENTS:

- Adds the `last_password` field to the static role [credential response](https://www.vaultproject.io/api-docs/secret/openldap#static-role-passwords)
- Adds the `userdn` and `userattr` configuration parameters to control how user LDAP
  search is performed for service account check-out and static roles.
- Adds the `upndomain` configuration parameter to allow construction of a userPrincipalName
  (UPN) string for authentication.

BUG FIXES:

- Fix config updates so that they retain prior values set in storage
- Fix `last_bind_password` client rotation retry that may occur after a root credential rotation
