---
layout: "intro"
page_title: "Dynamic Secrets - Getting Started"
sidebar_current: "gettingstarted-dynamicsecrets"
description: |-
  On this page we introduce dynamic secrets by showing you how to create AWS access keys with Vault.
---

# Dynamic Secrets

Now that you've experimented with the `kv` secrets engine, it is time to explore
another feature of Vault: _dynamic secrets_.

Unlike the `kv` secrets where you had to put data into the store yourself,
dynamic secrets are generated when they are accessed. Dynamic secrets do not
exist until they are read, so there is no risk of someone stealing them or
another client using the same secrets. Because Vault has built-in revocation
mechanisms, dynamic secrets can be revoked immediately after use, minimizing the
amount of time the secret existed.

-> **Note:** Before starting this page, please register for an
[AWS account](https://aws.amazon.com). We won't be using any features that
cost money, so you shouldn't be charged for anything. However, we are not
responsible for any charges you may incur.

## Enable the AWS Secrets Engine

Unlike the `kv` secrets engine which is enabled by default, the AWS secrets
engine must be enabled before use. This step is usually done via configuration
management.

```text
$ vault secrets enable -path=aws aws
Success! Enabled the aws secrets engine at: aws/
```

The AWS secrets engine is now enabled at `aws/`. As we covered in the previous
sections, different secrets engines allow for different behavior. In this case,
the AWS secrets engine generates dynamic, on-demand AWS access credentials.

## Configuring the AWS Secrets Engine

After enabling the AWS secrets engine, you must configure it to authenticate and
communicate with AWS. This requires privileged account credentials. If you are
unfamiliar with AWS, use your root account keys.

~> Do not use your root account keys in production. This is a getting started
guide and is not "best practices" for production installations.

```text
$ vault write aws/config/root \
    access_key=AKIAI4SGLQPBX6CSENIQ \
    secret_key=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
Success! Data written to: aws/config/root
```

These credentials are now stored in this AWS secrets engine. The engine will use
these credentials when communicating with AWS in future requests.

## Creating a Role

The next step is to configure a "role". A "role" in Vault is a human-friendly
identifier to an action. Think of it as a symlink.

Vault knows how to create an IAM user via the AWS API, but it does not know what
permissions, groups, and policies you want to attach to that user. This is where
roles come in - roles map your configuration options to those API calls.

For example, here is an IAM policy that enables all actions on EC2. When Vault
generates an access key, it will automatically attach this policy. The generated
access key will have full access to EC2 (as dictated by this policy), but not
IAM or other AWS services. If you are not familiar with AWS' IAM policies, that
is okay - just use this one for now.

```json
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

As mentioned above, we need to map this policy document to a named role. To do
that, write to `aws/roles/:name`:

```text
$ vault write aws/roles/my-role policy=-<<EOF
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
EOF
Success! Data written to: aws/roles/my-role
```

Again, we're using a special path here `aws/roles/:name` to write an IAM policy
to Vault. We just told Vault:

> When I ask for a credential for "my-role", create it and attach the IAM policy `{ "Version": "2012..." }`.

## Generating the Secret

Now that the AWS secrets engine is enabled and configured with a role, we can
ask Vault to generate an access key pair for that role by reading from
`aws/creds/:name` where `:name` corresponds to the name of an existing role:

```text
$ vault read aws/creds/my-role
Key                Value
---                -----
lease_id           aws/creds/my-role/0bce0782-32aa-25ec-f61d-c026ff22106e
lease_duration     768h
lease_renewable    true
access_key         AKIAJELUDIANQGRXCTZQ
secret_key         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
security_token     <nil>
```

Success! The access and secret key can now be used to perform any EC2 operations
within AWS. Notice that these keys are new, they are not the keys you entered
earlier. If you were to run the command a second time, you would get a new
access key pair. Each time you read from `aws/creds/:name`, Vault will connect
to AWS and generate a new IAM user and key pair.

Take careful note of the `lease_id` field in the output. This value is used for
renewal, revocation, and inspection. Copy this `lease_id` to your clipboard.
Note that the `lease_id` is the **full path**, not just the UUID at the end.

## Revoking the Secret

Vault will automatically revoke this credential after 768 hours (see
`lease_duration` in the output), but perhaps we want to revoke it early. Once
the secret is revoked, the access keys are no longer valid.

To revoke the secret, use `vault revoke` with the lease ID that was outputted
from `vault read` when you ran it:

```text
$ vault lease revoke aws/creds/my-role/0bce0782-32aa-25ec-f61d-c026ff22106
Success! Revoked lease: aws/creds/my-role/0bce0782-32aa-25ec-f61d-c026ff22106e
```

Done! If you login to your AWS account, you will see that no IAM users exist. If
you try to use the access keys that were generated, you will find that they no
longer work.

With such easy dynamic creation and revocation, you can hopefully begin to see
how easy it is to work with dynamic secrets and ensure they only exist for the
duration that they are needed.

## Next

On this page we experienced our first dynamic secret, and we also saw the
revocation system in action. Dynamic secrets are incredibly powerful. As time
goes on, we expect that more systems will support some sort of API to create
access credentials, and Vault will be ready to get the most value out of this
practice.

Before going further, we're going to take a quick detour to learn about the
[built-in help system](/intro/getting-started/help.html).
