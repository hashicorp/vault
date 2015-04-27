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
on every path, use `vault help` after mounting the backend.

## Quick Start

Mount the aws secret backend using the `vault mount` command:

```text
$ vault mount aws
Successfully mounted 'aws' at 'aws'!
```

Configure the root credentials that are used to manage IAM credentials:

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

Create an IAM policy:

```text
$ vault write aws/policy/deploy \
    name=deploy \
    policy=@policy.json
```

This path will generate a new, never before used key pair for
accessing AWS. The IAM policy used to back this key pair will be
the "name" parameter, which is "deploy" in this example.

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

For more information on IAM policies, please see the
[AWS IAM policy documentation](http://docs.aws.amazon.com/IAM/latest/UserGuide/PoliciesOverview.html).

Vault can now generate IAM credentials under the given policy:

```text
$ vault read aws/deploy
Key             Value
lease_id        aws/deploy/7cb8df71-782f-3de1-79dd-251778e49f58
lease_duration  3600
access_key      AKIAIOMYUTSLGJOGLHTQ
secret_key      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
```

If you run the command again, you will get a new set of credentials:

```text
$ vault read aws/deploy
Key             Value
lease_id        aws/deploy/82d89562-ff19-382e-6be9-cb45c8f6a42d
lease_duration  3600
access_key      AKIAJZ5YRPHFH3QHRRRQ
secret_key      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
```

If you get stuck at any time, simply run `vault help aws` or with a subpath for
interactive help output.

## API

### /aws/config/root
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the root IAM credentials used.
    This is a root protected endpoint.
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
    This is a root protected endpoint.
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

### /aws/policy/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates a named policy.
    This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/aws/policy/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">policy</span>
        <span class="param-flags">required</span>
        The IAM policy in JSON format.
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
    Queries a named policy.
    This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/aws/policy/<name>`</dd>

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

  </dd>
</dl>

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes a named policy.
    This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/aws/policy/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>


### /aws/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a dynamic IAM credential based on the named policy.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/aws/<name>`</dd>

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
            "secret_key": "..."
        }
    }
    ```

  </dd>
</dl>
