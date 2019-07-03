---
layout: "docs"
page_title: "KV - Secrets Engines"
sidebar_title: "K/V Version 2"
sidebar_current: "docs-secrets-kv-v2"
description: |-
  The KV secrets engine can store arbitrary secrets.
---

# KV Secrets Engine - Version 2

The `kv` secrets engine is used to store arbitrary secrets within the
configured physical storage for Vault.

Key names must always be strings. If you write non-string values directly via
the CLI, they will be converted into strings. However, you can preserve
non-string values by writing the key/value pairs to Vault from a JSON file or
using the HTTP API.

This secrets engine honors the distinction between the `create` and `update`
capabilities inside ACL policies.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

A v2 `kv` secrets engine can be enabled by:

```text
$ vault secrets enable -version=2 kv
```

Or, you can pass `kv-v2` as the secrets engine type:

```text
$ vault secrets enable kv-v2
```

Additionally, when running a dev-mode server, the v2 `kv` secrets engine is enabled by default at the
path `secret/` (for non-dev servers, it is currently v1). It can be disabled, moved, or enabled multiple times at
different paths. Each instance of the KV secrets engine is isolated and unique.

## Upgrading from Version 1

An existing version 1 kv store can be upgraded to a version 2 kv store via the CLI or API, as shown below. This will start an upgrade process to upgrade the existing key/value data to a versioned format. The mount will be inaccessible during this process. This process could take a long time, so plan accordingly.

