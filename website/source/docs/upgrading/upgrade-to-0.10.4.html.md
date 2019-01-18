---
layout: "docs"
page_title: "Upgrading to Vault 0.10.4 - Guides"
sidebar_title: "Upgrade to 0.10.4"
sidebar_current: "docs-upgrading-to-0.10.4"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.10.4. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.10.4 compared to 0.10.3. Please read it carefully.

### Revocations of dynamic secrets leases now asynchronous

Dynamic secret lease revocation are now queued/asynchronous rather
than synchronous. This allows Vault to take responsibility for revocation
even if the initial attempt fails. The previous synchronous behavior can be
attained via the `-sync` CLI flag or `sync` API parameter. When in
synchronous mode, if the operation results in failure it is up to the user
to retry.

### CLI Retries

The CLI will no longer retry commands on 5xx errors. This was a
source of confusion to users as to why Vault would "hang" before returning a
5xx error. The Go API client still defaults to two retries.

### Identity Entity Alias metadata

You can no longer manually set metadata on
entity aliases. All alias data (except the canonical entity ID it refers to)
is intended to be managed by the plugin providing the alias information, so
allowing it to be set manually didn't make sense.
