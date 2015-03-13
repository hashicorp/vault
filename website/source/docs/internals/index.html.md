---
layout: "docs"
page_title: "Internals"
sidebar_current: "docs-internals"
description: |-
  This section covers the internals of Vault and explains how plans are generated, the lifecycle of a provider, etc. The goal of this section is to remove any notion of "magic" from Vault. We want you to be able to trust and understand what Vault is doing to function.
---

# Vault Internals

This section covers the internals of Vault and explains how
plans are generated, the lifecycle of a provider, etc. The goal
of this section is to remove any notion of "magic" from Vault.
We want you to be able to trust and understand what Vault is
doing to function.

-> **Note:** Knowledge of Vault internals is not
required to use Vault. If you aren't interested in the internals
of Vault, you may safely skip this section.
