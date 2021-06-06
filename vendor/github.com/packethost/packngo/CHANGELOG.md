# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).
This project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

This release contains a bunch of fixes to the package api after some more real
world use. There a few breaks in backwards compatibility, but we are tying to
minimize them and move towards a 1.0 release.

### Added
- "acceptance" tests which run against production api (will incur charges)
- HardwareReservation to Device
- RootPassword to Device
- Spot market support
- Management and Manageable fields to discern between Elastic IPs and device unique IP
- Support for Volume attachments to Device and Volume
- Support for ProvisionEvents
- DoRequest sugar to Client
- Add ListProject function to the SSHKeys interface
- Operations for switching between Network Modes, aka "L2 support"
  Support for Organization, Payment Method and Billing address resources

### Fixed
- User.Emails json tag is fixed to match api response
- Single error object api response is now handled correctly

### Changed
- IPService was split to DeviceIPService and ProjectIPService
- Renamed Device.IPXEScriptUrl -> Device.IPXEScriptURL
- Renamed DeviceCreateRequest.HostName -> DeviceCreateRequest.Hostname
- Renamed DeviceCreateRequest.IPXEScriptUrl -> DeviceCreateRequest.IPXEScriptURL
- Renamed DeviceUpdateRequest.HostName -> DeviceUpdateRequest.Hostname
- Renamed DeviceUpdateRequest.IPXEScriptUrl -> DeviceUpdateRequest.IPXEScriptURL
- Sync with packet.net api change to /projects/{id}/ips which no longer returns
  the address in CIDR form
- Removed package level exported functions that should have never existed

## [0.1.0] - 2017-08-17

Initial release, supports most of the api for interacting with:

- Plans
- Users
- Emails
- SSH Keys
- Devices
- Projects
- Facilities
- Operating Systems
- IP Reservations
- Volumes
