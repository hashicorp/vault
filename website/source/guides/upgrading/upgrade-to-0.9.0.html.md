---
layout: "guides"
page_title: "Upgrading to Vault 0.9.0 - Guides"
sidebar_current: "guides-upgrading-to-0.9.0"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.9.0. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.9.0 compared to the most recent release. Please read it carefully.

## CouchDB Storage Changes

Vault may write values to storage that start with an underscore (`_`)
character. This is a reserved character in CouchDB, which can cause breakage.
As a result, this backend now stores each value prefixed with a `$` character.

If you are upgrading from existing CouchDB usage, you can turn off this
behavior by setting the `"prefixed"` configuration value to `"false"`.
Alternately, if you need to handle underscores at the start of keys, you can
rewrite your existing keys to start with a `$` character.
