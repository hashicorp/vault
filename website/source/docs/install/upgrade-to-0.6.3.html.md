---
layout: "install"
page_title: "Upgrading to Vault 0.6.3"
sidebar_current: "docs-install-upgrade-to-0.6.3"
description: |-
  Learn how to upgrade to Vault 0.63.
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
