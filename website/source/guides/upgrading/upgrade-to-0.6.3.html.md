---
layout: "guides"
page_title: "Upgrading to Vault 0.6.3 - Guides"
sidebar_current: "guides-upgrading-to-0.6.3"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.6.3. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.6.3. Please read it carefully.

## LDAP Null Binds Disabled By Default

When using the LDAP Auth Backend, `deny_null_bind` has a default value of
`true`, preventing a successful user authentication when an empty password
is provided. If you utilize passwordless LDAP binds, `deny_null_bind` must
be set to `false`. Upgrades will keep previous behavior until the LDAP
configuration information is rewritten, at which point the new behavior
will be utilized.

## Request Size Limitation

A maximum request size of 32MB is imposed to prevent a denial of service attack
with arbitrarily large requests.

## Any Audit Backend Successfully Activated Allows Active Duty

Previously, when a new Vault node was taking over service in an HA cluster, all
audit backends were required to be active successfully to take over active
duty. This behavior now matches the behavior of the audit logging system
itself: at least one audit backend must successfully be activated. The server
log contains an error when this occurs. This helps keep a Vault HA cluster
working when there is a misconfiguration on a standby node.
