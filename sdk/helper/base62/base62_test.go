package base62

import (
	"testing"
)

func TestRandom(t *testing.T) {
	strings := make(map[string]struct{})

	for i := 0; i < 100000; i++ {
		c, err := Random(16)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := strings[c]; ok {
			t.Fatalf("Unexpected duplicate string: %s", c)
		}
		strings[c] = struct{}{}

	}

	for i := 0; i < 3000; i++ {
		c, err := Random(i)
		if err != nil {
			t.Fatal(err)
		}
		if len(c) != i {
			t.Fatalf("Expected length %d, got: %d", i, len(c))
		}
	}
}

func TestDecode(t *testing.T) {
	str := "A fairly simple test"
	e := Encode([]byte(str))
	b, err := DecodeString(e)

	if err != nil {
		t.Fail()
	}
	if string(b) != str {
		t.Fail()
	}

	// A slightly harder test
	str,err = Random(200)
	if err != nil {
		t.Fail()
	}
	b, err = DecodeString(str)
	if err != nil {
		t.Fail()
	}
	e = Encode(b)
	if e != str {
		t.Fail()
	}
}

func BenchmarkEncodeDecode(b *testing.B) {
	str,_ := Random(200)
	for i := 0; i<b.N; i++ {
		o, _ := DecodeString(str)
		e := Encode(o)
		if e != str {
			b.Fail()
			break
		}
	}
}