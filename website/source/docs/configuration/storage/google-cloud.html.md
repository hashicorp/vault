---
layout: "docs"
page_title: "Google Cloud - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-google-cloud"
description: |-
  The Google Cloud storage backend is used to persist Vault's data in Google
  Cloud Storage.
---

# Google Cloud Storage Backend

The Google Cloud storage backend is used to persist Vault's data in
[Google Cloud Storage][gcs].

- **No High Availability** – the Google Cloud storage backend does not support
  high availability.

- **Community Supported** – the Google Cloud storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "gcs" {
  bucket           = "my-storage-bucket"
  credentials_file = "/tmp/credentials.json"
}
```

## `gcs` Parameters

- `bucket` `(string: <required>)` – Specifies the name of the Google Cloud
  Storage bucket to use. This bucket must already exist and the provided service
  account must have permission to read, write, and delete from the bucket. This
  can also be provided via the environment variable `GOOGLE_STORAGE_BUCKET`.

- `credentials_file` `(string: "<varies>")` – Specifies the path on disk to a
  Google Cloud Platform [service account][gcs-service-account] private key file
  in [JSON format][gcs-private-key]. The GCS client library will attempt to use
  the [application default credentials][adc] if this is not specified.

- `max_parallel` `(string: "128")` – Specifies the maximum number of concurrent
  requests.

## `gcs` Examples

### Default Example

This example shows a default configuration for the Google Cloud Storage backend.

```hcl
storage "gcs" {
  bucket           = "my-storage-bucket"
  credentials_file = "/tmp/credentials.json"
}
```

[adc]: https://developers.google.com/identity/protocols/application-default-credentials
[gcs]: https://cloud.google.com/storage/
[gcs-service-account]: https://cloud.google.com/compute/docs/access/service-accounts
[gcs-private-key]: https://cloud.google.com/storage/docs/authentication#generating-a-private-key
