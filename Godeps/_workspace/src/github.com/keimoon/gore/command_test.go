package gore

import (
	"testing"
)

func TestGetSetBasic(t *testing.T) {
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	rep, err := NewCommand("SET", "kirisame", "marisa").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
	rep, err = NewCommand("GET", "kirisame").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	if s, err := rep.String(); s != "marisa" {
		t.Fatal(s, err)
	}
	rep, err = NewCommand("DEL", "kirisame").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	if x, err := rep.Integer(); x != 1 {
		t.Fatal(x, err)
	}
	rep, err = NewCommand("GET", "kirisame").Run(conn)
	if err != nil || !rep.IsNil() {
		t.Fatal(err, "not nil")
	}
	rep, err = NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}

func TestValueConverting(t *testing.T) {
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	rep, err := NewCommand("SET", "int", -12345678987654321).Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
	rep, err = NewCommand("GET", "int").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	if x, err := rep.Int(); err != nil || x != -12345678987654321 {
		t.Fatal(x, err)
	}
	rep, err = NewCommand("SET", "float", -1234567.8765).Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
	rep, err = NewCommand("GET", "float").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	if x, err := rep.Float(); err != nil || x != -1234567.8765 {
		t.Fatal(x, err)
	}
	rep, err = NewCommand("SET", "bool", true).Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
	rep, err = NewCommand("GET", "bool").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	if x, err := rep.Bool(); err != nil || !x {
		t.Fatal(x, err)
	}
	rep, err = NewCommand("SET", "fixint", FixInt(-12345678987654321)).Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
	rep, err = NewCommand("GET", "fixint").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	if x, err := rep.FixInt(); err != nil || x != -12345678987654321 {
		t.Fatal(x, err)
	}
	rep, err = NewCommand("SET", "varint", VarInt(-12345678987654321)).Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
	rep, err = NewCommand("GET", "varint").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	if x, err := rep.VarInt(); err != nil || x != -12345678987654321 {
		t.Fatal(x, err)
	}
	rep, err = NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}

func TestArrayValue(t *testing.T) {
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	rep, err := NewCommand("ZADD", "int", 1, 1, 2, 2, 3, 3).Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	if x, err := rep.Integer(); x != 3 {
		t.Fatal(x, err)
	}
	rep, err = NewCommand("ZRANGE", "int", 0, -1).Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	sint := []int{}
	err = rep.Slice(&sint)
	if err != nil {
		t.Fatal(err)
	}
	if len(sint) != 3 || sint[0] != 1 || sint[1] != 2 || sint[2] != 3 {
		t.Fatal(sint)
	}
	rep, err = NewCommand("ZADD", "varint", 1, VarInt(1), 2, VarInt(2), 3, VarInt(3)).Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	if x, err := rep.Integer(); x != 3 {
		t.Fatal(x, err)
	}
	rep, err = NewCommand("ZRANGE", "varint", 0, -1).Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	svarint := []VarInt{}
	err = rep.Slice(&svarint)
	if err != nil {
		t.Fatal(err)
	}
	if len(svarint) != 3 || svarint[0] != 1 || svarint[1] != 2 || svarint[2] != 3 {
		t.Fatal(svarint)
	}
	rep, err = NewCommand("HMSET", "hash", "kazami", "yuuka", "shameimaru", "aya").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
	rep, err = NewCommand("HMGET", "hash", "kazami", "xxx", "shameimaru").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	sstring := []string{}
	err = rep.Slice(&sstring)
	if err != nil {
		t.Fatal(err)
	}
	if len(sstring) != 3 || sstring[0] != "yuuka" || sstring[1] != "" || sstring[2] != "aya" {
		t.Fatal(sstring)
	}
	rep, err = NewCommand("HGETALL", "hash").Run(conn)
	if err != nil {
		t.Fatal(err)
	}
	spair := []*Pair{}
	err = rep.Slice(&spair)
	if err != nil {
		t.Fatal(err)
	}
	if len(spair) != 2 || string(spair[0].First) != "kazami" || string(spair[0].Second) != "yuuka" ||
		string(spair[1].First) != "shameimaru" || string(spair[1].Second) != "aya" {
		t.Fatal(spair)
	}
	rep, err = NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}

func TestCommandGoroutine(t *testing.T) {
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := make(chan bool, 20)
	for i := 0; i < 2176; i++ {
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
	for i := 0; i < 2176; i++ {
		<-c
	}
	rep, err := NewCommand("FLUSHALL").Run(conn)
	if err != nil || !rep.IsOk() {
		t.Fatal(err, "not ok")
	}
}
