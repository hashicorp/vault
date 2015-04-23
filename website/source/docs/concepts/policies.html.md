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
path "sys" {
  policy = "deny"
}

path "secret" {
  policy = "write"
}

path "secret/foo" {
  policy = "read"
}
```

Policies use prefix-based routing to apply rules. They are deny by default,
so if a path isn't explicitly given, Vault will reject any access to it.
This works well due to Vault's architecture of being like a filesystem:
everything has a path associated with it, including the core configuration
mechanism under "sys".

## Policies

Allowed policies for a path are:

  * `write` - Read, write access to a path.

  * `read` - Read-only access to a path.

  * `deny` - No access allowed.

  * `sudo` - Read, write, and root access to a path.

The only non-obvious policy is "sudo". Some routes within Vault and mounted
backends are marked as _root_ paths. Clients aren't allowed to access root
paths unless they are a root user (have the special policy "root") or
have access to that path with the "sudo" policy.

For example, modifying the audit log backends is done via root paths.
Only root or "sudo" privilege users are allowed to do this.

## Root Policy

The "root" policy is special policy that can'be modified or removed.
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

For tokens, they are assocated at creation time with `vault token-create`
and the `-policy` flags. Child tokens can be associated with a subset of
a parent's policies. Root users can assign any policies.

There is no way to modify the policies associated with an active
identity. The identity must be revoked and reauthenticated to receive
the new policy list.

If an _existing_ policy is modified, the modifications propogate
to all associated users instantly. The above paragraph is more specifically
stating that you can't add new or remove policies associated with an
active identity.
