---
layout: "docs"
page_title: "Vault Enterprise MFA Support"
sidebar_current: "docs-vault-enterprise-mfa"
description: |-
  Vault Enterprise has support for Multi-factor Authentication (MFA), using different authentication types.

---

# Vault Enterprise MFA Support

Vault Enterprise has support for Multi-factor Authentication (MFA), using
different authentication types. MFA is built on top of the Identity system of
Vault.

## MFA Types

MFA in Vault can be of the following types.

- `Time-based One-time Password (TOTP)` - If configured and enabled on a path,
  this would require a TOTP passcode along with Vault token, to be presented
  while invoking the API request. The passcode will be validated against the
  TOTP key present in the identity of the caller in Vault.

- `Okta` - If Okta push is configured and enabled on a path, then the enrolled
  device of the user will get a push notification to approve or deny the access
  to the API. The Okta username will be derived from the caller identity's
  alias.

- `Duo` - If Duo push is configured and enabled on a path, then the enrolled
  device of the user will get a push notification to approve or deny the access
  to the API. The Duo username will be derived from the caller identity's
  alias.

- `PingID` - If PingID push is configured and enabled on a path, then the
  enrolled device of the user will get a push notification to approve or deny
  the access to the API. The PingID username will be derived from the caller
  identity's alias.

## Configuring MFA Methods

MFA methods are globally managed within the `System Backend` using the HTTP API.
Please see [MFA API](/api/system/mfa.html) for details on how to configure an MFA
method.

## MFA Methods In Policies

MFA requirements on paths are specified as `mfa_methods` along with other ACL
parameters.

### Sample Policy

```
path "secret/foo" {
    capabilities = ["read"]
    mfa_methods = ["dev_team_duo", "sales_team_totp"]
}
```

The above policy grants `read` access to `secret/foo` only after *both* the MFA
methods `dev_team_duo` and `sales_team_totp` are validated.

## Supplying MFA Credentials

MFA credentials are retrieved from the `X-Vault-MFA` HTTP header. The format of
the header is `mfa_method_name[:key[=value]]`. The items in the `[]` are
optional.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --header "X-Vault-MFA:my_totp:695452" \
    https://vault.rocks/v1/secret/foo
```

### API

MFA can be managed entirely over the HTTP API. Please see [MFA API](/api/system/mfa.html) for more details.
