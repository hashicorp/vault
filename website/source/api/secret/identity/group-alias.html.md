---
layout: "api"
page_title: "Identity Secret Backend: Group Alias - HTTP API"
sidebar_title: "Group Alias"
sidebar_current: "api-http-secret-identity-group-alias"
description: |-
  This is the API documentation for managing the group aliases in the identity store.
---

## Create a Group Alias

This endpoint creates or updates a group alias.

| Method   | Path                     |
| :----------------------- | :----------------------|
| `POST`   | `/identity/group-alias`  |

### Parameters

- `name` `(string: entity-<UUID>)` – Name of the group alias.

- `id` `(string: <optional>)` - ID of the group alias. If set, updates the
  corresponding existing group alias.

- `mount_accessor` `(string: "")` – Mount accessor which this alias belongs
  to.

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
    http://127.0.0.1:8200/v1/identity/group-alias
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

## Update Group Alias by ID

This endpoint is used to update an existing group alias.

| Method   | Path                              |
| :-------------------------------- | :--------------------- |
| `POST`    | `/identity/group-alias/id/:id`   |

### Parameters

- `id` `(string: <optional>)` - ID of the group alias.

- `name` `(string: entity-<UUID>)` – Name of the group alias.

- `mount_accessor` `(string: "")` – Mount accessor which this alias belongs
  to.

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
    http://127.0.0.1:8200/v1/identity/group-alias/id/ca726050-d8ac-6f1f-4210-3b5c5b613824
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

| Method   | Path                              |
| :-------------------------------- | :--------------------- |
| `GET`    | `/identity/group-alias/id/:id`    |

### Parameters

- `id` `(string: <required>)` – ID of the group alias.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/identity/group-alias/id/ca726050-d8ac-6f1f-4210-3b5c5b613824
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

| Method     | Path                             |
| :------------------------------- | :----------------------|
| `DELETE`   | `/identity/group-alias/id/:id`   |

## Parameters

- `id` `(string: <required>)` – ID of the group alias.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/identity/group-alias/id/ca726050-d8ac-6f1f-4210-3b5c5b613824
```

## List Group Alias by ID

This endpoint returns a list of available group aliases by their identifiers.

| Method   | Path                                      |
| :---------------------------------------- | :--------------------- |
| `LIST`   | `/identity/group-alias/id`                |
| `GET`    | `/identity/group-alias/id?list=true`      |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/identity/group-alias/id
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