Once upgraded to version 2, the former paths at which the data was accessible will no longer suffice. You will need to adjust user policies to add access to the version 2 paths as detailed in the [ACL Rules section below](/docs/secrets/kv/kv-v2.html#acl-rules). Similarly, users/applications will need to update the paths at which they interact with the kv data once it has been upgraded to version 2.

An existing version 1 kv can be upgraded to a version 2 KV store with the CLI command:

```text
$ vault kv enable-versioning secret/
Success! Tuned the secrets engine at: secret/
```

or via the API:

```text
$ cat payload.json
{
  "options": {
      "version": "2"
  }
}

$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/mounts/secret/tune
```

## ACL Rules

The version 2 kv store uses a prefixed API, which is different from the
version 1 API. Before upgrading from a version 1 kv the ACL rules
should be changed. Also different paths in the version 2 API can be ACL'ed
differently.

Writing and reading versions are prefixed with the `data/` path. This policy
that worked for the version 1 kv:

```
path "secret/dev/team-1/*" {
  capabilities = ["create", "update", "read"]
}
```

Should be changed to:

```
path "secret/data/dev/team-1/*" {
  capabilities = ["create", "update", "read"]
}
```

There are different levels of data deletion for this backend. To grant a policy
the permissions to delete the latest version of a key:

```
path "secret/data/dev/team-1/*" {
  capabilities = ["delete"]
}
```
To allow the policy to delete any version of a key:

```
path "secret/delete/dev/team-1/*" {
  capabilities = ["update"]
}
```

To allow a policy to undelete data:

```
path "secret/undelete/dev/team-1/*" {
  capabilities = ["update"]
}
```

To allow a policy to destroy versions:

```
path "secret/destroy/dev/team-1/*" {
  capabilities = ["update"]
}
```

To allow a policy to list keys:

```
path "secret/metadata/dev/team-1/*" {
  capabilities = ["list"]
}
```

To allow a policy to view metadata for each version:

```
path "secret/metadata/dev/team-1/*" {
  capabilities = ["read"]
}
```

To allow a policy to permanently remove all versions and metadata for a key:

```
path "secret/metadata/dev/team-1/*" {
  capabilities = ["delete"]
}
```

See the [API Specification](/api/secret/kv/kv-v2.html) for more
information.

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials. The `kv` secrets engine
allows for writing keys with arbitrary values.

### Writing/Reading arbitrary data

1. Write arbitrary data:

    ```text
    $ vault kv put secret/my-secret my-value=s3cr3t
	Key              Value
	---              -----
	created_time     2019-06-19T17:20:22.985303Z
	deletion_time    n/a
	destroyed        false
	version          1
    ```

1. Read arbitrary data:

    ```text
	$ vault kv get secret/my-secret
	====== Metadata ======
	Key              Value
	---              -----
	created_time     2019-06-19T17:20:22.985303Z
	deletion_time    n/a
	destroyed        false
	version          1

	====== Data ======
	Key         Value
	---         -----
	my-value    s3cr3t
    ```

1. Write another version, the previous version will still be accessible. The
   `-cas` flag can optionally be passed to perform a check-and-set operation. If
   not set the write will be allowed. If set to `-cas=0` a write will only be allowed
   if the key doesn’t exist. If the index is non-zero the write will only be
   allowed if the key’s current version matches the version specified in the
   cas parameter.

    ```text
    $ vault kv put -cas=1 secret/my-secret my-value=new-s3cr3t
	Key              Value
	---              -----
	created_time     2019-06-19T17:22:23.369372Z
	deletion_time    n/a
	destroyed        false
	version          2
    ```

1. Reading now will return the newest version of the data:

    ```text
    $ vault kv get secret/my-secret
	====== Metadata ======
	Key              Value
	---              -----
	created_time     2019-06-19T17:22:23.369372Z
	deletion_time    n/a
	destroyed        false
	version          2

	====== Data ======
	Key         Value
	---         -----
	my-value    new-s3cr3t
    ```

1. Previous versions can be accessed with the `-version` flag:

    ```text
    $ vault kv get -version=1 secret/my-secret
	====== Metadata ======
	Key              Value
	---              -----
	created_time     2019-06-19T17:20:22.985303Z
	deletion_time    n/a
	destroyed        false
	version          1

	====== Data ======
	Key         Value
	---         -----
	my-value    s3cr3t
    ```

### Deleting and Destroying Data

When deleting data the standard `vault kv delete` command will perform a
soft delete. It will mark the version as deleted and populate a `deletion_time`
timestamp. Soft deletes do not remove the underlying version data from storage,
which allows the version to be undeleted. The `vault kv undelete` command
handles undeleting versions.

A version's data is permanently deleted only when the key has more versions than
are allowed by the max-versions setting, or when using `vault kv destroy`. When
the destroy command is used the underlying version data will be removed and the
key metadata will be marked as destroyed. If a version is cleaned up by going
over max-versions the version metadata will also be removed from the key.

See the commands below for more information:

1. The latest version of a key can be deleted with the delete command, this also
   takes a `-versions` flag to delete prior versions:

    ```text
    $ vault kv delete secret/my-secret
	Success! Data deleted (if it existed) at: secret/my-secret
    ```

1. Versions can be undeleted:

    ```text
    $ vault kv undelete -versions=2 secret/my-secret
	Success! Data written to: secret/undelete/my-secret

    $ vault kv get secret/my-secret
	====== Metadata ======
	Key              Value
	---              -----
	created_time     2019-06-19T17:23:21.834403Z
	deletion_time    n/a
	destroyed        false
	version          2

	====== Data ======
	Key         Value
	---         -----
	my-value    short-lived-s3cr3t
    ```

1. Destroying a version permanently deletes the underlying data:

    ```text
    $ vault kv destroy -versions=2 secret/my-secret
	Success! Data written to: secret/destroy/my-secret
    ```

### Key Metadata

All versions and key metadata can be tracked with the metadata command & API.
Deleting the metadata key will cause all metadata and versions for that key to
be permanently removed.

See the commands below for more information:

1. All metadata and versions for a key can be viewed:

    ```text
    $ vault kv metadata get secret/my-secret
	========== Metadata ==========
	Key                     Value
	---                     -----
	cas_required            false
	created_time            2019-06-19T17:20:22.985303Z
	current_version         2
	delete_version_after    0s
	max_versions            0
	oldest_version          0
	updated_time            2019-06-19T17:22:23.369372Z

	====== Version 1 ======
	Key              Value
	---              -----
	created_time     2019-06-19T17:20:22.985303Z
	deletion_time    n/a
	destroyed        false

	====== Version 2 ======
	Key              Value
	---              -----
	created_time     2019-06-19T17:22:23.369372Z
	deletion_time    n/a
	destroyed        true
    ```

1. The metadata settings for a key can be configured:

    ```text
    $ vault kv metadata put -max-versions 2 -delete-version-after="3h25m19s" secret/my-secret
	Success! Data written to: secret/metadata/my-secret
    ```

	Delete-version-after settings will apply only to new versions. Max versions
	changes will be applied on next write:

    ```text
    $ vault kv put secret/my-secret my-value=newer-s3cr3t
	Key              Value
	---              -----
	created_time     2019-06-19T17:31:16.662563Z
	deletion_time    2019-06-19T20:56:35.662563Z
	destroyed        false
	version          4
    ```

	Once a key has more versions than max versions the oldest versions
	are cleaned up:

    ```text
    $ vault kv metadata get secret/my-secret
	========== Metadata ==========
	Key                     Value
	---                     -----
	cas_required            false
	created_time            2019-06-19T17:20:22.985303Z
	current_version         4
	delete_version_after    3h25m19s
	max_versions            2
	oldest_version          3
	updated_time            2019-06-19T17:31:16.662563Z

	====== Version 3 ======
	Key              Value
	---              -----
	created_time     2019-06-19T17:23:21.834403Z
	deletion_time    n/a
	destroyed        true

	====== Version 4 ======
	Key              Value
	---              -----
	created_time     2019-06-19T17:31:16.662563Z
	deletion_time    2019-06-19T20:56:35.662563Z
	destroyed        false
    ```

1. Permanently delete all metadata and versions for a key:

    ```text
    $ vault kv metadata delete secret/my-secret
	Success! Data deleted (if it existed) at: secret/metadata/my-secret
    ```

## API

The KV secrets engine has a full HTTP API. Please see the
[KV secrets engine API](/api/secret/kv/kv-v2.html) for more
details.
