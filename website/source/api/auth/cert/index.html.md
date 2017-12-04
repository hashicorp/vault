---
layout: "api"
page_title: "TLS Certificate Auth Backend - HTTP API"
sidebar_current: "docs-http-auth-cert"
description: |-
  This is the API documentation for the Vault TLS Certificate authentication
  backend.
---

# TLS Certificate Auth Backend HTTP API

This is the API documentation for the Vault TLS Certificate authentication 
backend. For general information about the usage and operation of the TLS
Certificate backend, please see the [Vault TLS Certificate backend documentation](/docs/auth/cert.html).

This documentation assumes the TLS Certificate backend is mounted at the
`/auth/cert` path in Vault. Since it is possible to mount auth backends at any
location, please update your API calls accordingly.

## Create CA Certificate Role

Sets a CA cert and associated parameters in a role name.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/cert/certs/:name`     | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` - The name of the certificate role.
- `certificate` `(string: <required>)` - The PEM-format CA certificate.
- `allowed_names` `(string: "")` - Constrain the Common and Alternative Names in 
  the client certificate with a [globbed pattern]
  (https://github.com/ryanuber/go-glob/blob/master/README.md#example). Value is 
  a comma-separated list of patterns.  Authentication requires at least one Name matching at least one pattern.  If not set, defaults to allowing all names.
- `policies` `(string: "")` - A comma-separated list of policies to set on tokens 
  issued when authenticating against this CA certificate.
- `display_name` `(string: "")` -   The `display_name` to set on tokens issued 
  when authenticating against this CA certificate. If not set, defaults to the 
  name of the role.
- `ttl` `(string: "")` - The TTL period of the token, provided as a number of 
  seconds. If not provided, the token is valid for the the mount or system 
  default TTL time, in that order.

### Sample Payload

```json
{
  "certificate": "-----BEGIN CERTIFICATE-----\nMIIEtzCCA5+.......ZRtAfQ6r\nwlW975rYa1ZqEdA=\n-----END CERTIFICATE-----",
  "display_name": "test"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json
    https://vault.rocks/v1/auth/cert/certs/test-ca
```

## Read CA Certificate Role

Gets information associated with the named role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/cert/certs/:name`     | `200 application/json` |

### Parameters

- `name` `(string: <required>)` - The name of the certificate role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/cert/certs/test-ca
```

### Sample Response

```json
{
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "certificate": "-----BEGIN CERTIFICATE-----\nMIIEtzCCA5+.......ZRtAfQ6r\nwlW975rYa1ZqEdA=\n-----END CERTIFICATE-----",
    "display_name": "test",
    "policies": "",
    "allowed_names": "",
    "ttl": 2764800
  },
  "warnings": null,
  "auth": null
}
```

## List Certificate Roles

Lists configured certificate names.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/cert/certs`           | `200 application/json` |
| `GET`   | `/auth/cert/certs?list=true`  | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/cert/certs

### Sample Response

```json
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "keys": [
      "cert1", 
      "cert2"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

## Delete Certificate Role

Deletes the named role and CA cert from the backend mount.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/cert/certs/:name`     | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` - The name of the certificate role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/cert/certs/cert1
```

## Create CRL

Sets a named CRL.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/cert/crls/:name`      | `204 (empty body)`     |


### Parameters 

- `name` `(string: <required>)` - The name of the CRL.
- `crl` `(string: <required>)` - The PEM format CRL.

### Sample Payload

```json
{
  "crl": "-----BEGIN X509 CRL-----\n...\n-----END X509 CRL-----"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --date @payload.json \
    https://vault.rocks/v1/auth/cert/crls/custom-crl
```

## Read CRL

Gets information associated with the named CRL (currently, the serial
numbers contained within).  As the serials can be integers up to an
arbitrary size, these are returned as strings.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/cert/crls/:name`      | `200 application/json` |

### Parameters

- `name` `(string: <required>)` - The name of the CRL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/cert/crls/custom-crl
```

### Sample Response

```json
{
  "auth": null,
  "data": {
    "serials": {
      "13": {}
    }
  },
  "lease_duration": 0,
  "lease_id": "",
  "renewable": false,
  "warnings": null
}
```

## Delete CRL

Deletes the named CRL from the backend mount.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/cert/crls/:name`      | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` - The name of the CRL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/cert/crls/cert1
```

## Configure TLS Certificate Backend

Configuration options for the backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/cert/config`          | `204 (empty body)`     |

### Parameters

- `disable_binding` `(boolean: false)` - If set, during renewal, skips the
  matching of presented client identity with the client identity used during
  login. 

### Sample Payload

```json
{
  "disable_binding": true
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --date @payload.json \
    https://vault.rocks/v1/auth/cert/certs/cert1
```

## Login with TLS Certiicate Backend

Log in and fetch a token. If there is a valid chain to a CA configured in the
backend and all role constraints are matched, a token will be issued.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/cert/login`           | `200 application/json` |

### Parameters

- `name` `(string: "")` - Authenticate against only the named certificate role, 
  returning its policy list if successful. If not set, defaults to trying all
  certificate roles and returning any one that matches.

### Sample Payload

```json
{
  "name": "cert1"
}
```

### Sample Request

```
$ curl \
    --request POST \
    --date @payload.json \
    https://vault.rocks/v1/auth/cert/login
```

### Sample Response

```json
{
  "auth": {
    "client_token": "cf95f87d-f95b-47ff-b1f5-ba7bff850425",
    "policies": [
      "web", 
      "stage"
    ],
    "lease_duration": 3600,
    "renewable": true,
  }
}
```