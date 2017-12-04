---
layout: "intro"
page_title: "Dynamic Secrets - Getting Started"
sidebar_current: "gettingstarted-dynamicsecrets"
description: |-
  On this page we introduce dynamic secrets by showing you how to create AWS access keys with Vault.
---

# Dynamic Secrets

Now that we've written basic secrets to Vault and we have an understanding
of the mount system, we're going to move on to the next core feature of
Vault: _dynamic secrets_.

Dynamic secrets are secrets that are generated when they're accessed,
and aren't statically written like we did in
[Your First Secret](/intro/getting-started/first-secret.html).
On this page, we'll use the built-in AWS secret backend to dynamically
generate AWS access keys.

The power of dynamic secrets is that they simply don't exist before
they're read, so there is no risk of someone stealing them or another
client using the same secrets. And because Vault has built-in revocation
mechanisms (covered later), the dynamic secret can be revoked right after
use, minimizing the amount of time the secret existed.

-> **Note:** Before starting this page, please register for an
[AWS account](https://aws.amazon.com). We won't be using any features that
cost money, so you shouldn't be charged for anything. However, we're not
responsible for any charges you may incur.

## Mounting the AWS Backend

Let's generate our first dynamic secret. We'll use the AWS backend to
dynamically generate an AWS access key pair. First, mount the AWS backend:

```
$ vault mount aws
Successfully mounted 'aws' at 'aws'!
```

The AWS backend is now mounted at `aws/`. As we covered in a previous
section: different secret backends allow for different behavior, and in this
case the AWS backend is a dynamic backend for generating AWS access credentials.

## Configuring the AWS Backend

With the AWS backend mounted, the first step is to configure it with
the AWS credentials that will be used to create the other credentials.
For now, use the root keys for your AWS account.

To configure the backend, we use `vault write` to a special path
`aws/config/root`:

```
$ vault write aws/config/root \
    access_key=AKIAI4SGLQPBX6CSENIQ \
    secret_key=z1Pdn06b3TnpG+9Gwj3ppPSOlAsu08Qw99PUW+eB
Success! Data written to: aws/config/root
```

Remember that secret backends can behave anyway they want when
reading/writing a path, so this path stores this configuration for
later. Notice you can't read it back:

```
$ vault read aws/config/root
Error reading aws/config/root: Error making API request.

URL: GET http://127.0.0.1:8200/v1/aws/config/root
Code: 405. Errors:

* unsupported operation
```

To help keep the credentials secure, the AWS backend doesn't let you
read them back even if you're using a root credential.

## Creating a Role

The next step is to configure the AWS backend with an IAM policy.
IAM is the system AWS uses for creating new credentials with limited
API permissions.

The AWS backend requires an IAM policy to associate created credentials
with. For this example, we'll write just one policy, but you can associate
many policies with the backend. Save a file named `policy.json` with the following contents:

```javascript
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "Stmt1426528957000",
      "Effect": "Allow",
      "Action": [
        "ec2:*"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
```

This is a basic IAM policy that lets the user perform any action within
Amazon EC2. With the policy saved, write it to Vault and create a new role:

```
$ vault write aws/roles/deploy policy=@policy.json
Success! Data written to: aws/roles/deploy
```

Again, we're using a special path here `aws/roles/<NAME>` to write
an IAM policy to Vault. We also used the special syntax `@filename` with
`vault write` to write the contents of a file.

## Generating the Secret

Now that we've configured the AWS backend and created a role, we can now
request an access key pair for that role. To do so, just read the
special path `aws/creds/<NAME>` where `NAME` is the role name:

```
$ vault read aws/creds/deploy
Key             Value
---             -----
lease_id        aws/creds/deploy/0d042c53-aa8a-7ce7-9dfd-310351c465e5
lease_duration  768h0m0s
lease_renewable true
access_key      AKIAJFN42DVCQWDHQYHQ
secret_key      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
security_token  <nil>
```

Success! The access and secret key can now be used to perform any EC2
operations within AWS. You can verify they work, if you want. Also notice
that these keys are new, they're not the keys you entered earlier.

The `lease_id` above is a special ID used for Vault for renewal,
revocation, etc. Copy and save your Lease ID now.

## Revoking the Secret

Let's complete the loop and revoke this secret now, purging it from
existence. Once the secret is revoked, the access keys will no longer
work.

To revoke the secret, use `vault revoke` with the lease ID that was
outputted from `vault read` when you ran it:

```
$ vault revoke aws/creds/deploy/0d042c53-aa8a-7ce7-9dfd-310351c465e5
Success! Revoked the secret with ID 'aws/creds/deploy/0d042c53-aa8a-7ce7-9dfd-310351c465e5', if it existed.
```

Done! If you look at your AWS account, you'll notice that no IAM users
exist. If you try to use the access keys that were generated, you'll
find that they no longer work.

With such easy dynamic creation and revocation, you can hopefully begin
to see how easy it is to work with dynamic secrets and ensure they only
exist for the duration that they're needed.

## Next

On this page we experienced our first dynamic secret, and we also saw
the revocation system in action. Dynamic secrets are incredibly powerful.
As time goes on, we expect that more systems will support some sort of
API to create access credentials, and Vault will be ready to get the
most value out of this practice.

Before going further, we're going to take a quick detour to learn
about the
[built-in help system](/intro/getting-started/help.html).
