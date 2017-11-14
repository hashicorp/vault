---
layout: "api"
page_title: "/sys/control-group - HTTP API"
sidebar_current: "docs-http-system-control-group"
description: |-
  The '/sys/control-group' endpoint handles the Control Group workflow.
---

## Authorize Control Group Request

This endpoint authorizes a control group request.

| Method   | Path                           | Produces               |
| :------- | :----------------------------- | :--------------------- |
| `POST`   | `/sys/control-group/authorize`   | `200 (application/json)`     |

### Parameters

- `accessor` `(string: <required>)` – The accessor for the control group wrapping token.

### Sample Payload

```json
{
  "accessor": "0ad21b78-e9bb-64fa-88b8-1e38db217bde",
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/control-group/authorize
```

### Sample Response

```json
{
    "data": {
        "approved": false
    }
}
```

## Check Control Group Request Status

This endpoint checks the status of a control group request.

| Method   | Path                           | Produces               |
| :------- | :----------------------------- | :--------------------- |
| `POST`   | `/sys/control-group/request`   | `200 (application/json)`     |

### Parameters

- `accessor` `(string: <required>)` – The accessor for the control group wrapping token.

### Sample Payload

```json
{
  "accessor": "0ad21b78-e9bb-64fa-88b8-1e38db217bde",
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/control-group/request
```

### Sample Response

```json
{
    "data": {
        "approved": false,
        "request_path": "secret/foo",
        "request_entity": {
                "id": "c8b6e404-de4b-50a4-2917-715ff8beec8e",
                "name": "Bob"
        },
        "authorizations": [
            {
                "entity_id": "6544a3ec-d3cd-443b-b87b-4fd2e889e0b7",
                "entity_name": "Abby Jones"
            },
            {
                "entity_id": "919084a4-417e-42ee-9d78-87fa2843af37",
                "entity_name": "James Franklin"
            }
        ]
    }
}
```
