---
layout: "guides"
page_title: "Upgrading to Vault 0.9.6 - Guides"
sidebar_current: "guides-upgrading-to-0.9.6"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.9.6. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.9.6 compared to 0.9.5. Please read it carefully.

### Change to AWS Role Output

The AWS authentication backend now allows binds for inputs as either a
comma-delimited string or a string array. However, to keep consistency with
input and output, when reading a role the binds will now be returned as string
arrays rather than strings.

### Change to AWS IAM Auth ARN Prefix Matching

In order to prefix-match IAM role and instance profile ARNs in AWS auth
backend, you now must explicitly opt-in by adding a `*` to the end of the ARN.
Existing configurations will be upgraded automatically, but when writing a new
role configuration the updated behavior will be used.
