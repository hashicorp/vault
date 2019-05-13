---
layout: "docs"
page_title: "AWS - Secrets Engines"
sidebar_title: "AWS"
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

Vault supports three different types of credentials to retrieve from AWS:

1. `iam_user`: Vault will create an IAM user for each lease, attach the managed
   and inline IAM policies as specified in the role to the user, and then return
   the access key and secret key to the caller. IAM users have no session tokens
   and so no session token will be returned.
2. `assumed_role`: Vault will call
   [sts:AssumeRole](https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html)
   and return the access key, secret key, and session token to the caller.
3. `federation_token`: Vault will call
   [sts:GetFederationToken](https://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html)
   passing in the supplied AWS policy document and return the access key, secret
   key, and session token to the caller.

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

1. Configure a Vault role that maps to a set of permissions in AWS as well as an
   AWS credential type. When users generate credentials, they are generated
   against this role. An example:

    ```text
    $ vault write aws/roles/my-role \
        credential_type=iam_user \
        policy_document=-<<EOF
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
    this role, Vault will create an IAM user and attach the specified policy
    document to the IAM user. Vault will then create an access key and secret
    key for the IAM user and return these credentials. You supply a
    user inline policy and/or provide references to an existing AWS policy's full
    ARN:

    ```text
    $ vault write aws/roles/my-other-role \
        policy_arns=arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess,arn:aws:iam::aws:policy/IAMReadOnlyAccess \
        credential_type=iam_user \
        policy_document=-<<EOF
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

The above demonstrated usage with `iam_user` credential types. As mentioned,
Vault also supports `assumed_role` and `federation_token` credential types.

### STS Federation Tokens

~> **Notice:** Due to limitations in AWS, in order to use the `federation_token`
credential type, Vault **must** be configured with IAM user credentials. AWS
does not allow temporary credentials (such as those from an IAM instance
profile) to be used.

An STS federation token inherits a set of permissions that are the combination
(intersection) of three sets of permissions:

1. The permissions granted to the `aws/config/root` credentials
2. The user inline policy configured for the `aws/role`
3. An implicit deny policy on IAM or STS operations.

Roles with a `credential_type` of `federation_token` can only specify a
`policy_document` in the Vault role. AWS does not support support managed
policies.

The `aws/config/root` credentials require IAM permissions for
`sts:GetFederationToken` and the permissions to delegate to the STS
federation token.  For example, this policy on the `aws/config/root` credentials
would allow creation of an STS federated token with delegated `ec2:*`
permissions (or any subset of `ec2:*` permissions):

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

An `ec2_admin` role would then assign an inline policy with the same `ec2:*`
permissions.

```text
$ vault write aws/roles/ec2_admin \
    credential_type=federation_token \
    policy_document=@policy.json
```

The policy.json file would contain an inline policy with similar permissions,
less the `sts:GetFederationToken` permission.  (We could grant
`sts:GetFederationToken` permissions, but STS attaches attach an implicit deny
that overrides the allow.)

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
$ vault write aws/sts/ec2_admin ttl=60m
Key            	Value
lease_id       	aws/sts/ec2_admin/31d771a6-fb39-f46b-fdc5-945109106422
lease_duration 	60m0s
lease_renewable	true
access_key     	ASIAJYYYY2AA5K4WIXXX
secret_key     	HSs0DYYYYYY9W81DXtI0K7X84H+OVZXK5BXXXX
security_token 	AQoDYXdzEEwasAKwQyZUtZaCjVNDiXXXXXXXXgUgBBVUUbSyujLjsw6jYzboOQ89vUVIehUw/9MreAifXFmfdbjTr3g6zc0me9M+dB95DyhetFItX5QThw0lEsVQWSiIeIotGmg7mjT1//e7CJc4LpxbW707loFX1TYD1ilNnblEsIBKGlRNXZ+QJdguY4VkzXxv2urxIH0Sl14xtqsRPboV7eYruSEZlAuP3FLmqFbmA0AFPCT37cLf/vUHinSbvw49C4c9WQLH7CeFPhDub7/rub/QU/lCjjJ43IqIRo9jYgcEvvdRkQSt70zO8moGCc7pFvmL7XGhISegQpEzudErTE/PdhjlGpAKGR3d5qKrHpPYK/k480wk1Ai/t1dTa/8/3jUYTUeIkaJpNBnupQt7qoaXXXXXXXXXX
```

### STS AssumeRole

The `assumed_role` credential type is typically used for cross-account
authentication or single sign-on (SSO) scenarios. In order to use an
`assumed_role` credential type, you must configure outside of Vault:

1. An IAM role
2. IAM inline policies and/or managed policies attached to the IAM role
3. IAM trust policy attached to the IAM role to grant privileges for Vault to
   assume the role

`assumed_role` credentials offer a few benefits over `federation_token`:

1. Assumed roles can invoke IAM and STS operations, if granted by the role's
   IAM policies.
2. Assumed roles support cross-account authentication
3. Temporary credentials (such as those granted by running Vault on an EC2
   instance in an IAM instance profile) can retrieve `assumed_role` credentials
   (but cannot retrieve `federation_token` credentials).

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

When specifying a Vault role with a `credential_type` of `assumed_role`, you can
specify more than one IAM role ARN. If you do so, Vault clients can select which
role ARN they would like to assume when retrieving credentials from that role.
You can further specify a `policy_document` which, if specified, acts as a
filter on the IAM permissions granted to the assumed role. For an action to be
allowed, it must be permitted by both the IAM policy on the AWS role that is
assumed as well as the `policy_document` specified on the Vault role. (The
`policy_document` parameter is passed in as the `Policy` parameter to the
[sts:AssumeRole](https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html)
API call.)

Note: When multiple `role_arns` are specified, clients requesting credentials
can specify any of the role ARNs that are defined on the Vault role in order to
retrieve credentials. However, when a `policy_document` is specified, that will
apply to ALL role credentials retrieved from AWS.

Let's create a "deploy" policy using the arn of our role to assume:

```text
$ vault write aws/roles/deploy \
    role_arns=arn:aws:iam::ACCOUNT-ID-WITHOUT-HYPHENS:role/RoleNameToAssume \
    credential_type=assumed_role
```

To generate a new set of STS assumed role credentials, we again write to
the role using the aws/sts endpoint:

```text
$ vault write aws/sts/deploy ttl=60m
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
