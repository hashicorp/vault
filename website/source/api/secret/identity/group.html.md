---
layout: "api"
page_title: "Identity Secret Backend: Group - HTTP API"
sidebar_current: "docs-http-secret-identity-group"
description: |-
  This is the API documentation for managing groups in the identity store.
---

## Create a Group

This endpoint creates or updates a Group.

| Method   | Path                | Produces               |
| :------- | :------------------ | :----------------------|
| `POST`   | `/identity/group`   | `200 application/json` |

### Parameters

- `name` `(string: entity-<UUID>)` – Name of the group.

- `id` `(string: <optional>)` - ID of the group. If set, updates the
  corresponding existing group.

- `type` `(string: "internal")` - Type of the group, `internal` or `external`.
  Defaults to `internal`.

- `metadata` `(key-value-map: {})` – Metadata to be associated with the
  group.

- `policies` `(list of strings: [])` – Policies to be tied to the group.

- `member_group_ids` `(list of strings: [])` -  Group IDs to be assigned as
  group members.

- `member_entity_ids` `(list of strings: [])` - Entity IDs to be assigned as
  group members.

### Sample Payload

```json
{
  "metadata": {
    "hello": "world"
  },
  "policies": ["grouppolicy1", "grouppolicy2"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/identity/group
```

### Sample Response

```json
{
  "data": {
    "id": "363926d8-dd8b-c9f0-21f8-7b248be80ce1",
    "name": "group_ab813d63"
  }
}
```

## Read Group by ID

This endpoint queries the group by its identifier.

| Method   | Path                        | Produces               |
| :------- | :-------------------------- | :--------------------- |
| `GET`    | `/identity/group/id/:id`    | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – Identifier of the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/identity/group/id/363926d8-dd8b-c9f0-21f8-7b248be80ce1
```

### Sample Response

```json
{
  "data": {
    "alias": {},
    "creation_time": "2017-11-13T19:36:47.102945Z",
    "id": "363926d8-dd8b-c9f0-21f8-7b248be80ce1",
    "last_update_time": "2017-11-13T19:36:47.102945Z",
    "member_entity_ids": [],
    "member_group_ids": null,
    "metadata": {
      "hello": "world"
    },
    "modify_index": 1,
    "name": "group_ab813d63",
    "policies": [
      "grouppolicy1",
      "grouppolicy2"
    ],
    "type": "internal"
  }
}
```

## Update Group by ID

This endpoint is used to update an existing group.

| Method   | Path                        | Produces               |
| :------- | :-------------------------- | :--------------------- |
| `POST`    | `/identity/group/id/:id`   | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – Identifier of the entity.

- `name` `(string: entity-<UUID>)` – Name of the group.

- `type` `(string: "internal")` - Type of the group, `internal` or `external`.
  Defaults to `internal`.

- `metadata` `(key-value-map: {})` – Metadata to be associated with the
  group.

- `policies` `(list of strings: [])` – Policies to be tied to the group.

- `member_group_ids` `(list of strings: [])` -  Group IDs to be assigned as
  group members.

- `member_entity_ids` `(list of strings: [])` - Entity IDs to be assigned as
  group members.

### Sample Payload

```json
{
  "name": "testgroupname",
    "metadata": {
      "hello": "everyone"
    },
  "policies": ["grouppolicy2", "grouppolicy3"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/identity/group/id/363926d8-dd8b-c9f0-21f8-7b248be80ce1
```

### Sample Response

```json
{
  "data": {
    "id": "363926d8-dd8b-c9f0-21f8-7b248be80ce1",
    "name": "testgroupname"
  }
}
```

## Delete Group by ID

This endpoint deletes a group.

| Method     | Path                       | Produces               |
| :--------- | :------------------------- | :----------------------|
| `DELETE`   | `/identity/group/id/:id`   | `204 (empty body)`     |

## Parameters

- `id` `(string: <required>)` – Identifier of the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/identity/group/id/363926d8-dd8b-c9f0-21f8-7b248be80ce1
```

## List Groups by ID

This endpoint returns a list of available groups by their identifiers.

| Method   | Path                           | Produces               |
| :------- | :----------------------------- | :--------------------- |
| `LIST`   | `/identity/group/id`           | `200 application/json` |
| `GET`    | `/identity/group/id?list=true` | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/identity/group/id
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "052567cf-1580-6f20-50c8-d38bc46dae6e",
      "26da8035-6691-b89e-67ac-ebf9ea7f9893",
      "363926d8-dd8b-c9f0-21f8-7b248be80ce1",
      "5c4a5720-7408-c113-1dcc-9ede725d0ac8",
      "d55e0f34-5c16-38ae-87af-324c9b656c43",
      "e4e56e04-0dec-9b68-9b20-a450975d898e"
    ]
  }
}
```
