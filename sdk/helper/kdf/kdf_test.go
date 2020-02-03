package kdf

import (
	"bytes"
	"testing"
)

func TestCounterMode(t *testing.T) {
	key := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	context := []byte("the quick brown fox")
	prf := HMACSHA256PRF
	prfLen := HMACSHA256PRFLen

	// Expect256 was generated in python with
	// import hashlib, hmac
	// hash = hashlib.sha256
	// context = "the quick brown fox"
	// key = "".join([chr(x) for x in range(1, 17)])
	// inp = "\x00\x00\x00\x00"+context+"\x00\x00\x01\x00"
	// digest = hmac.HMAC(key, inp, hash).digest()
	// print [ord(x) for x in digest]
	expect256 := []byte{219, 25, 238, 6, 185, 236, 180, 64, 248, 152, 251,
		153, 79, 5, 141, 222, 66, 200, 66, 143, 40, 3, 101, 221, 206, 163, 102,
		80, 88, 234, 87, 157}

	for _, l := range []uint32{128, 256, 384, 1024} {
		out, err := CounterMode(prf, prfLen, key, context, l)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if uint32(len(out)*8) != l {
			t.Fatalf("bad length: %#v", out)
		}

		if bytes.Contains(out, key) {
			t.Fatalf("output contains key")
		}

		if l == 256 && !bytes.Equal(out, expect256) {
			t.Fatalf("mis-match")
		}
	}

}

func TestHMACSHA256PRF(t *testing.T) {
	key := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	data := []byte("foobarbaz")
	out, err := HMACSHA256PRF(key, data)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if uint32(len(out)*8) != HMACSHA256PRFLen {
		t.Fatalf("Bad len")
	}

	// Expect was generated in python with:
	// import hashlib, hmac
	// hash = hashlib.sha256
	// msg = "foobarbaz"
	// key = "".join([chr(x) for x in range(1, 17)])
	// hm = hmac.HMAC(key, msg, hash)
	// print [ord(x) for x in hm.digest()]
	expect := []byte{9, 50, 146, 8, 188, 130, 150, 107, 205, 147, 82, 170,
		253, 183, 26, 38, 167, 194, 220, 111, 56, 118, 219, 209, 31, 52, 137,
		90, 246, 133, 191, 124}
	if !bytes.Equal(expect, out) {
		t.Fatalf("mis-matched output")
	}
}
