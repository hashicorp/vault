---
layout: "docs"
page_title: "Filesystem - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-filesystem"
description: |-
  The Filesystem storage backend stores Vault's data on the filesystem using a
  standard directory structure. It can be used for durable single server
  situations, or to develop locally where durability is not critical.
---

# Filesystem Storage Backend

The Filesystem storage backend stores Vault's data on the filesystem using a
standard directory structure. It can be used for durable single server
situations, or to develop locally where durability is not critical.

- **No High Availability** – the Filesystem backend does not support high
  availability.

- **HashiCorp Supported** – the Filesystem backend is officially supported by
  HashiCorp.

```hcl
storage "file" {
  path = "/mnt/vault/data"
}
```

Even though Vault's data is encrypted at rest, you should still take appropriate
measures to secure access to the filesystem.

## `file` Parameters

- `path` `(string: <required>)` – The absolute path on disk to the directory
  where the data will be stored. If the directory does not exist, Vault will
  create it.

## `file` Examples

This example shows the Filesytem storage backend being mounted at
`/mnt/vault/data`.

```hcl
storage "file" {
  path = "/mnt/vault/data"
}
```
