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

- **High Availability** – the Google Cloud storage backend supports high availability.
   Because GCS uses the time on the Vault node to implement
   the session lifetimes on its locks, significant clock skew across Vault nodes
   could cause contention issues on the lock.

- **Community Supported** – the Google Cloud storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "gcs" {
  bucket           = "my-storage-bucket"
  credentials_file = "/tmp/credentials.json"
  ha_enabled = "true"
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

- `ha_enabled` `(bool: false)` – Specifies whether this backend should be used
  to run Vault in high availability mode. This can also be provided via the
  environment variable `GCS_HA_ENABLED`.

This backend also supports the following high availability parameters. These are
discussed in more detail in the [HA concepts page](/docs/concepts/ha.html)

- `cluster_addr` `(string: "")` – Specifies the address to advertise to other
  Vault servers in the cluster for request forwarding. This can also be provided
  via the environment variable `VAULT_CLUSTER_ADDR`. This is a full URL, like
  `redirect_addr`, but Vault will ignore the scheme (all cluster members always
  use TLS with a private key/certificate)

- `disable_clustering` `(bool: false)` – Specifies whether clustering features
  such as request forwarding are enabled. Setting this to true on one Vault node
  will disable these features _only when that node is the active node_

- `redirect_addr` `(string: <required>)` – Specifies the address (full URL) to
  advertise to other Vault servers in the cluster for client redirection. This
  can also be provided via the environment variable `VAULT_REDIRECT_ADDR`.

## `gcs` Examples

### Default Example

This example shows a default configuration for the Google Cloud Storage backend.

```hcl
storage "gcs" {
  bucket           = "my-storage-bucket"
  credentials_file = "/tmp/credentials.json"
}
```

### Enabling High Availability

This example show enabling high availability for the GCS storage backend.

```hcl
storage "gcs" {
  ha_enabled    = "true"
  redirect_addr = "https://vault-leader.my-company.internal"
}
```

[adc]: https://developers.google.com/identity/protocols/application-default-credentials
[gcs]: https://cloud.google.com/storage/
[gcs-service-account]: https://cloud.google.com/compute/docs/access/service-accounts
[gcs-private-key]: https://cloud.google.com/storage/docs/authentication#generating-a-private-key
