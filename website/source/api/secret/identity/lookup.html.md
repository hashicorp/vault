---
layout: "api"
page_title: "Identity Secret Backend: Lookup - HTTP API"
sidebar_current: "docs-http-secret-identity-lookup"
description: |-
  This is the API documentation for entity and group lookups from identity
  store.
---

## Lookup an Entity

This endpoint queries the entity based on the given criteria. The criteria can
be `name`, `id`, `alias_id`, or a combination of `alias_name` and
`alias_mount_accessor`.

| Method   | Path                       | Produces               |
| :------- | :------------------------- | :----------------------|
| `POST`   | `/identity/lookup/entity`  | `200 application/json` |

### Parameters

- `name` `(string: "")` – Name of the entity.

- `id` `(string: "")` - ID of the entity.

- `alias_id` `(string: "")` - ID of the alias.

- `alias_name` `(string: "")` - Name of the alias. This should be supplied in
  conjunction with `alias_mount_accessor`.

- `alias_mount_accessor` `(string: "")` - Accessor of the mount to which the
  alias belongs to. This should be supplied in conjunction with `alias_name`.

### Sample Payload

```json
{
  "id": "043fedec-967d-b2c9-d3af-0c467b04e1fd"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/identity/lookup/entity
```

### Sample Response

```json
{
  "data": {
    "aliases": [],
    "creation_time": "2017-11-13T21:01:33.543497Z",
    "direct_group_ids": [],
    "group_ids": [],
    "id": "043fedec-967d-b2c9-d3af-0c467b04e1fd",
    "inherited_group_ids": [],
    "last_update_time": "2017-11-13T21:01:33.543497Z",
    "merged_entity_ids": null,
    "metadata": null,
    "name": "entity_43cc451b",
    "policies": null
  }
}
```

## Lookup a Group

This endpoint queries the group based on the given criteria. The criteria can
be `name`, `id`, `alias_id`, or a combination of `alias_name` and
`alias_mount_accessor`.

| Method   | Path                       | Produces               |
| :------- | :------------------------- | :----------------------|
| `POST`   | `/identity/lookup/group`   | `200 application/json` |

### Parameters

- `name` `(string: "")` – Name of the group.

- `id` `(string: "")` - ID of the group.

- `alias_id` `(string: "")` - ID of the alias.

- `alias_name` `(string: "")` - Name of the alias. This should be supplied in
  conjunction with `alias_mount_accessor`.

- `alias_mount_accessor` `(string: "")` - Accessor of the mount to which the
  alias belongs to. This should be supplied in conjunction with `alias_name`.

### Sample Payload

```json
{
  "id": "70a4bdef-9da3-4460-b524-bb08542eef25"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/identity/lookup/group
```

### Sample Response

```json
{
  "data": {
    "alias": {},
    "creation_time": "2017-11-13T21:06:44.475587Z",
    "id": "70a4bdef-9da3-4460-b524-bb08542eef25",
    "last_update_time": "2017-11-13T21:06:44.475587Z",
    "member_entity_ids": [],
    "member_group_ids": null,
    "metadata": null,
    "modify_index": 1,
    "name": "group_eaf2aab1",
    "policies": null,
    "type": "internal"
  }
}
```
