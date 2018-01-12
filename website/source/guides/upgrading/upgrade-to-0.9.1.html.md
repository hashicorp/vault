---
layout: "guides"
page_title: "Upgrading to Vault 0.9.1 - Guides"
sidebar_current: "guides-upgrading-to-0.9.1"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.9.1. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.9.1 compared to the 0.9.0. Please read it carefully.

### AppRole Case Sensitivity

In prior versions of Vault, `list` operations against AppRole roles would
require preserving case in the role name, even though most other operations
within AppRole are case-insensitive with respect to the role name. This has
been fixed; existing roles will behave as they have in the past, but new roles
will act case-insensitively in these cases.

### Token Auth Backend Roles Parameter Types

For `allowed_policies` and `disallowed_policies` in role definitions in the
token auth backend, input can now be a comma-separated string or an array of
strings. Reading a role will now return arrays for these parameters.

### Transit Key Exporting

You can now mark a key in the `transit` backend as `exportable` at any time,
rather than just at creation time; however, once this value is set, it still
cannot be unset.

### PKI Secret Backend Roles Parameter Types

For `allowed_domains` and `key_usage` in role definitions in the PKI secret
backend, input can now be a comma-separated string or an array of strings.
Reading a role will now return arrays for these parameters.

### SSH Dynamic Keys Method Defaults to 2048-bit Keys

When using the dynamic key method in the SSH backend, the default is now to use
2048-bit keys if no specific key bit size is specified.

### Consul Secret Backend Lease Handling

The `consul` secret backend can now accept both strings and integer numbers of
seconds for its lease value. The value returned on a role read will be an
integer number of seconds instead of a human-friendly string.

### Unprintable Characters Not Allowed in API Paths

Unprintable characters are no longer allowed in names in the API (paths and
path parameters), with an extra restriction on whitespace characters. Allowed
characters are those that are considered printable by Unicode plus spaces.
