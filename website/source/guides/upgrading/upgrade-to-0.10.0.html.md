---
layout: "guides"
page_title: "Upgrading to Vault 0.10.0 - Guides"
sidebar_current: "guides-upgrading-to-0.10.0"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.10.0. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.10.0 compared to 0.9.6. Please read it carefully.

### Database Plugin Compatibility

The database plugin interface was enhanced to support some additional
functionality related to root credential rotation and supporting templated 
URL strings. The changes were made in a backwards-compatible way and all 
builtin plugins were updated with the new features. Custom plugins not built
into Vault will need to be upgraded to support templated URL strings and 
root rotation. Additionally, the Initialize method was deprecated in favor 
of a new Init method that supports configuration modifications that occur in
the plugin back to the primary data store.
