---
layout: "api"
page_title: "Identity Secret Backend: Entity - HTTP API"
sidebar_current: "docs-http-secret-identity-entity"
description: |-
  This is the API documentation for managing entities in the identity store.
---

## Create an Entity

This endpoint creates or updates an Entity.

| Method   | Path                | Produces               |
| :------- | :------------------ | :----------------------|
| `POST`   | `/identity/entity`  | `200 application/json` |

### Parameters

- `name` `(string: entity-<UUID>)` – Name of the entity.

- `id` `(string: <optional>)` - ID of the entity. If set, updates the
  corresponding existing entity.

- `metadata` `(key-value-map: {})` – Metadata to be associated with the
  entity.

- `policies` `(list of strings: [])` – Policies to be tied to the entity.

### Sample Payload

```json
{
  "metadata": {
  "organization": "hashicorp",
    "team": "vault"
  },
  "policies": ["eng-dev", "infra-dev"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/identity/entity
```

### Sample Response

```json
{
  "data": {
    "id": "8d6a45e5-572f-8f13-d226-cd0d1ec57297",
    "aliases": null
  }
}
```

## Read Entity by ID

This endpoint queries the entity by its identifier.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/identity/entity/id/:id`    | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – Identifier of the entity.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/identity/entity/id/8d6a45e5-572f-8f13-d226-cd0d1ec57297
```

### Sample Response

```json
{
  "data": {
    "bucket_key_hash": "177553e4c58987f4cc5d7e530136c642",
    "creation_time": "2017-07-25T20:29:22.614756844Z",
    "id": "8d6a45e5-572f-8f13-d226-cd0d1ec57297",
    "last_update_time": "2017-07-25T20:29:22.614756844Z",
    "metadata": {
      "organization": "hashicorp",
      "team": "vault"
    },
    "name": "entity-c323de27-2ad2-5ded-dbf3-0c7ef98bc613",
    "aliases": [],
    "policies": [
      "eng-dev",
      "infra-dev"
    ]
  }
}
```

## Update Entity by ID

This endpoint is used to update an existing entity.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`    | `/identity/entity/id/:id`   | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – Identifier of the entity.

- `name` `(string: entity-<UUID>)` – Name of the entity.

- `metadata` `(key-value-map: {})` – Metadata to be associated with the entity.

- `policies` `(list of strings: [])` – Policies to be tied to the entity.


### Sample Payload

```json
{
  "name":"updatedEntityName",
  "metadata": {
  "organization": "hashi",
    "team": "nomad"
  },
  "policies": ["eng-developers", "infra-developers"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/identity/entity/id/8d6a45e5-572f-8f13-d226-cd0d1ec57297
```

### Sample Response

```json
{
  "data": {
    "id": "8d6a45e5-572f-8f13-d226-cd0d1ec57297",
    "aliases": null
  }
}
```

## Delete Entity by ID

This endpoint deletes an entity and all its associated aliases.

| Method     | Path                        | Produces               |
| :--------- | :-------------------------- | :----------------------|
| `DELETE`   | `/identity/entity/id/:id`   | `204 (empty body)`     |

## Parameters

- `id` `(string: <required>)` – Identifier of the entity.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/identity/entity/id/8d6a45e5-572f-8f13-d226-cd0d1ec57297
```

## List Entities by ID

This endpoint returns a list of available entities by their identifiers.

| Method   | Path                            | Produces               |
| :------- | :------------------------------ | :--------------------- |
| `LIST`   | `/identity/entity/id`           | `200 application/json` |
| `GET`    | `/identity/entity/id?list=true` | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/identity/entity/id
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "02fe5a88-912b-6794-62ed-db873ef86a95",
      "3bf81bc9-44df-8138-57f9-724a9ae36d04",
      "627fba68-98c9-c012-71ba-bfb349585ce1",
      "6c4c805b-b384-3d0e-4d51-44d349887b96",
      "70a72feb-35d1-c775-0813-8efaa8b4b9b5",
      "f1092a67-ce34-48fd-161d-c13a367bc1cd",
      "faedd89a-0d82-c197-c8f9-93a3e6cf0cd0"
    ]
  }
}
```
