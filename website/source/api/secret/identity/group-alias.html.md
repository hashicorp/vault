---
layout: "api"
page_title: "Identity Secret Backend: Group Alias - HTTP API"
sidebar_current: "docs-http-secret-identity-group-alias"
description: |-
  This is the API documentation for managing the group aliases in the identity store.
---

## Create a Group Alias

This endpoint creates or updates a group alias.

| Method   | Path                     | Produces               |
| :------- | :----------------------- | :----------------------|
| `POST`   | `/identity/group-alias`  | `200 application/json` |

### Parameters

- `name` `(string: entity-<UUID>)` – Name of the group alias.

- `id` `(string: <optional>)` - ID of the group alias. If set, updates the
  corresponding existing group alias.

- `mount_accessor` `(string: "")` – Mount accessor to which this alias belongs
  toMount accessor to which this alias belongs to.

- `canonical_id` `(string: "")` - ID of the group to which this is an alias.


### Sample Payload

```json
{
  "canonical_id": "b86920ea-2831-00ff-15c5-a3f923f1ee3b",
  "mount_accessor": "auth_github_232a90dc",
  "name": "dev-team"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/identity/group-alias
```

### Sample Response

```json
{
  "data": {
    "canonical_id": "b86920ea-2831-00ff-15c5-a3f923f1ee3b",
    "id": "ca726050-d8ac-6f1f-4210-3b5c5b613824"
  }
}
```

## Read Group Alias by ID

This endpoint queries the group alias by its identifier.

| Method   | Path                              | Produces               |
| :------- | :-------------------------------- | :--------------------- |
| `GET`    | `/identity/group-alias/id/:id`    | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – ID of the group alias.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/identity/group-alias/id/ca726050-d8ac-6f1f-4210-3b5c5b613824
```

### Sample Response

```json
{
  "data": {
    "canonical_id": "b86920ea-2831-00ff-15c5-a3f923f1ee3b",
    "creation_time": "2017-11-13T20:09:41.661694Z",
    "id": "ca726050-d8ac-6f1f-4210-3b5c5b613824",
    "last_update_time": "2017-11-13T20:09:41.661694Z",
    "merged_from_canonical_ids": null,
    "metadata": null,
    "mount_accessor": "auth_github_232a90dc",
    "mount_path": "",
    "mount_type": "github",
    "name": "dev-team"
  }
}
```

## Delete Group Alias by ID

This endpoint deletes a group alias.

| Method     | Path                             | Produces               |
| :--------- | :------------------------------- | :----------------------|
| `DELETE`   | `/identity/group-alias/id/:id`   | `204 (empty body)`     |

## Parameters

- `id` `(string: <required>)` – ID of the group alias.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/identity/group-alias/id/ca726050-d8ac-6f1f-4210-3b5c5b613824
```

## List Entities by ID

This endpoint returns a list of available group aliases by their identifiers.

| Method   | Path                                 | Produces               |
| :------- | :----------------------------------- | :--------------------- |
| `LIST`   | `/identity/group-alias/id`           | `200 application/json` |
| `GET`    | `/identity/entity/id?list=true`      | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/identity/group-alias/id
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "ca726050-d8ac-6f1f-4210-3b5c5b613824"
    ]
  }
}
```
