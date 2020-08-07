# Topology Package Design
This document outlines the design for this package.

## Topology
The `Topology` type handles monitoring the state of a MongoDB deployment and selecting servers.
Updating the description is handled by finite state machine which implements the server discovery
and monitoring specification. A `Topology` can be connected and fully disconnected, which enables
saving resources. The `Topology` type also handles server selection following the server selection
specification.

## Server
The `Server` type handles heartbeating a MongoDB server and holds a pool of connections.

## Connection
Connections are handled by two main types and an auxiliary type. The two main types are `connection`
and `Connection`. The first holds most of the logic required to actually read and write wire
messages. Instances can be created with the `newConnection` method. Inside the `newConnection`
method the auxiliary type, `initConnection` is used to perform the connection handshake. This is
required because the `connection` type does not fully implement `driver.Connection` which is
required during handshaking. The `Connection` type is what is actually returned to a consumer of the
`topology` package. This type does implement the `driver.Connection` type, holds a reference to a
`connection` instance, and exists mainly to prevent accidental continued usage of a connection after
closing it.

The connection implementations in this package are conduits for wire messages but they have no
ability to encode, decode, or validate wire messages. That must be handled by consumers.

## Pool
The `pool` type implements a connection pool. It handles caching idle connections and dialing
new ones, but it does not track a maximum number of connections. That is the responsibility of a
wrapping type, like `Server`.

The `pool` type has no concept of closing, instead it has concepts of connecting and disconnecting.
This allows a `Topology` to be disconnected,but keeping the memory around to be reconnected later.
There is a `close` method, but this is used to close a connection.

There are three methods related to getting and putting connections: `get`, `close`, and `put`. The
`get` method will either retrieve a connection from the cache or it will dial a new `connection`.
The `close` method will close the underlying socket of a `connection`. The `put` method will put a
connection into the pool, placing it in the cahce if there is space, otherwise it will close it.
