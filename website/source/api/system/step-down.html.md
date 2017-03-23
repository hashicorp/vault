---
layout: "api"
page_title: "/sys/step-down - HTTP API"
sidebar_current: "docs-http-system-step-down"
description: |-
  The `/sys/step-down` endpoint causes the node to give up active status.
---

# `/sys/step-down`

The `/sys/step-down` endpoint causes the node to give up active status.

## Step Down Leader

This endpoint forces the node to give up active status. If the node does not
have active status, this endpoint does nothing. Note that the node will sleep
for ten seconds before attempting to grab the active lock again, but if no
standby nodes grab the active lock in the interim, the same node may become the
active node again. Requires a token with `root` policy or `sudo` capability on
the path.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/step-down`             | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/step-down
```
