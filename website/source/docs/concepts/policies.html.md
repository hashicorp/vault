---
layout: "docs"
page_title: "Policies"
sidebar_current: "docs-concepts-policies"
description: |-
  Policies are how authorization is done in Vault, allowing you to restrict which parts of Vault a user can access.
---

# Access Control Policies

After [authenticating](/docs/concepts/auth.html) with Vault, the
next step is authorization. This is the process of determining what
a user is allowed to do. Authorization is unified in Vault in the form
of _policies_.

Policies are [HCL](https://github.com/hashicorp/hcl) or JSON documents
that describe what parts of Vault a user is allowed to access. An example
of a policy is shown below:

```javascript
path "sys/*" {
  policy = "deny"
}

path "secret/*" {
  policy = "write"
}

path "secret/foo" {
  policy = "read"
  capabilities = ["create", "sudo"]
}

path "secret/super-secret" {
  capabilities = ["deny"]
}

path "secret/bar" {
  capabilities = ["create"]
  allowed_parameters = {
    "*" = []
  }
  denied_parameters = {
    "foo" = ["bar"]
  }
}
```

Policies use path based matching to apply rules. A policy may be an exact
match, or might be a glob pattern which uses a prefix. Vault operates in a
whitelisting mode, so if a path isn't explicitly allowed, Vault will reject
access to it.  This works well due to Vault's architecture of being like a
filesystem: everything has a path associated with it, including the core
configuration mechanism under "sys".

~> Policy paths are matched using the most specific defined policy. This may
be an exact match or the longest-prefix match of a glob. This means if you
define a policy for `"secret/foo*"`, the policy would also match `"secret/foobar"`.
The glob character is only supported at the end of the path specification.

## Capabilities and Policies

Paths have an associated set of capabilities that provide fine-grained control
over operations. The capabilities are:

  * `create` - Create a value at a path. (At present, few parts of Vault
    distinguish between `create` and `update`, so most operations require
    `update`. Parts of Vault that provide such a distinction, such as
    the `generic` backend, are noted in documentation.)

  * `read` - Read the value at a path.

  * `update` - Change the value at a path. In most parts of Vault, this also
    includes the ability to create the initial value at the path.

  * `delete` - Delete the value at a path.

  * `list` - List key names at a path. Note that the keys returned by a
    `list` operation are *not* filtered by policies.  Do not encode sensitive
    information in key names.

  * `sudo` - Gain access to paths that are _root-protected_. This is _additive_
    to other capabilities, so a path that requires `sudo` access will also
    require `read`, `update`, etc. as appropriate.

  * `deny` - No access allowed. This always takes precedence regardless of any
    other defined capabilities, including `sudo`.

The only non-obvious capability is `sudo`. Some routes within Vault and mounted
backends are marked as _root-protected_ paths. Clients aren't allowed to access
root paths unless they are a root user (have the special policy "root" attached
to their token) or have access to that path with the `sudo` capability (in
addition to the other necessary capabilities for performing an operation
against that path, such as `read` or `delete`).

For example, modifying the audit log backends is done via root paths.
Only root or `sudo` privilege users are allowed to do this.

Prior to Vault 0.5, the `policy` keyword was used per path rather than a set of
`capabilities`. In Vault 0.5+ these are still supported as shorthand and to
maintain backwards compatibility, but internally they map to a set of
capabilities. These mappings are as follows:

  * `deny` - `["deny"]`

  * `sudo` - `["create", "read", "update", "delete", "list", "sudo"]`

  * `write` - `["create", "read", "update", "delete", "list"]`

  * `read` - `["read", "list"]`

## Fine-Grained Control

There are a few optional fields that allow for fine-grained control over client
behavior on a given path. The capabilities associated with this path take
precedence over permissions on parameters.

### Allowed and Disallowed Parameters

These parameters allow the administrator to restrict the keys (and optionally
values) that a user is allowed to specify when calling a path.

  * `allowed_parameters` - A map of keys to an array of values that acts as a
    whitelist. Setting a key with an `[]` value will allow changes to
    parameters with that name. Setting a key with a populated value array (e.g.
    `["foo", "bar"]`, `[3600, 7200]` or `[true]` will allow that parameter to
    only be set to one of the values in the array. If any keys exist in the
    `allowed_parameters` object all keys not specified will be denied unless
    there the key `"*"` is set (mapping to an empty array), which will allow
    all other parameters to be modified; parameters with specific values will
    still be restricted to those values.
  * `denied_parameters` - A map of keys to an array of values that acts as a
    blacklist, and any parameter set here takes precedence over
    `allowed_parameters`. Setting to "*" will deny any parameter (so only calls
    made without specifying any parameters will be allowed). Otherwise setting
    a key with an `[]` value will deny any changes to parameters with that
    name. Setting a key with a populated value array will deny any attempt to
    set a parameter with that name and value. If keys exist in the
    `denied_parameters` object all keys not specified will be allowed (unless
    `allowed_parameters` is also set, in which case normal rules will apply).

String values inside a populated value array support prefix/suffix globbing. 
Globbing is enabled by prepending or appending a `*` to the value (e.g. 
`["*foo", "bar*"]` would match `"...foo"` and `"bar..."`).

### Required Minimum/Maximum Response Wrapping TTLs

These parameters can be used to set minimums/maximums on TTLs set by clients
when requesting that a response be
[wrapped](/docs/concepts/response-wrapping.html), with a granularity of a second. These can either be specified as a number of seconds or a string with a `s`, `m`, or `h` suffix indicating seconds, minutes, and hours respectively.

In practice, setting a minimum TTL of one second effectively makes response
wrapping mandatory for a particular path.

  * `min_wrapping_ttl` - The minimum allowed TTL that clients can specify for a
    wrapped response. In practice, setting a minimum TTL of one second
    effectively makes response wrapping mandatory for a particular path. It can
    also be used to ensure that the TTL is not too low, leading to end targets
    being unable to unwrap before the token expires.
  * `max_wrapping_ttl` - The maximum allowed TTL that clients can specify for a
    wrapped response.

If both are specified, the minimum value must be less than the maximum. In
addition, if paths are merged from different stanzas, the lowest value
specified for each is the value that will result, in line with the idea of
keeping token lifetimes as short as possible.

## Root Policy

The "root" policy is a special policy that can not be modified or removed.
Any user associated with the "root" policy becomes a root user. A root
user can do _anything_ within Vault.

There always exists at least one root user (associated with the token
when initializing a new server). After this root user, it is recommended
to create more strictly controlled users. The original root token should
be protected accordingly.

## Managing Policies

Policy management can be done via the API or CLI. The CLI commands are
`vault policies` and `vault policy-write`. Please see the help associated
with these commands for more information. They are very easy to use.

## Associating Policies

To associate a policy with a user, you must consult the documentation for
the authentication backend you're using.

For tokens, they are associated at creation time with `vault token-create`
and the `-policy` flags. Child tokens can be associated with a subset of
a parent's policies. Root users can assign any policies.

There is no way to modify the policies associated with a token once the token
has been issued. The token must be revoked and a new one acquired to receive a
new set of policies.

However, the _contents_ of policies are parsed in real-time at every token use.
As a result, if a policy is modified, the modified rules will be in force the
next time a token with that policy attached is used to make a call to Vault.
