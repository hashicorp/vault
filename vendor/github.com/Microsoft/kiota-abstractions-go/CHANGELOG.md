# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Added

## [1.7.0] - 2024-07-09

-  Added accessors for headers and status to `ApiErrorable`  [#177](https://github.com/microsoft/kiota-abstractions-go/issues/177)

### Changed

## [1.6.1] - 2024-07-09

- Corrected two instances of `octet-steam` to `octet-stream` [#173](https://github.com/microsoft/kiota-abstractions-go/pull/173), [#174](https://github.com/microsoft/kiota-abstractions-go/pull/174)

## [1.6.0] - 2024-02-29

### Added

- Added support for untyped nodes. (https://github.com/microsoft/kiota/pull/4095)

## [1.5.6] - 2024-01-18

### Changed

- The input contains http or https which function will return an error. [#130](https://github.com/microsoft/kiota-abstractions-go/issues/130)

## [1.5.5] - 2024-01-17

### Changed

- Fixed a bug where reseting properties to null would be impossible with the in memory backing store. [microsoftgraph/msgraph-sdk-go#643](https://github.com/microsoftgraph/msgraph-sdk-go/issues/643)

## [1.5.4] - 2024-01-16

### Changed

- Fix bug where empty string query parameters are added to the request. [#133](https://github.com/microsoft/kiota-abstractions-go/issues/133)

## [1.5.3] - 2023-11-24

### Added

- Added support for multi valued query and path parameters of type other than string. [#124](https://github.com/microsoft/kiota-abstractions-go/pull/124)

## [1.5.2] - 2023-11-22

### Added

- Added ApiErrorable interface. [microsoft/kiota-http-go#110](https://github.com/microsoft/kiota-http-go/issues/110)

## [1.5.1] - 2023-11-15

### Added

- Added support for query an path parameters of enum type. [microsoft/kiota#3693](https://github.com/microsoft/kiota/issues/3693)

## [1.5.0] - 2023-11-08

### Added

- Added request information methods to reduce the amount of generated code.

## [1.4.0] - 2023-11-01

### Added

- Added serialization helpers. [microsoft/kiota#3406](https://github.com/microsoft/kiota/issues/3406)

## [1.3.1] - 2023-10-31

### Changed

- Fixed an issue where query parameters of type array of anything else than string would not be expanded properly. [#114](https://github.com/microsoft/kiota-abstractions-go/issues/114)

## [1.3.0] - 2023-10-12

### Added

- Added an overload method to set binary content with their content type.

## [1.2.3] - 2023-10-05

### Added

- A tryAdd method to RequestHeaders

## [1.2.2] - 2023-09-21

### Changed

- Switched the RFC 6570 implementation to std-uritemplate

## [1.2.1] - 2023-09-06

### Changed

- Fixed a bug where serialization registries would always replace existing values. [#95](https://github.com/microsoft/kiota-abstractions-go/issues/95)

## [1.2.0] - 2023-07-26

### Added

- Added support for multipart request body.

## [1.1.0] - 2023-05-04

### Added

- Added an interface to represent composed types.

## [1.0.0] - 2023-05-04

### Changed

- GA Release.

## [0.20.0] - 2023-04-12

### Added

- Adds response headers to Api Error class

### Changed

## [0.19.1] - 2023-04-12

### Added

### Changed

- Fixes concurrent map write panics when enabling backing stores.

## [0.19.0] - 2023-03-22

### Added

- Adds base request builder class to reduce generated code duplication.

## [0.18.0] - 2023-03-20

### Added

- Adds utility functions `CopyMap` and `CopyStringMap` that returns a copy of the passed map.

## [0.17.3] - 2023-03-15

### Changed

- Fixes panic when updating in-memory slices, maps or structs .

## [0.17.2] - 2023-03-01

### Added

- Adds ResponseStatusCode field in ApiError struct.

## [0.17.1] - 2023-01-28

### Added

- Adds a type qualifier for backing store instance type to be `BackingStoreFactory`.

### Changed

## [0.17.0] - 2023-01-23

### Added

- Added support for backing store.

## [0.16.0] - 2023-01-10

### Added

- Added a method to convert abstract requests to native requests in the request adapter interface.

## [0.15.2] - 2023-01-09

### Changed

- Fix bug where empty string query parameters are added to the request.

## [0.15.1] - 2022-12-15

### Changed

- Fix bug preventing adding authentication key to header requests.

## [0.15.0] - 2022-12-15

### Added

- Added support for multi-valued request headers.

## [0.14.0] - 2022-10-28

### Changed

- Fixed a bug where request bodies collections with single elements would not serialize properly

## [0.13.0] - 2022-10-18

### Added

- Added an API key authentication provider.

## [0.12.0] - 2022-09-27

### Added

- Added tracing support through OpenTelemetry.

## [0.11.0] - 2022-09-22

### Add
- Adds generic helper methods to reduce code duplication for serializer and deserializers
- Adds `WriteAnyValue` to support serialization of objects with undetermined properties at execution time e.g maps.
- Adds `GetRawValue` to allow returning an `interface{}` from the parse-node

## [0.10.1] - 2022-09-14

### Changed

- Fix: Add getter and setter on `ResponseHandler` pointer .

## [0.10.0] - 2022-09-02

### Added

- Added support for composed types serialization.

## [0.9.1] - 2022-09-01

### Changed

- Add `ResponseHandler` to request information struct

## [0.9.0] - 2022-08-24

### Changed

- Changes RequestAdapter contract passing a `Context` object as the first parameter for SendAsync

## [0.8.2] - 2022-08-11

### Added

- Add tests to verify DateTime and DateTimeOffsets default to ISO 8601.
- Adds check to return error when the baseUrl path parameter is not set when needed.

## [0.8.1] - 2022-06-07

### Changed

- Updated yaml package version through testify dependency.

## [0.8.0] - 2022-05-26

### Added

- Adds support for enum and enum collections responses.

## [0.7.0] - 2022-05-18

### Changed

- Breaking: adds support for continuous access evaluation.

## [0.6.0] - 2022-05-16

- Added a method to set the content from a scalar value in request information.

## [0.5.0] - 2022-04-21

### Added

- Added vanity methods to request options to add headers and options to simplify code generation.

## [0.4.0] - 2022-04-19

### Changed

- Upgraded uri template library for quotes in template fix.
- Upgraded to Go 18

## [0.3.0] - 2022-04-08

### Added

- Added support for query parameters with special characters in the name.

## [0.2.0] - 2022-04-04

### Changed

- Breaking: simplifies the field deserializers.

## [0.1.0] - 2022-03-30

### Added

- Initial tagged release of the library.
