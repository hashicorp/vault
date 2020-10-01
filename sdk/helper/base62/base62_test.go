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

	e := EncodeToString([]byte(str))
	b, err := DecodeString(e)

	if err != nil || string(b) != str {
		t.Fail()
	}

	_, err = DecodeString("invalid.string")
	if err == nil {
		t.Fail()
	}

	_, err = DecodeString(".also-bad")
	if err == nil {
		t.Fail()
	}

	input := make([]byte, 50)
	for i := 0; i < 100; i++ {
		rand.Read(input)
		str = EncodeToString(input)
		b, err := DecodeString(str)
		if err != nil || !bytes.Equal(b, input) {
			//e = EncodeToString(b)
			t.Fail()
		}
	}
}


func BenchmarkEncodeDecode(b *testing.B) {
	c := 64
	input := make([]byte, c)
	rand.Read(input)
	for i := 0; i < b.N; i++ {
		e := EncodeToString(input)
		o, _ := DecodeString(e)
		if !bytes.Equal(o,input) {
			b.Fail()
			break
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	c := 64
	input := make([]byte, c)
	rand.Read(input)
	e := EncodeToString(input)
	for i := 0; i < b.N; i++ {
		o, _ := DecodeString(e)
		if !bytes.Equal(o,input) {
			b.Fail()
			break
		}
	}
}
