# Release Notes

## Version 10.5.2 (24 September 2024)

### Fixed issues

* [GOCBC-1659](https://jira.issues.couchbase.com/browse/GOCBC-1659): Fixed nil pointer dereference in "ns-server mode" where an attempt was made to refresh the config using CCCP upon receiving a 0x0d (config-only) KV status.
* [GOCBC-1660](https://jira.issues.couchbase.com/browse/GOCBC-1660): Fixed potential data race that was occurring because the value of a lock was being logged in kvMux.

### New Features and Behavioral Changes

* [GOCBC-1632](https://jira.issues.couchbase.com/browse/GOCBC-1632): Added support for server groups in GetOneReplica & LookupIn, at uncommitted API stability level.

## Version 10.5.1 (18 July 2024)

### Fixed issues

* [GOCBC-1644](https://issues.couchbase.com/browse/GOCBC-1644):
  Fixed race that could occur when a request was retried and cancelled at the same time.

### New Features and Behavioral Changes

* [GOCBC-1645](https://issues.couchbase.com/browse/GOCBC-1645):
* [GOCBC-1640](https://issues.couchbase.com/browse/GOCBC-1640):
  Adjust logging for config management.

## Version 10.5.0 (18 June 2024)

### Fixed issues

* [GOCBC-1636](https://issues.couchbase.com/browse/GOCBC-1636):
  Fixed race in clusterAgent (only used by AgentGroup().

### New Features and Behavioral Changes

* [GOCBC-1631](https://issues.couchbase.com/browse/GOCBC-1631):
  Handle scope not found error from memcached, and treat it the same as collection not found for retry behaviour.
* [GOCBC-1626](https://issues.couchbase.com/browse/GOCBC-1626):
  Updated behaviour of ping and http requests on agent close to cancel in flight requests and prevent requests from being sent after close.

## Version 10.4.1 (17 April 2024)

### Fixed issues

* [GOCBC-1615](https://issues.couchbase.com/browse/GOCBC-1615):
  Removed superfluous validations from SCRAM client.
* [GOCBC-1622](https://issues.couchbase.com/browse/GOCBC-1622):
  Updated handling of degraded target state for `WaitUntilReady`.


## Version 10.4.0 (13 March 2024)

### New Features and Behavioral Changes

* [GOCBC-1614](https://issues.couchbase.com/browse/GOCBC-1614):
  Updated API stability level of range scan to committed 
  Updated API stability level of `UseClusterMapNotifications` to committed 

## Version 10.3.2 (20 February 2024)

### New Features and Behavioral Changes

* [GOCBC-1585](https://issues.couchbase.com/browse/GOCBC-1585):
  If non-idempotent requests fail due to the socket closing while they are in-flight, they are now exposed to the retry orchestrator, with the reason `SocketCloseInFlightRetryReason`.
* [GOCBC-1591](https://issues.couchbase.com/browse/GOCBC-1591):
  Added support for scoped search indexes.

## Version 10.3.1 (17 January 2024)

### New Features and Behavioral Changes

* [GOCBC-1494](https://issues.couchbase.com/browse/GOCBC-1494):
  Added handling for some missing query error codes.
* [GOCBC-1555](https://issues.couchbase.com/browse/GOCBC-1555):
  Added support for change history in DCP.
* [GOCBC-1558](https://issues.couchbase.com/browse/GOCBC-1558):
  Added support for the new `NOT_LOCKED` KV status.
  Exposed `ErrMemdNotLocked`.
* [GOCBC-1561](https://issues.couchbase.com/browse/GOCBC-1561):
  Use the KV error map description as the error message when receiving an unknown KV status code.

### Fixed issues

* [GOCBC-1550](https://issues.couchbase.com/browse/GOCBC-1550):
  Fixed issue where the SDK deadlocks if it is in ns_server mode and the server provides a config which does not contain a `thisNode` entry.
* [GOCBC-1569](https://issues.couchbase.com/browse/GOCBC-1569):
  Fixed issue where all direct dispatch retries fail when a pipeline closed error occurs, as the same pipeline is used.
  Reload the pipeline if direct dispatch fails with pipeline closed.
* [GOCBC-1573](https://issues.couchbase.com/browse/GOCBC-1573):
  Fixed issue where the SDK attempts to use a prepared query name that exists in a different query context.
  Updated the query cache to store both the query context and the statement.
* [GOCBC-1586](https://issues.couchbase.com/browse/GOCBC-1586):
  Fixed issue where unmarshalling the collections manifest fails when the server reports that a collection has `maxTTL` equal to -1.
  Changed the type of `MaxTTL` in `ManifestCollection` from `uint32` to `int32`.

## Version 10.3.0 (21 November 2023)

### New Features and Behavioral Changes

* [GOCBC-1439](https://issues.couchbase.com/browse/GOCBC-1439):
  Improvements for faster failover.
  Added support for SnappyEverywhere Hello.
  Added support for cluster config known versions.
  Added support for NMVB deduplicated response bodies
  Added support for brief cluster map notifications, see: `AgentConfig{}.UseClusterMapNotifications`
* [GOCBC-1451](https://issues.couchbase.com/browse/GOCBC-1451):
  Updated DCP agents to only use HTTP polling.
* [GOCBC-1542](https://issues.couchbase.com/browse/GOCBC-1542):
  Dropped "identical revision" log line down to debug level.

### Fixed Issues

* [GOCBC-1449](https://issues.couchbase.com/browse/GOCBC-1449):
  Fixed issue where `GetCollectionManifest`, `GetAllCollectionManifests` and a`GetCollectionID` options were missing user impersonation.
* [GOCBC-1471](https://issues.couchbase.com/browse/GOCBC-1471):
  Fixed issue where calling Close on an Agent before polling has started up would lead to the poller failing to stop.

## Version 10.2.9 (18 October 2023)

### Fixed Issues

* [GOCBC-1483](https://issues.couchbase.com/browse/GOCBC-1483):
  Retry requests when resetting cid cache queue.

### New Features and Behavioral Changes

* [GOCBC-1489](https://issues.couchbase.com/browse/GOCBC-1489):
  Expose `ErrCircuitBreakerOpen`.

## Version 10.2.8 (25 September 2023)

### New Features and Behavioral Changes

* [GOCBC-1479](https://issues.couchbase.com/browse/GOCBC-1479):
  Do not apply cluster configs during bootstrap if select bucket fails.

## Version 10.2.7 (30 August 2023)

### Fixed Issues

* [GOCBC-1458](https://issues.couchbase.com/browse/GOCBC-1458):
  Updated route config log output.
* [GOCBC-1465](https://issues.couchbase.com/browse/GOCBC-1465):
  Fixed issue where NoTLSSeedNode localhost ipv6 seed was parsed incorrectly.

## Version 10.2.6 (25 July 2023)

### New Features and Behavioral Changes

* [GOCBC-1434](https://issues.couchbase.com/browse/GOCBC-1434):
  Added support for `LookupIn` replica reads.

### Fixed Issues

* [GOCBC-1446](https://issues.couchbase.com/browse/GOCBC-1446):
  Fixed issue where calling `Commit` without any staged mutations would execute the callback with an error but not return from the function.
* [GOCBC-1441](https://issues.couchbase.com/browse/GOCBC-1429):
  Reverted [GOCBC-1429](https://issues.couchbase.com/browse/GOCBC-1429).

## Version 10.2.5 (21 June 2023)

### New Features and Behavioral Changes

* [GOCBC-1426](https://issues.couchbase.com/browse/GOCBC-1426):
  Improved logging during agent shutdown.
* [GOCBC-1432](https://issues.couchbase.com/browse/GOCBC-1432):
  Improved logging around when a cluster config is rejected.

### Fixed Issues

* [GOCBC-1413](https://issues.couchbase.com/browse/GOCBC-1413):
  Fixed issue where a panic would occur when UseTLS was set to false and NoTLSSeedNode set to true - no returns an error.
* [GOCBC-1429](https://issues.couchbase.com/browse/GOCBC-1429):
  Fixed issue where a node would incorrectly be identified as the seed node when running in ns_server mode.

## Version 10.2.4 (17 May 2023)

### New Features and Behavioral Changes

* [GOCBC-1409](https://issues.couchbase.com/browse/GOCBC-1409):
  Updated DCP OSO to use improved SeqNoAdvance

### Fixed Issues

* [GOCBC-1410](https://issues.couchbase.com/browse/GOCBC-1410):
  Fixed issue where applying a cluster config during SDK shutdown could panic.

## Version 10.2.3 (18 April 2023)

### Fixed Issues

* [GOCBC-1401](https://issues.couchbase.com/browse/GOCBC-1401):
  Exposed SeqNo on DCP rollback error.
* [GOCBC-1403](https://issues.couchbase.com/browse/GOCBC-1403):
  Fixed issue where cccp poller would wait for a cluster config before starting.

## Version 10.2.2 (22 March 2023)

### New Features and Behavioral Changes

* [GOCBC-1393](https://issues.couchbase.com/browse/GOCBC-1393):
  Altered the behaviour of retries for enhanced prepared statements.
* [GOCBC-1395](https://issues.couchbase.com/browse/GOCBC-1395):
  Improved timeout errors on http based services.

## Version 10.2.1 (22 February 2023)

### New Features and Behavioral Changes

* [GOCBC-1362](https://issues.couchbase.com/browse/GOCBC-1362):
  Added support for sending unsupported frames with `memd.Conn`.
* [GOCBC-1322](https://issues.couchbase.com/browse/GOCBC-1322):
  Added volatile stability support for kv range scan.
  Added volatile stability support for waiting for a config snapshot to be available.
* [GOCBC-1373](https://issues.couchbase.com/browse/GOCBC-1373):
  Added support for query error code 1197.

### Fixed Issues
* [GOCBC-1376](https://issues.couchbase.com/browse/GOCBC-1376):
  Fixed issue where lost cleanup would log an incorrectly formatted log line.
* [GOCBC-1387](https://issues.couchbase.com/browse/GOCBC-1387):
  Fixed issue where an edge case could trigger a race between releasing connection buffers and reading on the connection - leading to a panic.
* [GOCBC-1388](https://issues.couchbase.com/browse/GOCBC-1388):
  Fixed issue where the SDK could not connect to all nodes when `NoTLSSeedNode` is set in environments where multiple nodes are identifying as 127.0.0.1 (and so do not set a hostname in the cluster config).

## Version 10.2.0 (19 October 2022)

### New Features and Behavioral Changes

* [GOCBC-1159](https://issues.couchbase.com/browse/GOCBC-1159):
  Added support for refreshing the DNS SRV record when cluster becomes uncontactable.
* [GOCBC-1284](https://issues.couchbase.com/browse/GOCBC-1284):
  Significant refactoring work to kv bootstrap.
* [GOCBC-1303](https://issues.couchbase.com/browse/GOCBC-1303):
  Added `ServerWaitBackoff` to agent options.
* [GOCBC-1316](https://issues.couchbase.com/browse/GOCBC-1316):
  Added support for transactions ExtInsertExisting.
* [GOCBC-1328](https://issues.couchbase.com/browse/GOCBC-1328):
  Only fallback from cccp polling to http polling once all nodes tried.
* [GOCBC-1331](https://issues.couchbase.com/browse/GOCBC-1331):
  Added support for pipelining fetching a config into kv bootstrap.
* [GOCBC-1335](https://issues.couchbase.com/browse/GOCBC-1335):
  Updated logging to include address and pointer location in memdclient.
* [GOCBC-1351](https://issues.couchbase.com/browse/GOCBC-1351):
  Updated error message logged on auth failures.
* [GOCBC-1352](https://issues.couchbase.com/browse/GOCBC-1352):
  Added support for trusting the system cert store when TLS enabled and no cert provider registered.
* [GOCBC-1356](https://issues.couchbase.com/browse/GOCBC-1356):
  Updated the behaviour when `MutateIn` or `Add` returns `NOT_STORED` to return a `ErrDocumentExists`.


### Fixed Issues
* [GOCBC-1347](https://issues.couchbase.com/browse/GOCBC-1347):
  Fixed issue where a nil agent value could cause logging `TransactionATRLocation` to log a panic.
* [GOCBC-1348](https://issues.couchbase.com/browse/GOCBC-1348):
  Fixed issue where a race on creating a client record could lead to a panic.

## Version 10.1.5 (21 September 2022)

### New Features and Behavioral Changes

* [GOCBC-1293](https://issues.couchbase.com/browse/GOCBC-1293):
  Added support for resource units.
* [GOCBC-1332](https://issues.couchbase.com/browse/GOCBC-1332):
  Added deadlines to collections operations options.
* [GOCBC-1339](https://issues.couchbase.com/browse/GOCBC-1339):
  Removed support for `CleanupWatchATRs` from `TransactionsConfig`.
  Note that whilst this field still exists it is *not* used internally, it is included only for API level backward compatibility.
* [GOCBC-1340](https://issues.couchbase.com/browse/GOCBC-1340):
  Added support for automatically starting lost cleanup on `TransactionsConfig` `CustomATRLocation`.


### Fixed Issues
* [GOCBC-1338](https://issues.couchbase.com/browse/GOCBC-1338):
  Fixed issue where `lazyCircuitBreaker` was not using 64-bit aligned values.

### Known Issues
* [GOCBC-1347](https://issues.couchbase.com/browse/GOCBC-1347):
  Known issue where a nil agent value could cause logging `TransactionATRLocation` to log a panic.
* [GOCBC-1348](https://issues.couchbase.com/browse/GOCBC-1348):
  Known issue where a race on creating a client record can lead to a panic.

## Version 10.1.4 (20 July 2022)

### New Features and Behavioral Changes

* [GOCBC-1246](https://issues.couchbase.com/browse/GOCBC-1246):
  Added support for `TransactionLogger` to `TransactionOptions`.
* [GOCBC-1314](https://issues.couchbase.com/browse/GOCBC-1314):
  Improved logging in the lost transactions process.
* [GOCBC-1318](https://issues.couchbase.com/browse/GOCBC-1318):
  Changed `WaitUntilReady` to always wait for any explicitly defined services to be online.
* [GOCBC-1319](https://issues.couchbase.com/browse/GOCBC-1319):
  Added a `String` implemented to `memd.Packet`.


### Fixed Issues
* [GOCBC-1320](https://issues.couchbase.com/browse/GOCBC-1320):
  Fixed issue where vbucket hashing function wasn't masking out the 16th bit of the key.

## Version 10.1.3 (22 June 2022)

### New Features and Behavioral Changes

* [GOCBC-1264](https://issues.couchbase.com/browse/GOCBC-1264):
  Added more documentation to `AgentConfig`.
* [GOCBC-1298](https://issues.couchbase.com/browse/GOCBC-1298):
* [GOCBC-1299](https://issues.couchbase.com/browse/GOCBC-1299):
  Masked the underlying cause of `TransactionOperationFailedError`.
* [GOCBC-1159](https://issues.couchbase.com/browse/GOCBC-1159):
  Made improvements to handle a rebalance during a freeze in serverless environments.
* [GOCBC-1283](https://issues.couchbase.com/browse/GOCBC-1283):
  Update forward compatibility errors to include document details.

### Fixed Issues
* [GOCBC-1300](https://issues.couchbase.com/browse/GOCBC-1300):
  Added collection unknown check to `ProcessATR` to improve lost cleanup deleted collection handling.
* [GOCBC-1304](https://issues.couchbase.com/browse/GOCBC-1304):
  Fixed issue where lost cleanup would block the SDK response thread for a connection.
* [GOCBC-1301](https://issues.couchbase.com/browse/GOCBC-1301):
  Fixed issue where `addLostCleanupLocation` was left nil after `ResumeTransactionAttempt` called.

## Version 10.1.2 (26 April 2022)

### New Features and Behavioral Changes

* [GOCBC-1265](https://issues.couchbase.com/browse/GOCBC-1265):
  Bundle Capella CA certificates with the SDK.

## Version 10.1.1 (15 March 2022)

### New Features and Behavioral Changes

* [GOCBC-1221](https://issues.couchbase.com/browse/GOCBC-1221):
  Added support for improved query error handling.
* [GOCBC-1238](https://issues.couchbase.com/browse/GOCBC-1238):
  Add config option to set the connection read buffer size.
* [GOCBC-1242](https://issues.couchbase.com/browse/GOCBC-1242):
  Drain DCP queue on non-user initiated EOF.
* [GOCBC-1244](https://issues.couchbase.com/browse/GOCBC-1221):
  Updated dependencies.

### Fixed Issues

* [GOCBC-1248](https://issues.couchbase.com/browse/GOCBC-1248):
  Fixed issue where a hard close of a memdclient during a graceful close could trigger a panic.
* [GOCBC-1256](https://issues.couchbase.com/browse/GOCBC-1256):
  Fixed issue where config polling would fallback to using the http poller, when no http addresses are registered for use.
* [GOCBC-1258](https://issues.couchbase.com/browse/GOCBC-1258):
  Fixed issue where log redaction tags were not closed correctly.

## Version 10.1.0 (15 February 2022)

### New Features and Behavioral Changes

* [TXNG-127](https://issues.couchbase.com/browse/TXNG-127):
  Integrate transactions into SDK.

### Fixed Issues

* [GOCBC-1232](https://issues.couchbase.com/browse/GOCBC-1232):
  Fixed issue where DCP stream End could race with request cancellation (due to rebalance, etc...).
* [GOCBC-1233](https://issues.couchbase.com/browse/GOCBC-1233):
  Fixed issue where Agent close could hang if called whilst auth request in flight.

## Version 10.0.7 (24 January 2022)

### New Features and Behavioral Changes

* [GOCBC-1216](https://issues.couchbase.com/browse/GOCBC-1216):
  Add support for missing memcached status code 0x8d
* [GOCBC-1222](https://issues.couchbase.com/browse/GOCBC-1222):
  Updated memcached connections to use a `sync.Pool` for buffers for readers, to help reduce memory footprint.

### Fixed Issues

* [GOCBC-1214](https://issues.couchbase.com/browse/GOCBC-1214):
  Fixed issue where nodes "actual" IP could be used for internal config instead of seed address when `NoTLSSeedNode` in use.

## Version 10.0.6 (14 December 2021)

### New Features and Behavioral Changes

* [GOCBC-1190](https://issues.couchbase.com/browse/GOCBC-1190):
  Added internal stability support for sending queries to specific nodes.
* [GOCBC-1196](https://issues.couchbase.com/browse/GOCBC-1196):
* Added error body and status code to analytics, query, search, view errors.

### Fixed Issues

* [GOCBC-1205](https://issues.couchbase.com/browse/GOCBC-1205):
  Fixed issue where tracer spans were not always being finished.
* [GOCBC-1206](https://issues.couchbase.com/browse/GOCBC-1206):
  Fixed issue where metrics were always incorrectly reporting very short durations for operations.
* [GOCBC-1208](https://issues.couchbase.com/browse/GOCBC-1208):
  Fixed issue where cluster config polling would fallback to HTTP polling even when there was no bucket.
* [GOCBC-1209](https://issues.couchbase.com/browse/GOCBC-1209):
  Fixed issue where the ns server connection string scheme wouldn't work for DCP.


## Version 10.0.5 (16 November 2021)

### New Features and Behavioral Changes

* [GOCBC-1179](https://issues.couchbase.com/browse/GOCBC-1179):
  Gracefully close memdclients on pipeline shutdown/reconnect.
* [GOCBC-1180](https://issues.couchbase.com/browse/GOCBC-1180):
  Added support for the ns_server connection string scheme and seed (i.e. localhost) poller.
* [GOCBC-1181](https://issues.couchbase.com/browse/GOCBC-1181):
  Added support for `ReconfigureSecurity` function.
* [GOCBC-1182](https://issues.couchbase.com/browse/GOCBC-1182):
  Request error map v2 from the server.
* [GOCBC-1193](https://issues.couchbase.com/browse/GOCBC-1193):
  Added the response body to query errors.

### Fixed Issues

* [GOCBC-1194](https://issues.couchbase.com/browse/GOCBC-1194):
  Fixed issue where we wouldn't try to build a route config with all seed nodes for default network type before trying external network type.

## Version 10.0.4 (19 October 2021)

### New Features and Behavioral Changes

* [GOCBC-1178](https://issues.couchbase.com/browse/GOCBC-1178):
  Don't remove poller controller watcher from cluster config updates.

### Fixed Issues

* [GOCBC-1177](https://issues.couchbase.com/browse/GOCBC-1177):
  Fixed issue where a connection being closed by the server during bootstrap could cause the SDK to loop reconnect without backoff.


## Version 10.0.3 (21 September 2021)

###New Features and Behavioral Changes

* [GOCBC-1162](https://issues.couchbase.com/browse/GOCBC-1162):
  Added support for initially bootstrapping the SDK over nonTLS when TLS is in use.
* [GOCBC-1169](https://issues.couchbase.com/browse/GOCBC-1169):
  Updated query streamer so that additional calls to `NextRow` return nil rather than panic.

### Fixed Issues

* [GOCBC-1160](https://issues.couchbase.com/browse/GOCBC-1160):
  Fixed issue where HTTP header used for user impersonation was incorrect.
* [GOCBC-1163](https://issues.couchbase.com/browse/GOCBC-1163):
  Fixed issue where cluster config parsing would check existence of wrong ports for TLS (although then assign correct ports).

## Version 10.0.2 (17 August 2021)

###New Features and Behavioral Changes

* [GOCBC-1146](https://issues.couchbase.com/browse/GOCBC-1146):
  Added support for user impersonation to non-KV services.
* [GOCBC-1148](https://issues.couchbase.com/browse/GOCBC-1148):
  Added support for forcibly reconnecting all connections.
* [GOCBC-1150](https://issues.couchbase.com/browse/GOCBC-1150):
  Update user impersonation options for KV to use a string rather than []byte.

### Fixed Issues

* [GOCBC-1139](https://issues.couchbase.com/browse/GOCBC-1139):
  Fixed issue where DCP agent would try to use SCRAM auth with TLS enabled, causing LDAP usage to always fail bootstrap.
* [GOCBC-1147](https://issues.couchbase.com/browse/GOCBC-1147):
  Fixed issue where failing to fetch the error map during bootstrap would lead to bootstrap hanging.

## Version 10.0.1 (15 July 2021)

### Fixed Issues

* Fixed issue where modules file contained incorrect gocbcore version.

## Version 10.0.0 (15 July 2021) (Do not use, see v10.0.1)

###New Features and Behavioral Changes

* [GOCBC-901](https://issues.couchbase.com/browse/GOCBC-901):
  Broke the `AgentConfig` up into grouped components.
* [GOCBC-1008](https://issues.couchbase.com/browse/GOCBC-1008):
  Updated mutate in to return cas mismatch error rather than document exists when doing a replace.
* [GOCBC-1062](https://issues.couchbase.com/browse/GOCBC-1062):
  Added support for DCP snapshot marker v2 and v2.1.
* [GOCBC-1081](https://issues.couchbase.com/browse/GOCBC-1081):
  During CCCP polling don't retry request if the error is request cancelled.
* [GOCBC-1130](https://issues.couchbase.com/browse/GOCBC-1130):
  Updated Query error handling to return an authentication error on error code 13104.
* [GOCBC-1087](https://issues.couchbase.com/browse/GOCBC-1087):
  Added support for communicating with Eventing and Backup services.
* [GOCBC-1093](https://issues.couchbase.com/browse/GOCBC-1093):
  Added support for `RevEpoch` in bucket configs.
* [GOCBC-1044](https://issues.couchbase.com/browse/GOCBC-1044):
* [GOCBC-1128](https://issues.couchbase.com/browse/GOCBC-1128):
  Added `Meter` interface and operation level response latency metric.
* [GOCBC-1133](https://issues.couchbase.com/browse/GOCBC-1133):
  Remove `ViewQuery` from `AgentGroup`.

### Fixed Issues

* [GOCBC-1135](https://issues.couchbase.com/browse/GOCBC-1135):
  Fixed issue where cmd traces could be ended twice in some scenarios when operation was cancelled.

## Version 9.1.5 (15 June 2021)

### Fixed Issues

* [GOCBC-1095](https://issues.couchbase.com/browse/GOCBC-1095):
  Fixed issue where SDK was parsing view error contents incorrectly.
* [GOCBC-1102](https://issues.couchbase.com/browse/GOCBC-1102):
  Fixed issue where `WaitUntilReady` wouldn't recover if one of the HTTP based services returned an error.
* [GOCBC-1106](https://issues.couchbase.com/browse/GOCBC-1106):
* [GOCBC-1112](https://issues.couchbase.com/browse/GOCBC-1112):
  Fixed issues where fts responses were being parsed incorrectly.
* [GOCBC-1127](https://issues.couchbase.com/browse/GOCBC-1127):
  Fixed issue where query errors could be parsed incorrectly.

## Version 9.1.4 (20 April 2021)

###New Features and Behavioral Changes

* [GOCBC-1071](https://issues.couchbase.com/browse/GOCBC-1071):
  Updated SDK to use new protocol level changes for get collection id.
* [GOCBC-1068](https://issues.couchbase.com/browse/GOCBC-1068):
  Dropped log level to warn for when applying a cluster config object is preempted.
* [GOCBC-1079](https://issues.couchbase.com/browse/GOCBC-1079):
  During bootstrap don't retry authentication if the error is request cancelled.
* [GOCBC-1081](https://issues.couchbase.com/browse/GOCBC-1081):
  During CCCP polling don't retry request if the error is request cancelled.

### Fixed Issues

* [GOCBC-1080](https://issues.couchbase.com/browse/GOCBC-1080):
  Fixed issue where SDK would always rebuild connections on first cluster config fetched against server 7.0.
* [GOCBC-1082](https://issues.couchbase.com/browse/GOCBC-1082):
  Fixed issue where bootstrapping a node during an SDK wide reconnect would cause a delay in connecting to that node.
* [GOCBC-1088](https://issues.couchbase.com/browse/GOCBC-1088):
  Fixed issue where the poller controller could deadlock if a node reported a bucket not found at the same time as CCCP successfully fetched a cluster config for the first time.
  
## Version 9.1.3 (16 March 2021)

### New Features and Behavioral Changes

* [GOCBC-1056](https://issues.couchbase.com/browse/GOCBC-1056):
  Various performance improvements to reduce CPU level.
* [GOCBC-1068](https://issues.couchbase.com/browse/GOCBC-1068):
  Dropped the log level for preempted config updates.
* [GOCBC-940](https://issues.couchbase.com/browse/GOCBC-940):
  Updated the tracing interfaces and orphaned response logging output.

### Fixed Issues

* [GOCBC-1066](https://issues.couchbase.com/browse/GOCBC-1066):
  Fixed issue which could cause the config pollers to panic.

## Version 9.1.2 (16 February 2021)

### New Features and Behavioral Changes

* [GOCBC-1041](https://issues.couchbase.com/browse/GOCBC-1041):
  Dropped the log level for memdclient read failures to warn, from error.
* [GOCBC-1046](https://issues.couchbase.com/browse/GOCBC-1046):
  Added `MaxTTl` to `ManifestCollection`.

### Fixed Issues

* [GOCBC-1042](https://issues.couchbase.com/browse/GOCBC-1042):
  Fixed issue where bucket names were not being correctly escaped.
* [GOCBC-1050](https://issues.couchbase.com/browse/GOCBC-1050):
  Fixed issue where the diagnostics component could panic if an operation was cancelled by the user after it had already been internally cancelled.

## Version 9.1.1 (19 January 2021)

### New Features and Behavioral Changes

* [GOCBC-1032](https://issues.couchbase.com/browse/GOCBC-1032):
  Added support for bucket capability support verification to agent, at API stability internal.
* [GOCBC-1030](https://issues.couchbase.com/browse/GOCBC-1030):
  Added support for internal cancellation of bootstrap before completion, allowing pipeline clients to shutdown without waiting for bootstrap to complete (such as on connection takeover).

  Added support to fallback to http config fetching if select bucket fails with a valid fallback error, allowing for faster config fetching against non-kv nodes.

## Version 9.1.0 (15 December 2020)

### New Features and Behavioral Changes

* [GOCBC-854](https://issues.couchbase.com/browse/GOCBC-854):
Added support for user impersonation.
* [GOCBC-1013](https://issues.couchbase.com/browse/GOCBC-1013):
Added support for `StatsKeys` and `StatsChunks` to `SingleServerStats` to support responses for stats keys such as `connections` which contain complex objects per packet.

### Fixed Issues

* [GOCBC-1016](https://issues.couchbase.com/browse/GOCBC-1016):
Fixed issue where creating an agent with no bucket and a non-default port HTTP address could lead to a panic in `WaitForReady`.
(Note: `WaitForReady` will *never* return success in this scenario)
* [GOCBC-1028](https://issues.couchbase.com/browse/GOCBC-1028):
Fixed issue where bootstrapping against a non-kv node could never successfully fully connect.
