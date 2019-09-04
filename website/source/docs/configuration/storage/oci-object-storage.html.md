---
layout: "docs"
page_title: "OCI ObjectStorage - Storage Backends - Configuration"
sidebar_title: "OCI ObjectStorage"
sidebar_current: "docs-configuration-storage-oci-objectstorage"
description: |-
  The OCI ObjectStorage backend is used to persist Vault's data in OCI object storage.
---

# OCI ObjectStorage Storage Backend

The OCI ObjectStorage backend is used to persist Vault's data in OCI object storage.

- **High Availability** – the OCI ObjectStorage backend supports high availability.

- **Community Supported** – the OCI ObjectStorage backend is supported by the community. While it has undergone review by HashiCorp employees, they may not be as knowledgeable about the technology. If you encounter problems with them, you may be referred to the original author.

```hcl
storage "oci_objectstorage" {
   namespace_name = "<object_storage_namespace_name>"
   bucket_name = "<vault_data_bucket_name>"
   ha_enabled = "<boolean true/false>"
   lock_bucket_name = "<leader_lock_bucket_name>"
   auth_type_api_key = "<boolean setting for using api-key instead of instance principals>"
}
```

For more information on OCI object storage, please see the Oracle's [OCI object storage documentation][ocios-docs].


## `oci_objectstorage` Setup

To use the OCI ObjectStorage Vault storage backend, you must have a OCI account. Either using the API or web interface, create the data bucket and lock bucket if enabling high availability.

The OCI object storage backend does not support creating the buckets automatically at this time.


## `oci_objectstorage` Authentication

The OCI ObjectStorage Vault storage backend uses the official OCI Golang SDK. This means it supports the common ways of providing credentials to OCI.

For more information on service accounts, please see the [OCI Identity documentation] [oci-identity].

## `oci_objectstorage` Parameters

- `namespace_name` `(string: <required>)` – Specifies the name of the ObjectStorage namespaces containing the data bucket and the lock bucket.

- `bucket_name` `(string: <required>)` - Specifies the name of the bucket that will be used to store the vault data.


### High Availability Parameters

- `ha_enabled` `(string: "<required>")` - Specifies if high availability mode is
  enabled. This is a boolean value, but it is specified as a string like "true"
  or "false".

- `lock_bucket_name` `(string: "<required>")` - Specifies the name of the bucket that will be used to store the node lease data.

## `oci_objectstorage` Examples

### Standalone vault instance

This example shows configuring OCI ObjectStorage as a standalone instance.

```hcl
storage "oci_objectstorage" {
    namespace_name = "MyNamespace
    bucket_name = "DataBucket"
}
```

### High Availability

This example shows configuring OCI ObjectStorage with high availability enabled.

```hcl
storage "oci_objectstorage" {
   namespaceName = "MyNamespace
   bucketName = "DataBucket"
   ha_enabled = "true"
   lockBucketName = "LockBucket"
}
```

[oci-identity]: https://docs.cloud.oracle.com/iaas/Content/Identity/Concepts/overview.htm
[ocios-docs]: https://docs.cloud.oracle.com/iaas/Content/Object/Concepts/objectstorageoverview.htm
