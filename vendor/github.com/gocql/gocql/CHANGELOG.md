# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Fixed

## [1.7.0] - 2024-09-23

This release is the first after the donation of gocql to the Apache Software Foundation (ASF)

### Changed
- Update DRIVER_NAME parameter in STARTUP messages to a different value intended to clearly identify this
  driver as an ASF driver.  This should clearly distinguish this release (and future gocql-cassandra-driver
  releases) from prior versions. (#1824)
- Supported Go versions updated to 1.23 and 1.22 to conform to gocql's sunset model. (#1825)

## [1.6.0] - 2023-08-28

### Added
- Added the InstaclustrPasswordAuthenticator to the list of default approved authenticators. (#1711)
- Added the `com.scylladb.auth.SaslauthdAuthenticator` and `com.scylladb.auth.TransitionalAuthenticator`
  to the list of default approved authenticators. (#1712)
- Added transferring Keyspace and Table names to the Query from the prepared response and updating
  information about that every time this information is received. (#1714)

### Changed
- Tracer created with NewTraceWriter now includes the thread information from trace events in the output. (#1716)
- Increased default timeouts so that they are higher than Cassandra default timeouts.
  This should help prevent issues where a default configuration overloads a server using default timeouts
  during retries. (#1701, #1719)

## [1.5.2] - 2023-06-12

Same as 1.5.0. GitHub does not like gpg signed text in the tag message (even with prefixed armor),
so pushing a new tag.

## [1.5.1] - 2023-06-12

Same as 1.5.0. GitHub does not like gpg signed text in the tag message,
so pushing a new tag.

## [1.5.0] - 2023-06-12

### Added

- gocql now advertises the driver name and version in the STARTUP message to the server.
  The values are taken from the Go module's path and version
  (or from the replacement module, if used). (#1702)
  That allows the server to track which fork of the driver is being used.
- Query.Values() to retrieve the values bound to the Query.
  This makes writing wrappers around Query easier. (#1700)

### Fixed
- Potential panic on deserialization (#1695)
- Unmarshalling of dates outside of `[1677-09-22, 2262-04-11]` range. (#1692)

## [1.4.0] - 2023-04-26

### Added

### Changed

- gocql now refreshes the entire ring when it receives a topology change event and
  when control connection is re-connected.
  This simplifies code managing ring state. (#1680)
- Supported versions of Cassandra that we test against are now 4.0.x and 4.1.x. (#1685)
- Default HostDialer now uses already-resolved connect address instead of hostname when establishing TCP connections (#1683).

### Fixed

- Deadlock in Session.Close(). (#1688)
- Race between Query.Release() and speculative executions (#1684)
- Missed ring update during control connection reconnection (#1680)

## [1.3.2] - 2023-03-27

### Changed

- Supported versions of Go that we test against are now Go 1.19 and Go 1.20.

### Fixed

- Node event handling now processes topology events before status events.
  This fixes some cases where new nodes were missed. (#1682)
- Learning a new IP address for an existing node (identified by host ID) now triggers replacement of that host.
  This fixes some Kubernetes reconnection failures. (#1682)
- Refresh ring when processing a node UP event for an unknown host.
  This fixes some cases where new nodes were missed. (#1669)

## [1.3.1] - 2022-12-13

### Fixed

- Panic in RackAwareRoundRobinPolicy caused by wrong alignment on 32-bit platforms. (#1666)

## [1.3.0] - 2022-11-29

### Added

- Added a RackAwareRoundRobinPolicy that attempts to keep client->server traffic in the same rack when possible.

### Changed

- Supported versions of Go that we test against are now Go 1.18 and Go 1.19.

## [1.2.1] - 2022-09-02

### Changed

- GetCustomPayload now returns nil instead of panicking in case of query error. (#1385)

### Fixed

- Nil pointer dereference in events.go when handling node removal. (#1652)
- Reading peers from DataStax Enterprise clusters. This was a regression in 1.2.0. (#1646)
- Unmarshaling maps did not pre-allocate the map. (#1642)

## [1.2.0] - 2022-07-07

This release improves support for connecting through proxies and some improvements when using Cassandra 4.0 or later.

### Added
- HostDialer interface now allows customizing connection including TLS setup per host. (#1629)

### Changed
- The driver now uses `host_id` instead of connect address to identify nodes. (#1632)
- gocql reads `system.peers_v2` instead of `system.peers` when connected to Cassandra 4.0 or later and
  populates `HostInfo.Port` using the native port. (#1635)

### Fixed
- Data race in `HostInfo.HostnameAndPort()`. (#1631)
- Handling of nils when marshaling/unmarshaling lists and maps. (#1630)
- Silent data corruption in case a map was serialized into UDT and some fields in the UDT were not present in the map.
  The driver now correctly writes nulls instead of shifting fields. (#1626, #1639)

## [1.1.0] - 2022-04-29

### Added
- Changelog.
- StreamObserver and StreamObserverContext interfaces to allow observing CQL streams.
- ClusterConfig.WriteTimeout option now allows to specify a write-timeout different from read-timeout.
- TypeInfo.NewWithError method.

### Changed
- Supported versions of Go that we test against are now Go 1.17 and Go 1.18.
- The driver now returns an error if SetWriteDeadline fails. If you need to run gocql on
  a platform that does not support SetWriteDeadline, set WriteTimeout to zero to disable the timeout.
- Creating streams on a connection that is closing now fails early.
- HostFilter now also applies to control connections.
- TokenAwareHostPolicy now panics immediately during initialization instead of at random point later
  if you reuse the TokenAwareHostPolicy between multiple sessions. Reusing TokenAwareHostPolicy between
  sessions was never supported.

### Fixed
- The driver no longer resets the network connection if a write fails with non-network-related error.
- Blocked network write to a network could block other goroutines, this is now fixed.
- Fixed panic in unmarshalUDT when trying to unmarshal a user-defined-type to a non-pointer Go type.
- Fixed panic when trying to unmarshal unknown/custom CQL type.

## Deprecated
- TypeInfo.New, please use TypeInfo.NewWithError instead.

## [1.0.0] - 2022-03-04
### Changed
- Started tagging versions with semantic version tags
