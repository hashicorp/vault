---
layout: "intro"
page_title: "Starting the Server - Getting Started"
sidebar_current: "gettingstarted-devserver"
description: |-
  After installing Vault, the next step is to start the server.
---

# Starting the Vault Server

With Vault installed, the next step is to start a Vault server.

Vault operates as a client/server application. The Vault server is the
only piece of the Vault architecture that interacts with the data
storage and backends. All operations done via the Vault CLI interact
with the server over a TLS connection.

In this page, we'll start and interact with the Vault server to understand
how the server is started.

## Starting the Dev Server

First, we're going to start a Vault _dev server_. The dev server
is a built-in, pre-configured server that is not very
secure but useful for playing with Vault locally. Later in this guide
we'll configure and start a real server.

To start the Vault dev server, run:

```
$ vault server -dev
WARNING: Dev mode is enabled!

In this mode, Vault is completely in-memory and unsealed.
Vault is configured to only have a single unseal key. The root
token has already been authenticated with the CLI, so you can
immediately begin using the Vault CLI.

The only step you need to take is to set the following
environment variable since Vault will be talking without TLS:

    export VAULT_ADDR='http://127.0.0.1:8200'

The unseal key and root token are reproduced below in case you
want to seal/unseal the Vault or play with authentication.

Unseal Key: 2252546b1a8551e8411502501719c4b3
Root Token: 79bd8011-af5a-f147-557e-c58be4fedf6c

==> Vault server configuration:

         Log Level: info
           Backend: inmem
        Listener 1: tcp (addr: "127.0.0.1:8200", tls: "disabled")

...
```

You should see output similar to that above. Vault does not fork, so it will
continue to run in the foreground; to connect to it with later commands, open
another shell.

As you can see, when you start a dev server, Vault warns you loudly. The dev
server stores all its data in-memory (but still encrypted), listens on
`localhost` without TLS, and automatically unseals and shows you the unseal key
and root access key.  We'll go over what all this means shortly.

The important thing about the dev server is that it is meant for
development only.

-> **Note:** Do not run the dev server in production.

Even if the dev server was run in production, it wouldn't be very useful
since it stores data in-memory and every restart would clear all your
secrets.

With the dev server running, do the following three things before anything
else:

  1. Launch a new terminal session.

  2. Copy and run the `export VAULT_ADDR ...` command from the terminal
     output. This will configure the Vault client to talk to our dev server.

  3. Save the unseal key somewhere. Don't worry about _how_ to save this
     securely. For now, just save it anywhere.

  4. Do the same as step 3, but with the root token. We'll use this later.

## Verify the Server is Running

Verify the server is running by running the `vault status` command. This should
succeed and exit with exit code 0. If you see an error about opening
a connection, make sure you copied and executed the `export VAULT_ADDR...`
command from above properly.

If it ran successfully, the output should look like the below:

```
$ vault status
Sealed: false
Key Shares: 1
Key Threshold: 1
Unseal Progress: 0

High-Availability Enabled: false
```

If the output looks different, especially if the numbers are different
or the Vault is sealed, then restart the dev server and try again. The
only reason these would ever be different is if you're running a dev
server from going through this guide previously.

We'll cover what this output means later in the guide.

## Next

Congratulations! You've started your first Vault server. We haven't stored
any secrets yet, but we'll do that in the next section.

Next, we're going to
[read and write our first secrets](/intro/getting-started/first-secret.html).
