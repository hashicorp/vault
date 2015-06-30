// +build all integration

package gocql

import (
	"fmt"
	"testing"
)

func TestSessionAPI(t *testing.T) {

	cfg := &ClusterConfig{}
	pool, err := NewSimplePool(cfg)
	if err != nil {
		t.Fatal(err)
	}

	s := &Session{
		Pool: pool,
		cfg:  *cfg,
		cons: Quorum,
	}

	defer s.Close()

	s.SetConsistency(All)
	if s.cons != All {
		t.Fatalf("expected consistency 'All', got '%v'", s.cons)
	}

	s.SetPageSize(100)
	if s.pageSize != 100 {
		t.Fatalf("expected pageSize 100, got %v", s.pageSize)
	}

	s.SetPrefetch(0.75)
	if s.prefetch != 0.75 {
		t.Fatalf("expceted prefetch 0.75, got %v", s.prefetch)
	}

	trace := &traceWriter{}

	s.SetTrace(trace)
	if s.trace != trace {
		t.Fatalf("expected traceWriter '%v',got '%v'", trace, s.trace)
	}

	qry := s.Query("test", 1)
	if v, ok := qry.values[0].(int); !ok {
		t.Fatalf("expected qry.values[0] to be an int, got %v", qry.values[0])
	} else if v != 1 {
		t.Fatalf("expceted qry.values[0] to be 1, got %v", v)
	} else if qry.stmt != "test" {
		t.Fatalf("expected qry.stmt to be 'test', got '%v'", qry.stmt)
	}

	boundQry := s.Bind("test", func(q *QueryInfo) ([]interface{}, error) {
		return nil, nil
	})
	if boundQry.binding == nil {
		t.Fatal("expected qry.binding to be defined, got nil")
	} else if boundQry.stmt != "test" {
		t.Fatalf("expected qry.stmt to be 'test', got '%v'", boundQry.stmt)
	}

	itr := s.executeQuery(qry)
	if itr.err != ErrNoConnections {
		t.Fatalf("expected itr.err to be '%v', got '%v'", ErrNoConnections, itr.err)
	}

	testBatch := s.NewBatch(LoggedBatch)
	testBatch.Query("test")
	err = s.ExecuteBatch(testBatch)

	if err != ErrNoConnections {
		t.Fatalf("expected session.ExecuteBatch to return '%v', got '%v'", ErrNoConnections, err)
	}

	s.Close()
	if !s.Closed() {
		t.Fatal("expected s.Closed() to be true, got false")
	}
	//Should just return cleanly
	s.Close()

	err = s.ExecuteBatch(testBatch)
	if err != ErrSessionClosed {
		t.Fatalf("expected session.ExecuteBatch to return '%v', got '%v'", ErrSessionClosed, err)
	}
}

func TestQueryBasicAPI(t *testing.T) {
	qry := &Query{}

	if qry.Latency() != 0 {
		t.Fatalf("expected Query.Latency() to return 0, got %v", qry.Latency())
	}

	qry.attempts = 2
	qry.totalLatency = 4
	if qry.Attempts() != 2 {
		t.Fatalf("expected Query.Attempts() to return 2, got %v", qry.Attempts())
	}
	if qry.Latency() != 2 {
		t.Fatalf("expected Query.Latency() to return 2, got %v", qry.Latency())
	}

	qry.Consistency(All)
	if qry.GetConsistency() != All {
		t.Fatalf("expected Query.GetConsistency to return 'All', got '%s'", qry.GetConsistency())
	}

	trace := &traceWriter{}
	qry.Trace(trace)
	if qry.trace != trace {
		t.Fatalf("expected Query.Trace to be '%v', got '%v'", trace, qry.trace)
	}

	qry.PageSize(10)
	if qry.pageSize != 10 {
		t.Fatalf("expected Query.PageSize to be 10, got %v", qry.pageSize)
	}

	qry.Prefetch(0.75)
	if qry.prefetch != 0.75 {
		t.Fatalf("expected Query.Prefetch to be 0.75, got %v", qry.prefetch)
	}

	rt := &SimpleRetryPolicy{NumRetries: 3}
	if qry.RetryPolicy(rt); qry.rt != rt {
		t.Fatalf("expected Query.RetryPolicy to be '%v', got '%v'", rt, qry.rt)
	}

	qry.Bind(qry)
	if qry.values[0] != qry {
		t.Fatalf("expected Query.Values[0] to be '%v', got '%v'", qry, qry.values[0])
	}
}

