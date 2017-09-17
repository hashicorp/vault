---
layout: "docs"
page_title: "Auth Backend: Okta"
sidebar_current: "docs-auth-okta"
description: |-
  The Okta auth backend allows users to authenticate with Vault using Okta credentials.
---

# Auth Backend: Okta

Name: `okta`

The Okta auth backend allows authentication using Okta
and user/password credentials. This allows Vault to be integrated
into environments using Okta.

The mapping of groups in Okta to Vault policies is managed by using the
`users/` and `groups/` paths.

## Authentication

#### Via the CLI

```
$ vault auth -method=okta username=mitchellh
Password (will be hidden):
Successfully authenticated! The policies that are associated
with this token are listed below:

admins
```

#### Via the API

The endpoint for the login is `auth/okta/login/<username>`.

The password should be sent in the POST body encoded as JSON.

```shell
$ curl $VAULT_ADDR/v1/auth/okta/login/mitchellh \
    -d '{ "password": "foo" }'
```

The response will be in JSON. For example:

```javascript
{
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "auth": {
    "client_token": "c4f280f6-fdb2-18eb-89d3-589e2e834cdb",
    "policies": [
      "admins"
    ],
    "metadata": {
      "username": "mitchellh"
    },
    "lease_duration": 0,
    "renewable": false
  }
}
```

## Configuration

First, you must enable the Okta auth backend:

```
$ vault auth-enable okta
Successfully enabled 'okta' at 'okta'!
```

Now when you run `vault auth -methods`, the Okta backend is available:

```
Path       Type      Description
okta/      okta
token/     token     token based credentials
```

To use the Okta auth backend, it must first be configured for your Okta account.
The configuration options are categorized and detailed below.

Configuration is written to `auth/okta/config`.

### Connection parameters

* `org_name` (string, required) - The Okta organization.  This will be the first part of the url `https://XXX.okta.com` url.
* `api_token` (string, optional) - The Okta API token.  This is required to query Okta for user group membership. If this is not supplied only locally configured groups will be enabled. This can be generated from http://developer.okta.com/docs/api/getting_started/getting_a_token.html
* `base_url` (string, optional) - The Okta url. Examples: `oktapreview.com`, The default is `okta.com`
* `max_ttl` (string, optional) - Maximum duration after which authentication will be expired.
 Either number of seconds or in a format parsable by Go's [time.ParseDuration](https://golang.org/pkg/time/#ParseDuration)
* `ttl` (string, optional) - Duration after which authentication will be expired.
 Either number of seconds or in a format parsable by Go's [time.ParseDuration](https://golang.org/pkg/time/#ParseDuration)

Use `vault path-help` for more details.

## Examples:

### Scenario 1

* Okta organization `XXXTest`.
* With no token supplied only locally configured group membership will be available.  Groups will not be queried from Okta.

```
$ vault write auth/okta/config \
    org_name="XXXTest"
...
```

### Scenario 2

* Okta organization `dev-123456`.
* Okta base_url for developer account `oktapreview.com`
* API token `00KzlTNCqDf0enpQKYSAYUt88KHqXax6dT11xEZz_g`. This will allow group membership to be queried.

```
$ vault write auth/okta/config base_url="oktapreview.com" \
    org_name="dev-123456" \
    api_token="00KzlTNCqDf0enpQKYSAYUt88KHqXax6dT11xEZz_g" 
...
```

## Okta Group -> Policy Mapping

Next we want to create a mapping from an Okta group to a Vault policy:

```
$ vault write auth/okta/groups/scientists policies=foo,bar
```

This maps the Okta group "scientists" to the "foo" and "bar" Vault policies.

We can also add specific Okta users to additional (potentially non-Okta) groups:

```
$ vault write auth/okta/groups/engineers policies=foobar
$ vault write auth/okta/users/tesla groups=engineers
```

This adds the Okta user "tesla" to the "engineers" group, which maps to
the "foobar" Vault policy.

Finally, we can test this by authenticating:

```
$ vault auth -method=okta username=tesla
Password (will be hidden):
Successfully authenticated! The policies that are associated
with this token are listed below:

bar, foo, foobar
```

## Note on Okta Group's

Groups can only be pulled from Okta if an API token is configured via `token`

## Note on policy mapping

It should be noted that user -> policy mapping (via group membership) happens at token creation time. And changes in group membership in Okta will not affect tokens that have already been provisioned. To see these changes, old tokens should be revoked and the user should be asked to reauthenticate.

## API

The Okta authentication backend has a full HTTP API. Please see the
[Okta Auth API](/api/auth/okta/index.html) for more
details.