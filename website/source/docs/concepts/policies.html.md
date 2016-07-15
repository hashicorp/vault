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
