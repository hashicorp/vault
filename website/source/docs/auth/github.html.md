---
layout: "docs"
page_title: "GitHub - Auth Methods"
sidebar_current: "docs-auth-github"
description: |-
  The GitHub auth method allows authentication with Vault using GitHub.
---

# GitHub Auth Method

The `github` auth method can be used to authenticate with Vault using a GitHub
personal access token. This method of authentication is most useful for humans:
operators or developers using Vault directly via the CLI.

~> **IMPORTANT NOTE:** Vault does not support an OAuth workflow to generate
GitHub tokens, so does not act as a GitHub application. As a result, this method
uses personal access tokens. An important consequence is that any valid GitHub
access token with the `read:org` scope can be used for authentication. If such a
token is stolen from a third party service, and the attacker is able to make
network calls to Vault, they will be able to log in as the user that generated
the access token. When using this method it is a good idea to ensure that access
to Vault is restricted at a network level rather than public. If these risks are
unacceptable to you, you should use a different method.

## Authentication

### Via the CLI

The default path is `/github`. If this auth method was enabled at a different
path, specify `-path=/my-path` in the CLI.

```text
$ vault login -method=github token="MY_TOKEN"
```

### Via the API

The default endpoint is `auth/github/login`. If this auth method was enabled
at a different path, use that value instead of `github`.

```shell
$ curl \
    --request POST \
    --data '{"token": "MY_TOKEN"}' \
    https://vault.rocks/v1/auth/github/login
```

The response will contain a token at `auth.client_token`:

```json
{
  "auth": {
    "renewable": true,
    "lease_duration": 2764800,
    "metadata": {
      "username": "my-user",
      "org": "my-org"
    },
    "policies": [
      "default",
      "dev-policy"
    ],
    "accessor": "f93c4b2d-18b6-2b50-7a32-0fecf88237b8",
    "client_token": "1977fceb-3bfa-6c71-4d1f-b64af98ac018"
  }
}
```

## Configuration

Auth methods must be configured in advance before users or machines can
authenticate. These steps are usually completed by an operator or configuration
management tool.

1. Enable the GitHub auth method:

    ```text
    $ vault auth enable github
    ```

1. Use the `/config` endpoint to configure Vault to talk to GitHub. For the list
   of available configuration options, please see the API documentation.

    ```text
    $ vault write auth/github/config organization=hashicorp
    ```

    For the complete list of configuration options, please see the API
    documentation.

1. Map the users/teams of that GitHub organization to policies in Vault. Team
   names must be "slugified":

    ```text
    $ vault write auth/github/map/teams/dev value=dev-policy
    ```

    In this example, when members of the team "dev" in the organization
    "hashicorp" authenticate to Vault using a GitHub personal access token, they
    will be given a token with the "dev-policy" policy attached.

    ---

    You can also create mappings for a specific user `map/users/<user>`
    endpoint:

    ```text
    $ vault write auth/github/map/users/sethvargo value=sethvargo-policy
    ```

    In this example, a user with the GitHub username `sethvargo` will be
    assigned the `sethvargo-policy` policy **in addition to** any team policies.

## API

The GitHub auth method has a full HTTP API. Please see the
[GitHub Auth API](/api/auth/github/index.html) for more
details.
