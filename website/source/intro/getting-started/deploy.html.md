---
layout: "intro"
page_title: "Deploy Vault - Getting Started"
sidebar_current: "gettingstarted-deploy"
description: |-
  Learn how to deploy Vault into production, how to initialize it, configure it, etc.
---

# Deploy Vault

Up to this point, we've been working with the dev server, which
automatically authenticated us, setup in-memory storage, etc. Now that
you know the basics of Vault, it is important to learn how to deploy
Vault into a real environment.

On this page, we'll cover how to configure Vault, start Vault, the
seal/unseal process, and scaling Vault.

## Configuring Vault

Vault is configured using [HCL](https://github.com/hashicorp/hcl) files.
As a reminder, JSON files are also fully HCL-compatible; HCL is a superset of JSON.
The configuration file for Vault is relatively simple. An example is shown below:

```javascript
storage "consul" {
  address = "127.0.0.1:8500"
  path = "vault"
}

listener "tcp" {
 address = "127.0.0.1:8200"
 tls_disable = 1
}
```

Within the configuration file, there are two primary configurations:

  * `storage` - This is the physical backend that Vault uses for storage. Up to
    this point the dev server has used "inmem" (in memory), but in the example
    above we're using [Consul](https://www.consul.io), a much more
    production-ready backend.

  * `listener` - One or more listeners determine how Vault listens for API
    requests. In the example above we're listening on localhost port 8200
    without TLS. In your environment set `VAULT_ADDR=http://127.0.0.1:8200` so
    the Vault client will connect without TLS.

For now, copy and paste the configuration above to a file called
`example.hcl`. It will configure Vault to expect an instance of Consul
running locally.

Starting a local Consul instance takes only a few minutes. Just follow the
[Consul Getting Started Guide](https://www.consul.io/intro/getting-started/install.html)
up to the point where you have installed Consul and started it with this command:

```shell
$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul -bind 127.0.0.1
```

## Starting the Server

With the configuration in place, starting the server is simple, as
shown below. Modify the `-config` flag to point to the proper path
where you saved the configuration above.

```
$ vault server -config=example.hcl
==> Vault server configuration:

         Log Level: info
           Storage: consul
        Listener 1: tcp (addr: "127.0.0.1:8200", tls: "disabled")

==> Vault server started! Log data will stream in below:
```

Vault outputs some information about its configuration, and then blocks.
This process should be run using a resource manager such as systemd or
upstart.

You'll notice that you can't execute any commands. We don't have any
auth information! When you first setup a Vault server, you have to start
by _initializing_ it.

On Linux, Vault may fail to start with the following error:

```shell
$ vault server -config=example.hcl
Error initializing core: Failed to lock memory: cannot allocate memory

This usually means that the mlock syscall is not available.
Vault uses mlock to prevent memory from being swapped to
disk. This requires root privileges as well as a machine
that supports mlock. Please enable mlock on your system or
disable Vault from using it. To disable Vault from using it,
set the `disable_mlock` configuration option in your configuration
file.
```

For guidance on dealing with this issue, see the discussion of
`disable_mlock` in [Server Configuration](/docs/configuration/index.html).

## Initializing the Vault

Initialization is the process of first configuring the Vault. This
only happens once when the server is started against a new backend that
has never been used with Vault before.

During initialization, the encryption keys are generated, unseal keys
are created, and the initial root token is setup. To initialize Vault
use `vault init`. This is an _unauthenticated_ request, but it only works
on brand new Vaults with no data:

```
$ vault init
Key 1: 427cd2c310be3b84fe69372e683a790e01
Key 2: 0e2b8f3555b42a232f7ace6fe0e68eaf02
Key 3: 37837e5559b322d0585a6e411614695403
Key 4: 8dd72fd7d1af254de5f82d1270fd87ab04
Key 5: b47fdeb7dda82dbe92d88d3c860f605005
Initial Root Token: eaf5cc32-b48f-7785-5c94-90b5ce300e9b

Vault initialized with 5 keys and a key threshold of 3!
...
```

Initialization outputs two incredibly important pieces of information:
the _unseal keys_ and the _initial root token_. This is the
**only time ever** that all of this data is known by Vault, and also the
only time that the unseal keys should ever be so close together.

For the purpose of this getting started guide, save all of these keys
somewhere, and continue. In a real deployment scenario, you would never
save these keys together. Instead, you would likely use Vault's PGP and
Keybase.io support to encrypt each of these keys with the users' PGP keys.
This prevents one single person from having all the unseal keys. Please
see the documentation on [using PGP, GPG, and Keybase](/docs/concepts/pgp-gpg-keybase.html)
for more information.

## Seal/Unseal

Every initialized Vault server starts in the _sealed_ state. From
the configuration, Vault can access the physical storage, but it can't
read any of it because it doesn't know how to decrypt it. The process
of teaching Vault how to decrypt the data is known as _unsealing_ the
Vault.

Unsealing has to happen every time Vault starts. It can be done via
the API and via the command line. To unseal the Vault, you
must have the _threshold_ number of unseal keys. In the output above,
notice that the "key threshold" is 3. This means that to unseal
the Vault, you need 3 of the 5 keys that were generated.

-> **Note:** Vault does not store any of the unseal key shards. Vault
uses an algorithm known as
[Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)
to split the master key into shards. Only with the threshold number of keys
can it be reconstructed and your data finally accessed.

Begin unsealing the Vault with `vault unseal`:

```
$ vault unseal
Key (will be hidden):
Sealed: true
Key Shares: 5
Key Threshold: 3
Unseal Progress: 1
```

After pasting in a valid key and confirming, you'll see that the Vault
is still sealed, but progress is made. Vault knows it has 1 key out of 3.
Due to the nature of the algorithm, Vault doesn't know if it has the
_correct_ key until the threshold is reached.

Also notice that the unseal process is stateful. You can go to another
computer, use `vault unseal`, and as long as it's pointing to the same server,
that other computer can continue the unseal process. This is incredibly
important to the design of the unseal process: multiple people with multiple
keys are required to unseal the Vault. The Vault can be unsealed from
multiple computers and the keys should never be together. A single malicious
operator does not have enough keys to be malicious.

Continue with `vault unseal` to complete unsealing the Vault. To unseal
the vault you must use three _different_ keys, the same key repeated
will not work. As you use keys, as long as they are correct, you should
soon see output like this:

```
$ vault unseal
Key (will be hidden):
Sealed: false
Key Shares: 5
Key Threshold: 3
Unseal Progress: 0
```

The `Sealed: false` means the Vault is unsealed!

Feel free to play around with entering invalid keys, keys in different
orders, etc. in order to understand the unseal process. It is very important.
Once you're ready to move on, use `vault auth` to authenticate with
the root token.

As a root user, you can reseal the Vault with `vault seal`. A single
operator is allowed to do this. This lets a single operator lock down
the Vault in an emergency without consulting other operators.

When the Vault is sealed again, it clears all of its state (including
the encryption key) from memory. The Vault is secure and locked down
from access.

## Next

You now know how to configure, initialize, and unseal/seal Vault.
This is the basic knowledge necessary to deploy Vault into a real
environment. Once the Vault is unsealed, you access it as you have
throughout this getting started guide (which worked with an unsealed Vault).

Next, we have a [short tutorial](/intro/getting-started/apis.html) on using the HTTP APIs to authenticate and access secrets.
