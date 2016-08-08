---
layout: "docs"
page_title: "Auth Backend: GitHub"
sidebar_current: "docs-auth-github"
description: |-
  The GitHub auth backend allows authentication with Vault using GitHub.
---

# Auth Backend: GitHub

Name: `github`

The GitHub auth backend can be used to authenticate with Vault using
a GitHub personal access token.
This method of authentication is most useful for humans: operators or
developers using Vault directly via the CLI.

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
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "warnings": null,
  "auth": {
    "client_token": "c4f280f6-fdb2-18eb-89d3-589e2e834cdb",
    "policies": [
      "admins"
    ],
    "metadata": {
      "org": "test_org",
      "username": "rajanadar",
    },
    "lease_duration": 0,
    "renewable": false
  }
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
$ vault write auth/github/map/teams/admins value=admins
Success! Data written to: auth/github/map/teams/admins
```

The above would make anyone in the "admins" team receive tokens with the policy `admins`.

You can then auth with a user that is a member of the "admins" team using a Personal Access Token with the `read:org` scope.

GitHub token can also be supplied from the env variable `VAULT_AUTH_GITHUB_TOKEN`.

```
$ vault auth -method=github token=000000905b381e723b3d6a7d52f148a5d43c4b45
Successfully authenticated! The policies that are associated
with this token are listed below:

admins
```

