---
layout: "docs"
page_title: "Alibaba - Auth Methods"
sidebar_current: "docs-auth-alibaba"
description: |-
  The Alibaba auth method allows automated authentication of Alibaba entities.
---

# Alibaba Auth Method

The `alibaba` auth method provides an automated mechanism to retrieve
a Vault token for Alibaba ECS instances using IAM roles. Unlike most Vault
auth methods, this method does not require manual first-deploying, or
provisioning security-sensitive credentials (tokens, username/password, client
certificates, etc), by operators. It treats Alibaba as a Trusted Third Party 
and uses a special Alibaba request signed with Alibaba RAM credentials. RAM 
credentials are automatically supplied to Alibaba instances in RAM metadata, 
and it is this information that Vault can use to authenticate clients.

## Authentication Workflow

The Alibaba STS API includes a method,
[`sts:GetCallerIdentity`](https://www.alibabacloud.com/help/doc-detail/43767.htm),
which allows you to validate the identity of a client. The client signs
a `GetCallerIdentity` query using the [Alibaba signature
algorithm](https://www.alibabacloud.com/help/doc-detail/67332.htm). It then 
submits 2 pieces of information to the Vault server to recreate a valid signed 
request: the request URL, and the request headers. The Vault server then 
reconstructs the query and forwards it on to the Alibaba STS service and validates 
the result back.

Importantly, the credentials used to sign the GetCallerIdentity request can come
from the ECS instance metadata service for an ECS instance, which obviates the
need for an operator to manually provision some sort of identity material first.
However, the credentials can, in principle, come from anywhere, not just from
the locations Alibaba has provided for you.

Each signed Alibaba request includes the current timestamp and a nonce to mitigate 
the risk of replay attacks.

It's also important to note that Alibaba does NOT include any sort
of authorization around calls to `GetCallerIdentity`. For example, if you have
a RAM policy on your credential that requires all access to be MFA authenticated,
non-MFA authenticated credentials will still be able to authenticate to Vault 
using this method. It does not appear possible to enforce a RAM principal to be 
MFA authenticated while authenticating to Vault.

## Authorization Workflow

The basic mechanism of operation is per-role. 

Roles are associated with a role ARN that has been pre-created in Alibaba. 
Alibaba's console displays each role's ARN. A role in Vault has a 1:1 relationship
with a role in Alibaba, and must bear the same name.

When a client assumes that role and sends its `GetCallerIdentity` request to Vault,
Vault matches the arn of its assumed role with that of a pre-created role in Vault.
It then checks what policies have been associated with the role, and grants a
token accordingly.

## Authentication

### Via the CLI

#### Enable Alibaba authentication in Vault.

```
$ vault auth enable alibaba
```

#### Configure the policies on the role.

```
$ vault write auth/alibaba/role/elk arn='acs:ram::5138828231865461:role/elk'
```

#### Perform the login operation

```
$ vault write auth/alibaba/login role=dev-role \
        identity_request_url=$IDENTITY_REQUEST_URL_BASE_64 \
        identity_request_headers=$IDENTITY_REQUEST_HEADERS_BASE_64
```

For the RAM auth method, generating the signed request is a non-standard
operation. The Vault cli supports generating this for you:

```
$ vault login -method=alibaba access_key_id=... access_key_secret=... security_token=... region=...
```

This assumes you have the Alibaba credentials you would find on an ECS instance using the 
following call:
```
curl 'http://100.100.100.200/latest/meta-data/ram/security-credentials/$ROLE_NAME'
```
Please note the `$ROLE_NAME` above is case-sensitive and must be consistent with how it's reflected
on the instance.

An example of how to generate the required request values for the `login` method
can be found found in the 
[vault cli source code](https://github.com/hashicorp/vault-plugin-auth-alibaba/blob/master/cli.go).

## API

The Alibaba auth method has a full HTTP API. Please see the
[Alibaba Auth API](/api/auth/alibaba/index.html) for more
details.
