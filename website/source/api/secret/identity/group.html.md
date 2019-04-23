---
layout: "api"
page_title: "Identity Secret Backend: Group - HTTP API"
sidebar_title: "Group"
sidebar_current: "api-http-secret-identity-group"
description: |-
  This is the API documentation for managing groups in the identity store.
---

## Create a Group

This endpoint creates or updates a Group.

| Method   | Path                |
| :------------------ | :----------------------|
| `POST`   | `/identity/group`   |

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
    http://127.0.0.1:8200/v1/identity/group
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

| Method   | Path                        |
| :-------------------------- | :--------------------- |
| `GET`    | `/identity/group/id/:id`    |

### Parameters

- `id` `(string: <required>)` – Identifier of the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/identity/group/id/363926d8-dd8b-c9f0-21f8-7b248be80ce1
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

| Method   | Path                        |
| :-------------------------- | :--------------------- |
| `POST`    | `/identity/group/id/:id`   |

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
    http://127.0.0.1:8200/v1/identity/group/id/363926d8-dd8b-c9f0-21f8-7b248be80ce1
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

| Method     | Path                       |
| :------------------------- | :----------------------|
| `DELETE`   | `/identity/group/id/:id`   |

## Parameters

- `id` `(string: <required>)` – Identifier of the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/identity/group/id/363926d8-dd8b-c9f0-21f8-7b248be80ce1
```

## List Groups by ID

This endpoint returns a list of available groups by their identifiers.

| Method   | Path                           |
| :----------------------------- | :--------------------- |
| `LIST`   | `/identity/group/id`           |
| `GET`    | `/identity/group/id?list=true` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/identity/group/id
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

## Create/Update Group by Name

This endpoint is used to create or update a group by its name.

| Method   | Path                            |
| :------------------------------ | :--------------------- |
| `POST`    | `/identity/group/name/:name`   |


### Parameters

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
    http://127.0.0.1:8200/v1/identity/group/name/testgroupname
```

### Sample Response
```json
{
  "request_id": "b98b4a3d-a9f1-e151-11e1-ad91cfb08351",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "id": "5a3a04a0-0c3a-a4c3-74e8-26b1adbeaece",
    "name": "testgroupname"
  },
  "warnings": null
}
```

## Read Group by Name

This endpoint queries the group by its name.

| Method   | Path                            |
| :------------------------------ | :--------------------- |
| `GET`    | `/identity/group/name/:name`    |

### Parameters

- `name` `(string: <required>)` – Name of the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/identity/group/name/testgroupname
```

### Sample Response

```json
{
  "data": {
    "alias": {},
    "creation_time": "2018-09-19T22:02:04.395128091Z",
    "id": "5a3a04a0-0c3a-a4c3-74e8-26b1adbeaece",
    "last_update_time": "2018-09-19T22:02:04.395128091Z",
    "member_entity_ids": [],
    "member_group_ids": null,
    "metadata": {
      "foo": "bar"
    },
    "modify_index": 1,
    "name": "testgroupname",
    "parent_group_ids": null,
    "policies": [
      "grouppolicy1",
      "grouppolicy2"
    ],
    "type": "internal"
  }
}
```

## Delete Group by Name

This endpoint deletes a group, given its name.

| Method     | Path                           |
| :----------------------------- | :----------------------|
| `DELETE`   | `/identity/group/name/:name`   |

## Parameters

- `name` `(string: <required>)` – Name of the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/identity/group/name/testgroupname
```

## List Groups by Name

This endpoint returns a list of available groups by their names.

| Method   | Path                             |
| :------------------------------- | :--------------------- |
| `LIST`   | `/identity/group/name`           |
| `GET`    | `/identity/group/name?list=true` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/identity/group/name
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "testgroupname"
    ]
  }
}
```
