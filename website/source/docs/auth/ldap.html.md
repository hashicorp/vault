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

The mapping of groups and users in LDAP to Vault policies is managed by using
the `users/` and `groups/` paths.

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

admins
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

To use the ldap auth backend, it must first be configured with connection
details for your LDAP server, information on how to authenticate users, and
instructions on how to query for group membership.
The configuration options are categorized and detailed below.

Configuration is written to `auth/ldap/config`.

### Connection parameters

* `url` (string, required) - The LDAP server to connect to. Examples: `ldap://ldap.myorg.com`, `ldaps://ldap.myorg.com:636`. This can also be a comma-delineated list of URLs, e.g. `ldap://ldap.myorg.com,ldaps://ldap.myorg.com:636`, in which case the servers will be tried in-order if there are errors during the connection process.
* `starttls` (bool, optional) - If true, issues a `StartTLS` command after establishing an unencrypted connection.
* `insecure_tls` - (bool, optional) - If true, skips LDAP server SSL certificate verification - insecure, use with caution!
* `certificate` - (string, optional) - CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded.

### Binding parameters

There are two alternate methods of resolving the user object used to authenticate the end user: _Search_ or _User Principal Name_. When using _Search_, the bind can be either anonymous or authenticated. User Principal Name is method of specifying users supported by Active Directory. More information on UPN can be found [here](https://msdn.microsoft.com/en-us/library/ms677605(v=vs.85).aspx#userPrincipalName).

#### Binding - Authenticated Search

* `binddn` (string, optional) - Distinguished name of object to bind when performing user and group search. Example: `cn=vault,ou=Users,dc=example,dc=com`
* `bindpass` (string, optional) - Password to use along with `binddn` when performing user search.
* `userdn` (string, optional) - Base DN under which to perform user search. Example: `ou=Users,dc=example,dc=com`
* `userattr` (string, optional) - Attribute on user attribute object matching the username passed when authenticating. Examples: `sAMAccountName`, `cn`, `uid`

#### Binding - Anonymous Search

* `discoverdn` (bool, optional) - If true, use anonymous bind to discover the bind DN of a user
* `userdn` (string, optional) - Base DN under which to perform user search. Example: `ou=Users,dc=example,dc=com`
* `userattr` (string, optional) - Attribute on user attribute object matching the username passed when authenticating. Examples: `sAMAccountName`, `cn`, `uid`
* `deny_null_bind` (bool, optional) - This option prevents users from bypassing authentication when providing an empty password. The default is `true`.

#### Binding - User Principal Name (AD)

* `upndomain` (string, optional) - userPrincipalDomain used to construct the UPN string for the authenticating user. The constructed UPN will appear as `[username]@UPNDomain`. Example: `example.com`, which will cause vault to bind as `username@example.com`.

### Group Membership Resolution

Once a user has been authenticated, the LDAP auth backend must know how to resolve which groups the user is a member of. The configuration for this can vary depending on your LDAP server and your directory schema. There are two main strategies when resolving group membership - the first is searching for the authenticated user object and following an attribute to groups it is a member of. The second is to search for group objects of which the authenticated user is a member of. Both methods are supported.

* `groupfilter` (string, optional) - Go template used when constructing the group membership query. The template can access the following context variables: \[`UserDN`, `Username`\]. The default is `(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))`, which is compatible with several common directory schemas. To support nested group resolution for Active Directory, instead use the following query: `(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))`.
* `groupdn` (string, required) - LDAP search base to use for group membership search. This can be the root containing either groups or users. Example: `ou=Groups,dc=example,dc=com`
* `groupattr` (string, optional) - LDAP attribute to follow on objects returned by `groupfilter` in order to enumerate user group membership. Examples: for groupfilter queries returning _group_ objects, use: `cn`. For queries returning _user_ objects, use: `memberOf`. The default is `cn`.

*Note*: When using _Authenticated Search_ for binding parameters (see above) the distinguished name defined for `binddn` is used for the group search.  Otherwise, the authenticating user is used to perform the group search.

Use `vault path-help` for more details.

## Examples:

### Scenario 1

* LDAP server running on `ldap.example.com`, port 389.
* Server supports `STARTTLS` command to initiate encryption on the standard port.
* CA Certificate stored in file named `ldap_ca_cert.pem`
* Server is Active Directory supporting the userPrincipalName attribute. Users are identified as `username@example.com`.
* Groups are nested, we will use `LDAP_MATCHING_RULE_IN_CHAIN` to walk the ancestry graph.
* Group search will start under `ou=Groups,dc=example,dc=com`. For all group objects under that path, the `member` attribute will be checked for a match against the authenticated user.
* Group names are identified using their `cn` attribute.

```
$ vault write auth/ldap/config \
    url="ldap://ldap.example.com" \
    userdn="ou=Users,dc=example,dc=com" \
    groupdn="ou=Groups,dc=example,dc=com" \
    groupfilter="(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))" \
    groupattr="cn" \
    upndomain="example.com" \
    certificate=@ldap_ca_cert.pem \
    insecure_tls=false \
    starttls=true
...
```

### Scenario 2

* LDAP server running on `ldap.example.com`, port 389.
* Server supports `STARTTLS` command to initiate encryption on the standard port.
* CA Certificate stored in file named `ldap_ca_cert.pem`
* Server does not allow anonymous binds for performing user search.
* Bind account used for searching is `cn=vault,ou=users,dc=example,dc=com` with password `My$ecrt3tP4ss`.
* User objects are under the `ou=Users,dc-example,dc=com` organizational unit.
* Username passed to vault when authenticating maps to the `sAMAccountName` attribute.
* Group membership will be resolved via the `memberOf` attribute of _user_ objects. That search will begin under `ou=Users,dc=example,dc=com`.

```
$ vault write auth/ldap/config \
    url="ldap://ldap.example.com" \
    userattr=sAMAccountName \
    userdn="ou=Users,dc=example,dc=com" \
    groupdn="ou=Users,dc=example,dc=com" \
    groupfilter="(&(objectClass=person)(uid={{.Username}}))" \
    groupattr="memberOf" \
    binddn="cn=vault,ou=users,dc=example,dc=com" \
    bindpass='My$ecrt3tP4ss' \
    certificate=@ldap_ca_cert.pem \
    insecure_tls=false \
    starttls=true
...
```

### Scenario 3

* LDAP server running on `ldap.example.com`, port 636 (LDAPS)
* CA Certificate stored in file named `ldap_ca_cert.pem`
* User objects are under the `ou=Users,dc=example,dc=com` organizational unit.
* Username passed to vault when authenticating maps to the `uid` attribute.
* User bind DN will be auto-discovered using anonymous binding.
* Group membership will be resolved via any one of `memberUid`, `member`, or `uniqueMember` attributes. That search will begin under `ou=Groups,dc=example,dc=com`.
* Group names are identified using the `cn` attribute.

```
$ vault write auth/ldap/config \
    url="ldaps://ldap.example.com" \
    userattr="uid" \
    userdn="ou=Users,dc=example,dc=com" \
    discoverdn=true \
    groupdn="ou=Groups,dc=example,dc=com" \
    certificate=@ldap_ca_cert.pem \
    insecure_tls=false \
    starttls=true
...
```

## LDAP Group -> Policy Mapping

Next we want to create a mapping from an LDAP group to a Vault policy:

```
$ vault write auth/ldap/groups/scientists policies=foo,bar
```

This maps the LDAP group "scientists" to the "foo" and "bar" Vault policies.
We can also add specific LDAP users to additional (potentially non-LDAP)
groups. Note that policies can also be specified on LDAP users as well.

```
$ vault write auth/ldap/groups/engineers policies=foobar
$ vault write auth/ldap/users/tesla groups=engineers policies=zoobar
```

This adds the LDAP user "tesla" to the "engineers" group, which maps to
the "foobar" Vault policy. User "tesla" itself is associated with "zoobar"
policy.

Finally, we can test this by authenticating:

```
$ vault auth -method=ldap username=tesla
Password (will be hidden):
Successfully authenticated! The policies that are associated
with this token are listed below:

default, foobar, zoobar
```

## Note on policy mapping

It should be noted that user -> policy mapping happens at token creation time. And changes in group membership on the LDAP server will not affect tokens that have already been provisioned. To see these changes, old tokens should be revoked and the user should be asked to reauthenticate.

## API
### /auth/ldap/config
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Configures the LDAP authentication backend.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/config`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">url</span>
        <span class="param-flags">required</span>
        The LDAP server to connect to. Examples: `ldap://ldap.myorg.com`,
        `ldaps://ldap.myorg.com:636`
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">starttls</span>
        <span class="param-flags">optional</span>
        If true, issues a `StartTLS` command after establishing an unencrypted
        connection. Defaults to `false`.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">tls_min_version</span>
        <span class="param-flags">optional</span>
        Minimum TLS version to use. Accepted values are `tls10`, `tls11` or
        `tls12`. Defaults to `tls12`.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">tls_max_version</span>
        <span class="param-flags">optional</span>
        Maximum TLS version to use. Accepted values are `tls10`, `tls11` or
        `tls12`. Defaults to `tls12`.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">insecure_tls</span>
        <span class="param-flags">optional</span>
        If true, skips LDAP server SSL certificate verification - insecure, use
        with caution! Defaults to `false`.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">certificate</span>
        <span class="param-flags">optional</span>
        CA certificate to use when verifying LDAP server certificate, must be
        x509 PEM encoded.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">binddn</span>
        <span class="param-flags">optional</span>
        Distinguished name of object to bind when performing user search.
        Example: `cn=vault,ou=Users,dc=example,dc=com`
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bindpass</span>
        <span class="param-flags">optional</span>
        Password to use along with `binddn` when performing user search.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">userdn</span>
        <span class="param-flags">optional</span>
        Base DN under which to perform user search. Example:
        `ou=Users,dc=example,dc=com`
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">userattr</span>
        <span class="param-flags">optional</span>
        Attribute on user attribute object matching the username passed when
        authenticating. Examples: `sAMAccountName`, `cn`, `uid`
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">discoverdn</span>
        <span class="param-flags">optional</span>
        Use anonymous bind to discover the bind DN of a user. Defaults to
        `false`.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">deny_null_bind</span>
        <span class="param-flags">optional</span>
        This option prevents users from bypassing authentication when providing
        an empty password. Defaults to `true`.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">upndomain</span>
        <span class="param-flags">optional</span>
        userPrincipalDomain used to construct the UPN string for the
        authenticating user. The constructed UPN will appear as
        `[username]@UPNDomain`. Example: `example.com`, which will cause
        vault to bind as `username@example.com`.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">groupfilter</span>
        <span class="param-flags">optional</span>
        Go template used when constructing the group membership query. The
        template can access the following context variables:
        \[`UserDN`, `Username`\]. The default is `(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))`,
        which is compatible with several common directory schemas. To support
        nested group resolution for Active Directory, instead use the following
        query: `(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))`.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">groupdn</span>
        <span class="param-flags">optional</span>
        LDAP search base to use for group membership search. This can be the
        root containing either groups or users.
        Example: `ou=Groups,dc=example,dc=com`
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">groupattr</span>
        <span class="param-flags">optional</span>
        LDAP attribute to follow on objects returned by `groupfilter` in order
        to enumerate user group membership. Examples: for groupfilter queries
        returning _group_ objects, use: `cn`. For queries returning _user_
        objects, use: `memberOf`. The default is `cn`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
  Retrieves the LDAP configuration.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/config`</dd>

  <dt>Parameters</dt>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": null,
      "warnings": null,
      "wrap_info": null,
      "data": {
        "binddn": "cn=vault,ou=Users,dc=example,dc=com",
        "bindpass": "",
        "certificate": "",
        "deny_null_bind": true,
        "discoverdn": false,
        "groupattr": "cn",
        "groupdn": "ou=Groups,dc=example,dc=com",
        "groupfilter": "(\u0026(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))",
        "insecure_tls": false,
        "starttls": false,
        "tls_max_version": "tls12",
        "tls_min_version": "tls12",
        "upndomain": "",
        "url": "ldaps://ldap.myorg.com:636",
        "userattr": "samaccountname",
        "userdn": "ou=Users,dc=example,dc=com"
      },
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>

### /auth/ldap/groups
#### LIST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Lists the existing groups in the backend.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/groups` (LIST) or `/auth/ldap/groups?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
  None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": null,
      "warnings": null,
      "wrap_info": null,
      "data": {
        "keys": [
          "scientists",
          "engineers"
        ]
      },
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>

### /auth/ldap/groups/[group_name]
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Creates and updates the LDAP group policy associations.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/groups/[group_name]`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">required</span>
        Comma-separated list of policies associated to the group.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
  Reads the LDAP group policy mappings.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/groups/[group_name]`</dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "policies": "admin,default"
      },
      "renewable": false,
      "lease_id": ""
      "lease_duration": 0,
      "warnings": null
    }
    ```

  </dd>
</dl>

#### DELETE
<dl class="api">
  <dt>Description</dt>
  <dd>
  Deletes an LDAP group.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/groups/[group_name]`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/ldap/users
#### LIST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Lists the existing users in the backend.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/users` (LIST) or `/auth/ldap/users?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
  None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": null,
      "warnings": null,
      "wrap_info": null,
      "data": {
        "keys": [
          "tesla"
        ]
      },
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>

### /auth/ldap/users/[username]
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Creates and updates the LDAP user group and policy mappings.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/users/[username]`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">groups</span>
        <span class="param-flags">optional</span>
        Comma-separated list of groups associated to the user.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
        Comma-separated list of policies associated to the user.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
  Reads the LDAP user.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/users/[username]`</dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "policies": "admins,default",
        "groups": ""
      },
      "renewable": false,
      "lease_id": ""
      "lease_duration": 0,
      "warnings": null
    }
    ```

  </dd>
</dl>

#### DELETE
<dl class="api">
  <dt>Description</dt>
  <dd>
  Deletes an LDAP user.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/users/[username]`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/ldap/login/[username]
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Creates and updates the LDAP user group and policy associations.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/ldap/login/[username]`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">password</span>
        <span class="param-flags">required</span>
        Password for the user.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": null,
      "auth": {
        "client_token": "c4f280f6-fdb2-18eb-89d3-589e2e834cdb",
        "policies": [
          "admins",
          "default"
        ],
        "metadata": {
          "username": "mitchellh"
        },
        "lease_duration": 0,
        "renewable": false
      }
    }
    ```

  </dd>
</dl>
