---
layout: "intro"
page_title: "Access Control Policies"
sidebar_current: "gettingstarted-acl"
description: |-
  Access control policies in Vault control what a user can access.
---

# Access Control Policies (ACLs)

Access control policies in Vault control what a user can access. In
the last section, we learned about _authentication_. This section is
about _authorization_.

Whereas for authentication Vault has multiple options or backends that
can be enabled and used, the authorization or policies of Vault are always
the same format. All authentication backends must map identities back to
the core policies that are configured with Vault.

When initializing Vault, there is always one special policy created
that can't be removed: the "root" policy. This policy is a special policy
that gives superuser access to everything in Vault. An identity mapped to
the root policy can do anything.

## Policy Format

Policies in Vault are formatted with
[HCL](https://github.com/hashicorp/hcl). HCL is a human-readable configuration
format that is also JSON-compatible, so you can use JSON as well. An example
policy is shown below:

```
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

The policy format uses a longest matching prefix system on the API path
to determine access control. Since everything in Vault must be accessed
via the API, this gives strict control over every aspect of Vault, including
mounting backends, authenticating, as well as secret access.

In the policy above, a user could write any secret to `secret/`, except
to `secret/foo`, where only read access is allowed. Policies default to
deny, so any access to an unspecified path is not allowed.

Save the above policy as `acl.hcl`.

## Writing the Policy

To write a policy, use the `vault policy-write` command:

```
$ vault policy-write secret acl.hcl
Policy 'secret' written.
```

You can see the policies that are available with `vault policies`, and you
can see the contents of a policy with `vault policy <name>`. Only users with
root access can do this.

## Testing the Policy

To use the policy, let's create a token and assign it to that policy.
Make sure to save your root token somewhere so you can authenticate
back to a root user later.

```
$ vault token-create -policy="secret"
d97ef000-48cf-45d9-1907-3ea6ce298a29

$ vault auth
```
