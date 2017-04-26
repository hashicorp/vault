---
layout: "docs"
page_title: "Azure - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-azure"
description: |-
  The Azure storage backend is used to persist Vault's data in an Azure Storage
  Container. The storage container must already exist and the provided account
  credentials must have read and write  permissions to the storage container.
---

# Azure Storage Backend

The Azure storage backend is used to persist Vault's data in an
[Azure Storage Container][azure-storage]. The storage container must already
exist and the provided account credentials must have read and write permissions
to the storage container.

- **No High Availability** – the Azure storage backend does not support high
  availability.

- **Community Supported** – the Azure storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "azure" {
  accountName = "my-storage-account"
  accountKey  = "abcd1234"
  container   = "container-efgh5678"
}
```

The current implementation is limited to a maximum of 4 megabytes per blob.

## `azure` Parameters

- `accountName` `(string: <required>)` – Specifies the Azure Storage account
  name.

- `accountKey` `(string: <required>)` – Specifies the Azure Storage account key.

- `container` `(string: <required>)` – Specifies the Azure Storage Blob
  container name.

- `max_parallel` `(string: "128")` – Specifies The maximum number of concurrent
  requests to Azure.

## `azure` Examples

This example shows configuring the Azure storage backend with a custom number of
maximum parallel connections.

```hcl
storage "azure" {
  accountName  = "my-storage-account"
  accountKey   = "abcd1234"
  container    = "container-efgh5678"
  max_parallel = 512
}
```

[azure-storage]: https://azure.microsoft.com/en-us/services/storage/
