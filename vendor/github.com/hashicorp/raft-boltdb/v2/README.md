raft-boltdb/v2
===========

This implementation uses the maintained version of BoltDB, [BBolt](https://github.com/etcd-io/bbolt). This is the primary version of `raft-boltdb` and should be used whenever possible. 

There is no breaking API change to the library. However, there is the potential for disk format incompatibilities so it was decided to be conservative and making it a separate import path. This separate import path will allow both versions (original and v2) to be imported to perform a safe in-place upgrade of old files read with the old version and written back out with the new one. 