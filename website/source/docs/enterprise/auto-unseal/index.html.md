---
layout: "docs"
page_title: "Vault Enterprise Auto Unseal"
sidebar_current: "docs-vault-enterprise-auto-unseal"
description: |-
  Vault Enterprise supports automatic unsealing via cloud technologies like KMS.
---

# Vault Enterprise Auto Unseal

As of version 0.9, Vault Enterprise supports opt-in automatic unsealing via
cloud technologies such Amazon KMS or Google Cloud KMS. This feature enables
operators to delegate the unsealing process to trusted cloud providers to ease
operations in the event of partial failure and to aid in the creation of new or
ephemeral clusters.

## Enabling Auto Unseal

Automatic unsealing is not enabled by default. To enable automatic unsealing,
specify the `seal` stanza in your Vault configuration file:

```hcl
seal "awskms" {
  aws_region = "us-east-1"
  access_key = "..."
  secret_key = "..."
  kms_key_id = "..."
}
```

For a complete list of examples and supported technologies, please see the
[seal documentation](/docs/configuration/seal/index.html).
