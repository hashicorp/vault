---
layout: "docs"
page_title: "KV - Secrets Engines"
sidebar_current: "docs-secrets-kv"
description: |-
  The KV secrets engine can store arbitrary secrets.
---

# KV Secrets Engine

The `kv` secrets engine is used to store arbitrary secrets within the
configured physical storage for Vault.

Writing to a key in the `kv` backend will replace the old value; sub-fields are
not merged together.

Key names must always be strings. If you write non-string values directly via
the CLI, they will be converted into strings. However, you can preserve
non-string values by writing the key/value pairs to Vault from a JSON file or
using the HTTP API. 

This secrets engine honors the distinction between the `create` and `update`
capabilities inside ACL policies.

~> **Note**: Path and key names are _not_ obfuscated or encrypted; only the
values set on keys are. You should not store sensitive information as part of a
secret's path.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

The `kv` secrets engine is enabled by default at the path `secret/`. It can
be disabled, moved, or enabled multiple times at different paths. Each instance
of the KV secrets engine is isolated and unique.

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials. The `kv` secrets engine
allows for writing keys with arbitrary values.

1. Write arbitrary data:

    ```text
    $ vault write secret/my-secret my-value=s3cr3t
    Success! Data written to: secret/my-secret
    ```

1. Read arbitrary data:

    ```text
    $ vault read secret/my-secret
    Key                 Value
    ---                 -----
    refresh_interval    768h
    my-value            s3cr3t
    ```

## TTLs

Unlike other secrets engines, the KV secrets engine does not enforce TTLs
for expiration. Instead, the `lease_duration` is a hint for how often consumers
should check back for a new value. This is commonly displayed as
`refresh_interval` instead of `lease_duration` to clarify this in output.

If provided a key of `ttl`, the KV secrets engine will utilize this value
as the lease duration:

```text
$ vault write secret/my-secret ttl=30m my-value=s3cr3t
Success! Data written to: secret/my-secret
```

Even will a `ttl` set, the secrets engine _never_ removes data on its own. The
`ttl` key is merely advisory.

When reading a value with a `ttl`, both the `ttl` key _and_ the refresh interval
will reflect the value:

```text
$ vault read secret/my-secret
Key                 Value
---                 -----
refresh_interval    30m
my-value            s3cr3t
ttl                 30m
```

## API

The KV secrets engine has a full HTTP API. Please see the
[KV secrets engine API](/api/secret/kv/index.html) for more
details.
