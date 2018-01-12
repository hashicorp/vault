---
layout: "docs"
page_title: "Nomad Secret Backend"
sidebar_current: "docs-secrets-nomad"
description: |-
  The Nomad secret backend for Vault generates tokens for Nomad dynamically.
---

# Nomad Secret Backend

Name: `Nomad`

The Nomad secret backend for Vault generates
[Nomad](https://www.nomadproject.io)
API tokens dynamically based on pre-existing Nomad ACL policies.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

~> **Version information** ACLs are only available on Nomad 0.7.0 and above.

## Quick Start

The first step to using the vault backend is to mount it.
Unlike the `generic` backend, the `nomad` backend is not mounted by default.

```
$ vault mount nomad
Successfully mounted 'nomad' at 'nomad'!
```

Optionally, we can configure the lease settings for credentials generated
by Vault. This is done by writing to the `config/lease` key:

```
$ vault write nomad/config/lease ttl=3600 max_ttl=86400
Success! Data written to: nomad/config/lease
```

For a quick start, you can use the SecretID token provided by the [Nomad ACL bootstrap
process](https://www.nomadproject.io/guides/acl.html#generate-the-initial-token), although this
is discouraged for production deployments. 

```
$ nomad acl bootstrap
Accessor ID  = 95a0ee55-eaa6-2c0a-a900-ed94c156754e
Secret ID    = c25b6ca0-ea4e-000f-807a-fd03fcab6e3c
Name         = Bootstrap Token
Type         = management
Global       = true
Policies     = n/a
Create Time  = 2017-09-20 19:40:36.527512364 +0000 UTC
Create Index = 7
Modify Index = 7
```
The suggested pattern is to generate a token specifically for Vault, following the 
[Nomad ACL guide](https://www.nomadproject.io/guides/acl.html)

Next, we must configure Vault to know how to contact Nomad.
This is done by writing the access information:

```
$ vault write nomad/config/access \
    address=http://127.0.0.1:4646 \
    token=adf4238a-882b-9ddc-4a9d-5b6758e4159e
Success! Data written to: nomad/config/access
```

In this case, we've configured Vault to connect to Nomad
on the default port with the loopback address. We've also provided
an ACL token to use with the `token` parameter. Vault must have a management
type token so that it can create and revoke ACL tokens.

The next step is to configure a role. A role is a logical name that maps
to a set of policy names used to generate those credentials. For example, lets create
an "monitoring" role that maps to a "readonly" policy:

```
$ vault write nomad/role/monitoring policies=readonly
Success! Data written to: nomad/role/monitoring
```

The backend expects either a single or a comma separated list of policy names.

To generate a new Nomad ACL token, we simply read from that role:

```
$ vault read nomad/creds/monitoring
Key              Value
---              -----
lease_id         nomad/creds/monitoring/78ec3ef3-c806-1022-4aa8-1dbae39c760c
lease_duration   768h0m0s
lease_renewable  true
accessor_id      a715994d-f5fd-1194-73df-ae9dad616307
secret_id        b31fb56c-0936-5428-8c5f-ed010431aba9
```

Here we can see that Vault has generated a new Nomad ACL token for us.
We can test this token out, by reading it in Nomad (by it's accesor):

```
$ nomad acl token info a715994d-f5fd-1194-73df-ae9dad616307
Accessor ID  = a715994d-f5fd-1194-73df-ae9dad616307
Secret ID    = b31fb56c-0936-5428-8c5f-ed010431aba9
Name         = Vault example root 1505945527022465593
Type         = client
Global       = false
Policies     = [readonly]
Create Time  = 2017-09-20 22:12:07.023455379 +0000 UTC
Create Index = 138
Modify Index = 138
```

## API

The Nomad secret backend has a full HTTP API. Please see the
[Nomad secret backend API](/api/secret/nomad/index.html) for more
details.
