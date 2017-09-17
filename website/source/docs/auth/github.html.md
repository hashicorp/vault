---
layout: "docs"
page_title: "Auth Backend: GitHub"
sidebar_current: "docs-auth-github"
description: |-
  The GitHub auth backend allows authentication with Vault using GitHub.
---

# Auth Backend: GitHub

Name: `github`

The GitHub auth backend can be used to authenticate with Vault using a GitHub
personal access token. This method of authentication is most useful for humans:
operators or developers using Vault directly via the CLI.

**N.B.**: Vault does not support an OAuth workflow to generate GitHub tokens,
so does not act as a GitHub application. As a result, this backend uses
personal access tokens. An important consequence is that any valid GitHub
access token with the `read:org` scope can be used for authentication. If such
a token is stolen from a third party service, and the attacker is able to make
network calls to Vault, they will be able to log in as the user that generated
the access token. When using this backend it is a good idea to ensure that
access to Vault is restricted at a network level rather than public. If these
risks are unacceptable to you, you should use a different backend.

## Authentication

#### Via the CLI

```
$ vault auth -method=github token=<api token>
...
```

#### Via the API

The endpoint for the GitHub login is `auth/github/login`. 

The `github` mountpoint value in the url is the default mountpoint value.
If you have mounted the `github` backend with a different mountpoint, use that value.

The `token` should be sent in the POST body encoded as JSON.

```shell
$ curl $VAULT_ADDR/v1/auth/github/login \
    -d '{ "token": "your_github_personal_access_token" }'
```

The response will be in JSON. For example:

```javascript
{
  "auth": {
    "renewable": true,
    "lease_duration": 2764800,
    "metadata": {
      "username": "vishalnayak",
      "org": "hashicorp"
    },
    "policies": [
      "default",
      "dev-policy"
    ],
    "accessor": "f93c4b2d-18b6-2b50-7a32-0fecf88237b8",
    "client_token": "1977fceb-3bfa-6c71-4d1f-b64af98ac018"
  },
  "warnings": null,
  "wrap_info": null,
  "data": null,
  "lease_duration": 0,
  "renewable": false,
  "lease_id": "",
  "request_id": "3c346f3b-e089-39ab-a953-a349f2284e3c"
}
```

## Configuration

First, you must enable the GitHub auth backend:

```
$ vault auth-enable github
Successfully enabled 'github' at 'github'!
```

Now when you run `vault auth -methods`, the GitHub backend is available:

```
Path       Type      Description
github/    github
token/     token     token based credentials
```

Prior to using the GitHub auth backend, it must be configured. To
configure it, use the `/config` endpoint with the following arguments:

  * `organization` (string, required) - The organization name a user must
     be a part of to authenticate.
  * `base_url` (string, optional) - For GitHub Enterprise or other API-compatible
     servers, the base URL to access the server.
  * `max_ttl` (string, optional) - Maximum duration after which authentication will be expired.
     This must be a string in a format parsable by Go's [time.ParseDuration](https://golang.org/pkg/time/#ParseDuration)
  * `ttl` (string, optional) - Duration after which authentication will be expired.
     This must be a string in a format parsable by Go's [time.ParseDuration](https://golang.org/pkg/time/#ParseDuration)

###Generate a GitHub Personal Access Token
Access your Personal Access Tokens in GitHub at [https://github.com/settings/tokens](https://github.com/settings/tokens).
Generate a new Token that has the scope `read:org`. Save the generated token. This is what you will provide to vault.

For example:

```
$ vault write auth/github/config organization=hashicorp
Success! Data written to: auth/github/config
```

After configuring that, you must map the teams of that organization to
policies within Vault. Use the `map/teams/<team>` endpoints to do that.
Team names must be slugified, so if your team name is: `Some Amazing Team`, 
you will need to include it as: `some-amazing-team`. 
Example:

```
$ vault write auth/github/map/teams/dev value=dev-policy
Success! Data written to: auth/github/map/teams/dev
```

The above would make anyone in the `dev` team receive tokens with the policy
`dev-policy`.

You can then auth with a user that is a member of the `dev` team using a
Personal Access Token with the `read:org` scope.

You can also create mappings for specific users in a similar fashion with the 
`map/users/<user>` endpoint.
Example:

```
$ vault write auth/github/map/users/user1 value=user1-policy
Success! Data written to: auth/github/map/teams/user1
```

Now a user with GitHub username `user1` will be assigned the `user1-policy` on authentication, 
in addition to any team policies.

GitHub token can also be supplied from the env variable `VAULT_AUTH_GITHUB_TOKEN`.

```
$ vault auth -method=github token=000000905b381e723b3d6a7d52f148a5d43c4b45
Successfully authenticated! You are now logged in.
The token below is already saved in the session. You do not
need to "vault auth" again with the token.
token: 0d9ab511-bc25-4fb6-a58b-94ce12b8da9c
token_duration: 2764800
token_policies: [default dev-policy]
```

Clients can use this token to perform an allowed set of operations on all the
paths contained by the policy set.

## API

The GitHub authentication backend has a full HTTP API. Please see the
[GitHub Auth API](/api/auth/github/index.html) for more
details.
