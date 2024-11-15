raft-boltdb/v2
===========

This implementation uses the maintained version of BoltDB, [BBolt](https://github.com/etcd-io/bbolt). This is the primary version of `raft-boltdb` and should be used whenever possible. 

There is no breaking API change to the library. However, there is the potential for disk format incompatibilities so it was decided to be conservative and making it a separate import path. This separate import path will allow both versions (original and v2) to be imported to perform a safe in-place upgrade of old files read with the old version and written back out with the new one. 

## Metrics

The raft-boldb library emits a number of metrics utilizing github.com/armon/go-metrics. Those metrics are detailed in the following table. One note is that the application which pulls in this library may add its own prefix to the metric names. For example within [Consul](https://github.com/hashicorp/consul), the metrics will be prefixed with `consul.`.

| Metric                              | Unit         | Type    | Description           |
| ----------------------------------- | ------------:| -------:|:--------------------- |
| `raft.boltdb.freelistBytes`         | bytes        | gauge   | Represents the number of bytes necessary to encode the freelist metadata. When [`raft_boltdb.NoFreelistSync`](/docs/agent/options#NoFreelistSync) is set to `false` these metadata bytes must also be written to disk for each committed log. |
| `raft.boltdb.freePageBytes`         | bytes        | gauge   | Represents the number of bytes of free space within the raft.db file. |
| `raft.boltdb.getLog`                | ms           | timer   | Measures the amount of time spent reading logs from the db. |
| `raft.boltdb.logBatchSize`          | bytes        | sample  | Measures the total size in bytes of logs being written to the db in a single batch. |
| `raft.boltdb.logsPerBatch`          | logs         | sample  | Measures the number of logs being written per batch to the db. |
| `raft.boltdb.logSize`               | bytes        | sample  | Measures the size of logs being written to the db. |
| `raft.boltdb.numFreePages`          | pages        | gauge   | Represents the number of free pages within the raft.db file. |
| `raft.boltdb.numPendingPages`       | pages        | gauge   | Represents the number of pending pages within the raft.db that will soon become free. |
| `raft.boltdb.openReadTxn`           | transactions | gauge   | Represents the number of open read transactions against the db |
| `raft.boltdb.storeLogs`             | ms           | timer   | Measures the amount of time spent writing logs to the db. |
| `raft.boltdb.totalReadTxn`          | transactions | gauge   | Represents the total number of started read transactions against the db |
| `raft.boltdb.txstats.cursorCount`   | cursors      | counter | Counts the number of cursors created since Consul was started. |
| `raft.boltdb.txstats.nodeCount`     | allocations  | counter | Counts the number of node allocations within the db since Consul was started. |
| `raft.boltdb.txstats.nodeDeref`     | dereferences | counter | Counts the number of node dereferences in the db since Consul was started. |
| `raft.boltdb.txstats.pageAlloc`     | bytes        | gauge   | Represents the number of bytes allocated within the db since Consul was started. Note that this does not take into account space having been freed and reused. In that case, the value of this metric will still increase. |
| `raft.boltdb.txstats.pageCount`     | pages        | gauge   | Represents the number of pages allocated since Consul was started. Note that this does not take into account space having been freed and reused. In that case, the value of this metric will still increase. |
| `raft.boltdb.txstats.rebalance`     | rebalances   | counter | Counts the number of node rebalances performed in the db since Consul was started. |
| `raft.boltdb.txstats.rebalanceTime` | ms           | timer   | Measures the time spent rebalancing nodes in the db. |
| `raft.boltdb.txstats.spill`         | spills       | counter | Counts the number of nodes spilled in the db since Consul was started. |
| `raft.boltdb.txstats.spillTime`     | ms           | timer   | Measures the time spent spilling nodes in the db. |
| `raft.boltdb.txstats.split`         | splits       | counter | Counts the number of nodes split in the db since Consul was started. |
| `raft.boltdb.txstats.write`         | writes       | counter | Counts the number of writes to the db since Consul was started. |
| `raft.boltdb.txstats.writeTime`     | ms           | timer   | Measures the amount of time spent performing writes to the db. |
| `raft.boltdb.writeCapacity`         | logs/second  | sample  | Theoretical write capacity in terms of the number of logs that can be written per second. Each sample outputs what the capacity would be if future batched log write operations were similar to this one. This similarity encompasses 4 things: batch size, byte size, disk performance and boltdb performance. While none of these will be static and its highly likely individual samples of this metric will vary, aggregating this metric over a larger time window should provide a decent picture into how this BoltDB store can perform |
