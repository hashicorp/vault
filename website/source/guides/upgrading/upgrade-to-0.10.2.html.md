---
layout: "guides"
page_title: "Upgrading to Vault 0.10.2 - Guides"
sidebar_current: "guides-upgrading-to-0.10.2"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.10.2. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.10.2 compared to 0.10.1. Please read it carefully.

### Convergent Encryption version 3

If you are using `transit`'s convergent encryption feature, which prior to this
release was at version 2, we recommend
[rotating](https://www.vaultproject.io/api/secret/transit/index.html#rotate-key)
your encryption key (the new key will use version 3) and
[rewrapping](https://www.vaultproject.io/api/secret/transit/index.html#rewrap-data)
your data to mitigate the chance of offline plaintext-confirmation attacks.

### PKI duration return types

The PKI backend now returns durations (e.g. when reading a role) as an integer
number of seconds instead of a Go-style string.
