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

## A Note on Escaping

**It is up to the administrator** to provide properly escaped DNs. This
includes the user DN, bind DN for search, and so on.

The only DN escaping performed by this backend is on usernames given at login
time when they are inserted into the final bind DN, and uses escaping rules
defined in RFC 4514.

Additionally, Active Directory has escaping rules that differ slightly from the
RFC; in particular it requires escaping of '#' regardless of position in the DN
(the RFC only requires it to be escaped when it is the first character), and
'=', which the RFC indicates can be escaped with a backslash, but does not
contain in its set of required escapes. If you are using Active Directory and
these appear in your usernames, please ensure that they are escaped, in
addition to being properly escaped in your configured DNs.

For reference, see [RFC 4514](https://www.ietf.org/rfc/rfc4514.txt) and this
[TechNet post on characters to escape in Active
Directory](http://social.technet.microsoft.com/wiki/contents/articles/5312.active-directory-characters-to-escape.aspx).

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
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "auth": {
    "client_token": "c4f280f6-fdb2-18eb-89d3-589e2e834cdb",
    "policies": [
      "root"
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

If your users are not located directly below the "userdn", e.g. in several
OUs like
```
    ou=users,dc=example,dc=com
ou=people    ou=external     ou=robots
```
you can also specify a `binddn` and `bindpass` for vault to search for the DN
of a user. This also works for the AD where a typical setup is to have user
DNs in the form `cn=Firstname Lastname,ou=Users,dc=example,dc=com` but you
want to login users using the `sAMAccountName` attribute. For that specify
```
$ vault write auth/ldap/config url="ldap://ldap.forumsys.com" \
    userattr=sAMAccountName \
    userdn="ou=users,dc=example,dc=com" \
    groupdn="dc=example,dc=com" \
    binddn="cn=vault,ou=users,dc=example,dc=com" \
    bindpass='My$ecrt3tP4ss' \
    certificate=@ldap_ca_cert.pem \
    insecure_tls=false \
    starttls=true
...
```
To discover the bind dn for a user with an anonymous bind, use the `discoverdn=true`
parameter and leave the `binddn` / `bindpass` empty.

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

