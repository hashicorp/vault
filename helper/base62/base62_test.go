package base62

import (
	"testing"
)

func TestValid(t *testing.T) {
	tCases := []struct {
		in  string
		out string
	}{
		{
			"",
			"0",
		},
		{
			"foo",
			"sapp",
		},
		{
			"5d5746d044b9a9429249966c9e3fee178ca679b91487b11d4b73c9865202104c",
			"cozMP2pOYdDiNGeFQ2afKAOGIzO0HVpJ8OPFXuVPNbHasFyenK9CzIIPuOG7EFWOCy4YWvKGZa671N4kRSoaxZ",
		},
		{
			"5ba33e16d742f3c785f6e7e8bb6f5fe82346ffa1c47aa8e95da4ddd5a55bb334",
			"cotpEJPnhuTRofLi4lDe5iKw2fkSGc6TpUYeuWoBp8eLYJBWLRUVDZI414OjOCWXKZ0AI8gqNMoxd4eLOklwYk",
		},
		{
			" ",
			"w",
		},
		{
			"-",
			"J",
		},
		{
			"0",
			"M",
		},
		{
			"1",
			"N",
		},
		{
			"-1",
			"30B",
		},
		{
			"11",
			"3h7",
		},
		{
			"abc",
			"qMin",
		},
		{
			"1234598760",
			"1a0AFzKIPnihTq",
		},
		{
			"abcdefghijklmnopqrstuvwxyz",
			"hUBXsgd3F2swSlEgbVi2p0Ncr6kzVeJTLaW",
		},
	}

	for _, c := range tCases {
		e := Encode([]byte(c.in))
		d := string(Decode(e))

		if d != c.in {
			t.Fatalf("decoded value didn't match input %#v %#v", c.in, d)
		}

		if e != c.out {
			t.Fatalf("encoded value didn't match expected %#v, %#v", e, c.out)
		}
	}
}

func TestInvalid(t *testing.T) {
	d := Decode("!0000/")
	if len(d) != 0 {
		t.Fatalf("Decode of invalid string should be empty, got %#v", d)
	}
}

func TestRandom(t *testing.T) {
	a, err1 := Random(16, true)
	b, err2 := Random(16, true)

	if err1 != nil || err2 != nil {
		t.Fatalf("Unexpected errors: %v, %v", err1, err2)
	}

	if a == b {
		t.Fatalf("Expected different random values. Got duplicate: %s", a)
	}

	for i := 0; i < 3000; i++ {
		c, _ := Random(i, true)
		if len(c) != i {
			t.Fatalf("Expected length %d, got: %d", i, len(c))
		}
	}

	d, _ := Random(100, false)
	if len(d) < 133 || len(d) > 135 {
		t.Fatalf("Expected length 133-135, got: %d", len(d))
	}
}
