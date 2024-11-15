# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

## [1.1.0] - 2024-08-08

### Changed

- Continuous Access Evaluation is now enabled by default.

## [1.0.2] - 2024-01-19

### Changed

- Validates that provided valid hosts don't start with a scheme.

## [1.0.1] - 2023-10-13

### Changed

- Allow http on localhost.

## [1.0.0] - 2023-05-04

### Changed

- GA Release.

## [0.6.0] - 2023-01-17

### Changed

- Removes the Microsoft Graph specific default values.

## [0.5.0] - 2022-09-27

### Added

- Added tracing through OpenTelemetry.

## [0.4.1] - 2022-09-02

### Changed

- Upgraded abstractions and yaml dependencies.

### Changed

## [0.4.0] - 2022-08-31

### Changed

- Pass `context.Context` for on `GetAuthorizationToken` method.

## [0.3.1] - 2022-06-07

### Changed

- Upgraded abstractions and yaml dependencies.

## [0.3.0] - 2022-05-18

### Added

- Added preliminary work to support continuous access evaluation.

## [0.2.1] - 2022-04-19

### Changed

- Upgraded abstractions to 0.4.0.

## [0.2.0] - 2022-04-18

### Changed

- Bumped required go version to 1.18 as Azure Identity now requires it.

## [0.1.0] - 2022-03-30

### Added

- Initial tagged release of the library.