func TestQueryShouldPrepare(t *testing.T) {
	toPrepare := []string{"select * ", "INSERT INTO", "update table", "delete from", "begin batch"}
	cantPrepare := []string{"create table", "USE table", "LIST keyspaces", "alter table", "drop table", "grant user", "revoke user"}
	q := &Query{}

	for i := 0; i < len(toPrepare); i++ {
		q.stmt = toPrepare[i]
		if !q.shouldPrepare() {
			t.Fatalf("expected Query.shouldPrepare to return true, got false for statement '%v'", toPrepare[i])
		}
	}

	for i := 0; i < len(cantPrepare); i++ {
		q.stmt = cantPrepare[i]
		if q.shouldPrepare() {
			t.Fatalf("expected Query.shouldPrepare to return false, got true for statement '%v'", cantPrepare[i])
		}
	}
}

func TestBatchBasicAPI(t *testing.T) {

	cfg := &ClusterConfig{RetryPolicy: &SimpleRetryPolicy{NumRetries: 2}}
	pool, err := NewSimplePool(cfg)
	if err != nil {
		t.Fatal(err)
	}

	s := &Session{
		Pool: pool,
		cfg:  *cfg,
		cons: Quorum,
	}
	defer s.Close()

	b := s.NewBatch(UnloggedBatch)
	if b.Type != UnloggedBatch {
		t.Fatalf("expceted batch.Type to be '%v', got '%v'", UnloggedBatch, b.Type)
	} else if b.rt != cfg.RetryPolicy {
		t.Fatalf("expceted batch.RetryPolicy to be '%v', got '%v'", cfg.RetryPolicy, b.rt)
	}

	b = NewBatch(LoggedBatch)
	if b.Type != LoggedBatch {
		t.Fatalf("expected batch.Type to be '%v', got '%v'", LoggedBatch, b.Type)
	}

	b.attempts = 1
	if b.Attempts() != 1 {
		t.Fatalf("expceted batch.Attempts() to return %v, got %v", 1, b.Attempts())
	}

	if b.Latency() != 0 {
		t.Fatalf("expected batch.Latency() to be 0, got %v", b.Latency())
	}

	b.totalLatency = 4
	if b.Latency() != 4 {
		t.Fatalf("expected batch.Latency() to return %v, got %v", 4, b.Latency())
	}

	b.Cons = One
	if b.GetConsistency() != One {
		t.Fatalf("expected batch.GetConsistency() to return 'One', got '%s'", b.GetConsistency())
	}

	b.Query("test", 1)
	if b.Entries[0].Stmt != "test" {
		t.Fatalf("expected batch.Entries[0].Stmt to be 'test', got '%v'", b.Entries[0].Stmt)
	} else if b.Entries[0].Args[0].(int) != 1 {
		t.Fatalf("expected batch.Entries[0].Args[0] to be 1, got %v", b.Entries[0].Args[0])
	}

	b.Bind("test2", func(q *QueryInfo) ([]interface{}, error) {
		return nil, nil
	})

	if b.Entries[1].Stmt != "test2" {
		t.Fatalf("expected batch.Entries[1].Stmt to be 'test2', got '%v'", b.Entries[1].Stmt)
	} else if b.Entries[1].binding == nil {
		t.Fatal("expected batch.Entries[1].binding to be defined, got nil")
	}
	r := &SimpleRetryPolicy{NumRetries: 4}

	b.RetryPolicy(r)
	if b.rt != r {
		t.Fatalf("expected batch.RetryPolicy to be '%v', got '%v'", r, b.rt)
	}

	if b.Size() != 2 {
		t.Fatalf("expected batch.Size() to return 2, got %v", b.Size())
	}

}

func TestConsistencyNames(t *testing.T) {
	names := map[fmt.Stringer]string{
		Any:         "ANY",
		One:         "ONE",
		Two:         "TWO",
		Three:       "THREE",
		Quorum:      "QUORUM",
		All:         "ALL",
		LocalQuorum: "LOCAL_QUORUM",
		EachQuorum:  "EACH_QUORUM",
		Serial:      "SERIAL",
		LocalSerial: "LOCAL_SERIAL",
		LocalOne:    "LOCAL_ONE",
	}

	for k, v := range names {
		if k.String() != v {
			t.Fatalf("expected '%v', got '%v'", v, k.String())
		}
	}
}
