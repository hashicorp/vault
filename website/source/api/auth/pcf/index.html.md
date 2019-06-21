---
layout: "api"
page_title: "PCF - Auth Methods - HTTP API"
sidebar_title: "PCF"
sidebar_current: "api-http-auth-pcf"
description: |-
  This is the API documentation for the Vault PCF auth method.
---

# Pivotal Cloud Foundry (PCF) Auth Method (API)

This is the API documentation for the Vault PCF auth method. For
general information about the usage and operation of the PCF method, please
see the [Vault PCF method documentation](/docs/auth/pcf.html).

This documentation assumes the PCF method is mounted at the `/auth/pcf`
path in Vault. Since it is possible to enable auth methods at any location,
please update your API calls accordingly.

## Create Configuration

Configure the root CA certificate to be used for verifying instance identity
certificates, and configure access to the PCF API. For detailed instructions
on how to obtain these values, please see the [Vault PCF method 
documentation](/docs/auth/pcf.html).

| Method   | Path                  |
| :--------|---------------------- |
| `POST`   | `/auth/pcf/config`    |

### Parameters

- `identity_ca_certificates` `(array: [], required)` - The root CA certificate(s) 
to be used for verifying that the `CF_INSTANCE_CERT` presented for logging in was
issued by the proper authority.
- `pcf_api_addr` `(string: required)`: PCF's full API address, to be used for verifying
that a given `CF_INSTANCE_CERT` shows an application ID, space ID, and organization ID 
that presently exist.
- `pcf_username` `(string: required)`: The username for authenticating to the PCF API.
- `pcf_password` `(string: required)`: The password for authenticating to the PCF API.
- `pcf_api_trusted_certificates` `(array: [])`: The certificate that's presented by the
PCF API. This configures Vault to trust this certificate when making API calls, resolving
`x509: certificate signed by unknown authority` errors.
- `login_max_seconds_old` `(int: 300)`: The maximum number of seconds in the past when a 
signature could have been created. The lower the value, the lower the risk of replay
attacks.
- `login_max_seconds_ahead` `(int: 60)`: In case of clock drift, the maximum number of
seconds in the future when a signature could have been created. The lower the value, 
the lower the risk of replay attacks.

### Sample Payload

```json
{
  "identity_ca_certificates": ["-----BEGIN CERTIFICATE-----\nMIIEtzCCA5+.......ZRtAfQ6r\nwlW975rYa1ZqEdA=\n-----END CERTIFICATE-----"],
  "pcf_api_addr": "https://api.sys.somewhere.cf-app.com",
  "pcf_username": "vault",
  "pcf_password": "pa55w0rd",
  "pcf_api_trusted_certificates": ["-----BEGIN CERTIFICATE-----\nMIIEtzCCA5+.......ZRtAfQ6r\nwlW975rYa1ZqEdA=\n-----END CERTIFICATE-----"],
  "login_max_seconds_old": 5,
  "login_max_seconds_ahead": 1
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/pcf/config
```

## Read Config

Returns the present PCF configuration.

| Method   | Path                  |
| :--------|---------------------- |
| `GET`    | `/auth/pcf/config`    |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/pcf/config
```

### Sample Response

```json
{
  "identity_ca_certificates": ["-----BEGIN CERTIFICATE-----\nMIIEtzCCA5+.......ZRtAfQ6r\nwlW975rYa1ZqEdA=\n-----END CERTIFICATE-----"],
  "pcf_api_addr": "https://api.sys.somewhere.cf-app.com",
  "pcf_username": "vault",
  "pcf_api_trusted_certificates": ["-----BEGIN CERTIFICATE-----\nMIIEtzCCA5+.......ZRtAfQ6r\nwlW975rYa1ZqEdA=\n-----END CERTIFICATE-----"],
  "login_max_seconds_old": 5,
  "login_max_seconds_ahead": 1
}
```

## Delete Config

Deletes the present PCF configuration.

| Method   | Path                  |
| :--------|---------------------- |
| `DELETE` | `/auth/pcf/config`    |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/pcf/config
```

## Create Role

Create a role in Vault granting a particular level of access to a particular group
of PCF instances. We recommend using the PCF API or the CF CLI to gain the IDs you
wish to target. 

If you list no `bound` parameters, then any entity with a valid 
`CF_INSTANCE_CERT` that's been issued by any configured `identity_ca_certificates`
will be able to authenticate against this role.

| Method   | Path                   |
| :--------|----------------------- |
| `POST`   | `/auth/pcf/roles/:role`|

### Parameters

- `role` `(string: required)` - The name of the role.
- `bound_application_ids` `(array: [])` - An optional list of application IDs 
an instance must be a member of to qualify as a member of this role.
- `bound_space_ids` `(array: [])` - An optional list of space IDs 
an instance must be a member of to qualify as a member of this role.
- `bound_organization_ids` `(array: [])` - An optional list of organization IDs 
an instance must be a member of to qualify as a member of this role.
- `bound_instance_ids` `(array: [])` - An optional list of instance IDs 
an instance must be a member of to qualify as a member of this role. Please note that 
every time you use `cf push` on an app, its instance ID changes. Also, instance IDs 
are not verifiable as being presently alive using the PCF API. Thus, we recommend against
using this setting for most use cases.
- `bound_cidrs` `(array: [])` - Comma separated string or list of CIDR blocks. 
If set, specifies the blocks of IP addresses which can perform the login operation.
- `policies` `(array: [])` - Policies to be set on tokens issued using this role.
- `disable_ip_matching` `(bool: false)` - If set to true, disables the default behavior 
that logging in must be performed from an acceptable IP address described by the 
certificate presented. Should only be set to true if required, generally when a proxy
is used to perform logins.
- `ttl` `(string: "")` - The TTL period of tokens issued using this role,
  provided as "1h", where hour is the largest suffix.
- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens issued using
  this role.
- `period` `(string: "")` - If set, indicates that the token generated using
  this role should never expire. The token should be renewed within the duration
  specified by this value. At each renewal, the token's TTL will be set to the
  value of this parameter.

### Sample Payload

```json
{
  "bound_application_ids": ["09d7eb6a-afc2-49a0-bb32-858c22f2b346"],
  "bound_space_ids": ["21005ebb-8943-433e-84e6-d9d9d7338853"],
  "bound_organization_ids": ["9785a884-5e93-49bd-97ee-57bf7c2b20e0"],
  "bound_instance_ids": ["f3e0f176-3f83-4efb-5842-2ff4"],
  "bound_cidrs": ["127.0.0.1/32", "128.252.0.0/16"],
  "policies": ["default"],
  "ttl": "1h",
  "max_ttl": "1h",
  "period": "1h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/pcf/roles/:role
```

## Read Role

Returns a PCF role.

| Method   | Path                   |
| :--------|----------------------- |
| `GET`    | `/auth/pcf/roles/:role`|

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/pcf/roles/:role
```

### Sample Response

```json
{
  "bound_application_ids": ["09d7eb6a-afc2-49a0-bb32-858c22f2b346"],
  "bound_space_ids": ["21005ebb-8943-433e-84e6-d9d9d7338853"],
  "bound_organization_ids": ["9785a884-5e93-49bd-97ee-57bf7c2b20e0"],
  "bound_instance_ids": ["f3e0f176-3f83-4efb-5842-2ff4"],
  "bound_cidrs": ["127.0.0.1/32", "128.252.0.0/16"],
  "policies": ["default"],
  "ttl": 2764790,
  "max_ttl": 2764790,
  "period": 2764790
}
```

## Delete Role

Deletes a PCF role.

| Method   | Path                   |
| :--------|----------------------- |
| `DELETE` | `/auth/pcf/roles/:role`|

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/pcf/roles/:role
```

## List Roles

Returns a PCF role.

| Method   | Path                   |
| :--------|----------------------- |
| `LIST`   | `/auth/pcf/roles`      |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST
    http://127.0.0.1:8200/v1/auth/pcf/roles
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "role1",
      "role2"
    ]
  }
}
```

## Login

Log in to PCF. 

Vault provides both an agent and a CLI tool for logging in that
eliminates the need to build a signature yourself. However, if you do wish to 
build the signature, its signing algorithm is viewable [here](https://github.com/hashicorp/vault-plugin-auth-pcf/tree/master/signatures).
The [plugin repo](https://github.com/hashicorp/vault-plugin-auth-pcf) also contains
a command-line tool (`generate-signature`) that can be compiled as a binary for generating a signature,
and a test that outputs steps in generating the signature so they can be duplicated.

However, at a high level, these are the steps for generating a signature:
- Get and format the current time, ex. `2006-01-02T15:04:05Z`.
- Get the full body of the file located at `CF_INSTANCE_CERT`.
- Get the name of the role.
- Concatenate them together in the above order, with no extra string used for joining them.
- Create a SHA256 checksum of the resulting string (`checksum` below).
- Sign the string using the key located at `CF_INSTANCE_KEY`. In Go, this is performed using
the following line of code which you can more deeply inspect: 
```
rsa.SignPSS(rand.Reader, rsaPrivateKey, crypto.SHA256, checksum, nil)
```
- Convert the signature to a string.

| Method   | Path                   |
| :--------|----------------------- |
| `POST`   | `/auth/pcf/login`      |

### Parameters

- `role` `(string: required)` - The name of the role.
- `cf_instance_cert` `(string: required)` - The full body of the file available at
the path denoted by `CF_INSTANCE_CERT`.
- `signing_time` `(string: required)` - The date and time used to construct the signature.
- `signature` `(string: required)` - The signature generated by the algorithm described
above using the `CF_INSTANCE_KEY`.

### Sample Payload

```json
{
  "role": "test-role",
  "cf_instance_cert": "-----BEGIN CERTIFICATE-----\nMIIEtzCCA5+.......ZRtAfQ6r\nwlW975rYa1ZqEdA=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIEtzCCA5+.......ZRtAfQ6r\nwlW975rYa1ZqEdA=\n-----END CERTIFICATE-----",
  "signing_time": "2006-01-02T15:04:05Z",
  "signature": "MmyUjQ1OqxQmF0W6raVQDL-hlqqe1oG-7abA6Oi3NHwT-9lMfrYxsCwMnd2HKGMly2tCgetcoA2orfquoe6MkMuksx_KGH_KLObcAykt53z4rHceHKGvm7eGj60cjWFYtiNPic-lzUGERLbUeKLMi6NlThm9ueb7hhpyTUpEYtphV3gorbVxvlkrnuYSbgy2NGpOUY1N8dRzcxmHkYjh12XoWEw4Is5aFr6eFKbZ0vmLWBzhJ7_w20CFyTpRYB-6heGz1iR9qEG8mZk3_x4rZpT5mejJ5zmH2xlUjBJMndfcz47btEi2BO9pFVxK2wK-tKeUUFgx6RcomAopTskkmg=="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/pcf/login
```

### Sample Response
```json
{
  "auth": {
    "renewable": true,
    "lease_duration": 1800000,
    "policies": [
      "default",
      "dev"
    ],
    "accessor": "20b89871-e6f2-1160-fb29-31c2f6d4645e",
    "client_token": "c9368254-3f21-aded-8a6f-7c818e81b17a"
  }
}
```
