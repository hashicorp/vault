---
layout: "docs"
page_title: "Secret Backend: AWS"
sidebar_current: "docs-secrets-aws"
description: |-
  The AWS secret backend for Vault generates access keys dynamically based on IAM policies.
---

# AWS Secret Backend

Name: `aws`

The AWS secret backend for Vault generates AWS access credentials dynamically
based on IAM policies. This makes IAM much easier to use: credentials could
be generated on the fly, and are automatically revoked when the Vault
lease is expired.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Quick Start

The first step to using the aws backend is to mount it.
Unlike the `generic` backend, the `aws` backend is not mounted by default.

```text
$ vault mount aws
Successfully mounted 'aws' at 'aws'!
```

Next, we must configure the root credentials that are used to manage IAM credentials:

```text
$ vault write aws/config/root \
    access_key=AKIAJWVN5Z4FOFT7NLNA \
    secret_key=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 \
    region=us-east-1
```

The following parameters are required:

- `access_key` - the AWS access key that has permission to manage IAM
  credentials.
- `secret_key` - the AWS secret key that has permission to manage IAM
  credentials.
- `region` the AWS region for API calls.

Note: the client uses the official AWS SDK and will use environment variable or IAM 
role-provided credentials if available.

The next step is to configure a role. A role is a logical name that maps
to a policy used to generated those credentials.
You can either supply a user inline policy (via the policy argument), or
provide a reference to an existing AWS policy by supplying the full ARN
reference (via the arn argument).

For example, lets first create a "deploy" role using an user inline policy as an example:

```text
$ vault write aws/roles/deploy \
    policy=@policy.json
```

This path will create a named role along with the IAM policy used
to restrict permissions for it. This is used to dynamically create
a new pair of IAM credentials when needed.

The `@` tells Vault to load the policy from the file named `policy.json`. Here
is an example IAM policy to get started:

```javascript
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "iam:*",
    "Resource": "*"
  }
}
```

As a second example, lets create a "readonly" role using an existing AWS policy as an example:

```text
$ vault write aws/roles/readonly arn=arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess
```

This path will create a named role pointing to an existing IAM policy used
to restrict permissions for it. This is used to dynamically create
a new pair of IAM credentials when needed.

For more information on IAM policies, please see the
[AWS IAM policy documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/PoliciesOverview.html).


To generate a new set of IAM credentials, we simply read from that role:

```text
$ vault read aws/creds/deploy
Key             Value
lease_id        aws/creds/deploy/7cb8df71-782f-3de1-79dd-251778e49f58
lease_duration  3600
access_key      AKIAIOMYUTSLGJOGLHTQ
secret_key      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
security_token  <nil>
```

If you run the command again, you will get a new set of credentials:

```text
$ vault read aws/creds/deploy
Key             Value
lease_id        aws/creds/deploy/82d89562-ff19-382e-6be9-cb45c8f6a42d
lease_duration  3600
access_key      AKIAJZ5YRPHFH3QHRRRQ
secret_key      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
security_token  <nil>
```

## Dynamic IAM users

The `aws/creds` endpoint will dynamically create a new IAM user and respond
with an IAM access key for the newly created user.

