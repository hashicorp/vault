---
layout: "docs"
page_title: "AliCloud - Auth Methods"
sidebar_title: "AliCloud"
sidebar_current: "docs-auth-alicloud"
description: |-
  The AliCloud auth method allows automated authentication of AliCloud entities.
---

# AliCloud Auth Method

The `alicloud` auth method provides an automated mechanism to retrieve
a Vault token for AliCloud entities. Unlike most Vault auth methods, this 
method does not require manual first-deploying, or provisioning 
security-sensitive credentials (tokens, username/password, client certificates, 
etc), by operators. It treats AliCloud as a Trusted Third Party and uses a 
special AliCloud request signed with private credentials. A variety of credentials 
can be used to construct the request, but AliCloud offers 
[instance metadata](https://www.alibabacloud.com/help/faq-detail/49122.htm) 
that's ideally suited for the purpose. By launching an instance with a role,
the role's STS credentials under instance metadata can be used to securely 
build the request.

## Authentication Workflow

The AliCloud STS API includes a method,
[`sts:GetCallerIdentity`](https://www.alibabacloud.com/help/doc-detail/43767.htm),
which allows you to validate the identity of a client. The client signs
a `GetCallerIdentity` query using the [AliCloud signature
algorithm](https://www.alibabacloud.com/help/doc-detail/67332.htm). It then 
submits 2 pieces of information to the Vault server to recreate a valid signed 
request: the request URL, and the request headers. The Vault server then 
reconstructs the query and forwards it on to the AliCloud STS service and validates 
the result back.

Importantly, the credentials used to sign the GetCallerIdentity request can come
from the ECS instance metadata service for an ECS instance, which obviates the
need for an operator to manually provision some sort of identity material first.
However, the credentials can, in principle, come from anywhere, not just from
the locations AliCloud has provided for you.

Each signed AliCloud request includes the current timestamp and a nonce to mitigate 
the risk of replay attacks.

It's also important to note that AliCloud does NOT include any sort
of authorization around calls to `GetCallerIdentity`. For example, if you have
a RAM policy on your credential that requires all access to be MFA authenticated,
non-MFA authenticated credentials will still be able to authenticate to Vault 
using this method. It does not appear possible to enforce a RAM principal to be 
MFA authenticated while authenticating to Vault.

## Authorization Workflow

The basic mechanism of operation is per-role. 

Roles are associated with a role ARN that has been pre-created in AliCloud. 
AliCloud's console displays each role's ARN. A role in Vault has a 1:1 relationship
with a role in AliCloud, and must bear the same name.

When a client assumes that role and sends its `GetCallerIdentity` request to Vault,
Vault matches the arn of its assumed role with that of a pre-created role in Vault.
It then checks what policies have been associated with the role, and grants a
token accordingly.

## Authentication

### Via the CLI

#### Enable AliCloud authentication in Vault.

```
$ vault auth enable alicloud
```

#### Configure the policies on the role.

```
$ vault write auth/alicloud/role/dev-role arn='acs:ram::5138828231865461:role/dev-role'
```

#### Perform the login operation

```
$ vault write auth/alicloud/login \
        role=dev-role \
        identity_request_url=$IDENTITY_REQUEST_URL_BASE_64 \
        identity_request_headers=$IDENTITY_REQUEST_HEADERS_BASE_64
```

For the RAM auth method, generating the signed request is a non-standard
operation. The Vault CLI supports generating this for you:

```
$ vault login -method=alicloud access_key=... secret_key=... security_token=... region=...
```

This assumes you have the AliCloud credentials you would find on an ECS instance using the 
following call:
```
curl 'http://100.100.100.200/latest/meta-data/ram/security-credentials/$ROLE_NAME'
```
Please note the `$ROLE_NAME` above is case-sensitive and must be consistent with how it's reflected
on the instance.

An example of how to generate the required request values for the `login` method
can be found found in the 
[Vault CLI source code](https://github.com/hashicorp/vault-plugin-auth-alicloud/blob/master/tools/tools.go).

## API

The AliCloud auth method has a full HTTP API. Please see the
[AliCloud Auth API](/api/auth/alicloud/index.html) for more
details.
