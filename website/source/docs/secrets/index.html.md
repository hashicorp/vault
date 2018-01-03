---
layout: "docs"
page_title: "Secrets Engines"
sidebar_current: "docs-secrets"
description: |-
  Secrets engines are mountable engines that store or generate secrets in Vault.
---

# Secrets Engines

Secrets engines are components which store, generate, or encrypt data. Secrets
engines are incredibly flexible, so it is easiest to think about them in terms
of their function. Secrets engines are provided some set of data, they take some
action on that data, and they return a result.

Some secrets engines simply store and read data - like encrypted
Redis/Memcached. Other secrets engines connect to other services and generate
dynamic credentials on demand. Other secrets engines provide encryption as a
service, totp generation, certificates, and much more.

Secrets engines are enabled at a "path" in Vault. When a request comes to Vault,
the router automatically routes anything with the route prefix to the secrets
engine. In this way, each secrets engine defines its own paths and properties.
To the user, secrets engines behave similar to a virtual filesystem, supporting
operations like read, write, and delete.

## Secrets Engines Lifecycle

Most secrets engines can be enabled, disabled, tuned, and moved via the CLI or
API. Previous versions of Vault called these "mounts", but that term was
overloaded.

- **Enable** - This enables a secrets engine at a given path. With few
  exceptions, secrets engines can be enabled at multiple paths. Each secrets
  engine is isolated to its path. By default, they are enabled at their "type"
  (e.g. "aws" enables at "aws/").

- **Disable** - This disables an existing secrets engine. When a secrets engine
  is disabled, all of its secrets are revoked (if they support it), and all of
  the data stored for that engine in the physical storage layer is deleted.

- **Move** - This moves the path for an existing secrets engine. This process
  revokes all secrets, since secret leases are tied to the path they were
  created at. The configuration data stored for the engine persists through the
  move.

- **Tune** - This tunes global configuration for the secrets engine such as the
  TTLs.

Once a secrets engine is enabled, you can interact with it directly at its path
according to its own API. Use `vault path-help` to determine the paths it
responds to.

Note that mount points cannot conflict with each other in Vault. There are
two broad implications of this fact. The first is that you cannot have
a mount which is prefixed with an existing mount. The second is that you
cannot create a mount point that is named as a prefix of an existing mount.
As an example, the mounts `foo/bar` and `foo/baz` can peacefully coexist
with each other whereas `foo` and `foo/baz` cannot

## Barrier View

Secrets engines receive a _barrier view_ to the configured Vault physical
storage. This is a lot like a [chroot](https://en.wikipedia.org/wiki/Chroot).

When a secrets engine is enabled, a random UUID is generated. This becomes the
data root for that engine. Whenever that engine writes to the physical storage
layer, it is prefixed with that UUID folder. Since the Vault storage layer
doesn't support relative access (such as `../`), this makes it impossible for a
enabled secrets engine to access other data.

This is an important security feature in Vault - even a malicious engine
cannot access the data from any other engine.
