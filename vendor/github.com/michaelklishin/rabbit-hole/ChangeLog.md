## Changes Between 1.4.0 and 1.5.0 (unreleased)

### More Binding Management Functions

`ListExchangeBindings`, `ListExchangeBindingsWithSource`, `ListExchangeBindingsWithDestination`,
and `ListExchangeBindingsBetween` are new functions that list bindings,
in particular between exchanges.

GitHub issue: [#109](https://github.com/michaelklishin/rabbit-hole/pull/109).

### Password Hash Generation Helpers

It is now possible to specify a `password_hash` when creating a user.
Helper functions `GenerateSalt` and `SaltedPasswordHashSHA256` make this more
straightforward compared to implementing [the algorithm](http://www.rabbitmq.com/passwords.html#computing-password-hash)
directly.

GitHub issue: [#119](https://github.com/michaelklishin/rabbit-hole/pull/119)

### Paginated Queue Listing

A new function, `PagedListQueuesWithParameters`, can list queues with pagination support.

GitHub issue: [#118](https://github.com/michaelklishin/rabbit-hole/pull/118)

### More `NodeInfo` and `QueueInfo` Attributes

GitHub issue: [#115](https://github.com/michaelklishin/rabbit-hole/issues/115)

### URL.Opaque Left to Its Own Devices

The client no longer messes with `URL.Opaque` as it doesn't seem to
be necessary any more for correct %-encoding of URL path.

GitHub issue: [#121](https://github.com/michaelklishin/rabbit-hole/issues/121)


## Changes Between 1.0.0 and 1.1.0 (Dec 1st, 2015)

### More Complete Message Stats Information

Message stats now include fields such as `deliver_get` and `redeliver`.

GH issue: [#73](https://github.com/michaelklishin/rabbit-hole/pull/73).

Contributed by Edward Wilde.


## 1.0 (first tagged release, Dec 25th, 2015)

### TLS Support

`rabbithole.NewTLSClient` is a new function which works
much like `NewClient` but additionally accepts a transport.

Contributed by @[GrimTheReaper](https://github.com/GrimTheReaper).

### Federation Support

It is now possible to create federation links
over HTTP API.

Contributed by [Ryan Grenz](https://github.com/grenzr-bskyb).

### Core Operations Support

Most common HTTP API operations (listing and management of
vhosts, users, permissions, queues, exchanges, and bindings)
are supported by the client.
