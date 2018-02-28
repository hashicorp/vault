---
layout: "intro"
page_title: "Authentication - Getting Started"
sidebar_current: "gettingstarted-auth"
description: |-
  Authentication to Vault gives a user access to use Vault. Vault can authenticate using multiple methods.
---

# Authentication

Now that we know how to use the basics of Vault, it is important to understand
how to authenticate to Vault itself. Up to this point, we have not "logged in"
to Vault. When starting the Vault server in `dev` mode, it automatically logs
you in as the root user with admin permissions. In a non-dev setup, you would
have had to authenticate first.

On this page, we'll talk specifically about authentication. On the next page, we
talk about [authorization](/intro/getting-started/policies.html). Authentication
is the mechanism of assigning an identity to a Vault user. The access control
and permissions associated with an identity are authorization, and will not be
covered on this page.

Vault has pluggable auth methods, making it easy to authenticate with Vault
using whatever form works best for your organization. On this page we will use
the token auth method and the GitHub auth method.

## Background

Authentication is the process by which user or machine-supplied information is
verified and converted into a Vault token with matching policies attached. The
easiest way to think about Vault's authentication is to compare it to a website.

When a user authenticates to a website, they enter their username, password, and
maybe 2FA code. That information is verified against external sources (a
database most likely), and the website responds with a success or failure. On
success, the website also returns a signed cookie that contains a session id
which uniquely identifies that user for this session. That cookie and session id
are automatically carried by the browser to future requests so the user is
authenticated. Can you imagine how terrible it would be to require a user to
enter their login credentials on each page?

Vault behaves very similarly, but it is much more flexible and pluggable than a
standard website. Vault supports many different authentication mechanisms, but
they all funnel into a single "session token", which we call the "Vault token".

Authentication is simply the process by which a user or machine gets a Vault
token.

## Tokens

Token authentication is enabled by default in Vault and cannot be disabled. When
you start a dev server with `vault server -dev`, it outputs your _root token_.
The root token is the initial access token to configure Vault. It has root
privileges, so it can perform any operation within Vault.

You can create more tokens:

```text
$ vault token create
Key                Value
---                -----
token              463763ae-0c3b-ff77-e137-af668941465c
token_accessor     57b6b540-57c8-64c4-e9c6-0b18ab058144
token_duration     ∞
token_renewable    false
token_policies     [root]
```

By default, this will create a child token of your current token that inherits
all the same policies. The "child" concept here is important: tokens always have
a parent, and when that parent token is revoked, children can also be revoked
all in one operation. This makes it easy when removing access for a user, to
remove access for all sub-tokens that user created as well.

After a token is created, you can revoke it:

```text
$ vault token revoke 463763ae-0c3b-ff77-e137-af668941465c
Success! Revoked token (if it existed)
```

In a previous section, we use the `vault lease revoke` command. This command
is only used for revoking _leases_. For revoking _tokens_, use
`vault token revoke`.

To authenticate with a token:

```text
$ vault login d08e2bd5-ffb0-440d-6486-b8f650ec8c0c
Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                Value
---                -----
token              a402d075-6d59-6129-1ac7-3718796d4346
token_accessor     7636b2f8-0cf1-e110-9b18-8f8b5ecf8351
token_duration     ∞
token_renewable    false
token_policies     [root]
```

This authenticates with Vault. It will verify your token and let you know what
access policies the token is associated with. If you want to test the `vault
login` command, create a new token first.

### Best Practice

In practice, operators should not use the `token create` command to generate
Vault tokens for users or machines. Instead, those users or machines should
authenticate to Vault using any of Vault's configured auth methods such as
GitHub, LDAP, AppRole, etc. For legacy applications which cannot generate their
own token, operators may need to create a token in advance. Auth methods are
discussed in more detail in the next section.

## Auth Methods

Vault supports many auth methods, but they must be enabled before use. Auth
methods give you flexibility. Enabling and configuration are typically performed
by a Vault operator or security team. As an example of a human-focused auth
method, let's authenticate via GitHub.

First, enable the GitHub auth method:

```text
$ vault auth enable -path=github github
Success! Enabled github auth method at: github/
```

