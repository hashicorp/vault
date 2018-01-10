---
layout: "docs"
page_title: "Okta - Auth Methods"
sidebar_current: "docs-auth-okta"
description: |-
  The Okta auth method allows users to authenticate with Vault using Okta
  credentials.
---

# Okta Auth Method

The `okta` auth method allows authentication using Okta and user/password
credentials. This allows Vault to be integrated into environments using Okta.

The mapping of groups in Okta to Vault policies is managed by using the
`users/` and `groups/` paths.

## Authentication

### Via the CLI

The default path is `/okta`. If this auth method was enabled at a different
path, specify `-path=/my-path` in the CLI.

```text
$ vault login -method=okta username=my-username
```

### Via the API

The default endpoint is `auth/okta/login`. If this auth method was enabled
at a different path, use that value instead of `okta`.

```shell
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data '{"password": "MY_PASSWORD"}' \
    https://vault.rocks/v1/auth/okta/login/my-username
```

The response will contain a token at `auth.client_token`:

```json
{
  "auth": {
    "client_token": "c4f280f6-fdb2-18eb-89d3-589e2e834cdb",
    "policies": [
      "admins"
    ],
    "metadata": {
      "username": "mitchellh"
    }
  }
}
```

## Configuration

Auth methods must be configured in advance before users or machines can
authenticate. These steps are usually completed by an operator or configuration
management tool.

### Via the CLI

1. Enable the Okta auth method:

    ```text
    $ vault auth enable okta
    ```

1. Configure Vault to communicate with your Okta account:

    ```text
    $ vault write auth/okta/config \
        base_url="okta.com" \
        organization="dev-123456" \
        token="00KzlTNCqDf0enpQKYSAYUt88KHqXax6dT11xEZz_g"
    ```

    **If no token is supplied, Vault will function, but only locally configured
    group membership will be available. Without a token, groups will not be
    queried.**

    For the complete list of configuration options, please see the API
    documentation.

1. Map an Okta group to a Vault policy:

    ```text
    $ vault write auth/okta/groups/scientists policies=nuclear-reactor
    ```

    In this example, anyone who successfully authenticates via Okta who is a
    member of the "scientists" group will receive a Vault token with the
    "nuclear-reactor" policy attached.

    ---

    It is also possible to add users directly:

    ```text
    $ vault write auth/okta/groups/engineers policies=autopilot
    $ vault write auth/okta/users/tesla groups=engineers
    ```

    This adds the Okta user "tesla" to the "engineers" group, which maps to
    the "autopilot" Vault policy.

      **The user-policy mapping via group membership happens at token _creation
      time_. Any changes in group membership in Okta will not affect existing
      tokens that have already been provisioned. To see these changes, users
      will need to re-authenticate. You can force this by revoking the
      existing tokens.**

## API

The Okta auth method has a full HTTP API. Please see the
[Okta Auth API](/api/auth/okta/index.html) for more details.
