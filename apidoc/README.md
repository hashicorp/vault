# apidoc

## Goals
1. Deliver an OpenAPI spec for the `/sys` endpoints.
2. Build a foundation to allow as much auto-generation of documentation as possible, both for /sys and other endpoints, including plugin backends.
3. Allow generation of different documentation formats. Don’t lock into OpenAPI.


## Design
* Generating docs from code was chosen. Code from docs was investigated but is much more difficult to fit into our existing implementations.
* The basic design is to build an executable that will `import` things to be documented (i.e. actual go dependencies), and the run this application to generate the documentation.
* The code currently lives in Vault under `/apidocs` for development convenience, but I think it should really live in its own repo, importing Vault and plugins alike.
* A simple document structure (Go constructs) is defined to represent the api spec. The structure is similar to OpenAPI but it doesn’t need to be. The idea is that all documentation will be parsed into this format and then rendered into one of the target formats.
* A `*framework.Backend` can be parsed into basic documentation automatically using defined operations, descriptions and help.
* Functions are available to define or extend documentation. This was needed for a number of endpoints that aren’t part of a mounted backend, and also for adding documentation not part of the current help system (e.g. example responses). This approach is hopefully just a stopgap. I’d like to see `*framework.Backend` or another structure be modified to accommodate richer documentation as opposed to having sidecar structures.
There are some paths that only have path help, for example.
