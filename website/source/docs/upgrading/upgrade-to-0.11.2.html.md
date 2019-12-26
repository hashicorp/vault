---
layout: "docs"
page_title: "Upgrading to Vault 0.11.2 - Guides"
sidebar_title: "Upgrade to 0.11.2"
sidebar_current: "docs-upgrading-to-0.11.2"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.11.2. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.11.2 compared to 0.11.1. Please read it carefully.

### `sys/seal-status` Behavior Change

The `sys/seal-status` endpoint now includes an initialized boolean in the
output. If Vault is not initialized, it will return a 200 with this value
set false instead of a 400

### Mount Config Passthrough Headers

The mount config option for `passthrough_request_headers` will now deny
certain headers from being provided to backends based on a global denylist.
