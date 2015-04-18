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

Prior to using the GitHub auth backend, it must be configured. To
configure it, use the `/config` endpoint and pass in the following arguments:

  * `organization` (string, required) - The organization name a user must
       be a part of to authenticate.

After configuring that, you must map the teams of that organization to
policies within Vault. Use the `map/teams/<team>` endpoints to do that.
Example:

```
$ vault write auth/github/map/teams/owners value=root
...
```

The above would make anyone in the "owners" team a root user in Vault
(not recommended).
