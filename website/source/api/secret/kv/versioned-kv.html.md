---
layout: "api"
page_title: "KV - Secrets Engines - HTTP API"
sidebar_current: "docs-http-secret-kv"
description: |-
  This is the API documentation for the Vault KV secrets engine.
---

# Versioned KV Secrets Engine (API)

This is the API documentation for the Vault KV secrets engine while running in
versioned mode. For general information about the usage and operation of the kv
secrets engine, please see the [Vault kv
documentation](/docs/secrets/kv/index.html).

This documentation assumes the kv secrets engine is enabled at the
`/secret` path in Vault and that versioning has been enabled. Since it is
possible to enable secrets engines at any location, please update your API calls
accordingly.

## Read Secret Version

This endpoint retrieves the secret at the specified location.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/secret/data/:path`              | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to read.
  This is specified as part of the URL.
- `version` `(int: 0)` - Specifies the version to return. If not set the latest
  version is returned.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/secret/data/my-secret
```

### Sample Response

```json
{
  "data": {
    "data": {
        "foo": "bar"
    },
    "metadata": {
      "created_time": "2018-03-22T02:24:06.945319214Z",
      "deletion_time": "",
      "destroyed": false,
      "version": 1
    }
  },
}
```

## Create/Update Secret

This endpoint creates a new version of a secret at the specified location. If
the value does not yet exist, the calling token must have an ACL policy granting
the `create` capability. If the value already exists, the calling token must
have an ACL policy granting the `update` capability.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/secret/data/:path`         | `204 (empty body)`     |
| `PUT`    | `/secret/data/:path`         | `204 (empty body)`     |

### Parameters

- `options` `(Map: <optional>)` – An object that holds option settings.
    - `cas` `(int: <optional>)` - Set the "cas" value to use a Check-And-Set
      operation. If not set the write will be allowed. If set to 0 a write will
      only be allowed if the key doesn’t exist. If the index is non-zero the
      write will only be allowed if the key’s current version matches the
      version specified in the cas parameter.  

- `data` `(Map: <required>)` – The contents of the data map will be stored and
  returned on read.

### Sample Payload

```json
{
  "data": {
	  "foo": "bar",
	  "zip": "zap"
	}
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

### Sample Response
```
{
  "data": {
    "created_time": "2018-03-22T02:36:43.986212308Z",
    "deletion_time": "",
    "destroyed": false,
    "version": 1
  }
}
```

## Delete Secret

This endpoint issues a soft delete of the secret's latest version at the
specified location. This marks the version as deleted and will stop it from
being returned from reads, but the underlying data will not be removed. A
delete can be undone using the `undelete` path.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/secret/data/:path`              | `204 (empty body)`     |

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

## Read Secret Metadata

This endpoint retrieves the metadata and versions for the secret at the
specified location.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/secret/metadata/:path`     | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to read.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/secret/metadata/my-secret
```

### Sample Response

```json
{
  "data": {
    "created_time": "2018-03-22T02:24:06.945319214Z",
    "current_version": 3,
    "max_versions": 0,
    "oldest_version": 0,
    "updated_time": "2018-03-22T02:36:43.986212308Z",
    "versions": {
      "1": {
        "created_time": "2018-03-22T02:24:06.945319214Z",
        "deletion_time": "",
        "destroyed": false
      },
      "2": {
        "created_time": "2018-03-22T02:36:33.954880664Z",
        "deletion_time": "",
        "destroyed": false
      },
      "3": {
        "created_time": "2018-03-22T02:36:43.986212308Z",
        "deletion_time": "",
        "destroyed": false
      }
    }
  }
}
```

## List Secrets

This endpoint returns a list of key names at the specified location. Folders are
suffixed with `/`. The input must be a folder; list on a file will not return a
value. Note that no policy-based filtering is performed on keys; do not encode
sensitive information in key names. The values themselves are not accessible via
this command.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/secret/metadata/:path`              | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secrets to list.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/secret/metadata/my-secret
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

## Update Metadata

This endpoint creates a new version of a secret at the specified location. If
the value does not yet exist, the calling token must have an ACL policy granting
the `create` capability. If the value already exists, the calling token must
have an ACL policy granting the `update` capability.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/secret/metadata/:path`         | `204 (empty body)`     |
| `PUT`    | `/secret/metadata/:path`         | `204 (empty body)`     |

### Parameters

- `max_versions` `(int: 0)` – The number of versions to keep. If not
  set, the backend’s configured max version is used. 

- `cas_required` `(bool: false)` – If true the key will require the cas
  parameter to be set on all write requests. If false, the backend’s
  configuration will be used. 

### Sample Payload

```json
{
  "data": {
	  "max_versions": 5,
	  "cas_required": false
	}
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/secret/metadata/my-secret
```
