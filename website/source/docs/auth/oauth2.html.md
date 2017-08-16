---
layout: "docs"
page_title: "Auth Backend: Oauth2"
sidebar_current: "docs-auth-oauth2"
description: |-
  The Oauth2 auth backend allows users to authenticate with Vault using an external Oauth2 provider.
---

# Auth Backend: Oauth2

Name: `oauth2`

The Oauth2 auth backend allows authentication using an external Oauth2
provider using the [Resource Owner Password flow](https://tools.ietf.org/html/rfc6749#section-1.3.3).

The mapping of users to Vault policies is managed by using the
`users/` and `groups/` paths.  A user can be directly assigned to
policies and groups in Vault, or group membership can be queried
from a userinfo API.

## Configuration

First, you must enable the Oauth2 auth backend:

```
$ vault auth-enable oauth2
Successfully enabled 'oauth2' at 'oauth2'!
```

Now when you run `vault auth -methods`, the Oauth2 backend is available:

```
Path       Type      Description
oauth2/    oauth2
token/     token     token based credentials
```

To use the Oauth2 auth backend, it must first be configured for your Oauth2 account.
The configuration options are categorized and detailed below.

Configuration is written to `auth/oauth2/config`.

### Connection parameters

* `provider_url` (string, required) - The URL of the external Oauth2 provider token endpoint to be used to authenticate user-supplied credentials.
* `client_id` (string, optional) - The Oauth2 client ID that Vault should use to identify itself with the Oauth2 provider.  This is passed in the HTTP authorization header.  Whether or not this is required depends on how your oauth provider authenticates these requests.
* `client_secret` (string, optional) - The Oauth2 client secret that Vault should use in conjunction with the `client_id` to identify itself with the Oauth2 provider.
* `userinfo_url` (string, optional) - The URL of an HTTP/JSON API that provides data regarding group assignments for a particular user.  Queried following a successful authentication in order to determine which policies to grant for a user.
* `userinfo_group_key` (string, optional) - The field name in the returned userinfo JSON object that contains a comma-separated list of groups to which the authenticated user belongs.
* `scope` (string, optional) - An Oauth2 scope that allows access to the group assignment data in the above `userinfo_url` API.

Use `vault path-help` for more details.

## Authentication

#### Via the CLI

```
$ vault auth -method=oauth2 username=bob
Password (will be hidden):
Successfully authenticated! You are now logged in.
The token below is already saved in the session. You do not
need to "vault auth" again with the token.
token: 0f0a1dc3-c1f0-2a15-d59b-b091d6699644
token_duration: 2764799
token_policies: [default admins]
```

#### Via the API

The endpoint for the login is `auth/oauth2/login/<username>`.

The password should be sent in the POST body encoded as JSON.

```shell
$ curl $VAULT_ADDR/v1/auth/oauth2/login/bob \
    -d '{ "password": "foo" }'
```

The response will be in JSON. For example:

```javascript
{
  "request_id": "e3a8fd7e-0299-ca17-b689-30c0de17325a",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {},
  "wrap_info": null,
  "warnings": [],
  "auth": {
    "client_token": "0787bd87-511a-bb43-4776-9cf686c02e78",
    "accessor": "0b0f0173-a692-a000-0514-4ddf6328042c",
    "policies": [
      "default",
      "admins"
    ],
    "metadata": {
      "policies": ",default,admins",
      "username": "bob"
    },
    "lease_duration": 2764800,
    "renewable": true
  }
}
```

## Examples:

### Scenario 1

* Vault configured with a client ID & Secret to identify itself with provider.
* With no userinfo URL supplied only locally configured group membership will be available.  Groups will not be queried from the provider.

```
$ vault write auth/oauth2/config \
    provider_url='https://www.example.com/oauth/token' \
    client_id='VAULT-CLIENT' \
    client_secret='1234567890zxcvbnmASDFGHJKL'
```

### Scenario 2

* No client secret necessary from Vault
* A userinfo URL is configured to allow querying for group membership from Oauth2 provider.

```
$ vault write auth/oauth2/config \
    provider_url='https://www.example.com/oauth/token' \
    client_id='VAULT-CLIENT' \
    userinfo_url='https://www.example.com/oauth/userinfo' \
    userinfo_group_key='vault_entitlements' \
    scope='vault'
```

## Oauth2 Group -> Policy Mapping

Next we want to create a mapping from an Oauth2 group to a Vault policy:

```
$ vault write auth/oauth2/groups/scientists policies=foo,bar
```

This maps the group "scientists" coming back in the userinfo group list to the
"foo" and "bar" Vault policies.

We can also add specific Oauth2 users to additional (potentially non-Oauth2) groups:

```
$ vault write auth/oauth2/groups/engineers policies=foobar
$ vault write auth/oauth2/users/tesla groups=engineers
```

This adds the Oauth2 user "tesla" to the "engineers" group, which maps to
the "foobar" Vault policy.

Finally, we can test this by authenticating:

```
$ vault auth -method=oauth2 username=tesla
Password (will be hidden):
Successfully authenticated! You are now logged in.
The token below is already saved in the session. You do not
need to "vault auth" again with the token.
token: 0f0a1dc3-c1f0-2a15-d59b-b091d6699644
token_duration: 2764799
token_policies: [default bar foo foobar]
```

## Note on groups from Oauth2

Groups can only be pulled from Oauth2 if a userinfo endpoint is configured.

## Note on policy mapping

It should be noted that user -> policy mapping (via group membership) happens at
token creation time. And changes in group membership in Oauth2 will not affect
tokens that have already been provisioned. To see these changes, old tokens
should be revoked and the user should be asked to re-authenticate.
