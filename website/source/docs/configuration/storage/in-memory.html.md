---
layout: "docs"
page_title: "In-Memory - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-in-memory"
description: |-
  The In-Memory storage backend is used to persist Vault's data entirely
  in-memory on the same machine in which Vault is running. This is useful for
  development and experimentation, but use of this backend is highly discouraged
  in production except in very specific use-cases.
---

# In-Memory Storage Backend

The In-Memory storage backend is used to persist Vault's data entirely in-memory
on the same machine in which Vault is running. This is useful for development
and experimentation, but use of this backend is **highly discouraged in
production**. All data is lost when Vault or the machine on which it is running
is restarted.

- **No High Availability** – the In-Memory backend does not support high
  availability.

- **Not Production Recommended** – the In-Memory backend is not recommended for
  production installations as data does not persist beyond restarts.

- **HashiCorp Supported** – the In-Memory backend is officially supported by
  HashiCorp.

```hcl
storage "inmem" {}
```

## `inmem` Parameters

The In-Memory storage backend has no configuration parameters.

## `inmem` Examples

This example shows activating the In-Memory storage backend.

```hcl
storage "inmem" {}
```
