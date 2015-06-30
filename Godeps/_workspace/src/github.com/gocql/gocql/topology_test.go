// +build all unit

package gocql

import (
	"testing"
)

// fakeNode is used as a simple structure to test the RoundRobin API
type fakeNode struct {
	conn   *Conn
	closed bool
}

// Pick is needed to satisfy the Node interface
func (n *fakeNode) Pick(qry *Query) *Conn {
	if n.conn == nil {
		n.conn = &Conn{}
	}
	return n.conn
}

//Close is needed to satisfy the Node interface
func (n *fakeNode) Close() {
	n.closed = true
}

//TestRoundRobinAPI tests the exported methods of the RoundRobin struct
//to make sure the API behaves accordingly.
func TestRoundRobinAPI(t *testing.T) {
	node := &fakeNode{}
	rr := NewRoundRobin()
	rr.AddNode(node)

	if rr.Size() != 1 {
		t.Fatalf("expected size to be 1, got %v", rr.Size())
	}

	if c := rr.Pick(nil); c != node.conn {
		t.Fatalf("expected conn %v, got %v", node.conn, c)
	}

	rr.Close()
	if rr.pool != nil {
		t.Fatalf("expected rr.pool to be nil, got %v", rr.pool)
	}

	if !node.closed {
		t.Fatal("expected node.closed to be true, got false")
	}
}
