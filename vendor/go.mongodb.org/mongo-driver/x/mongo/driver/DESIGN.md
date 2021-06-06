# Driver Library Design
This document outlines the design for this package.

## Deployment, Server, and Connection
Acquiring a `Connection` from a `Server` selected from a `Deployment` enables sending and receiving
wire messages. A `Deployment` represents an set of MongoDB servers and a `Server` represents a
member of that set. These three types form the operation execution stack.

### Compression
Compression is handled by Connection type while uncompression is handled automatically by the
Operation type. This is done because the compressor to use for compressing a wire message is
chosen by the connection during handshake, while uncompression can be performed without this
information. This does make the design of compression non-symmetric, but it makes the design simpler
to implement and more consistent.

## Operation
The `Operation` type handles executing a series of commands using a `Deployment`. For most uses
`Operation` will only execute a single command, but the main use case for a series of commands is
batch split write commands, such as insert. The type itself is heavily documented, so reading the
code and comments together should provide an understanding of how the type works.

This type is not meant to be used directly by callers. Instead an wrapping type should be defined
using the IDL and an implementation generated using `operationgen`.
