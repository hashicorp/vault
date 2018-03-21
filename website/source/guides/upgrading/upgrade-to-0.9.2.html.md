---
layout: "guides"
page_title: "Upgrading to Vault 0.9.2 - Guides"
sidebar_current: "guides-upgrading-to-0.9.2"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.9.2. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.9.2 compared to 0.9.1. Please read it carefully.

### Backwards Compatible CLI Changes

This upgrade guide is typically reserved for breaking changes, however it
is worth calling out that the CLI interface to Vault has been completely
revamped while maintaining backwards compatibility. This could lead to
potential confusion  while browsing the latest version of the Vault
documentation on vaultproject.io.

All previous CLI commands should continue to work and are backwards
compatible in almost all cases.

Documentation for previous versions of Vault can be accessed using
the GitHub interface by browsing tags (eg [0.9.1 website tree](https://github.com/hashicorp/vault/tree/v0.9.1/website)) or by
[building the Vault website locally](https://github.com/hashicorp/vault/tree/v0.9.1/website#running-the-site-locally).

### `sys/health` DR Secondary Reporting

The `replication_dr_secondary` bool returned by `sys/health` could be
misleading since it would be `false` both when a cluster was not a DR secondary
but also when the node is a standby in the cluster and has not yet fully
received state from the active node. This could cause health checks on LBs to
decide that the node was acceptable for traffic even though DR secondaries
cannot handle normal Vault traffic. (In other words, the bool could only convey
"yes" or "no" but not "not sure yet".) This has been replaced by
`replication_dr_mode` and `replication_perf_mode` which are string values that
convey the current state of the node; a value of `disabled` indicates that
replication is disabled or the state is still being discovered. As a result, an
LB check can positively verify that the node is both not `disabled` and is not
a DR secondary, and avoid sending traffic to it if either is true.


### PKI Secret Backend Roles Parameter Types

For `ou` and `organization` in role definitions in the PKI secret backend,
input can now be a comma-separated string or an array of strings. Reading a
role will now return arrays for these parameters.


### Plugin API Changes

The plugin API has been updated to utilize golang's context.Context package.
Many function signatures now accept a context object as the first parameter.
Existing plugins will need to pull in the latest Vault code and update their
function signatures to begin using context and the new gRPC transport.
