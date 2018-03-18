---
layout: "intro"
page_title: "Your First Secret - Getting Started"
sidebar_current: "gettingstarted-firstsecret"
description: |-
  With the Vault server running, let's read and write our first secret.
---

# Your First Secret

Now that the dev server is up and running, let's get straight to it and
read and write our first secret.

One of the core features of Vault is the ability to read and write
arbitrary secrets securely. On this page, we'll do this using the CLI,
but there is also a complete
[HTTP API](/api/index.html)
that can be used to programmatically do anything with Vault.

Secrets written to Vault are encrypted and then written to backend
storage. For our dev server, backend storage is in-memory, but in production
this would more likely be on disk or in [Consul](https://www.consul.io).
Vault encrypts the value before it is ever handed to the storage driver.
The backend storage mechanism _never_ sees the unencrypted value and doesn't
have the means necessary to decrypt it without Vault.

## Writing a Secret

Let's start by writing a secret. Before we do that though, we need to authenticate with the Vault server. We will be using the root token the server provided us in the previous step:

```text
$ vault login
Token (will be hidden):
```

Upon successful authentication, you should see something like the following:

```text
Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                Value
---                -----
token              deb4c4bc-fed3-cea5-8cac-bc1fba790ea5
token_accessor     cb8f39f7-8eb7-06a6-a738-3768246b41ed
token_duration     ∞
token_renewable    false
token_policies     [root]
```

We should now be ready to writing our first secret by using the `vault write` command as shown below:

```text
$ vault write secret/hello value=world
Success! Data written to: secret/hello
```

This writes the pair `value=world` to the path `secret/hello`. We'll
cover paths in more detail later, but for now it is important that the
path is prefixed with `secret/`, otherwise this example won't work. The
`secret/` prefix is where arbitrary secrets can be read and written.

You can even write multiple pieces of data, if you want:

```text
$ vault write secret/hello value=world excited=yes
Success! Data written to: secret/hello
```

`vault write` is a very powerful command. In addition to writing data
directly from the command-line, it can read values and key pairs from
`STDIN` as well as files. For more information, see the
[command documentation](/docs/commands/index.html).

~> **Warning:** The documentation uses the `key=value` based entry
throughout, but it is more secure to use files if possible. Sending
data via the CLI is often logged in shell history. For real secrets,
please use files. See the link above about reading in from `STDIN` for more information.

## Reading a Secret

As you might expect, secrets can be read with `vault read`:

```text
$ vault read secret/hello
Key                 Value
---                 -----
refresh_interval    768h
excited             yes
value               world
```

As you can see, the values we wrote are given back to us. Vault reads
the data from storage and decrypts it.

The output format is purposefully whitespace separated to make it easy
to pipe into a tool like `awk`.

This contains some extra information. Many secrets engines create leases for
secrets that allow time-limited access to other systems, and in those cases
`lease_id` would contain a lease identifier and `lease_duration` would contain
the length of time for which the lease is valid, in seconds.

Optional JSON output is very useful for scripts. For example below we use the
`jq` tool to extract the value of the `excited` secret:

```text
$ vault read -format=json secret/hello | jq -r .data.excited
yes
```

When supported, you can also get a field directly:

```text
$ vault read -field=excited secret/hello
yes
```

## Deleting a Secret

Now that we've learned how to read and write a secret, let's go ahead
and delete it. We can do this with `vault delete`:

```text
$ vault delete secret/hello
Success! Data deleted (if it existed) at: secret/hello
```

## Next

In this section we learned how to use the powerful CRUD features of
Vault to store arbitrary secrets. On its own this is already a useful
but basic feature.

Next, we'll learn the basics about [secrets engines](/intro/getting-started/secrets-engines.html).
