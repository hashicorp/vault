---
layout: "api"
page_title: "/identity/groups - HTTP API"
sidebar_current: "api-http-secret-identity-groups"
description: |-
  This is the API documentation for the identity groups.
---

## Create/Update Group

This endpoint creates or updates a group.

| Method   | Path                |
| :------------------ | :----------------------|
| `POST`   | `/identity/group`   |

### Parameters

- `name` `(string: group-<UUID>)` – Name of the group.

- `id` `(string: "")` - ID of the group. If this is set, this endpoint will
  update the corresponding group.

- `metadata` `(list of strings: [])` – Metadata to be associated with the group. Format should be a list of `key=value` pairs.

- `policies` `(list of strings: [])` – Policies to be tied to the group. Comma separated list of strings.

- `member_group_ids` `(list of strings: [])` - Group IDs to be assigned as group members.

- `member_entity_ids` `(list of strings: [])` - Entity IDs to be assigned as group members.

### Sample Payload

```json
{
    "name": "engineering-group",
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
    http://127.0.0.1:8200/v1/identity/group
```

### Sample Response

```json
{
        "data": {
                "id": "454ceeb5-76d7-a131-b92a-7ecfb15523e8",
                "name": "engineering-group"
        }
}
```

## Update Group by ID

This endpoint updates the group by its ID.

| Method   | Path                       |
| :------------------------- | :----------------------|
| `POST`   | `/identity/group/id/:id`   |

### Parameters

- `id` `(string: "")` - ID of the group.

- `name` `(string: group-<UUID>)` – Name of the group.

- `metadata` `(list of strings: [])` – Metadata to be associated with the group. Format should be a list of `key=value` pairs.

- `policies` `(list of strings: [])` – Policies to be tied to the group. Comma separated list of strings.

- `member_group_ids` `(list of strings: [])` - Group IDs to be assigned as group members.

- `member_entity_ids` `(list of strings: [])` - Entity IDs to be assigned as group members.

### Sample Payload

```json
{
	    "metadata": ["organization=updatedorg", "team=updatedteam"],
	    "policies": ["updatedpolicy"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/identity/group/id/454ceeb5-76d7-a131-b92a-7ecfb15523e8
```

### Sample Response

```json
{
        "data": {
                "id": "454ceeb5-76d7-a131-b92a-7ecfb15523e8",
                "name": "engineering-group"
        }
}
```

## Read Group by ID

This endpoint reads the group by its ID.

| Method   | Path                       |
| :------------------------- | :--------------------- |
| `GET`    | `/identity/group/id/:id`   |

### Parameters

- `id` `(string: "")` - ID of the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
        http://127.0.0.1:8200/v1/identity/group/id/454ceeb5-76d7-a131-b92a-7ecfb15523e8
```

### Sample Response

```json
{
        "data": {
                "creation_time": "2017-09-13T01:17:26.755474204Z",
                "id": "454ceeb5-76d7-a131-b92a-7ecfb15523e8",
                "last_update_time": "2017-09-13T01:17:26.755474204Z",
                "member_entity_ids": [],
                "member_group_ids": null,
                "metadata": {
                        "organization": "hashicorp",
                        "team": "vault"
                },
                "modify_index": 1,
                "name": "engineering-group",
                "policies": [
                        "dev-policy"
                ]
        }
}
```

## Delete Group by ID

This endpoint deleted the group by its ID.

| Method     | Path                       |
| :------------------------- | :----------------------|
| `DELETE`   | `/identity/group/id/:id`   |

### Parameters

- `id` `(string: "")` - ID of the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/identity/group/id/454ceeb5-76d7-a131-b92a-7ecfb15523e8
```


## List Groups by ID

This endpoint lists all the groups by their ID.

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
      "454ceeb5-76d7-a131-b92a-7ecfb15523e8",
      "7b2fb80c-9516-68d1-35fc-11450f6477ab"
    ]
  }
}
```

## Lookup Group by ID

This endpoint queries the group by its ID.

| Method   | Path                       |
| :------------------------- | :----------------------|
| `POST`   | `/identity/lookup/group`   |

### Parameters

- `type` `(string: "")` - Type of query. Supported values are `by_id` and `by_name`.

- `group_name` `(string: "")` - Name of the group.

- `group_id` `(string: "")` - ID of the group.

### Sample Payload

```json
{
    "type": "by_id",
    "group_id": "454ceeb5-76d7-a131-b92a-7ecfb15523e8"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/identity/lookup/group
```

### Sample Response

```json
{
        "data": {
                "creation_time": "2017-09-13T01:17:26.755474204Z",
                "id": "454ceeb5-76d7-a131-b92a-7ecfb15523e8",
                "last_update_time": "2017-09-13T01:17:26.755474204Z",
                "member_entity_ids": [],
                "member_group_ids": null,
                "metadata": {
                        "organization": "hashicorp",
                        "team": "vault"
                },
                "modify_index": 1,
                "name": "engineering-group",
                "policies": [
                        "dev-policy"
                ]
        }
}
```
