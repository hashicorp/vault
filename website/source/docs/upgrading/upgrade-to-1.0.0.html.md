---
layout: "docs"
page_title: "Upgrading to Vault 1.0.0 - Guides"
sidebar_title: "Upgrade to 1.0.0"
sidebar_current: "docs-upgrading-to-1.0.0"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 1.0.0. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.11.5 compared to 1.0.0. Please read it carefully.

## Token Format Change

Tokens are now prefixed by a designation to indicate what type of token they are. Service tokens start with s. and batch tokens start with b.. Existing tokens will still work (they are all of service type and will be considered as such). Prefixing allows us to be more efficient when consuming a token, which keeps the critical path of requests faster.

## Key Encoding Enforcement

Vault will no longer accept updates when the storage key has invalid UTF-8 character encoding.

## Upsert in Mount Tuning Endpoints

Mount/Auth tuning the options map on backends will now upsert any provided values, and keep any of the existing values in place if not provided. The options map itself cannot be unset once it's set, but the keypairs within the map can be unset if an empty value is provided, with the exception of the version keypair which is handled differently for KVv2 purposes.

## Agent Reauthentication

Agent no longer automatically reauthenticates when new credentials are detected. It's not strictly necessary and in some cases was causing reauthentication much more often than intended.

## HSM Key Regeneration Removed

Vault no longer supports destroying and regenerating encryption keys on an HSM; it only supports creating them. Although this has never been a source of a customer incident, it is simply a code path that is too trivial to activate, especially by mistyping regenerate_key instead of generate_key.

## Seal Upgrade before Migration

When upgrading from Vault 0.8.x, the seal type in the barrier config storage entry will be upgraded from "hsm-auto" to "awskms" or "pkcs11" upon unseal if using AWSKMS or HSM seals. If performing seal migration, the barrier config should first be upgraded prior to starting migration.

## Pooled API Client

The Go API client now uses a connection-pooling HTTP client by default. For CLI operations this makes no difference but it should provide significant performance benefits for those writing custom clients using the Go API library. As before, this can be changed to any custom HTTP client by the caller.

## Plugin Catalog Changes

Builtin Secret Engines and Auth Methods are integrated deeper into the plugin system. The plugin catalog can now override builtin plugins with custom versions of the same name. Additionally the plugin system now requires a plugin type field when configuring plugins, this can be "auth", "database", or "secret".

## Removed Endpoints

Paths within auth/token that allow specifying a token or accessor in the URL have been removed. These have been deprecated since March 2016 and undocumented, but were retained for backwards compatibility. They shouldn't be used due to the possibility of those paths being logged, so at this point they are simply being removed.
