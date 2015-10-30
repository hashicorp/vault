// +build all unit

package gocql

import (
	"bytes"
	"testing"

	"github.com/golang/snappy"
)

func TestSnappyCompressor(t *testing.T) {
	c := SnappyCompressor{}
	if c.Name() != "snappy" {
		t.Fatalf("expected name to be 'snappy', got %v", c.Name())
	}

	str := "My Test String"
	//Test Encoding
	expected := snappy.Encode(nil, []byte(str))
	if res, err := c.Encode([]byte(str)); err != nil {
		t.Fatalf("failed to encode '%v' with error %v", str, err)
	} else if bytes.Compare(expected, res) != 0 {
		t.Fatal("failed to match the expected encoded value with the result encoded value.")
	}

	val, err := c.Encode([]byte(str))
	if err != nil {
		t.Fatalf("failed to encode '%v' with error '%v'", str, err)
	}

	//Test Decoding
	if expected, err := snappy.Decode(nil, val); err != nil {
		t.Fatalf("failed to decode '%v' with error %v", val, err)
	} else if res, err := c.Decode(val); err != nil {
		t.Fatalf("failed to decode '%v' with error %v", val, err)
	} else if bytes.Compare(expected, res) != 0 {
		t.Fatal("failed to match the expected decoded value with the result decoded value.")
	}
}
