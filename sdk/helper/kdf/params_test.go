package kdf

import (
	"bytes"
	"crypto/sha256"
	"hash"
	"testing"
)

func TestCounterVariable(t *testing.T) {
	var x CounterVariable
	x.LittleEndian = false

	x.Width = 0
	if x.Validate() == nil {
		t.Fatal("unexpected pass from CounterVariable.Validate()")
	}

	x.Width = 1
	if x.Validate() == nil {
		t.Fatal("unexpected pass from CounterVariable.Validate()")
	}

	x.Width = 8
	if err := x.Validate(); err != nil {
		t.Fatalf("unexpected failure from CounterVariable.Validate(): %v", err)
	}

	x.Width = 72
	if x.Validate() == nil {
		t.Fatal("unexpected pass from CounterVariable.Validate()")
	}

	type testCase struct {
		LittleEndian bool
		Width        uint8
		Counter      uint64
		Expected     []byte
	}

	var testCases = []testCase{
		testCase{false, 8, 0, []byte{0}},
		testCase{false, 8, 1, []byte{1}},
		testCase{true, 8, 0, []byte{0}},
		testCase{true, 8, 1, []byte{1}},
		testCase{false, 16, 0, []byte{0, 0}},
		testCase{false, 16, 1, []byte{0, 1}},
		testCase{true, 16, 0, []byte{0, 0}},
		testCase{true, 16, 1, []byte{1, 0}},
		testCase{false, 24, 0, []byte{0, 0, 0}},
		testCase{false, 24, 1, []byte{0, 0, 1}},
		testCase{true, 24, 1, []byte{1, 0, 0}},
	}

	for index, test := range testCases {
		x.LittleEndian = test.LittleEndian
		x.Width = test.Width

		var actual = x.Encode(test.Counter)
		if !bytes.Equal(actual, test.Expected) {
			t.Fatalf("test %d failed: got %v, expected %v", index, actual, test.Expected)
		}
	}
}

func TestDKMLengthVariable(t *testing.T) {
	var y DKMLength
	y.Method = 3
	y.Width = 8
	y.LittleEndian = false

	if y.Validate() == nil {
		t.Fatal("unexpected pass from DKMLength.Validate()")
	}

	type testCase struct {
		Method       DKMLengthMethod
		LittleEndian bool
		Width        uint8
		PRF          hash.Hash
		Keys         []int
		Expected     []byte
	}

	var testCases = []testCase{
		testCase{SumOfKeys, false, 16, sha256.New(), []int{5}, []byte{0, 5}},
		testCase{SumOfSegments, false, 16, sha256.New(), []int{5}, []byte{1, 0}},
		testCase{SumOfKeys, true, 16, sha256.New(), []int{5}, []byte{5, 0}},
		testCase{SumOfSegments, true, 16, sha256.New(), []int{5}, []byte{0, 1}},
		testCase{SumOfKeys, false, 16, sha256.New(), []int{5, 5}, []byte{0, 10}},
		testCase{SumOfSegments, false, 16, sha256.New(), []int{5, 5}, []byte{2, 0}},
		testCase{SumOfKeys, false, 16, sha256.New(), []int{255, 255}, []byte{1, 254}},
		testCase{SumOfSegments, false, 16, sha256.New(), []int{255, 255}, []byte{2, 0}},
		testCase{SumOfKeys, false, 16, sha256.New(), []int{256, 256}, []byte{2, 0}},
		testCase{SumOfSegments, false, 16, sha256.New(), []int{256, 256}, []byte{2, 0}},
	}

	for index, test := range testCases {
		y.Method = test.Method
		y.LittleEndian = test.LittleEndian
		y.Width = test.Width

		var actual = y.Encode(test.PRF, test.Keys)
		if !bytes.Equal(actual, test.Expected) {
			t.Fatalf("test %d failed: got %v, expected %v", index, actual, test.Expected)
		}
	}
}
