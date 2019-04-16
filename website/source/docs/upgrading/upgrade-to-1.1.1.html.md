---
layout: "docs"
page_title: "Upgrading to Vault 1.1.1 - Guides"
sidebar_title: "Upgrade to 1.1.1"
sidebar_current: "docs-upgrading-to-1.1.1"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 1.1.1. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 1.1.0 compared to 1.1.1. Please read it carefully.

## Known Issues

There is a known issue that could cause the upgrade to 1.1.1 to fail under
certain circumstances. This issue occurs when a KV version 2 mount exists but
contains no data. This will be fixed in 1.1.2. Addtionally a work around does
exist: prior to upgrading ensure all KV v2 mounts have at least one key written
to it. 

## JWT/OIDC Plugin

Logins of role_type "oidc" via the /login path are no longer allowed. 

## ACL Wildcards

New ordering defines which policy wins when there are multiple inexact matches
and at least one path contains `+`. `+*` is now illegal in policy paths. The
previous behavior simply selected any matching segment-wildcard path that
matched.

## Replication

Due to technical limitations, mounting and unmounting was not previously
possible from a performance secondary. These have been resolved, and these
operations may now be run from a performance secondary.
