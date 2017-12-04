---
layout: "intro"
page_title: "Built-in Help - Getting Started"
sidebar_current: "gettingstarted-help"
description: |-
  Vault has a built-in help system to learn about the available paths in Vault and how to use them.
---

# Built-in Help

You've now worked with `vault write` and `vault read` for multiple
paths: the kv secret backend with `secret/` and dynamic AWS
credentials with the AWS backend provider at `aws/`. In both cases, the
structure and usage of each backend differed, for example the AWS
backend has special paths like `aws/config`.

Instead of having to memorize or reference documentation constantly
to determine what paths to use, we've built a help system directly into
Vault. This help system can be accessed via the API or the command-line and
generates human-readable help for any mounted backend.

On this page, we'll learn how to use this help system. It is an invaluable
tool as you continue to work with Vault.

## Backend Overview

For this, we'll assume you have the AWS backend mounted. If not, mount
it with `vault mount aws`. Even if you don't have an AWS account, you
can still mount the AWS backend.

With the backend mounted, let's learn about it with the `vault
path-help` command:

```
$ vault path-help aws
## DESCRIPTION

The AWS backend dynamically generates AWS access keys for a set of
IAM policies. The AWS access keys have a configurable lease set and
are automatically revoked at the end of the lease.

After mounting this backend, credentials to generate IAM keys must
be configured with the "root" path and policies must be written using
the "roles/" endpoints before any access keys can be generated.

## PATHS

The following paths are supported by this backend. To view help for
any of the paths below, use the help command with any route matching
the path pattern. Note that depending on the policy of your auth token,
you may or may not be able to access certain paths.

    ^config/lease$
        Configure the default lease information for generated credentials.

    ^config/root$
        Configure the root credentials that are used to manage IAM.

    ^creds/(?P<name>\w+)$
        Generate an access key pair for a specific role.

    ^roles/(?P<name>\w+)$
        Read and write IAM policies that access keys can be made for.
```

The `vault path-help` command takes a path. By specifying the root path for
a mount, it will give us the overview of that mount. Notice how the help
not only contains a description, but also the exact regular expressions
used to match routes for this backend along with a brief description
of what the route is for.

## Path Help

After seeing the overview, we can continue to dive deeper by getting
help for an individual path. For this, just use `vault path-help` with a path
that would match the regular expression for that path. Note that the path
doesn't need to actually _work_. For example, we'll get the help below
for accessing `aws/creds/operator`, even though we never created the `operator`
role:

```
$ vault path-help aws/creds/operator
Request:        creds/operator
Matching Route: ^creds/(?P<name>\w+)$

Generate an access key pair for a specific role.

## PARAMETERS

    name (string)
        Name of the role

## DESCRIPTION

This path will generate a new, never before used key pair for
accessing AWS. The IAM policy used to back this key pair will be
the "name" parameter. For example, if this backend is mounted at "aws",
then "aws/creds/deploy" would generate access keys for the "deploy" role.

The access keys will have a lease associated with them. The access keys
can be revoked by using the lease ID.
```

Within a path, we are given the parameters that this path requires.
Some parameters come from the route itself. In this case, the `name`
parameter is a named capture from the route regular expression.

There is also a description of what that path does.

Go ahead and explore more paths! Mount other backends, traverse their
help systems and learn about what they do. For example, learn about the
kv `secret/` path.

## Next

The help system may not be the most exciting feature of Vault, but it
is indispensable in day-to-day usage of Vault. The help system lets you
learn about how to use any backend within Vault without leaving the command
line.

Next, we'll learn about
[authentication](/intro/getting-started/authentication.html).
