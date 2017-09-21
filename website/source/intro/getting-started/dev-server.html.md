---
layout: "intro"
page_title: "Starting the Server - Getting Started"
sidebar_current: "gettingstarted-devserver"
description: |-
  After installing Vault, the next step is to start the server.
---

# Starting the Vault Server

With Vault installed, the next step is to start a Vault server.

Vault operates as a client/server application. The Vault server is the only
piece of the Vault architecture that interacts with the data storage and
backends. All operations done via the Vault CLI interact with the server over a
TLS connection.

In this page, we'll start and interact with the Vault server to understand how
the server is started.

## Starting the Dev Server

First, we're going to start a Vault _dev server_. The dev server is a built-in,
pre-configured server that is not very secure but useful for playing with Vault
locally. Later in this guide we'll configure and start a real server.

To start the Vault dev server, run:

```text
$ vault server -dev
==> Vault server configuration:

                     Cgo: disabled
         Cluster Address: https://127.0.0.1:8201
              Listener 1: tcp (addr: "127.0.0.1:8200", cluster address: "127.0.0.1:8201", tls: "disabled")
               Log Level: info
                   Mlock: supported: false, enabled: false
        Redirect Address: http://127.0.0.1:8200
                 Storage: inmem
                 Version: Vault v1.2.3
             Version Sha: ...

WARNING! dev mode is enabled! In this mode, Vault runs entirely in-memory
and starts unsealed with a single unseal key. The root token is already
authenticated to the CLI, so you can immediately begin using Vault.

You may need to set the following environment variable:

    $ export VAULT_ADDR='http://127.0.0.1:8200'

The unseal key and initial root token are displayed below in case you want to
seal/unseal the Vault or re-authenticate.

Unseal Key: 1aKM7rNnyW+7Jx1XDAXFswgkRVe+78JB28k/bel90jY=
Root Token: root

Development mode should NOT be used in production installations!

==> Vault server started! Log data will stream in below:

# ...
```

You should see output similar to that above. Vault does not fork, so it will
continue to run in the foreground. Open another shell or terminal tab to run the
remaining commands.

The dev server stores all its data in-memory (but still encrypted), listens on
`localhost` without TLS, and automatically unseals and shows you the unseal key
and root access key. **Do not run a dev server in production!**

With the dev server running, do the following three things before anything else:

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

```text
$ vault status
Key             Value
---             -----
Sealed          false
Total Shares    1
Version         (version unknown)
Cluster Name    vault-cluster-81109a1a
Cluster ID      f6e0aa8a-700e-38b8-5dc5-4265c880b2a1
HA Enabled      false
```

If the output looks different, especially if the numbers are different or the
Vault is sealed, then restart the dev server and try again. The only reason
these would ever be different is if you're running a dev server from going
through this guide previously.

We'll cover what this output means later in the guide.

## Next

Congratulations! You've started your first Vault server. We haven't stored
any secrets yet, but we'll do that in the next section.

Next, we're going to
[read and write our first secrets](/intro/getting-started/first-secret.html).
