# Changelog
This file will contain all notable changes to this repository. Any new releases with changes will be documented here.

## 0.2.5 (May 7 2024)

### Changed
- Log untriggered action if provider isn't running

## 0.2.4 (January 23 2024)

### Changed
- Support for updating the SCADA provider configuration

## 0.2.3 (April 7 2023)
Minor maintenance release.

### Changed
- Change "connect requested" log from info to debug
- Updated `golang.org/x/net` to `v0.7.0` to address numerous CVEs
- Updated `golang.org/x/text` to `v0.3.8` to address numerous CVEs
- Updated `go.mongodb.org/mongo-driver` to `v1.5.1` to address `CVE-2021-20329`

## 0.2.2 (January 16 2023)
Minor maintenance release.

### Changed
- Ensure `New()` has proper semantics

## 0.2.1 (January 11 2023)
Minor maintenance release.

### Fixed
- Two of the accessors (`LastError` & `SessionStatus`) had a data race

### Changed
- Added data race detection to the github workflow

## 0.2.0 (October 28 2022)
### Changed
- Updated `net-rpc-msgpackrp` to `v2` to address a diamond dependency problem for projects importing hcp-scada-provider and an older version of hashicorp/go-msgpack

## 0.1.1 (September 22 2022)
Minor update mainly for the new scada hostname.

### Changed
- Updated to hcp-sdk-go v0.23.0 to get the new scada public address

### Added
- Documented the `SessionStatus` function with better descriptions of the possible statuses
- Support for the `ErrInvalidCredentials` error in the `LastError` function

## 0.1.0 (September 9 2022)
First version of the library containing capability to connect to SCADA broker as SCADA provider using HCP identity.
