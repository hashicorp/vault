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
for Vault 0.10.0 compared to 0.9.0. Please read it carefully.

## Changes Since 0.9.6

### Database Plugin Compatibility

The database plugin interface was enhanced to support some additional
functionality related to root credential rotation and supporting templated 
URL strings. The changes were made in a backwards-compatible way and all 
builtin plugins were updated with the new features. Custom plugins not built
into Vault will need to be upgraded to support templated URL strings and 
root rotation. Additionally, the Initialize method was deprecated in favor 
of a new Init method that supports configuration modifications that occur in
the plugin back to the primary data store.

### Removal of Returned Secret Information

For a long time Vault has returned configuration given to various secret
engines and auth methods with secret values (such as secret API keys or
passwords) still intact, and with a warning to the user on write that anyone
with read access could see the secret. This was mostly done to make it easy for
tools like Terraform to judge whether state had drifted. However, it also feels
quite un-Vault-y to do this and we've never felt very comfortable doing so. In
0.10 we have gone through and removed this bevhavior from the various backends;
fields which contained secret values are simply no longer returned on read. We
are working with the Terraform team to make changes to their provider to
accommodate this as best as possible, and users of other tools may have to make
adjustments, but in the end we felt that the ends did not justify the means and
we needed to prioritize security over operational convenience.

### LDAP Auth Method Case Sensitivity

We now treat usernames and groups configured locally for policy assignment in a
case insensitive fashion by default. Existing configurations will continue to
work as they do now; however, the next time a configuration is written
`case_sensitive_names` will need to be explicitly set to `true`.

### TTL Handling Moved to Core

All lease TTL handling has been centralized within the core of Vault to ensure
consistency across all backends. Since this was previously delegated to
individual backends, there may be some slight differences in TTLs generated
from some backends.

### Default `secret/` Mount is Deprecated

In 0.12 we will stop mounting `secret/` by default at initialization time (it
will still be available in `dev` mode).

## Full List Since 0.9.0

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
### Backwards Compatible CLI Changes

This upgrade guide is typically reserved for breaking changes, however it
is worth calling out that the CLI interface to Vault has been completely
revamped while maintaining backwards compatibility. This could lead to
potential confusion  while browsing the latest version of the Vault
documentation on vaultproject.io.

All previous CLI commands should continue to work and are backwards
compatible in almost all cases.

Documentation for previous versions of Vault can be accessed using
the GitHub interface by browsing tags (eg [0.9.1 website tree](https://github.com/hashicorp/vault/tree/v0.9.1/website)) or by
[building the Vault website locally](https://github.com/hashicorp/vault/tree/v0.9.1/website#running-the-site-locally).

### `sys/health` DR Secondary Reporting

The `replication_dr_secondary` bool returned by `sys/health` could be
misleading since it would be `false` both when a cluster was not a DR secondary
but also when the node is a standby in the cluster and has not yet fully
received state from the active node. This could cause health checks on LBs to
decide that the node was acceptable for traffic even though DR secondaries
cannot handle normal Vault traffic. (In other words, the bool could only convey
"yes" or "no" but not "not sure yet".) This has been replaced by
`replication_dr_mode` and `replication_perf_mode` which are string values that
convey the current state of the node; a value of `disabled` indicates that
replication is disabled or the state is still being discovered. As a result, an
LB check can positively verify that the node is both not `disabled` and is not
a DR secondary, and avoid sending traffic to it if either is true.


### PKI Secret Backend Roles Parameter Types

For `ou` and `organization` in role definitions in the PKI secret backend,
input can now be a comma-separated string or an array of strings. Reading a
role will now return arrays for these parameters.

### Plugin API Changes

The plugin API has been updated to utilize golang's context.Context package.
Many function signatures now accept a context object as the first parameter.
Existing plugins will need to pull in the latest Vault code and update their
function signatures to begin using context and the new gRPC transport.

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
