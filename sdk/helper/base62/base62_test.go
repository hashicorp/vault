package base62

import (
	"bytes"
	"crypto/rand"
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
	str := "A fairly simple test case"

	e := Encode([]byte(str))
	b, err := Decode(nil, e)

	if err != nil {
		t.Fail()
	}
	if string(b) != str {
		t.Fail()
	}

	input := make([]byte, 4)
	for i := 0; i < 100; i++ {
		output := make([]byte, 4)
		rand.Read(input)
		str = Encode(input)
		b, err = Decode(output, str)
		if !bytes.Equal(b, input) {
			//e = Encode(b)
			t.Fail()
		}
	}
}