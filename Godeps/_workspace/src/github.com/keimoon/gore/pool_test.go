package gore

import (
	"testing"
)

func TestPool(t *testing.T) {
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	pool := &Pool{
		InitialConn: 5,
		MaximumConn: 5,
	}
	err = pool.Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}

	c := make(chan bool, 20)
	for i := 0; i < 10000; i++ {
		go func(pool *Pool, c chan bool, x int64) {
			defer func() {
				c <- true
			}()
			conn, err := pool.Acquire()
			if err != nil || conn == nil {
				t.Fatal(err, "nil")
			}
			defer pool.Release(conn)
			rep, err := NewCommand("SET", x, x).Run(conn)
			if err != nil || !rep.IsOk() {
				t.Fatal(err, "not ok")
			}
		}(pool, c, int64(i))
	}
	for i := 0; i < 10000; i++ {
		<-c
	}
	for i := 0; i < 10000; i++ {
		go func(pool *Pool, c chan bool, x int64) {
			defer func() {
				c <- true
			}()
			conn, err := pool.Acquire()
			if err != nil || conn == nil {
				t.Fatal(err, "nil")
			}
			defer pool.Release(conn)
			rep, err := NewCommand("GET", x).Run(conn)
			if err != nil {
				t.Fatal(err)
			}
			y, err := rep.Int()
			if err != nil || y != x {
				t.Fatal(err, x, y)
			}
		}(pool, c, int64(i))
	}
	for i := 0; i < 10000; i++ {
		<-c
	}
	rep, err := NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}

func TestPoolClose(t *testing.T) {
	pool := &Pool{
		InitialConn: 20,
		MaximumConn: 20,
	}
	err := pool.Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	c := make(chan bool, 20)
	ready := make(chan bool, 10)
	for i := 0; i < 1000; i++ {
		go func() {
			defer func() {
				c <- true
			}()
			conn, err := pool.Acquire()			
			if err != nil {
				t.Fatal(err)
			}
			if conn != nil {
				ready <- true
			}
		}()
	}
	for i := 0; i < 20; i++ {
		<- ready
	}
	pool.Close()
	for i := 0; i < 1000; i++ {
		<-c
	}
}
