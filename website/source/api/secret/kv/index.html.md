---
layout: "api"
page_title: "Key/Value Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-kv"
description: |-
  This is the API documentation for the Vault Key/Value secret backend.
---

# Key/Value Secret Backend HTTP API

This is the API documentation for the Vault Key/Value secret backend. For general
information about the usage and operation of the Key/Value backend, please see
the [Vault Key/Value backend documentation](/docs/secrets/kv/index.html).

This documentation assumes the Key/Value backend is mounted at the `/secret`
path in Vault. Since it is possible to mount secret backends at any location,
please update your API calls accordingly.

## Read Secret

This endpoint retrieves the secret at the specified location.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/secret/:path`              | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to read.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/secret/my-secret
```

### Sample Response

```json
{
  "auth": null,
  "data": {
    "foo": "bar"
  },
  "lease_duration": 2764800,
  "lease_id": "",
  "renewable": false
}
```

_Note_: the `lease_duration` field (which on the CLI shows as
`refresh_interval`) is advisory. No lease is created. This is a way for writers
to indicate how often a given value shold be re-read by the client. See the
[Vault Key/Value backend documentation](/docs/secrets/kv/index.html) for
more details.

## List Secrets

This endpoint returns a list of key names at the specified location. Folders are
suffixed with `/`. The input must be a folder; list on a file will not return a
value. Note that no policy-based filtering is performed on keys; do not encode
sensitive information in key names. The values themselves are not accessible via
this command.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/secret/:path`              | `200 application/json` |
| `GET`    | `/secret/:path?list=true`    | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secrets to list.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/secret/my-secret
```

### Sample Response

The example below shows output for a query path of `secret/` when there are
secrets at `secret/foo` and `secret/foo/bar`; note the difference in the two
entries.

```json
{
  "auth": null,
  "data": {
    "keys": ["foo", "foo/"]
  },
  "lease_duration": 2764800,
  "lease_id": "",
  "renewable": false
}
```

## Create/Update Secret

This endpoint stores a secret at the specified location. If the value does not
yet exist, the calling token must have an ACL policy granting the `create`
capability. If the value already exists, the calling token must have an ACL
policy granting the `update` capability.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/secret/:path`              | `204 (empty body)`     |
| `PUT`    | `/secret/:path`              | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secrets to
  create/update. This is specified as part of the URL.

- `:key` `(string: "")` – Specifies a key, paired with an associated value, to
  be held at the given location. Multiple key/value pairs can be specified, and
  all will be returned on a read operation. A key called `ttl` will trigger
  some special behavior; see the [Vault Key/Value backend
  documentation](/docs/secrets/kv/index.html) for details.

### Sample Payload

```json
{
  "foo": "bar",
  "zip": "zap"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/secret/my-secret
```

## Delete Secret

This endpoint deletes the secret at the specified location.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/secret/:path`              | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to delete.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/secret/my-secret
```
