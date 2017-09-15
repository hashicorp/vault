---
layout: "docs"
page_title: "Consul Secret Backend"
sidebar_current: "docs-secrets-consul"
description: |-
  The Consul secret backend for Vault generates tokens for Consul dynamically.
---

# Consul Secret Backend

Name: `consul`

The Consul secret backend for Vault generates
[Consul](https://www.consul.io)
API tokens dynamically based on Consul ACL policies.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Quick Start

The first step to using the consul backend is to mount it.
Unlike the `kv` backend, the `consul` backend is not mounted by default.

```
$ vault mount consul
Successfully mounted 'consul' at 'consul'!
```

[Acquire a management token from
Consul](https://www.consul.io/docs/agent/http/acl.html#acl_create), using the
`acl_master_token` from your Consul configuration file or any other management
token:

```shell
$ curl \
    -H "X-Consul-Token: secret" \
    -X PUT \
    -d '{"Name": "sample", "Type": "management"}' \
    http://127.0.0.1:8500/v1/acl/create
```
```javascript
{
  "ID": "adf4238a-882b-9ddc-4a9d-5b6758e4159e"
}
```

Next, we must configure Vault to know how to contact Consul.
This is done by writing the access information:

```
$ vault write consul/config/access \
    address=127.0.0.1:8500 \
    token=adf4238a-882b-9ddc-4a9d-5b6758e4159e
Success! Data written to: consul/config/access
```

In this case, we've configured Vault to connect to Consul
on the default port with the loopback address. We've also provided
an ACL token to use with the `token` parameter. Vault must have a management
type token so that it can create and revoke ACL tokens.

The next step is to configure a role. A role is a logical name that maps
to a role used to generate those credentials. For example, lets create
a "readonly" role:

```
POLICY='key "" { policy = "read" }'
$ echo $POLICY | base64 | vault write consul/roles/readonly policy=-
Success! Data written to: consul/roles/readonly
```

The backend expects the policy to be base64 encoded, so we need to encode it
properly before writing. The policy language is [documented by
Consul](https://www.consul.io/docs/internals/acl.html), but we've defined a
read-only policy.

To generate a new set Consul ACL token, we simply read from that role:

```
$ vault read consul/creds/readonly
Key           	Value
lease_id      	consul/creds/readonly/c7a3bd77-e9af-cfc4-9cba-377f0ef10e6c
lease_duration	3600
token         	973a31ea-1ec4-c2de-0f63-623f477c2510
```

Here we can see that Vault has generated a new Consul ACL token for us.
We can test this token out, and verify that it is read-only:

```
$ curl 127.0.0.1:8500/v1/kv/foo?token=973a31ea-1ec4-c2de-0f63-623f477c2510
[{"CreateIndex":12,"ModifyIndex":53,"LockIndex":4,"Key":"foo","Flags":3304740253564472344,"Value":"YmF6"}]

$ curl -X PUT -d 'test' 127.0.0.1:8500/v1/kv/foo?token=973a31ea-1ec4-c2de-0f63-623f477c2510
Permission denied
```

## API

The Consul secret backend has a full HTTP API. Please see the
[Consul secret backend API](/api/secret/consul/index.html) for more
details.
