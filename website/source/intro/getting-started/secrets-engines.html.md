---
layout: "intro"
page_title: "Secrets Engines - Getting Started"
sidebar_title: "Secrets Engines"
sidebar_current: "gettingstarted-secret-backends"
description: |-
  secrets engines are what create, read, update, and delete secrets.
---

# Secrets Engines

Previously, we saw how to read and write arbitrary secrets to Vault. You may
have noticed all requests started with `secret/`. Try using a different prefix -
Vault will return an error:

```text
$ vault write foo/bar a=b
# ...
* no handler for route 'foo/bar'
```

The path prefix tells Vault which secrets engine it should route
traffic to. When a request comes to Vault, it matches the initial path part using a
longest prefix match and then passes the request to the corresponding secrets
engine enabled at that path.

By default, Vault enables a secrets engine called `kv` at the path `secret/`.
The kv secrets engine reads and writes raw data to the backend storage.

Vault supports many other secrets engines besides `kv`, and this feature makes
Vault flexible and unique. For example, the `aws` secrets engine generates AWS
IAM access keys on demand. The `database` secrets engine generates on-demand,
time-limited database credentials. These are just a few examples of the many
available secrets engines.

For simplicity and familiarity, Vault presents these secrets engines similar to
a filesystem. A secrets engine is enabled at a path. Vault itself performs
prefix routing on incoming requests and routes the request to the correct
secrets engine based on the path at which they were enabled.

This page discusses secrets engines and the operations they support. This
information is important to both operators who will configure Vault and users
who will interact with Vault.

## Enable a Secrets Engine

To get started, enable another instance of the `kv` secrets engine at a
different path. Just like a filesystem, Vault can enable a secrets engine at
many different paths. Each path is completely isolated and cannot talk to other
paths. For example, a `kv` secrets engine enabled at `foo` has no ability to
communicate with a `kv` secrets engine enabled at `bar`.

```text
$ vault secrets enable -path=kv kv
Success! Enabled the kv secrets engine at: kv/
```

The path where the secrets engine is enabled defaults to the name of the secrets engine. Thus, the following commands are actually equivalent:

```text
$ vault secrets enable -path=kv kv
$ vault secrets enable kv
```

To verify our success and get more information about the secrets engine, use the
`vault secrets list` command:

```text
$ vault secrets list
Path          Type         Description
----          ----         -----------
cubbyhole/    cubbyhole    per-token private secret storage
kv/           kv           n/a
secret/       kv           key/value secret storage
sys/          system       system endpoints used for control, policy and debugging
```

This shows there are 4 enabled secrets engines on this Vault server. You can see
the type of the secrets engine, the corresponding path, and an optional
description (or "n/a" if none was given).

~> The `sys/` path corresponds to the system backend. While the system backend
is not specifically discussed in this guide, there is plentiful documentation on
the system backend. Many of these operations interact with Vault's core system
and is not required for beginners.

Take a few moments to read and write some data to the new `kv` secrets engine
enabled at `kv/`. Here are a few ideas to get started:

```text
$ vault write kv/my-secret value="s3c(eT"

$ vault write kv/hello target=world

$ vault write kv/airplane type=boeing class=787

$ vault list kv
```

## Disable a Secrets Engine

When a secrets engine is no longer needed, it can be disabled. When a secrets
engine is disabled, all secrets are revoked and the corresponding Vault data and
configuration is removed. Any requests to route data to the original path would
result in an error, but another secrets engine could now be enabled at that
path.

If, for some reason, Vault is unable to delete the data or revoke the leases,
the disabling operation will fail. If this happens, the secrets engine will
remain enabled and available, but the request will return an error.

```text
$ vault secrets disable kv/
Success! Disabled the secrets engine (if it existed) at: kv/
```

Note that this command takes a PATH to the secrets engine as an argument, not
the TYPE of the secrets engine.

In addition to disabling a secrets engine, it is also possible to "move" a
secrets engine to a new path. This is still a disruptive command. All
configuration data is retained, but any secrets are revoked, since secrets are
closely tied to their engine's paths.

## What is a Secrets Engine?

Now that you've successfully enabled and disabled a secrets engine... what is
it? What is the point of a secrets engine?

As mentioned above, Vault behaves similarly to a [virtual filesystem][vfs]. The
read/write/delete/list operations are forwarded to the corresponding secrets
engine, and the secrets engine decides how to react to those operations.

This abstraction is incredibly powerful. It enables Vault to interface directly
with physical systems, databases, HSMs, etc. But in addition to these physical
systems, Vault can interact with more unique environments like AWS IAM, dynamic
SQL user creation, etc. all while using the same read/write interface.

## Next

You now know about secrets engines and how to operate on them. This is important
knowledge to move forward and learn about other secrets engines.

Next, we'll use the AWS backend to
[generate dynamic secrets](/intro/getting-started/dynamic-secrets.html).

[vfs]: https://en.wikipedia.org/wiki/Virtual_file_system