The [Quick Start](#quick-start) describes how to setup the `aws/creds` endpoint.

## Root Credentials for Dynamic IAM users

The `aws/config/root` credentials need permission to manage dynamic IAM users.
Here is an example IAM policy that would grant these permissions:

```javascript
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iam:AttachUserPolicy",
        "iam:CreateAccessKey",
        "iam:CreateUser",
        "iam:DeleteAccessKey",
        "iam:DeleteUser",
        "iam:DeleteUserPolicy",
        "iam:DetachUserPolicy",
        "iam:ListAccessKeys",
        "iam:ListAttachedUserPolicies",
        "iam:ListGroupsForUser",
        "iam:ListUserPolicies",
        "iam:PutUserPolicy",
        "iam:RemoveUserFromGroup"
      ],
      "Resource": [
        "arn:aws:iam::ACCOUNT-ID-WITHOUT-HYPHENS:user/vault-*"
      ]
    }
  ]
}
```

Note that this policy example is unrelated to the policy you wrote to `aws/roles/deploy`.
This policy example should be applied to the IAM user (or role) associated with 
the root credentials that you wrote to `aws/config/root`. You have to apply it
yourself in IAM. The policy you wrote to `aws/roles/deploy` is the policy you
want the AWS secret backend to apply to the temporary credentials it returns
from `aws/creds/deploy`.

Unfortunately, IAM credentials are eventually consistent with respect to other
Amazon services. If you are planning on using these credential in a pipeline,
you may need to add a delay of 5-10 seconds (or more) after fetching
credentials before they can be used successfully.

If you want to be able to use credentials without the wait, consider using the STS
method of fetching keys. IAM credentials supported by an STS token are available for use
as soon as they are generated.

## STS credentials

Vault also supports an STS credentials instead of creating a new IAM user.

The `aws/sts` endpoint will always fetch credentials with a 1hr ttl.
Unlike the `aws/creds` endpoint, the ttl is enforced by STS.

Vault supports two of the [STS APIs](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_request.html),
[STS federation tokens](http://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html) and
[STS AssumeRole](http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html).

### STS Federation Tokens

An STS federation token inherits a set of permissions that are the combination
(intersection) of three sets of permissions:

1. The permissions granted to the `aws/config/root` credentials
2. The user inline policy configured for the `aws/role`
3. An implicit deny policy on IAM or STS operations.

STS federation token credentials can only be generated for user inline
policies; the AWS GetFederationToken API does not support managed policies.

The `aws/config/root` credentials require IAM permissions for
`sts:GetFederationToken` and the permissions to delegate to the STS
federation token.  For example, this policy on the `aws/config/root` credentials
would allow creation of an STS federated token with delegated `ec2:*` permissions:

```javascript
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": [
      "ec2:*",
      "sts:GetFederationToken"
    ],
    "Resource": "*"
  }
}
```

Our "deploy" role would then assign an inline user policy with the same `ec2:*`
permissions.

```text
$ vault write aws/roles/deploy \
    policy=@policy.json
```

The policy.json file would contain an inline policy with similar permissions,
less the `sts:GetFederationToken` permission.  (We could grant `sts` permissions,
but STS would attach an implict deny on `sts` that overides the allow.)

```javascript
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "ec2:*",
    "Resource": "*"
  }
}
```

To generate a new set of STS federation token credentials, we simply read from
the role using the aws/sts endpoint:

```text
$vault read aws/sts/deploy
Key            	Value
lease_id       	aws/sts/deploy/31d771a6-fb39-f46b-fdc5-945109106422
lease_duration 	3600
lease_renewable	true
access_key     	ASIAJYYYY2AA5K4WIXXX
secret_key     	HSs0DYYYYYY9W81DXtI0K7X84H+OVZXK5BXXXX
security_token 	AQoDYXdzEEwasAKwQyZUtZaCjVNDiXXXXXXXXgUgBBVUUbSyujLjsw6jYzboOQ89vUVIehUw/9MreAifXFmfdbjTr3g6zc0me9M+dB95DyhetFItX5QThw0lEsVQWSiIeIotGmg7mjT1//e7CJc4LpxbW707loFX1TYD1ilNnblEsIBKGlRNXZ+QJdguY4VkzXxv2urxIH0Sl14xtqsRPboV7eYruSEZlAuP3FLmqFbmA0AFPCT37cLf/vUHinSbvw49C4c9WQLH7CeFPhDub7/rub/QU/lCjjJ43IqIRo9jYgcEvvdRkQSt70zO8moGCc7pFvmL7XGhISegQpEzudErTE/PdhjlGpAKGR3d5qKrHpPYK/k480wk1Ai/t1dTa/8/3jUYTUeIkaJpNBnupQt7qoaXXXXXXXXXX
```

### STS AssumeRole

STS AssumeRole is typically used for cross-account authentication or single sign-on (SSO)
scenarios.  AssumeRole has additional complexity compared STS federation tokens:

1. The ARN of a IAM role to assume
2. IAM inline policies and/or managed policies attached to the IAM role
3. IAM trust policy attached to the IAM role to grant privileges for one identity
   to assume the role.

AssumeRole adds a few benefits over federation tokens:

1. Assumed roles can invoke IAM and STS operations, if granted by the role's
   IAM policies.
2. Assumed roles support cross-account authenication

The `aws/config/root` credentials must have an IAM policy that allows `sts:AssumeRole`
against the target role:

```javascript
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "sts:AssumeRole",
    "Resource": "arn:aws:iam::ACCOUNT-ID-WITHOUT-HYPHENS:role/RoleNameToAssume"
  }
}
```

You must attach a trust policy to the target IAM role to assume, allowing
the aws/root/config credentials to assume the role.

```javascript
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::ACCOUNT-ID-WITHOUT-HYPHENS:user/VAULT-AWS-ROOT-CONFIG-USER-NAME"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

Finally, let's create a "deploy" policy using the arn of our role to assume:

```text
$ vault write aws/roles/deploy \
    arn=arn:aws:iam::ACCOUNT-ID-WITHOUT-HYPHENS:role/RoleNameToAssume
```

To generate a new set of STS assumed role credentials, we again read from
the role using the aws/sts endpoint:

```text
$vault read aws/sts/deploy
Key            	Value
lease_id       	aws/sts/deploy/31d771a6-fb39-f46b-fdc5-945109106422
lease_duration 	3600
lease_renewable	true
access_key     	ASIAJYYYY2AA5K4WIXXX
secret_key     	HSs0DYYYYYY9W81DXtI0K7X84H+OVZXK5BXXXX
security_token 	AQoDYXdzEEwasAKwQyZUtZaCjVNDiXXXXXXXXgUgBBVUUbSyujLjsw6jYzboOQ89vUVIehUw/9MreAifXFmfdbjTr3g6zc0me9M+dB95DyhetFItX5QThw0lEsVQWSiIeIotGmg7mjT1//e7CJc4LpxbW707loFX1TYD1ilNnblEsIBKGlRNXZ+QJdguY4VkzXxv2urxIH0Sl14xtqsRPboV7eYruSEZlAuP3FLmqFbmA0AFPCT37cLf/vUHinSbvw49C4c9WQLH7CeFPhDub7/rub/QU/lCjjJ43IqIRo9jYgcEvvdRkQSt70zO8moGCc7pFvmL7XGhISegQpEzudErTE/PdhjlGpAKGR3d5qKrHpPYK/k480wk1Ai/t1dTa/8/3jUYTUeIkaJpNBnupQt7qoaXXXXXXXXXX
```


## Troubleshooting

### Dynamic IAM user errors

If you get an error message similar to either of the following, the root credentials that you wrote to `aws/config/root` have insufficient privilege:

```text
$ vault read aws/creds/deploy
* Error creating IAM user: User: arn:aws:iam::000000000000:user/hashicorp is not authorized to perform: iam:CreateUser on resource: arn:aws:iam::000000000000:user/vault-root-1432735386-4059

$ vault revoke aws/creds/deploy/774cfb27-c22d-6e78-0077-254879d1af3c
Revoke error: Error making API request.

URL: PUT http://127.0.0.1:8200/v1/sys/revoke/aws/creds/deploy/774cfb27-c22d-6e78-0077-254879d1af3c
Code: 400. Errors:

* invalid request
```

If you get stuck at any time, simply run `vault path-help aws` or with a subpath for
interactive help output.

### STS federated token errors

Vault generates STS tokens using the IAM credentials passed to `aws/config`.

Those credentials must have two properties:

- They must have permissions to call `sts:GetFederationToken`.
- The capabilities of those credentials have to be at least as permissive as those requested
by policies attached to the STS creds.

If either of those conditions are not met, a "403 not-authorized" error will be returned.

See http://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html for more details.

Vault 0.5.1 or later is recommended when using STS tokens to avoid validation
errors for exceeding the AWS limit of 32 characters on STS token names.

## API

### /aws/config/root
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the root IAM credentials used.
    If static credentials are not provided using
    this endpoint, then the credentials will be retrieved from the
    environment variables `AWS_ACCESS_KEY`, `AWS_SECRET_KEY` and `AWS_REGION`
    respectively. If the credentials are still not found and if the
    backend is configured on an EC2 instance with metadata querying
    capabilities, the credentials are fetched automatically.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/aws/config/root`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">access_key</span>
        <span class="param-flags">required</span>
        The AWS Access Key
      </li>
      <li>
        <span class="param">secret_key</span>
        <span class="param-flags">required</span>
        The AWS Secret Key
      </li>
      <li>
        <span class="param">region</span>
        <span class="param-flags">required</span>
        The AWS region for API calls
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /aws/config/lease
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the lease settings for generated credentials.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/aws/config/lease`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">lease</span>
        <span class="param-flags">required</span>
        The lease value provided as a string duration
        with time suffix. Hour is the largest suffix.
      </li>
      <li>
        <span class="param">lease_max</span>
        <span class="param-flags">required</span>
        The maximum lease value provided as a string duration
        with time suffix. Hour is the largest suffix.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /aws/roles/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates a named role.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/aws/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">policy</span>
        <span class="param-flags">required (unless arn specified)</span>
        The IAM policy in JSON format.
      </li>
      <li>
        <span class="param">arn</span>
        <span class="param-flags">required (unless policy specified)</span>
        The full ARN reference to the desired existing policy
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Queries a named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/aws/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    ```javascript
    {
      "data": {
        "policy": "..."       
      }
    }
    ```
    ```javascript
    {
      "data": {
        "arn": "..."       
      }
    }
    ```
  </dd>
</dl>

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes a named role.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/aws/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

#### LIST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns a list of existing roles in the backend
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/aws/roles` (LIST) or `/aws/roles/?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
    ```javascript
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "keys": [
      "devrole",
      "prodrole",
      "testrole"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
    ```
  </dd>
</dl>


### /aws/creds/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a dynamic IAM credential based on the named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/aws/creds/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    ```javascript
    {
      "data": {
        "access_key": "...",
        "secret_key": "...",
        "security_token": null
      }
    }
    ```
  </dd>
</dl>


### /aws/sts/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
      Generates a dynamic IAM credential with an STS token based on the named
      role. The TTL will be 3600 seconds (one hour).
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/aws/sts/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
      None
  </dd>

  <dt>Returns</dt>
  <dd>
    ```javascript
    {
        "data": {
            "access_key": "...",
            "secret_key": "...",
            "security_token": "..."
        }
    }
    ```
    </dd>
</dl>

#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
      Generates a dynamic IAM credential with an STS token based on the named
      role and the given TTL (defaults to 3600 seconds, or one hour).
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/aws/sts/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
        The TTL to use for the STS token.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    ```javascript
    {
        "data": {
            "access_key": "...",
            "secret_key": "...",
            "security_token": "..."
        }
    }
    ```
  </dd>
</dl>
