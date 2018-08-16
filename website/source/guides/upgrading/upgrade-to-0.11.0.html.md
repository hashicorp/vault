---
layout: "guides"
page_title: "Upgrading to Vault 0.11.0 - Guides"
sidebar_current: "guides-upgrading-to-0.11.0"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.11.0. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.11.0 Beta compared to 0.10.0. Please read it carefully.

**NOTE** This beta release does not have a forward compatibility guarantee and 
certain functionality may change that will be incompatible with the General 
Availability release. Please only use the beta releases to test functionality 
and upgrades with clusters that can be lost.

## Changes Since 0.10.4

### Request Timeouts

A default request timeout of 90s is now enforced. This setting can be
overwritten in the config file. If you anticipate requests taking longer than
90s this setting should be configured before upgrading.

### `sys/` Top Level Injection

For the last two years for backwards compatibility data for various `sys/`
routes has been injected into both the Secret's Data map and into the top level
of the JSON response object. However, this has some subtle issues that pop up
from time to time and is becoming increasingly complicated to maintain, so it's
finally being removed.

## Full List Since 0.10.0

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
