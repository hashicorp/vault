---
layout: "docs"
page_title: "Upgrading to Vault 0.11.6 - Guides"
sidebar_title: "Upgrade to 0.11.6"
sidebar_current: "docs-upgrading-to-0.11.6"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.11.6. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.11.6 compared to 0.11.5. Please read it carefully.

### Database Secret Engine Role Reads

On role read, empty statements will be returned as empty
slices instead of potentially being returned as JSON null values. This makes it
more in line with other parts of Vault and makes it easier for statically typed
languages to interpret the values.
