---
layout: "api"
page_title: "KV - Secrets Engines - HTTP API"
sidebar_title: "K/V Version 2"
sidebar_current: "api-http-secret-kv-v2"
description: |-
  This is the API documentation for the Vault KV secrets engine.
---

# KV Secrets Engine - Version 2 (API)

This is the API documentation for the Vault KV secrets engine while running in
versioned mode. For general information about the usage and operation of the kv
secrets engine, please see the [Vault kv
documentation](/docs/secrets/kv/index.html).

This documentation assumes the kv secrets engine is enabled at the
`/secret` path in Vault and that versioning has been enabled. Since it is
possible to enable secrets engines at any location, please update your API calls
accordingly.


## Configure the KV Engine 

This path configures backend level settings that are applied to every key in the
key-value store.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/secret/config`             |

### Parameters

- `max_versions` `(int: 0)` – The number of versions to keep per key. This value
  applies to all keys, but a key's metadata setting can overwrite this value.
  Once a key has more than the configured allowed versions the oldest version
  will be permanently deleted. Defaults to 10. 

- `cas_required` `(bool: false)` – If true all keys will require the cas
  parameter to be set on all write requests.

### Sample Payload

```json
{
  "max_versions": 5,
  "cas_required": false
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/secret/config
```

## Read KV Engine configuration

This path retrieves the current configuration for the secrets backend at the
given path.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`   | `/secret/config`             |


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://127.0.0.1:8200/v1/secret/config
```

### Sample Response

```json
{
  "data": {
    "cas_required": false,
    "max_versions": 0
  }
}
```


## Read Secret Version

This endpoint retrieves the secret at the specified location.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/secret/data/:path`         |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to read.
  This is specified as part of the URL.
- `version` `(int: 0)` - Specifies the version to return. If not set the latest
  version is returned.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://127.0.0.1:8200/v1/secret/data/my-secret
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

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/secret/data/:path`         |

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
  "options": {
      "cas": 0
  },
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
    https://127.0.0.1:8200/v1/secret/data/my-secret
```

### Sample Response

```json
{
  "data": {
    "created_time": "2018-03-22T02:36:43.986212308Z",
    "deletion_time": "",
    "destroyed": false,
    "version": 1
  }
}
```

## Delete Latest Version of Secret

This endpoint issues a soft delete of the secret's latest version at the
specified location. This marks the version as deleted and will stop it from
being returned from reads, but the underlying data will not be removed. A
delete can be undone using the `undelete` path.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/secret/data/:path`         |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to delete.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://127.0.0.1:8200/v1/secret/data/my-secret
```

## Delete Secret Versions

This endpoint issues a soft delete of the specified versions of the secret. This
marks the versions as deleted and will stop them from being returned from reads,
but the underlying data will not be removed. A delete can be undone using the
`undelete` path.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/secret/delete/:path`       |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to delete.
  This is specified as part of the URL.
- `versions` `([]int: <required>)` - The versions to be deleted. The versioned
  data will not be deleted, but it will no longer be returned in normal get
  requests.

### Sample Payload

```json
{
    "versions": [1, 2]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/secret/delete/my-secret
```

## Undelete Secret Versions

Undeletes the data for the provided version and path in the key-value store.
This restores the data, allowing it to be returned on get requests.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/secret/undelete/:path`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to undelete.
  This is specified as part of the URL.
- `versions` `([]int: <required>)` - The versions to undelete. The versions will
  be restored and their data will be returned on normal get requests.

### Sample Payload

```json
{
    "versions": [1, 2]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/secret/undelete/my-secret
```

## Destroy Secret Versions

Permanently removes the specified version data for the provided key and version
numbers from the key-value store.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/secret/destroy/:path`      |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to destroy.
  This is specified as part of the URL.
- `versions` `([]int: <required>)` - The versions to destroy. Their data will be
  permanently deleted. 

### Sample Payload

```json
{
    "versions": [1, 2]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/secret/destroy/my-secret
```

## List Secrets

This endpoint returns a list of key names at the specified location. Folders are
suffixed with `/`. The input must be a folder; list on a file will not return a
value. Note that no policy-based filtering is performed on keys; do not encode
sensitive information in key names. The values themselves are not accessible via
this command.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/secret/metadata/:path`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secrets to list.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://127.0.0.1:8200/v1/secret/metadata/my-secret
```

### Sample Response

The example below shows output for a query path of `secret/` when there are
secrets at `secret/foo` and `secret/foo/bar`; note the difference in the two
entries.

```json
{
  "data": {
    "keys": ["foo", "foo/"]
  },
}
```

## Read Secret Metadata

This endpoint retrieves the metadata and versions for the secret at the
specified path.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/secret/metadata/:path`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to read.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://127.0.0.1:8200/v1/secret/metadata/my-secret
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


## Update Metadata

This endpoint creates a new version of a secret at the specified location. If
the value does not yet exist, the calling token must have an ACL policy granting
the `create` capability. If the value already exists, the calling token must
have an ACL policy granting the `update` capability.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/secret/metadata/:path`     |

### Parameters

- `max_versions` `(int: 0)` – The number of versions to keep per key. If not
  set, the backend’s configured max version is used. Once a key has more than
  the configured allowed versions the oldest version will be permanently
  deleted. 

- `cas_required` `(bool: false)` – If true the key will require the cas
  parameter to be set on all write requests. If false, the backend’s
  configuration will be used. 

### Sample Payload

```json
{
  "max_versions": 5,
  "cas_required": false	
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/secret/metadata/my-secret
```

## Delete Metadata and All Versions

This endpoint permanently deletes the key metadata and all version data for the
specified key. All version history will be removed.  

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE`   | `/secret/metadata/:path`   |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to delete.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://127.0.0.1:8200/v1/secret/metadata/my-secret
```
