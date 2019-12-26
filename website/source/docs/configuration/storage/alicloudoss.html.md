---
layout: "docs"
page_title: "Alicloud OSS - Storage Backends - Configuration"
sidebar_title: "AliCloud OSS"
sidebar_current: "docs-configuration-storage-alicloudoss"
description: |-
  The Alicloud OSS storage backend is used to persist Vault's data in
  an Alicloud OSS bucket.
---

# Alicloud OSS Storage Backend

The Alicloud OSS storage backend is used to persist Vault's data in
an [Alicloud OSS][alicloudoss] bucket.

- **No High Availability** – the Alicloud OSS storage backend does not support
  high availability.

- **Community Supported** – the Alicloud OSS storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "alicloudoss" {
  access_key = "abcd1234"
  secret_key = "defg5678"
  endpoint   = "oss-us-west-1.aliyuncs.com"
  bucket     = "my-bucket"
}
```

## `alicloudoss` Parameters

- `bucket` `(string: <required>)` – Specifies the name of the OSS bucket. This
  can also be provided via the environment variable `ALICLOUD_OSS_BUCKET`.

- `endpoint` `(string: "")` – Specifies the OSS endpoint. This can also be
 provided via the environment variable `ALICLOUD_OSS_ENDPOINT`.

The following settings are used for authenticating to Alicloud.

- `access_key` – Specifies the Alicloud access key. This can also be provided via
  the environment variable `ALICLOUD_ACCESS_KEY`.

- `secret_key` – Specifies the Alicloud secret key. This can also be provided via
  the environment variable `ALICLOUD_SECRET_KEY`.

- `max_parallel` `(string: "128")` – Specifies the maximum number of concurrent
  requests to Alicloud OSS.

## `alicloudoss` Examples

### Default Example

This example shows using Alicloud OSS as a storage backend.

```hcl
storage "alicloudoss" {
  access_key = "abcd1234"
  secret_key = "defg5678"
  endpoint   = "oss-us-west-1.aliyuncs.com"
  bucket     = "my-bucket"
}
```

[alicloudoss]: https://www.alibabacloud.com/product/oss
