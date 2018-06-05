---
layout: "docs"
page_title: "AWS - Secrets Engines"
sidebar_current: "docs-secrets-aws"
description: |-
  The AWS secrets engine for Vault generates access keys dynamically based on
  IAM policies.
---

# AWS Secrets Engine

The AWS secrets engine generates AWS access credentials dynamically based on IAM
policies. This generally makes working with AWS IAM easier, since it does not
involve clicking in the web UI. Additionally, the process is codified and mapped
to internal auth methods (such as LDAP). The AWS IAM credentials are time-based
and are automatically revoked when the Vault lease expires.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1. Enable the AWS secrets engine:

    ```text
    $ vault secrets enable aws
    Success! Enabled the aws secrets engine at: aws/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Configure the credentials that Vault uses to communicate with AWS to generate
the IAM credentials:

    ```text
    $ vault write aws/config/root \
        access_key=AKIAJWVN5Z4FOFT7NLNA \
        secret_key=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 \
        region=us-east-1
    ```

    Internally, Vault will connect to AWS using these credentials. As such,
    these credentials must be a superset of any policies which might be granted
    on IAM credentials. Since Vault uses the official AWS SDK, it will use the
    specified credentials. You can also specify the credentials via the standard
    AWS environment credentials, shared file credentials, or IAM role/ECS task
    credentials.  (Note that you can't authorize vault with IAM role credentials if you plan
    on using STS Federation Tokens, since the temporary security credentials
    associated with the role are not authorized to use GetFederationToken.)

    ~> **Notice:** Even though the path above is `aws/config/root`, do not use
    your AWS root account credentials. Instead generate a dedicated user or
    role.

1. Configure a role that maps a name in Vault to a policy or policy file in AWS.
When users generate credentials, they are generated against this role:

    ```text
    $ vault write aws/roles/my-role \
        policy=-<<EOF
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:*",
          "Resource": "*"
        }
      ]
    }
    EOF
    ```

    This creates a role named "my-role". When users generate credentials against
    this role, the resulting IAM credential will have the permissions specified
    in the policy provided as the argument.

    You can either supply a user inline policy or provide a reference to an
    existing AWS policy's full ARN:

    ```text
    $ vault write aws/roles/my-other-role \
        arn=arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess
    ```

    For more information on IAM policies, please see the
    [AWS IAM policy documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/PoliciesOverview.html).

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials.

1. Generate a new credential by reading from the `/creds` endpoint with the name
of the role:

    ```text
    $ vault read aws/creds/my-role
    Key                Value
    ---                -----
    lease_id           aws/creds/my-role/f3e92392-7d9c-09c8-c921-575d62fe80d8
    lease_duration     768h
    lease_renewable    true
    access_key         AKIAIOWQXTLW36DV7IEA
    secret_key         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
    security_token     <nil>
    ```

    Each invocation of the command will generate a new credential.

    Unfortunately, IAM credentials are eventually consistent with respect to
    other Amazon services. If you are planning on using these credential in a
    pipeline, you may need to add a delay of 5-10 seconds (or more) after
    fetching credentials before they can be used successfully.

    If you want to be able to use credentials without the wait, consider using
    the STS method of fetching keys. IAM credentials supported by an STS token
    are available for use as soon as they are generated.

## Example IAM Policy for Vault

The `aws/config/root` credentials need permission to manage dynamic IAM users.
Here is an example AWS IAM policy that grants the most commonly required
permissions Vault needs:

```json
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

## STS credentials

Vault also supports an STS credentials instead of creating a new IAM user.

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
but STS would attach an implicit deny on `sts` that overrides the allow.)

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

To generate a new set of STS federation token credentials, we simply write to
the role using the aws/sts endpoint:

```text
$vault write aws/sts/deploy ttl=60m
Key            	Value
lease_id       	aws/sts/deploy/31d771a6-fb39-f46b-fdc5-945109106422
lease_duration 	60m0s
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
2. Assumed roles support cross-account authentication

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

To generate a new set of STS assumed role credentials, we again write to
the role using the aws/sts endpoint:

```text
$vault write aws/sts/deploy ttl=60m
Key            	Value
lease_id       	aws/sts/deploy/31d771a6-fb39-f46b-fdc5-945109106422
lease_duration 	60m0s
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

The AWS secrets engine has a full HTTP API. Please see the
[AWS secrets engine API](/api/secret/aws/index.html) for more
details.
