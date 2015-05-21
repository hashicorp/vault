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

The endpoint for the GitHub login is `/login`.

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
Example:

```
$ vault write auth/github/map/teams/owners value=root
Success! Data written to: auth/github/map/teams/owners
```

The above would make anyone in the "owners" team a root user in Vault
(not recommended).

You can then auth with a user that is a member of the "owners" team using a Personal Access Token with the `read:org` scope.

```
$ vault auth -method=github token=000000905b381e723b3d6a7d52f148a5d43c4b45
Successfully authenticated! The policies that are associated
with this token are listed below:

root
```

