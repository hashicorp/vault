package gore

import (
	"testing"
)

func TestTransaction(t *testing.T) {
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tr := NewTransaction(conn)
	tr.Add(NewCommand("SET", "kirisame", "marisa"))
	tr.Add(NewCommand("INCR", "X"))
	tr.Add(NewCommand("MGET", "kirisame", "X"))
	replies, err := tr.Commit()
	if err != nil {
		t.Fatal(err)
	}
	if len(replies) != 3 {
		t.Fatal("len error", replies)
	}
	if !replies[0].IsOk() {
		t.Fatal("not ok")
	}
	x, err := replies[1].Integer()
	if err != nil || x != 1 {
		t.Fatal(x, err)
	}
	s := []string{}
	err = replies[2].Slice(&s)
	if err != nil {
		t.Fatal(err)
	}
	if s[0] != "marisa" || s[1] != "1" {
		t.Fatal(s)
	}

	rep, err := NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}

func TestTransactionGoroutine(t *testing.T) {
	c := make(chan bool, 20)

	for i := 0; i < 50; i++ {
		go func(c chan bool) {
			defer func() {
				c <- true
			}()
			conn, err := Dial("localhost:6379")
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			for {
				tr := NewTransaction(conn)
				tr.Watch("X")
				rep, err := NewCommand("GET", "X").Run(conn)
				if err != nil {
					t.Fatal(err)
					return
				}
				x, err := rep.Int()
				if err != nil && err != ErrNil {
					t.Fatal(err)
					return
				}
				x++
				tr.Add(NewCommand("SET", "X", x))
				_, err = tr.Commit()
				if err == nil {
					break
				}
				if err != ErrKeyChanged {
					t.Fatal(err)
					break
				}
			}
		}(c)
	}

	for i := 0; i < 50; i++ {
		<-c
	}

	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	rep, err := NewCommand("GET", "X").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	x, err := rep.Int()
	if err != nil || x != 50 {
		t.Fatal(err, x)
	}
	rep, err = NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}
