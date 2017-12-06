---
layout: "docs"
page_title: "Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage"
description: |-
  The storage stanza configures the storage backend, which represents the
  location for the durable storage of Vault's information. Each backend has
  pros, cons, advantages, and trade-offs. For example, some backends support
  high availability while others provide a more robust backup and restoration
  process.
---

# `storage` Stanza

The `storage` stanza configures the storage backend, which represents the
location for the durable storage of Vault's information. Each backend has pros,
cons, advantages, and trade-offs. For example, some backends support high
availability while others provide a more robust backup and restoration process.
For information about a specific backend, choose one from the navigation on the
left.

## Configuration

Storage backend configuration is done through the Vault configuration file using
the `storage` stanza:

```hcl
storage [NAME] {
  [PARAMETERS...]
}
```

For example:

```hcl
storage "file" {
  path = "/mnt/vault/data"
}
```

For configuration options which also read an environment variable, the
environment variable will take precedence over values in the configuration
file.
