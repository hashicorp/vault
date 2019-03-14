---
layout: "docs"
page_title: "Upgrading to Vault 1.1.0 - Guides"
sidebar_title: "Upgrade to 1.1.0"
sidebar_current: "docs-upgrading-to-1.1.0"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 1.1.0. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 1.0.3 compared to 1.1.0. Please read it carefully.

## JWT Backend Changes

Specifying the group claims parameter has changed to use a standards based lookup.  The groups_claim_delimiter_pattern 
has been removed and if the groups claim is not at the top level, it can now be specified as a JSONPointer.

Additionally, roles now have a "role type" parameter with a default type of "oidc". To configure new JWT roles, a role 
type of "jwt" must be explicitly specified.

## Deprecated CLI Commands Removed

CLI commands deprecated in 0.9.2 are now removed. Please see the CLI help output for updated commands.

## Additional Changes

* Vault no longer automatically mounts a k/v backend at the "secret/" path when initalizing Vault.
* Vault's cluster port will now be opened on HA standby nodes.
* Vault no longer supports running netRPC plugins. These were deprecated in favor of gRPC based plugins and any plugin built since 0.9.4 defaults to gRPC. Older plugins may need to be recompiled against the latest Vault dependencies. 