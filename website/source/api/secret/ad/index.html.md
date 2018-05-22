---
layout: "api"
page_title: "Active Directory - Secrets Engines - HTTP API"
sidebar_current: "docs-http-secret-active-directory"
description: |-
  This is the API documentation for the Vault Active Directory secrets engine.
---

# Active Directory Secrets Engine (API)

This is the API documentation for the Vault AD secrets engine. For general
information about the usage and operation of the AD secrets engine, please see
the [Vault Active Directory documentation](/docs/secrets/ad/index.html).

This documentation assumes the AD secrets engine is enabled at the `/ad` path
in Vault. Since it is possible to enable secrets engines at any location, please
update your API calls accordingly.

## Configuration

The `config` endpoint configures the LDAP connection and binding parameters, as well as the password rotation configuration.

### Password parameters

* `ttl` (string, optional) - The default password time-to-live in seconds. Once the ttl has passed, a password will be rotated the next time it's requested.
* `max_ttl` (string, optional) - The maximum password time-to-live in seconds. No role will be allowed to set a custom ttl greater than the `max_ttl`.
* `password_length` (string, optional) - The desired password length. Defaults to 64. Minimum is 14. Note: to meet complexity requirements, all passwords begin with "?@09AZ".

### Connection parameters

* `url` (string, required) - The LDAP server to connect to. Examples: `ldap://ldap.myorg.com`, `ldaps://ldap.myorg.com:636`. This can also be a comma-delineated list of URLs, e.g. `ldap://ldap.myorg.com,ldaps://ldap.myorg.com:636`, in which case the servers will be tried in-order if there are errors during the connection process.
* `starttls` (bool, optional) - If true, issues a `StartTLS` command after establishing an unencrypted connection.
* `insecure_tls` - (bool, optional) - If true, skips LDAP server SSL certificate verification - insecure, use with caution!
* `certificate` - (string, optional) - CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded.

### Binding parameters

* `binddn` (string, required) - Distinguished name of object to bind when performing user and group search. Example: `cn=vault,ou=Users,dc=example,dc=com`
* `bindpass` (string, required) - Password to use along with `binddn` when performing user search.
* `userdn` (string, optional) - Base DN under which to perform user search. Example: `ou=Users,dc=example,dc=com`
* `upndomain` (string, optional) - userPrincipalDomain used to construct the UPN string for the authenticating user. The constructed UPN will appear as `[username]@UPNDomain`. Example: `example.com`, which will cause vault to bind as `username@example.com`.

## Config management

At present, this endpoint does not confirm that the provided AD credentials are
valid AD credentials with proper permissions.

| Method   | Path                   | Produces               |
| :------- | :--------------------- | :--------------------- |
| `POST`   | `/ad/config`           | `204 (empty body)`     |
| `GET`    | `/ad/config`           | `200 (see body below)` |
| `DELETE` | `/ad/config`           | `204 (empty body)`     |

### Sample Payload

```json
{
  "binddn": "domain-admin",
  "bindpass": "pa$$w0rd",
  "url": "ldap://127.0.0.11",
  "userdn": "dc=example,dc=com"
}
```

### Sample Post Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/ad/config
```
### Sample Get Response

```json
{
  "binddn": "domain-admin",
  "url": "ldaps://example.com",
  "certificate": "-----BEGIN CERTIFICATE-----....",
  "userdn": "dc=example,dc=com",
  "insecure_tls": false,
  "start_tls": true,
  "tls_min_version": "tls12",
  "tls_max_version": "tls12",
  "ttl": 2764800,
  "max_ttl": 2764800,
  "password_length": 64
}

```

## Role management

The `roles` endpoint configures how Vault will manage the passwords for individual service accounts.

### Parameters

* `service_account_name` (string, required) - The name of a pre-existing service account in Active Directory that maps to this role.
* `ttl` (string, optional) - The password time-to-live in seconds. Defaults to the configuration `ttl` if not provided.

When adding a role, Vault verifies its associated service account exists.

| Method   | Path                   | Produces               |
| :------- | :--------------------- | :--------------------- |
| `GET`    | `/ad/roles`            | `204 (see body below)` |
| `POST`   | `/ad/roles/:role_name` | `204 (empty body)`     |
| `GET`    | `/ad/roles/:role_name` | `200 (see body below)` |
| `DELETE` | `/ad/roles/:role_name` | `204 (empty body)`     |

### Sample Payload

```json
{
  "service_account_name": "my-application@example.com",
  "ttl": 100
}
```

### Sample Post Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/ad/roles/my-application
```

### Sample Get Response

```json
{
    "service_account_name": "my-application@example.com",
    "last_vault_rotation": "2018-03-29T14:33:24Z07:00",
    "password_last_set": "2018-03-30T14:33:24Z07:00",
    "ttl": 100
}
```

### Sample List Response

Performing a `GET` on the `/ad/roles` endpoint will list the names of all the roles Vault contains.

```json
[
  "my-application",
  "another-application"
]
```

## Retrieving passwords

The `creds` endpoint offers the credential information for a given role.

### Sample Get Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    --data @payload.json \
    http://127.0.0.1:8200/v1/ad/creds/my-application
```

### Sample Get Response

```json
{
    "username": "my-application",
    "current_password": "?@09AZ2wR37xJf",
    "last_password": "?@09AZ7WEu9fu8"
}
```