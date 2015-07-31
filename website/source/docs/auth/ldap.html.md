---
layout: "docs"
page_title: "Auth Backend: LDAP"
sidebar_current: "docs-auth-ldap"
description: |-
  The "ldap" auth backend allows users to authenticate with Vault using LDAP credentials.
---

# Auth Backend: LDAP

Name: `ldap`

The "ldap" auth backend allows authentication using an existing LDAP
server and user/password credentials. This allows Vault to be integrated
into environments using LDAP without duplicating the user/pass configuration
in multiple places.

The mapping of groups in LDAP to Vault policies is managed by using the
`users/` and `groups/` paths.

## Authentication

#### Via the CLI

```
$ vault auth -method=ldap username=mitchellh
Password (will be hidden):
Successfully authenticated! The policies that are associated
with this token are listed below:

root
```

#### Via the API

The endpoint for the login is `auth/ldap/login/<username>`.

The password should be sent in the POST body encoded as JSON.

```shell
$ curl $VAULT_ADDR/v1/auth/ldap/login/mitchellh \
    -d '{ "password": "foo" }'
```

The response will be in JSON. For example:

```javascript
{
  "lease_id":"",
  "renewable":false,
  "lease_duration":0,
  "data":null,
  "auth":{
    "client_token":"c4f280f6-fdb2-18eb-89d3-589e2e834cdb",
    "policies":[
      "root"
    ],
    "metadata":{
      "username":"mitchellh"
    },
    "lease_duration":0,
    "renewable":false
  }
}
```

## Configuration

First, you must enable the ldap auth backend:

```
$ vault auth-enable ldap
Successfully enabled 'ldap' at 'ldap'!
```

Now when you run `vault auth -methods`, the ldap backend is available:

```
Path       Type      Description
ldap/      ldap
token/     token     token based credentials
```

To use the "ldap" auth backend, an operator must configure it with
the address of the LDAP server that is to be used. An example is shown below.
Use `vault path-help` for more details.

```
$ vault write auth/ldap/config url="ldap://ldap.forumsys.com" \
		userattr=uid \
        userdn="dc=example,dc=com" \
        groupdn="dc=example,dc=com" \
        upndomain="forumsys.com" \
        certificate=@ldap_ca_cert.pem \
        insecure_tls=false \
        starttls=true
...
```

The above configures the target LDAP server, along with the parameters
specifying how users and groups should be queried from the LDAP server.

Next we want to create a mapping from an LDAP group to a Vault policy:

```
$ vault write auth/ldap/groups/scientists policies=foo,bar
```

This maps the LDAP group "scientists" to the "foo" and "bar" Vault policies.

We can also add specific LDAP users to additional (potentially non-LDAP) groups:

```
$ vault write auth/ldap/groups/engineers policies=foobar
$ vault write auth/ldap/users/tesla groups=engineers
```

This adds the LDAP user "tesla" to the "engineers" group, which maps to
the "foobar" Vault policy.

Finally, we can test this by authenticating:

```
$ vault auth -method=ldap username=tesla
Password (will be hidden):
Successfully authenticated! The policies that are associated
with this token are listed below:

bar, foo, foobar
```

