package gore

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	p := NewPipeline()
	for i := 0; i < 100; i++ {
		p.Add(NewCommand("SET", i, i))
	}
	replies, err := p.Run(conn)
	if err != nil || len(replies) != 100 {
		t.Fatal(err, len(replies))
	}
	for _, rep := range replies {
		if !rep.IsOk() {
			t.Fatal("not ok", rep)
		}
	}

	p = NewPipeline()
	for i := 0; i < 100; i++ {
		p.Add(NewCommand("GET", i))
	}
	replies, err = p.Run(conn)
	if err != nil || len(replies) != 100 {
		t.Fatal(err, len(replies))
	}
	for i, rep := range replies {
		x, err := rep.Int()
		if err != nil || x != int64(i) {
			t.Fatal(err, i, x)
		}
	}

	rep, err := NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}

func TestPipelineGoroutine(t *testing.T) {
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := make(chan bool, 20)
	for i := 0; i < 1000; i++ {
		go func(conn *Conn, c chan bool, x int64) {
			defer func() {
				c <- true
			}()
			_, err := NewCommand("SET", x, x).Run(conn)
			if err != nil {
				t.Fatal(err)
			}
			rep, err := NewCommand("GET", x).Run(conn)
			if err != nil {
				t.Fatal(err)
			}
			y, err := rep.Int()
			if err != nil || y != x {
				t.Fatal(err, x, y)
			}
		}(conn, c, int64(i))
	}
	for i := 0; i < 100; i++ {
		go func(conn *Conn, c chan bool, x int64) {
			defer func() {
				c <- true
			}()
			p := NewPipeline()
			for j := 0; j < 100; j++ {
				val := 1000 + x*100 + int64(j)
				p.Add(NewCommand("SET", val, val))
			}
			replies, err := p.Run(conn)
			if err != nil || len(replies) != 100 {
				t.Fatal(err, len(replies))
			}
			for _, r := range replies {
				if !r.IsOk() {
					t.Fatal("not ok", r)
				}
			}
			p = NewPipeline()
			for j := 0; j < 100; j++ {
				val := 1000 + x*100 + int64(j)
				p.Add(NewCommand("GET", val))
			}
			replies, err = p.Run(conn)
			if err != nil || len(replies) != 100 {
				t.Fatal(err, len(replies))
			}
			for j, r := range replies {
				y, err := r.Int()
				if err != nil || y != 1000+x*100+int64(j) {
					t.Fatal(err, x, j, y)
				}
			}
		}(conn, c, int64(i))
	}
	for i := 0; i < 1100; i++ {
		<-c
	}

	rep, err := NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}
