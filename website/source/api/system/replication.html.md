---
layout: "api"
page_title: "/sys/replication - HTTP API"
sidebar_current: "docs-http-system-replication"
description: |-
  The '/sys/replication' endpoint focuses on managing general operations in Vault Enterprise replication
---

# `/sys/replication`

~> **Enterprise Only** – These endpoints require Vault Enterprise.

## Attempt Recovery

This endpoint attempts recovery if replication is in an adverse state. For
example: an error has caused replication to stop syncing.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/recover`   | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/sys/replication/recover
```

### Sample Response

```json
{
  "warnings": ["..."]
}
```

## Reindex Replication

This endpoint reindexes the local data storage. This can cause a very long delay
depending on the number and size of objects in the data store.

**This endpoint requires 'sudo' capability.**

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/reindex`   | `200 application/json` |

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/sys/replication/reindex
```

### Sample Response

```json
{
  "warnings": ["..."]
}
```

## Check Status

This endpoint print information about the status of replication (mode,
sync progress, etc).

This is an authenticated endpoint.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/replication/status`    | `200 application/json` |

### Sample Request

```
$ curl \
    https://vault.rocks/v1/sys/replication/status
```

### Sample Response

The printed status of the replication environment. As an example, for a
primary, it will look something like:

```json
{
	"request_id": "d13e9665-d610-fea0-357f-8d652aa308cb",
	"lease_id": "",
	"lease_duration": 0,
	"renewable": false,
	"data": {
		"dr": {
			"cluster_id": "a876f38b-7577-25ac-6007-277528c99a1a",
			"known_secondaries": [
				"2"
			],
			"last_wal": 43,
			"merkle_root": "86d67839f47045f7d24beb4f39b14504d15a146c",
			"mode": "dr-primary",
			"primary_cluster_addr": ""
		},
		"performance": {
			"cluster_id": "11ab01df-32ea-1d79-b4bc-8bc973c1b749",
			"known_secondaries": [
				"1"
			],
			"last_wal": 43,
			"merkle_root": "e0531d566b23403101b0868e85b63d6774ba0ef2",
			"mode": "perf-primary",
			"primary_cluster_addr": ""
		}
	},
	"warnings": null
}
```
