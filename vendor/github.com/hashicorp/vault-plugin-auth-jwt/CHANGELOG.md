## Unreleased

## v0.22.0

IMPROVEMENTS:
* Updated dependencies:
   * `github.com/go-test/deep` v1.1.0 -> v1.1.1
   * `github.com/hashicorp/vault/api` v1.12.0 -> v1.15.0
   * `github.com/hashicorp/vault/sdk` v0.11.0 -> v0.14.0
   * `golang.org/x/oauth2` v0.21.0 -> v0.23.0
   * `golang.org/x/sync` v0.7.0 -> v0.8.0
   * `google.golang.org/api` v0.163.0 -> v0.197.0

## v0.21.0

NO CHANGES

## v0.20.3

BUG FIXES:
* Invalidate JWT with single non-empty string aud on empty bound audiences https://github.com/hashicorp/vault-plugin-auth-jwt/pull/295

## v0.20.2

IMPROVEMENTS:
* Updated dependencies:
  * `gopkg.in/square/go-jose.v2` v2.6.0 -> `gopkg.in/go-jose/go-jose.v2` v2.6.3
  * `github.com/docker/docker` v24.0.7+incompatible -> v24.0.9+incompatible
  * `golang.org/x/net` v0.22.0 -> v0.24.0
  * `golang.org/x/sys` v0.18.0 -> v0.19.0

BUG FIXES:
* Prevent error writing plugin config when run as a Vault built-in plugin for Vault version 1.16.1 https://github.com/hashicorp/vault-plugin-auth-jwt/pull/290

## v0.20.1

IMPROVEMENTS:
* Make redirect uri validation case-insensitive https://github.com/hashicorp/vault-plugin-auth-jwt/pull/282

## v0.20.0

IMPROVEMENTS:
* auth/jwt: adds the ability to specify more than one JWKS URL used to verify tokens https://github.com/hashicorp/vault-plugin-auth-jwt/pull/277
* Updated dependencies:
  * `github.com/hashicorp/vault/api` v1.10.0 -> v1.12.0
  * `github.com/hashicorp/vault/sdk` v0.10.2 -> v0.11.0
  * `golang.org/x/oauth2` v0.15.0 -> v0.17.0
  * `golang.org/x/sync` v0.5.0 -> v0.6.0
  * `google.golang.org/api` v0.154.0 -> v0.163.0

## v0.19.0

IMPROVEMENTS:
* Add support for numeric claims in `bound_claims` https://github.com/hashicorp/vault-plugin-auth-jwt/pull/265

## v0.18.0

IMPROVEMENTS:
* Include role name in Entity Alias metadata https://github.com/hashicorp/vault-plugin-auth-jwt/pull/160
* Updated dependencies:
  * `github.com/hashicorp/cap` v0.3.4 -> v0.4.0
  * `github.com/hashicorp/go-sockaddr` v1.0.2 -> v1.0.5
  * `github.com/hashicorp/vault/api` v1.9.2 -> v1.10.0
  * `github.com/hashicorp/vault/sdk` v0.9.2 -> v0.10.0
  * `golang.org/x/oauth2` v0.11.0 -> v0.12.0
  * `google.golang.org/api` v0.138.0 -> v0.143.0

FIXES:
* Add missing error check for parsing CLI flags https://github.com/hashicorp/vault-plugin-auth-jwt/pull/245

## 0.17.2

FIXES:
* Ensure SIGTSTP is only used in unix builds [[GH-255](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/255)]

## 0.17.1

IMPROVEMENTS:
* Close HTTP listener if stop or kill signal is received [[GH-251](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/251)]

## 0.17.0

IMPROVEMENTS:
* Support ADC for Google Workspace [[GH-240](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/240)]
* Updated dependencies:
   * `github.com/hashicorp/cap` v0.3.1 -> v0.3.4
   * `github.com/hashicorp/vault/sdk` v0.9.1 -> v0.9.2
   * `golang.org/x/oauth2 v0.9.0` -> v0.10.0
   * `google.golang.org/api v0.129.0` -> v0.134.0

## 0.16.0

IMPROVEMENTS:
* Updated dependencies:
   * `github.com/hashicorp/cap` v0.2.1-0.20230221194157-7894fed1633d -> v0.3.0
   * `github.com/hashicorp/vault/api` v1.9.0 -> v1.9.1
   * `github.com/hashicorp/vault/sdk` v0.8.1 -> v0.9.0
   * `github.com/stretchr/testify` v1.8.2 -> v1.8.3
   * `golang.org/x/oauth2` v0.6.0 -> v0.8.0
   * `golang.org/x/sync` v0.1.0 -> v0.2.0
   * `google.golang.org/api` v0.114.0 -> v0.124.0

## 0.15.2

IMPROVEMENTS:
* Make error response less verbose [[GH-233](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/233)]

## 0.15.1

IMPROVEMENTS:

* enable plugin multiplexing [GH-225](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/225)
* update dependencies
   * `github.com/hashicorp/vault/api` v1.9.0
   * `github.com/hashicorp/vault/sdk` v0.8.1
   * `github.com/go-test/deep` v1.0.8 -> v1.1.0
   * `github.com/hashicorp/cap` v0.2.1-0.20220727210936-60cd1534e220 -> v0.2.1-0.20230221194157-7894fed1633d
   * `github.com/hashicorp/go-hclog` v1.0.0 -> v1.5.0
   * `github.com/mitchellh/pointerstructure` v1.2.0 -> v1.2.1
   * `github.com/stretchr/testify` v1.7.0 -> v1.8.2
   * `golang.org/x/oauth2` v0.0.0-20220524215830-622c5d57e401 -> v0.6.0
   * `golang.org/x/sync` v0.0.0-20220722155255-886fb9371eb4 -> v0.1.0
   * `google.golang.org/api` v0.83.0 -> v0.114.0

## 0.15.0

IMPROVEMENTS:

* Adds `abort_on_error` parameter to CLI login command to help in non-interactive contexts [[GH-214]](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/214)
* Adds ability to set Google Workspace domain for groups search [[GH-220]](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/220)

## 0.14.0

* Updates dependency `google.golang.org/api@v0.83.0` [[GH-205](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/205)]
* Add Custom Provider for SecureAuth IdP [[GH-196](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/196)]
* Improves detection of Windows Subsystem for Linux (WSL) in CLI [[GH-209](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/209)]
* Adds support for Microsoft US Gov L4 to the Azure provider for groups fetching [[GH-211](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/211)]

## 0.13.0

* Adds ability to use JSON pointer syntax for the `user_claim` value [[GH-204](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/204)]

## 0.12.0

* Uses Proof Key for Code Exchange (PKCE) in OIDC flow [[GH-188](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/188)]

## 0.11.4

* Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode [[GH-192](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/192)]

## 0.11.3

* Uses Proof Key for Code Exchange (PKCE) in OIDC flow [[GH-191](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/191)]

## 0.11.2

* Add a skip_browser argument to make auto-launching of the default browser optional [[GH-182](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/182)]

## 0.10.2

* Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode [[GH-192](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/192)]

## 0.9.6

* Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode [[GH-192](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/192)]

## 0.8.1

BUG FIXES:

* Fixes `bound_claims` validation for provider-specific group and user info fetching [[GH-149](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/149)]
