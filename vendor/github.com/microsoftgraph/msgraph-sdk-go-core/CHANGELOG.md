# Changelog

All notable changes to this project will be documented in this file.

## [1.2.1](https://github.com/microsoftgraph/msgraph-sdk-go-core/compare/v1.2.0...v1.2.1) (2024-08-26)


### Bug Fixes

* repeated slice uploading on large file upload task ([cb329cc](https://github.com/microsoftgraph/msgraph-sdk-go-core/commit/cb329cc395946a619cda5501da88dcda15d84d9b))

## [1.2.0](https://github.com/microsoftgraph/msgraph-sdk-go-core/compare/v1.1.0...v1.2.0) (2024-07-15)


### Features

* add git release config ([69234a2](https://github.com/microsoftgraph/msgraph-sdk-go-core/commit/69234a236c1d212941e742593ce43d2a35a1212b))


### Bug Fixes

* allows registration of page iterator headers ([#309](https://github.com/microsoftgraph/msgraph-sdk-go-core/issues/309)) ([d4b0806](https://github.com/microsoftgraph/msgraph-sdk-go-core/commit/d4b0806dadcc3ccdf07a8eca8ca7b93150094d7f))
* content range order during upload ([#304](https://github.com/microsoftgraph/msgraph-sdk-go-core/issues/304)) ([f241e94](https://github.com/microsoftgraph/msgraph-sdk-go-core/commit/f241e947b28de38e8f7bc8c3d4eb6eb95b9afbdb))

## [1.1.0](https://github.com/microsoftgraph/msgraph-sdk-go-core/compare/v1.0.2...v1.1.0) (2024-07-10)


### Features

* add git release config ([69234a2](https://github.com/microsoftgraph/msgraph-sdk-go-core/commit/69234a236c1d212941e742593ce43d2a35a1212b))


### Bug Fixes

* content range order during upload ([#304](https://github.com/microsoftgraph/msgraph-sdk-go-core/issues/304)) ([f241e94](https://github.com/microsoftgraph/msgraph-sdk-go-core/commit/f241e947b28de38e8f7bc8c3d4eb6eb95b9afbdb))

## [1.1.0] - 2024-02-02

### Added

- Added support for large file uploads.

## [1.0.2] - 2023-12-01

### Changed

- Fixed a bug where GetBatchResponseById failed to deserialize error response bodies.

## [1.0.1] - 2023-11-24

### Changed

- Fixed a bug where page iterator would panic if it couldn't find the GetValue method on the collection.

## [1.0.0] - 2023-05-04

### Changed

- GA Release.

## [0.36.2] - 2023-05-01

### Added

- `PageIterator` exposes `odata.nextLink` and `odata.deltaLink` of most recent page.

## [0.36.1] - 2023-04-17

### Added

- Adds url token replacement to batch requests.

## [0.36.0] - 2023-03-27

### Added

- Adds `BatchRequestCollection` support.

## [0.35.0] - 2023-03-23

### Added

- `PageIterator` uses generics to define return type.

## [0.34.1] - 2023-03-06

### Changed

- Change `PageIterator` to use `GetValue` method instead of `value` field to access response.

## [0.34.0] - 2023-02-23

### Added

- Adds `UrlReplaceHandler` to default middleware.

## [0.33.1] - 2023-01-26

### Added

- Upgrade dependencies to support backing store.

## [0.33.0] - 2023-01-17

### Added

- Added authentication provider with Microsoft Graph defaults.

## [0.32.0] - 2023-01-11

### Changed

- Upgraded abstractions and http dependencies.

## [0.31.1] - 2022-12-15

### Changed

- Fixes path parameters missing when sending batch requests.
- Fixes appending items when sending batch requests.
- Fixes `Send` url when sending batch requests

## [0.31.0] - 2022-12-13

### Changed

- Updated references to core libraries for multi-valued request headers.

## [0.30.1] - 2022-10-21

### Changed

- Fix: Remove error swallowing in page iterator `fetchNextPage`.

## [0.30.0] - 2022-09-29

### Added

- Adds ability to batch requests.
- Adds tracing support via Open Telemetry.

## [0.29.0] - 2022-09-27

### Changed

- Updated dependencies for additional serialization methods.

## [0.28.1] - 2022-09-09

### Changed

- Updates references to kiota packages.

## [0.28.0] - 2022-08-24

### Changed

- Upgrade to library `kiota-abstraction` breaking change
- Introduces `context.Context` object to Page Iterator

## [0.27.0] - 2022-07-21

### Changed

- Fixes PageIterator to use updated nextLink property

### Changed

## [0.26.2] - 2022-06-12

### Changed

- Updated reference to kiota serialization json
- Updated reference to kiota http

## [0.26.1] - 2022-06-07

### Changed

- Updated references to kiota libraries and yaml dependencies.

## [0.26.0] - 2022-05-27

### Changed

- Updated references to kiota libraries to add support for enum and enum collections responses.

## [0.25.1] - 2022-05-25

### Changed

- Updated kiota http library reference.

## [0.25.0] - 2022-05-19

### Changed

- Upgraded kiota dependencies for preliminary continuous access evaluation support.

## [0.24.0] - 2022-04-28

### Changed

- Updated references to kiota libraries for request configuration revamp

## [0.23.0] - 2022-04-19

### Changed

- Upgraded kiota libraries to address quote in url template issue.
- Upgraded to go 18.

## [0.22.1] - 2022-04-14

### Changed

- Fixed an issue with date serialization in JSON.

## [0.22.0] - 2022-04-12

### Changed

- Updated references to kiota libraries for special character in parameter names support.
- Breaking: removed the odata parameter names handler.

## [0.21.0] - 2022-04-06

### Changed

- Updated reference to kiota libraries for deserialization simplification.

## [0.20.0] - 2022-03-31

### Changed

- Updated reference to kiota libraries that were moved to their own repository.

## [0.0.17] - 2022-03-30

### Added

- Added support for vendor specific content types
- Added support for 204 no content responses

### Changed

- Updated kiota libraries reference.

## [0.0.16] - 2022-03-21

### Changed

- Breaking: updates PageIterator to receive a RequestAdapter interface instead of GraphRequestAdapterBase concrete type
- Breaking: removed IsNil method from models

## [0.0.15] - 2022-03-15

### Changed

- Updated references to kiota libraries for new supported types (byte, unit8, ...)

## [0.0.14] - 2022-03-11

### Changed

- Publishes a version retraction for v0.11.0 that was wrongfully published and causes issues during upgrades

## [0.0.13] - 2022-03-04

### Changed

- Breaking: updates kiota dependencies for parsable interface split.

## [0.0.12] - 2022-03-03

### Changed

- Breaking: updates kiota dependencies to pass request information by reference and not by copy (request adapter, authentication provider).

## [0.0.11] - 2022-03-02

### Changed

- Breaking: updates kiota dependencies references to prepare for type discriminator support.

## [0.0.10] - 2022-02-28

### Changed

- Fixed a bug where http client configuration would impact the default client configuration for other usages.

## [0.0.9] - 2022-02-16

### Added

- Added support for deserializing error responses (will return error)

### Changed

- Fixed a bug where response body compression would send empty bodies

## [0.0.8] - 2022-02-08

### Added

- Added support for request body compression (gzip)
- Added support for response body decompression (gzip)

### Changed

- Fixes a bug where resuming the page iterator wouldn't work
- Fixes a bug where OData query parameters would be added twice in some cases

## [0.0.7] - 2022-02-03

### Changed

- Updated references to Kiota packages to fix a [bug where the access token would never be attached to the request](https://github.com/microsoft/kiota/pull/1116). 

## [0.0.6] - 2022-02-02

### Added

- Adds missing delta token for OData query parameters dollar sign injection.
- Adds PageIterator task

## [0.0.5] - 2021-12-02

### Changed

- Fixes a bug where the middleware pipeline would run only on the first request of the client/adapter/http client.

## [0.0.4] - 2021-12-01

### Changed

- Adds the missing github.com/microsoft/kiota/authentication/go/azure dependency

## [0.0.3] - 2021-11-30

### Changed

- Updated dependencies and switched to Go 17.

## [0.0.2] - 2021-11-08

### Changed

- Updated kiota abstractions and http to provide support for setting the base URL

## [0.0.1] - 2021-10-22

### Added

- Initial release
