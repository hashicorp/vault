---
layout: "docs"
page_title: "KV - Secrets Engines"
sidebar_current: "docs-secrets-kv"
description: |-
  The KV secrets engine can store arbitrary secrets.
---

# KV Secrets Engine

The `kv` secrets engine is used to store arbitrary secrets within the
configured physical storage for Vault. This backend can be run in one of two
modes. It can be a generic Key-Value store that stores one value for a key.
Versioning can be enabled and a configurable number of versions for each key
will be stored.

## KV Version 1

When running the `kv` secrets backend non-versioned only the most recently
written value for a key will be preserved. The benefits of non-versioned `kv`
is a reduced storage size for each key since no additional metadata or history
is stored. Additionally, requests going to a backend configured this way will be
more performant because for any given request there will be fewer storage calls
and no locking.

More information about running in this mode can be found in the [K/V Version 1
Docs](/docs/secrets/kv/kv-v1.html)

## KV Version 2

When running v2 of the `kv` backend a key can retain a configurable number of
versions. This defaults to 10 versions. The older versions' metadata and data
can be retrieved. Additionally, Check-and-Set operations can be used to avoid
overwritting data unintentionally.  

When a version is deleted the underlying data is not removed, rather it is
marked as deleted. Deleted versions can be undeleted. To permanently remove a
version's data the destroy command or API endpoint can be used. Additionally all
versions and metadata for a key can be deleted by deleting on the metadata
command or API endpoint. Each of these operations can be ACL'ed differently,
restricting who has permissions to soft delete, undelete, or fully remove data.

More information about running in this mode can be found in the [K/V Version 2
Docs](/docs/secrets/kv/kv-v2.html)
