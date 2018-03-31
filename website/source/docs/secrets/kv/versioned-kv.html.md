---
layout: "docs"
page_title: "KV - Secrets Engines"
sidebar_current: "docs-secrets-kv-versioned"
description: |-
  The KV secrets engine can store arbitrary secrets.
---

# KV Secrets Engine

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

The `kv` secrets engine is enabled by default at the path `secret/`. It can
be disabled, moved, or enabled multiple times at different paths. Each instance
of the KV secrets engine is isolated and unique.

## Upgrading from Non-Versioned

An existing non-versioned kv can be easily upgraded to a versioned key/value
store with the CLI command:

```
$ vault kv enable-versioning secret/
Success! Tuned the secrets engine at: secret/
```

or via the API:

```
$ cat payload.json
{
  "options": "versioned=true"
}

$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/mounts/secret/tune
```

This will start an upgrade process to upgrade the existing key/value data to
a versioned format. The mount will be inaccessible during this process. This
process could take a long time, so plan accordingly.  

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
    created_time     2018-03-30T22:11:48.589157362Z
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
    created_time     2018-03-30T22:11:48.589157362Z
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
    created_time     2018-03-30T22:18:37.124228658Z
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
    created_time     2018-03-30T22:18:37.124228658Z
    deletion_time    n/a
    destroyed        false
    version          2

    ====== Data ======
    Key         Value
    ---         -----
    my-value    new-s3cr3t
    ```

1. Previous versions can be accessed with the `-version` flag:

    ```
    $ vault kv get -version=1 secret/my-secret
    ====== Metadata ======
    Key              Value
    ---              -----
    created_time     2018-03-30T22:16:39.808909557Z
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
soft-delete. It will mark the version as deleted and populate a `deletion_time`
timestamp. Soft-deletes do not remove the underlying version data from storage,
this allows the version to be undeleted. The `vault kv undelete` commmand
handles undeleting versions. 

Version's data is permanently deleted only when the key has more versions than
are allowed by the max-versions setting, or when using `vault kv destroy`. When
the destroy command is used the underlying version data will be removed and the
key metadata will be marked as destroyed. If a version is cleaned up by going
over max-versions the version metadata will also be removed from the key.

See the commands below for more information:

1. The latest version of a key can be deleted with the delete command, this also
   takes a `-versions` flag to delete prior versions:
    
    ```
    $ vault kv delete secret/my-secret
    Success! Data deleted (if it existed) at: secret/my-secret
    ```

1. Versions can be undeleted:

    ```
    $ vault kv undelete -versions=2 secret/my-secret
    Success! Data written to: secret/undelete/my-secret

    $ vault kv get secret/my-secret
    ====== Metadata ======
    Key              Value
    ---              -----
    created_time     2018-03-30T22:18:37.124228658Z
    deletion_time    n/a
    destroyed        false
    version          2
    
    ====== Data ======
    Key         Value
    ---         -----
    my-value    new-s3cr3t
    ```

1. Destroying a version permanently deletes the underlying data:

    ```
    $ vault kv destroy -versions=2 secret/my-secret
    Success! Data written to: secret/destroy/my-secret
    ```

### Key Metadata

All versions and key metadata can be tracked with the metadata command & API.
Deleting the metadata key will cause all metadata and versions for that key to
be permanently removed.

See the commands below for more information:

1. All metadata and versions for a key can be viewed:
    
    ```
    $ vault kv metadata get secret/my-secret
    ======= Metadata =======
    Key                Value
    ---                -----
    created_time       2018-03-30T22:16:39.808909557Z
    current_version    2
    max_versions       0
    oldest_version     0
    updated_time       2018-03-30T22:18:37.124228658Z

    ====== Version 1 ======
    Key              Value
    ---              -----
    created_time     2018-03-30T22:16:39.808909557Z
    deletion_time    n/a
    destroyed        false
    
    ====== Version 2 ======
    Key              Value
    ---              -----
    created_time     2018-03-30T22:18:37.124228658Z
    deletion_time    n/a
    destroyed        true
    ```

1. The metadata settings for a key can be configured:

    ```
    $ vault kv metadata put -max-versions 1 secret/my-secret
    Success! Data written to: secret/metadata/my-secret
    ```

    Max versions changes will be applied on next write:

    ```
    $ vault kv put secret/my-secret my-value=newer-s3cr3t
    Key              Value
    ---              -----
    created_time     2018-03-30T22:41:09.193643571Z
    deletion_time    n/a
    destroyed        false
    version          3
    ```

    Once a key has more versions than max versions the oldest versions are cleaned
    up:

    ```
    $ vault kv metadata get secret/my-secret
    ======= Metadata =======
    Key                Value
    ---                -----
    created_time       2018-03-30T22:16:39.808909557Z
    current_version    3
    max_versions       1
    oldest_version     3
    updated_time       2018-03-30T22:41:09.193643571Z

    ====== Version 3 ======
    Key              Value
    ---              -----
    created_time     2018-03-30T22:41:09.193643571Z
    deletion_time    n/a
    destroyed        false
    ```

1. Permanently delete all metadata and versions for a key:

    ```
    $ vault kv metadata delete secret/my-secret
    Success! Data deleted (if it existed) at: secret/metadata/my-secret
    ```

## API

The KV secrets engine has a full HTTP API. Please see the
[KV secrets engine API](/api/secret/kv/versioned-kv.html) for more
details.
