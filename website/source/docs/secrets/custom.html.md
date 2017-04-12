---
layout: "docs"
page_title: "Custom Secret Backend"
sidebar_current: "docs-secrets-custom"
description: |-
  Create custom secret backends for Vault.
---

# Custom Secret Backends

Vault doesn't currently support the creation of custom secret backends.
The primary reason is because we want to ensure the core of Vault is
secure before attempting any sort of plug-in system. We're interested
in supporting custom secret backends, but don't yet have a clear strategy
or timeline to do.

In the mean time, you can use the
[generic backend](/docs/secrets/generic/index.html) to support custom
data with custom leases.
