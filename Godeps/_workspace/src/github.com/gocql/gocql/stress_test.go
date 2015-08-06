// +build all integration

package gocql

import (
	"sync/atomic"

	"testing"
)

func BenchmarkConnStress(b *testing.B) {
	const workers = 16

	cluster := createCluster()
	cluster.NumConns = 1
	cluster.NumStreams = workers
	session := createSessionFromCluster(cluster, b)
	defer session.Close()

	if err := createTable(session, "CREATE TABLE IF NOT EXISTS conn_stress (id int primary key)"); err != nil {
		b.Fatal(err)
	}

	var seed uint64
	writer := func(pb *testing.PB) {
		seed := atomic.AddUint64(&seed, 1)
		var i uint64 = 0
		for pb.Next() {
			if err := session.Query("insert into conn_stress (id) values (?)", i*seed).Exec(); err != nil {
				b.Error(err)
				return
			}
			i++
		}
	}

	b.SetParallelism(workers)
	b.RunParallel(writer)

}
