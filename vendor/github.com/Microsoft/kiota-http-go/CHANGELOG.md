# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.5] - 2024-09-03

### Changed
- Fixed a bug in compression middleware which caused empty body to send on retries

## [1.4.4] - 2024-08-13

### Changed

- Added `http.request.resend_delay` as a span attribute for the retry handler
- Changed the `http.retry_count` span attribute to `http.request.resend_count` to conform to OpenTelemetry specs.

## [1.4.3] - 2024-07-22

### Changed

- Fixed a bug to prevent double request compression by the compression handler.

## [1.4.2] - 2024-07-16

### Changed

- Prevent compression if Content-Range header is present.
- Fix bug which leads to a missing Content-Length header.

## [1.4.1] - 2024-05-09

### Changed

- Allow custom response handlers to return nil result values.

## [1.4.0] - 2024-05-09

- Support retry after as a date.

## [1.3.3] - 2024-03-19

- Fix bug where overriding http.DefaultTransport with an implementation other than http.Transport would result in an interface conversion panic

### Changed

## [1.3.2] - 2024-02-28

### Changed

- Fix bug with headers inspection handler using wrong key.

## [1.3.1] - 2024-02-09

### Changed

- Fix bug that resulted in the error "content is empty" being returned instead of HTTP status information if the request returned no content and an unsuccessful status code.

## [1.3.0] - 2024-01-22

### Added

- Added support to override default middleware with function `GetDefaultMiddlewaresWithOptions`.

## [1.2.1] - 2023-01-22

### Changed

- Fix bug passing no timeout in client as 0 timeout in context  .

## [1.2.0] - 2024-01-22

### Added

- Adds support for XXX status code.

## [1.1.2] - 2024-01-20

### Changed

- Changed the code by replacing ioutil.ReadAll and ioutil.NopCloser with io.ReadAll and io.NopCloser, respectively, due to their deprecation.

## [1.1.1] - 2023-11-22

### Added

- Added response headers and status code to returned error in `throwIfFailedResponse`.

## [1.1.0] - 2023-08-11

### Added

- Added headers inspection middleware and option.

## [1.0.1] - 2023-07-19

### Changed

- Bug Fix: Update Host for Redirect URL in go client.

## [1.0.0] - 2023-05-04

### Changed

- GA Release.

## [0.17.0] - 2023-04-26

### Added

- Adds Response Headers to the ApiError returned on Api requests errors.

## [0.16.2] - 2023-04-17

### Added

- Exit retry handler earlier if context is done.
- Adds exported method `ReplacePathTokens` that can be used to process url replacement logic globally.

## [0.16.1] - 2023-03-20

### Added

- Context deadline for requests defaults to client timeout when not provided.

## [0.16.0] - 2023-03-01

### Added

- Adds ResponseStatusCode to the ApiError returned on Api requests errors.

## [0.15.0] - 2023-02-23

### Added

- Added UrlReplaceHandler that replaces segments of the URL.

## [0.14.0] - 2023-01-25

### Added

- Added implementation methods for backing store.

## [0.13.0] - 2023-01-10

### Added

- Added a method to convert abstract requests to native requests in the request adapter interface.

## [0.12.0] - 2023-01-05

### Added

- Added User Agent handler to add the library information as a product to the header.

## [0.11.0] - 2022-12-20

### Changed

- Fixed a bug where retry handling wouldn't rewind the request body before retrying.

## [0.10.0] - 2022-12-15

### Added

- Added support for multi-valued request headers.

### Changed

- Fixed http.request_content_length attribute name for tracing

## [0.9.0] - 2022-09-27

### Added

- Added support for tracing via OpenTelemetry.

## [0.8.1] - 2022-09-26

### Changed

- Fixed bug for http go where response handler was overwritten in context object.

## [0.8.0] - 2022-09-22

### Added

- Added support for constructing a proxy authenticated client.

## [0.7.2] - 2022-09-09

### Changed

- Updated reference to abstractions.

## [0.7.1] - 2022-09-07

### Added

- Added support for additional status codes.

## [0.7.0] - 2022-08-24

### Added

- Adds context param in send async methods

## [0.6.2] - 2022-08-30

### Added

- Default 100 secs timeout for all request with a default context.

## [0.6.1] - 2022-08-29

### Changed

- Fixed a bug where an error would be returned for a 201 response with described response.

## [0.6.0] - 2022-08-17

### Added

- Adds a chaos handler optional middleware for tests

## [0.5.2] - 2022-06-27

### Changed

- Fixed an issue where response error was ignored for Patch calls

## [0.5.1] - 2022-06-07

### Changed

- Updated abstractions and yaml dependencies.

## [0.5.0] - 2022-05-26

### Added

- Adds support for enum or enum collections responses

## [0.4.1] - 2022-05-19

### Changed

- Fixed a bug where CAE support would leak connections when retrying.

## [0.4.0] - 2022-05-18

### Added

- Adds support for continuous access evaluation.

## [0.3.0] - 2022-04-19

### Changed

- Upgraded to abstractions 0.4.0.
- Upgraded to go 18.

## [0.2.0] - 2022-04-08

### Added

- Added support for decoding special characters in query parameters names.

## [0.1.0] - 2022-03-30

### Added

- Initial tagged release of the library.
