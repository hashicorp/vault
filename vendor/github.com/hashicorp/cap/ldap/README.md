# ldap

[![Go Reference](https://pkg.go.dev/badge/github.com/hashicorp/cap/ldap.svg)](https://pkg.go.dev/github.com/hashicorp/cap/ldap)

ldap is a package for writing clients that authenticate using Active Directory
or LDAP.

Primary types provided by the package:

* `ldap.Client`
* `ldap.ClientConfig`

<hr>

## Examples

* [CLI example](examples/cli/) which implements an ldap
  user authentication CLI.  

An abbreviated example of authenticating a user:

```go
client, err := ldap.NewClient(ctx, &clientConfig)
if err != nil { 
  // handle error appropriately
}

// authenticate and get the user's groups as well.
result, err := client.Authenticate(ctx, username, passwd, ldap.WithGroups())
if err != nil { 
  // handle error appropriately
}

if result.Success {
  // user successfully authenticated...
  if len(result.Groups) > 0 {
    // we found some groups associated with the authenticated user...
  } 
}
```

## Configuration

`ldap.ClientConfig` provides connection details for your LDAP server,
information on how to authenticate users, and instructions on how to query for
group membership. The configuration options are categorized and detailed below.

### Connection parameters

* `URLS` ([]string, required) - The LDAP server to connect to. Examples:
  `ldap://ldap.myorg.com`, `ldaps://ldap.myorg.com:636`. If there's more than one
  URL configured, the directories will be tried in-order if there are errors
  during the connection process.
* `StartTLS` (bool, optional) - If true, issues a StartTLS command after
  establishing an unencrypted connection.
* `InsecureTLS` (bool, optional) - If true, skips LDAP server SSL certificate
  verification - insecure, use with caution!
* `Certificate` (string, optional) - CA certificate to use when verifying LDAP
  server certificate, must be x509 PEM encoded.
* `ClientTLSCert` (string, optional) - Client certificate to provide to the LDAP
  server, must be x509 PEM encoded.
* `ClientTLSKey` (string, optional) - Client certificate key to provide to the
  LDAP server, must be x509 PEM encoded.

### Binding parameters

There are two alternate methods of resolving the user object used to
authenticate the end user: *Search* or *User Principal Name*. When using
*Search*, the bind can be either anonymous or authenticated. *User Principal
Name* is a method of specifying users supported by Active Directory. More
information on UPN can be found
[here](https://docs.microsoft.com/en-us/windows/win32/ad/naming-properties?redirectedfrom=MSDN#userPrincipalName).

#### Binding - Authenticated Search

* `BindDN` (string, optional) - Distinguished name of object to bind when
  performing user and group search. Example: `cn=application-acct,ou=Users,dc=example,dc=com`
* `BindPassword` (string, optional) - Password to use along with binddn when
  performing user search.
* `UserDN` (string, optional) - Base DN under which to perform user search.
  Example: `ou=Users,dc=example,dc=com`
* `UserAttr`  (string, optional) - Attribute on user attribute object matching
  the username passed when authenticating.  Examples: "cn", "uid"
* `UserFilter` (string, optional) - Go template used to construct a ldap user
  search filter. The template can access the following context variables:
  [UserAttr, Username]. The default userfilter is
  `({{.UserAttr}}={{.Username}})` or
  `(userPrincipalName={{.Username}}@UPNDomain)` if the upndomain parameter is
  set. The user search filter can be used to  restrict what user can attempt to
  log in. For example, to limit login to users that are not contractors, you
  could write
  `(&(objectClass=user)({{.UserAttr}}={{.Username}})(!(employeeType=Contractor)))`.

#### Binding - Anonymous Search

* `DiscoverDN` (bool, optional) - If true, use anonymous bind to discover the bind DN of a user
* `UserDN` (string, optional) - Base DN under which to perform user search.
  Example: `ou=Users,dc=example,dc=com`
* `UserAttr` (string, optional) - Attribute on user attribute object matching the username passed when authenticating. Examples: `cn`, `uid`
* `UserFilter` (string, optional) - Go template used to construct a ldap user search filter. The template can access the following context variables: [UserAttr, Username]. The default UserFilter is `({{.UserAttr}}={{.Username}})` or `(userPrincipalName={{.Username}}@UPNDomain)` if the UPNDomain parameter is set. The user search filter can be used to restrict what user can attempt to log in. For example, to limit login to users that are not contractors, you could write `(&(objectClass=user)({{.UserAttr}}={{.Username}})(!(employeeType=Contractor)))`.
* `AllowEmptyPasswordBinds` (bool, optional) - This option prevents users from bypassing authentication when providing an empty password. The default is `false`.
* `AnonymousGroupSearch` (bool, optional) - Use anonymous binds when performing LDAP group searches. Defaults to `false`.

#### Binding - User Principal Name (AD)

* `UPNDomain` (string, optional) - userPrincipalDomain used to construct the UPN
  string for the authenticating user. The constructed UPN will appear as
  `[username]@UPNDomain`.  Example: `example.com`, which will result in binding as `username@example.com`.

#### Alias dereferencing

* `DerefAliases` (string, optional) - Will control how aliases are dereferenced
  when performing the search. Possible values are: `never`, `finding`,
  `searching`, and `always`. If unset, a default of `never` is used. When set to
  `finding`, it will only dereference aliases during name resolution of the
  base. When set to `searching`, it will dereference aliases after name
  resolution.

### Group Membership Resolution

Once a user has been authenticated, the LDAP auth method must know how to resolve which groups the user is a member of. The configuration for this can vary depending on your LDAP server and your directory schema. There are two main strategies when resolving group membership - the first is searching for the authenticated user object and following an attribute to groups it is a member of. The second is to search for group objects of which the authenticated user is a member of. Both methods are supported.

* `GroupFilter` (string, optional) - Go template used when constructing the group membership query. The template can access the following context variables: [UserDN, Username]. The default is `(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))`, which is compatible with several common directory schemas. To support nested group resolution for Active Directory, instead use the following query: `(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))`.
* `GroupDN` (string, required) - LDAP search base to use for group membership search. This can be the root containing either groups or users. Example: `ou=Groups,dc=example,dc=com`
* `GroupAttr` (string, optional) - LDAP attribute to follow on objects returned by GroupFilter in order to enumerate user group membership. Examples: for GroupFilter queries returning group objects, use: `cn`. For queries returning user objects, use: `memberOf`. The default is `cn`.
Note: When using Authenticated Search for binding parameters (see above) the
distinguished name defined for `BindDN` is used for the group search. Otherwise,
the authenticating user is used to perform the group search.

### User Attributes

Using configuration you can choose to optionally include an authenticated
user's DN and entry attributes in the results of an authentication request.  

* `IncludeUserAttributes` (bool, optional) - If true, specifies that the
  authenticating user's DN and attributes be included an authentication
  AuthResult. Note: the default password attribute for both openLDAP
  (userPassword) and AD (unicodePwd) will always be excluded.

* `ExcludeUserAttributes` ([]string, optional) - If specified, optionally
  defines a set of user attributes to be excluded when an authenticating user's
  attributes are included in an AuthResult.Note: the default password attribute
  for both openLDAP (userPassword) and AD (unicodePwd) will always be excluded.

### Other

* `MaximumPageSize` (int, optional) - If set to a value greater than 0, the LDAP backend will use the LDAP server's paged search control to request pages of up to the given size. This can be used to avoid hitting the LDAP server's maximum result size limit. Otherwise, the LDAP backend will not use the paged search control.
