package gore

import (
	"testing"
)

func TestSharding(t *testing.T) {
	c := NewCluster()
	c.AddShard("127.0.0.1:6379", "127.0.0.1:6380")
	err := c.Dial()
	if err != nil {
		t.Fatal(err)
	}
	for x := 0; x < 10000; x++ {
		rep, err := c.Execute(NewCommand("SET", x, x))
		if err != nil || !rep.IsOk() {
			t.Fatal(err, rep)
		}
	}
	for x := 0; x < 10000; x++ {
		rep, err := c.Execute(NewCommand("GET", x))
		if err != nil {
			t.Fatal(err)
		}
		y, err := rep.Int()
		if err != nil || int64(x) != y {
			t.Fatal(err, x, y)
		}
	}
}
