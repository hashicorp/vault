## Unreleased

## 1.8.5 (Aug 19 2020)

- Added delegate_dataset support to instance creation

## 1.8.4 (May 11 2020)

- Fix panic when testing images without TRITON_TEST set [#186]

## 1.8.3 (May 6 2020)

- Add support for `brand`, `flexible_disk` and `disks` for packages [#182].

## 1.8.2 (May 5 2020)

- Fix panic when TRITON_TRACE_HTTP is set and resp is nil [#180].

## 1.8.1 (May 4 2020)

- Add ability to trace CloudAPI HTTP requests [#180]. You can set the
  TRITON_TRACE_HTTP environment variable to have the CloudAPI HTTP requests and
  responses be printed to stderr.

## 1.8.0 (April 27 2020)

- Update Triton acceptance tests to work with any Triton [#178]

## 1.7.2 (April 27 2020)

- Unable to connect to CloudAPI using wildcard certificates [#117]. You can now
  set the TRITON_SKIP_TLS_VERIFY environment variable to skip TLS checking.

## 1.7.1 (April 24 2020)

- Add support for volume tags [#175]

## 1.7.0 (June 26 2019)

- Expected instance to have Tags

## 1.6.1 (June 26 2019)

- compute/networks: support network objects for AddNIC [#169]

## 1.6.0 (June 24 2019)

- compute/networks: added support for network objects [#158]
- compute/instances: added instances().get support for deleted instances [#167]
- storage: added support for multipart upload [#160]
- storage: fixed directory list marker filtering [#156]

## 1.3.1 (April 27 2018)

- client: Fixing an issue where private Triton installations were marked as invalid DC [#152]

## 1.3.0 (April 17 2018)

- identity/roles: Add support for SetRoleTags [#112]
- Add support for Triton Service Groups endpoint [#148]

## 1.2.0 (March 20 2018)

- compute/instance: Instance Deletion status now included in the GET instance response [#138]

## 1.1.1 (March 13 2018)

- client: Adding the rbac user support to the SSHAgentSigner [BUG!]

## 1.1.0 (March 13 2018)

- client: Add support for Manta RBAC http signatures

## 1.0.0 (February 28 2018)

- client: Add support for querystring in client/ExecuteRequestRaw [#121]
- client: Introduce SetHeader for overriding API request header [#125]
- compute/instances: Add support for passing a list of tags to filter List instances [#116]
- compute/instances: Add support for getting a count of current instances from the CloudAPI [#119]
- compute/instances: Add ability to support name-prefix [#129]
- compute/instances: Add support for Instance Deletion Protection [#131]
- identity/user: Add support for ChangeUserPassword [#111]
- expose GetTritonEnv as a root level func [#126]

## 0.9.0 (January 23 2018)

**Please Note:** This is a precursor release to marking triton-go as 1.0.0. We are going to wait and fix any bugs that occur from this large set of changes that has happened since 0.5.2

- Add support for managing volumes in Triton [#100]
- identity/policies: Add support for managing policies in Triton [#86]
- addition of triton-go errors package to expose unwrapping of internal errors
- Migration from hashicorp/errwrap to pkg/errors
- Using path.Join() for URL structures rather than fmt.Sprintf()

## 0.5.2 (December 28 2017)

- Standardise the API SSH Signers input casing and naming

## 0.5.1 (December 28 2017)

- Include leading '/' when working with SSH Agent signers

## 0.5.0 (December 28 2017)

- Add support for RBAC in triton-go [#82]
This is a breaking change. No longer do we pass individual parameters to the SSH Signer funcs, but we now pass an input Struct. This will guard from from additional parameter changes in the future. 
We also now add support for using `SDC_*` and `TRITON_*` env vars when working with the Default agent signer

## 0.4.2 (December 22 2017)

- Fixing a panic when the user loses network connectivity when making a GET request to instance [#81]

## 0.4.1 (December 15 2017)

- Clean up the handling of directory sanitization. Use abs paths everywhere [#79]

## 0.4.0 (December 15 2017)

- Fix an issue where Manta HEAD requests do not return an error resp body [#77]
- Add support for recursively creating child directories [#78]

## 0.3.0 (December 14 2017)

- Introduce CloudAPI's ListRulesMachines under networking
- Enable HTTP KeepAlives by default in the client.  15s idle timeout, 2x
  connections per host, total of 10x connections per client.
- Expose an optional Headers attribute to clients to allow them to customize
  HTTP headers when making Object requests.
- Fix a bug in Directory ListIndex [#69](https://github.com/joyent/issues/69)
- Inputs to Object inputs have been relaxed to `io.Reader` (formerly a
  `io.ReadSeeker`) [#73](https://github.com/joyent/issues/73).
- Add support for ForceDelete of all children of a directory [#71](https://github.com/joyent/issues/71)
- storage: Introduce `Objects.GetInfo` and `Objects.IsDir` using HEAD requests [#74](https://github.com/joyent/triton-go/issues/74)

## 0.2.1 (November 8 2017)

- Fixing a bug where CreateUser and UpdateUser didn't return the UserID

## 0.2.0 (November 7 2017)

- Introduce CloudAPI's Ping under compute
- Introduce CloudAPI's RebootMachine under compute instances
- Introduce CloudAPI's ListUsers, GetUser, CreateUser, UpdateUser and DeleteUser under identity package
- Introduce CloudAPI's ListMachineSnapshots, GetMachineSnapshot, CreateSnapshot, DeleteMachineSnapshot and StartMachineFromSnapshot under compute package
- tools: Introduce unit testing and scripts for linting, etc.
- bug: Fix the `compute.ListMachineRules` endpoint

## 0.1.0 (November 2 2017)

- Initial release of a versioned SDK