Just like secrets engines, auth methods default to their TYPE as the PATH, so
the following commands are equivalent:

```text
$ vault auth enable -path=github github
$ vault auth enable github
```

Unlike secrets engines which are enabled at the root router, auth methods are
always prefixed with `auth/` in their path. So the GitHub auth method we just
enabled is accessible at `auth/github`. As another example:

```text
$ vault auth enable -path=my-github github
Success! Enabled github auth method at: my-github/
```

This would make the GitHub auth method accessible at `auth/my-github`. You can
use `vault path-help` to learn more about the paths.

Next, configure the GitHub auth method. Each auth method has different
configuration options, so please see the documentation for the full details. In
this case, the minimal set of configuration is to map teams to policies.

```text
$ vault write auth/github/config organization=hashicorp
```

With the GitHub auth method enabled, we first have to configure it. For GitHub,
we tell it what organization users must be a part of, and map a team to a
policy:

```text
$ vault write auth/github/config organization=hashicorp
Success! Data written to: auth/github/config

$ vault write auth/github/map/teams/my-team value=default,my-policy
Success! Data written to: auth/github/map/teams/my-team
```

The first command configures Vault to pull authentication data from the
"hashicorp" organization on GitHub. The next command tells Vault to map any
users who are members of the team "my-team" (in the hashicorp organization) to
map to the policies "default" and "my-policy". These policies do not have to
exist in the system yet - Vault will just produce a warning when we login.

---

As a user, you may want to find which auth methods are enabled and available:

```text
$ vault auth list
Path       Type      Description
----       ----      -----------
github/    github    n/a
token/     token     token based credentials
```

The `vault auth list` command will list all enabled auth methods. To learn more
about how to authenticate to a particular auth method via the CLI, use the
`vault auth help` command with the PATH or TYPE of an auth method:

```text
$ vault auth help github
Usage: vault login -method=github [CONFIG K=V...]

  The GitHub auth method allows users to authenticate using a GitHub
  personal access token. Users can generate a personal access token from the
  settings page on their GitHub account.

  Authenticate using a GitHub token:

      $ vault login -method=github token=abcd1234

Configuration:

  mount=<string>
      Path where the GitHub credential method is mounted. This is usually
      provided via the -path flag in the "vault login" command, but it can be
      specified here as well. If specified here, it takes precedence over the
      value for -path. The default value is "github".

  token=<string>
      GitHub personal access token to use for authentication.
```

Similarly, you can ask for help information about any CLI auth method, _even if
it is not enabled_:

```text
$ vault auth help aws
$ vault auth help userpass
$ vault auth help token
```

As per the help output, authenticate to GitHub using the `vault login` command.
Enter your [GitHub personal access token][gh-pat] and Vault will authenticate
you.

```text
$ vault login -method=github
GitHub Personal Access Token (will be hidden):
Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                    Value
---                    -----
token                  7efb3969-8743-880f-e234-afca6e12d790
token_accessor         f7bfb6a3-c41e-eb87-5317-88a0aad200ae
token_duration         768h
token_renewable        true
token_policies         [default my-policy]
token_meta_org         hashicorp
token_meta_username    my-user
```

Success! As the output indicates, Vault has already saved the resulting token in
its token helper, so you do not need to run `vault login` again. However, this
new user we just created does not have many permissions in Vault. To continue,
re-authenticate as the root token:

```text
$ vault login <initial-root-token>
```

You can revoke any logins from an auth method using `vault token revoke` with
the `-mode` argument. For example:

```text
$ vault token revoke -mode path auth/github
```

Alternatively, if you want to complete disable the GitHub auth method:

```text
$ vault auth disable github
Success! Disabled the auth method (if it existed) at: github/
```

This will also revoke any logins for that auth method.

## Next

In this page you learned about how Vault authenticates users. You learned about
the built-in token system as well as enabling other auth methods. At this point
you know how Vault assigns an _identity_ to a user.

The multiple auth methods Vault provides let you choose the most appropriate
authentication mechanism for your organization.

In this next section, we'll learn about
[authorization and policies](/intro/getting-started/policies.html).

[gh-pat]: https://help.github.com/articles/creating-an-access-token-for-command-line-use/
