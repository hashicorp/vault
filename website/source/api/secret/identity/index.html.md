---
layout: "api"
page_title: "Identity Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-identity"
description: |-
  This is the API documentation for the Vault Identity secret backend.
---

# Identity Secret Backend HTTP API

This is the API documentation for the Vault Identity secret backend. For
general information about the usage and operation of the Identity backend,
please see the
[Vault Identity backend documentation](/docs/secrets/identity/index.html).

## Register Entity

This endpoint creates or updates an Entity.

| Method   | Path                | Produces               |
| :------- | :------------------ | :----------------------|
| `POST`   | `/identity/entity`  | `200 application/json` |

### Parameters

- `name` `(string: entity-<UUID>)` – Name of the entity.

- `metadata` `(list of strings: [])` – Metadata to be associated with the entity. Format should be a list of `key=value` pairs.

- `policies` `(list of strings: [])` – Policies to be tied to the entity. Comma separated list of strings.

### Sample Payload

```json
{
	"metadata": ["organization=hashicorp", "team=vault"],
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
    "personas": null
  }
}
```

## Read Entity by ID

This endpoint queries the entity by its identifier.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/identity/entity/id/:id`    | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – Specifies the identifier of the entity.

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
    "personas": [],
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

- `id` `(string: <required>)` – Specifies the identifier of the entity.

- `name` `(string: entity-<UUID>)` – Name of the entity.

- `metadata` `(list of strings: [])` – Metadata to be associated with the entity. Format should be a list of `key=value` pairs.

- `policies` `(list of strings: [])` – Policies to be tied to the entity. Comma separated list of strings.


### Sample Payload

```json
{
	"name":"updatedEntityName",
	"metadata": ["organization=hashi", "team=nomad"],
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

```
{
  "data": {
    "id": "8d6a45e5-572f-8f13-d226-cd0d1ec57297",
    "personas": null
  }
}
```

## Delete Entity by ID

This endpoint deletes an entity and all its associated personas.

| Method     | Path                        | Produces               |
| :--------- | :-------------------------- | :----------------------|
| `DELETE`   | `/identity/entity/id/:id`   | `204 (empty body)`     |

## Parameters

- `id` `(string: <required>)` – Specifies the identifier of the entity.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/identity/entity/id/8d6a45e5-572f-8f13-d226-cd0d1ec57297
```

## List Entities by ID

This endpoint returns a list of available entities by their identifiers.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/identity/entity/id`        | `200 application/json` |

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

## Register Persona

This endpoint creates a new persona and attaches it to the entity with the
given identifier.

| Method   | Path                | Produces               |
| :------- | :------------------ | :----------------------|
| `POST`   | `/identity/persona`  | `200 application/json` |

### Parameters

- `name` (string: Required) - Name of the persona. Name should be the
  identifier of the client in the authentication source. For example, if the
  persona belongs to userpass backend, the name should be a valid username
  within userpass backend. If persona belongs to GitHub, it should be the
  GitHub username.

- `entity_id` (string: required) - Entity ID to which this persona belongs to.

- `mount_accessor` (string: required) - Accessor of the mount to which the
  persona should belong to.

- `metadata` `(list of strings: [])` – Metadata to be associated with the persona. Format should be a list of `key=value` pairs.

### Sample Payload

```
{
	"name": "testuser",
	"metadata": ["group=san_francisco", "region=west"],
	"entity_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
	"mount_accessor": "auth_userpass_e50b1a44"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/identity/persona
```

### Sample Response

```
{
  "data": {
    "entity_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
    "id": "34982d3d-e3ce-5d8b-6e5f-b9bb34246c31"
  }
}
```

## Read Persona by ID

This endpoint queries the persona by its identifier.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/identity/persona/id/:id`   | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – Specifies the identifier of the persona.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/identity/persona/id/34982d3d-e3ce-5d8b-6e5f-b9bb34246c31
```

### Sample Response

```
{
  "data": {
    "creation_time": "2017-07-25T21:41:09.820717636Z",
    "entity_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
    "id": "34982d3d-e3ce-5d8b-6e5f-b9bb34246c31",
    "last_update_time": "2017-07-25T21:41:09.820717636Z",
    "metadata": {
      "group": "san_francisco",
      "region": "west"
    },
    "mount_accessor": "auth_userpass_e50b1a44",
    "mount_path": "userpass/",
    "mount_type": "userpass",
    "name": "testuser"
  }
}
```

## Update Persona by ID

This endpoint is used to update an existing persona.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`    | `/identity/persona/id/:id`  | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – Specifies the identifier of the entity.

- `name` (string: Required) - Name of the persona. Name should be the
  identifier of the client in the authentication source. For example, if the
  persona belongs to userpass backend, the name should be a valid username
  within userpass backend. If persona belongs to GitHub, it should be the
  GitHub username.

- `entity_id` (string: required) - Entity ID to which this persona belongs to.

- `mount_accessor` (string: required) - Accessor of the mount to which the
  persona should belong to.

- `metadata` `(list of strings: [])` – Metadata to be associated with the
  persona. Format should be a list of `key=value` pairs.

### Sample Payload

```
{
	"name": "testuser",
	"metadata": ["group=philadelphia", "region=east"],
	"entity_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
	"mount_accessor": "auth_userpass_e50b1a44"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/identity/persona/id/34982d3d-e3ce-5d8b-6e5f-b9bb34246c31
```

### Sample Response

```
{
  "data": {
    "entity_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
    "id": "34982d3d-e3ce-5d8b-6e5f-b9bb34246c31"
  }
}
```

### Delete Persona by ID

This endpoint deletes a persona from its corresponding entity.

| Method     | Path                        | Produces               |
| :--------- | :-------------------------- | :----------------------|
| `DELETE`   | `/identity/persona/id/:id`  | `204 (empty body)`     |

## Parameters

- `id` `(string: <required>)` – Specifies the identifier of the persona.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/identity/persona/id/34982d3d-e3ce-5d8b-6e5f-b9bb34246c31
```

### List Personas by ID

This endpoint returns a list of available personas by their identifiers.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/identity/persona/id`       | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/identity/persona/id
```

### Sample Response

```
{
  "data": {
    "keys": [
      "2e8217fa-8cb6-8aec-9e22-3196d74ca2ba",
      "91ebe973-ec86-84db-3c7c-f760415326de",
      "92308b08-4139-3ec6-7af2-8e98166b4e0c",
      "a3b042e6-5cc1-d5a9-8874-d53a51954de2",
      "d5844921-017f-e496-2a9a-23d4a2f3e8a3"
    ]
  }
}
```

