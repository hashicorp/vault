---
layout: "intro"
page_title: "Policies - Getting Started"
sidebar_current: "gettingstarted-policies"
description: |-
  Policies in Vault control what a user can access.
---

# Policies

Policies in Vault control what a user can access. In the last section, we
learned about _authentication_. This section is about _authorization_.

For authentication Vault has multiple options or methods that can be enabled and
used. For authorization and policies Vault always uses the same format. All auth
methods map identities back to the core policies that are configured with Vault.

There are some built-in policies that cannot be removed. For example, the `root`
and `default` policies are required policies and cannot be deleted. The
`default` policy provides a common set of permissions and is included on all
tokens by default. The `root` policy gives a token super admin permissions,
similar to a root user on a linux machine.

## Policy Format

Policies are authored in [HCL][hcl], but it is JSON compatible. Here is an
example policy:

```hcl
# Normal servers have version 1 of KV mounted by default, so will need these
# paths:
path "secret/*" {
  capabilities = ["create"]
}
path "secret/foo" {
  capabilities = ["read"]
}

# Dev servers have version 2 of KV mounted by default, so will need these
# paths:
path "secret/data/*" {
  capabilities = ["create"]
}
path "secret/data/foo" {
  capabilities = ["read"]
}
```

With this policy, a user could write any secret to `secret/`, except to
`secret/foo`, where only read access is allowed. Policies default to deny, so
any access to an unspecified path is not allowed.

Do not worry about getting the exact policy format correct. Vault includes a
command that will format the policy automatically according to specification. It
also reports on any syntax errors.

```text
$ vault policy fmt my-policy.hcl
```

The policy format uses a prefix matching system on the API path to determine
access control. The most specific defined policy is used, either an exact match
or the longest-prefix glob match. Since everything in Vault must be accessed via
the API, this gives strict control over every aspect of Vault, including
enabling secrets engines, enabling auth methods, authenticating, as well as
secret access.

## Writing the Policy

To write a policy using the command line, specify the path to a policy file to
upload.

```text
$ vault policy write my-policy my-policy.hcl
Success! Uploaded policy: my-policy
```

Here is an example you can copy-paste in the terminal:

```text
$ vault policy write my-policy -<<EOF
# Normal servers have version 1 of KV mounted by default, so will need these
# paths:
path "secret/*" {
  capabilities = ["create"]
}
path "secret/foo" {
  capabilities = ["read"]
}

# Dev servers have version 2 of KV mounted by default, so will need these
# paths:
path "secret/data/*" {
  capabilities = ["create"]
}
path "secret/data/foo" {
  capabilities = ["read"]
}
EOF
```

To see the list of policies, run:

```text
$ vault policy list
default
my-policy
root
```

To view the contents of a policy, run:

```text
$ vault policy read my-policy
# Normal servers have version 1 of KV mounted by default, so will need these
# paths:
path "secret/*" {
  capabilities = ["create"]
}
...
```

## Testing the Policy

To use the policy, create a token and assign it to that policy:

```text
$ vault token create -policy=my-policy
Key                Value
---                -----
token              a4ebda12-23bf-5cf4-f80e-803ee2f37aab
token_accessor     aba6256e-401e-9591-31b2-a27048cb15ed
token_duration     768h
token_renewable    true
token_policies     [default my-policy]

$ vault login a4ebda12-23bf-5cf4-f80e-803ee2f37aab
Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                Value
---                -----
token              a4ebda12-23bf-5cf4-f80e-803ee2f37aab
token_accessor     aba6256e-401e-9591-31b2-a27048cb15ed
token_duration     767h59m18s
token_renewable    true
token_policies     [default my-policy]
```

Verify that you can write any data to `secret/`, but only read from
`secret/foo`:

### Dev servers

```text
$ vault kv put secret/bar robot=beepboop
Key              Value
---              -----
created_time     2018-05-22T18:05:42.537496856Z
deletion_time    n/a
destroyed        false
version          1

$ vault kv put secret/foo robot=beepboop
Error writing data to secret/data/foo: Error making API request.

URL: PUT http://127.0.0.1:8200/v1/secret/data/foo
Code: 403. Errors:

* permission denied
```

### Non-dev servers

```text
$ vault kv put secret/bar robot=beepboop
Success! Data written to: secret/bar

$ vault kv put secret/foo robot=beepboop
Error writing data to secret/foo: Error making API request.

URL: PUT http://127.0.0.1:8200/v1/secret/foo
Code: 403. Errors:

* permission denied
```

You also do not have access to `sys` according to the policy, so commands like
`vault policy list` or `vault secrets list` will not work. Re-authenticate as
the initial root token to continue:

```text
$ vault login <initial-root-token>
```

## Mapping Policies to Auth Methods

Vault is the single policy authority, unlike auth where you can enable multiple
auth methods. Any enabled auth method must map identities to these core
policies.

We use the `vault path-help` system with your auth method to determine how the
mapping is done, since it is specific to each auth method. For example, with
GitHub, it is done by team using the `map/teams/<team>` path:

```text
$ vault write auth/github/map/teams/default value=my-policy
Success! Data written to: auth/github/map/teams/default
```

For GitHub, the `default` team is the default policy set that everyone is
assigned to no matter what team they're on.

Other auth methods use alternate, but likely similar mechanisms for mapping
policies to identity.

## Next

Policies are an important part of Vault. While using the root token is easiest
to get up and running, you will want to restrict access to Vault very quickly,
and the policy system is the way to do this.

The syntax and function of policies is easy to understand and work with, and
because auth methods all must map to the central policy system, you only have to
learn this policy system.

Next, we will cover how to [deploy Vault](/intro/getting-started/deploy.html).

[HCL]: https://github.com/hashicorp/hcl
