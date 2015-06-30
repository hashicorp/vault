package gocql

import (
	"bytes"
	"testing"
)

func TestFuzzBugs(t *testing.T) {
	// these inputs are found using go-fuzz (https://github.com/dvyukov/go-fuzz)
	// and should cause a panic unless fixed.
	tests := [][]byte{
		[]byte("00000\xa0000"),
		[]byte("\x8000\x0e\x00\x00\x00\x000"),
		[]byte("\x8000\x00\x00\x00\x00\t0000000000"),
		[]byte("\xa0\xff\x01\xae\xefqE\xf2\x1a"),
		[]byte("\x8200\b\x00\x00\x00c\x00\x00\x00\x02000\x01\x00\x00\x00\x03" +
			"\x00\n0000000000\x00\x14000000" +
			"00000000000000\x00\x020000" +
			"\x00\a000000000\x00\x050000000" +
			"\xff0000000000000000000" +
			"0000000"),
		[]byte("\x82\xe600\x00\x00\x00\x000"),
		[]byte("\x8200\b\x00\x00\x00\b0\x00\x00\x00\x040000"),
		[]byte("\x8200\x00\x00\x00\x00\x100\x00\x00\x12\x00\x00\x0000000" +
			"00000"),
		[]byte("\x83000\b\x00\x00\x00\x14\x00\x00\x00\x020000000" +
			"000000000"),
		[]byte("\x83000\b\x00\x00\x000\x00\x00\x00\x04\x00\x1000000" +
			"00000000000000e00000" +
			"000\x800000000000000000" +
			"0000000000000"),
	}

	for i, test := range tests {
		t.Logf("test %d input: %q", i, test)

		var bw bytes.Buffer

		r := bytes.NewReader(test)
		head, err := readHeader(r, make([]byte, 9))
		if err != nil {
			continue
		}

		framer := newFramer(r, &bw, nil, byte(head.version))
		err = framer.readFrame(&head)
		if err != nil {
			continue
		}

		_, err = framer.parseFrame()
		if err != nil {
			continue
		}

		t.Errorf("(%d) expected to fail for input %q", i, test)
	}
}

func TestFrameWriteTooLong(t *testing.T) {
	w := &bytes.Buffer{}
	framer := newFramer(nil, w, nil, 2)

	framer.writeHeader(0, opStartup, 1)
	framer.writeBytes(make([]byte, maxFrameSize+1))
	err := framer.finishWrite()
	if err != ErrFrameTooBig {
		t.Fatalf("expected to get %v got %v", ErrFrameTooBig, err)
	}
}

func TestFrameReadTooLong(t *testing.T) {
	r := &bytes.Buffer{}
	r.Write(make([]byte, maxFrameSize+1))
	// write a new header right after this frame to verify that we can read it
	r.Write([]byte{0x02, 0x00, 0x00, opReady, 0x00, 0x00, 0x00, 0x00})

	framer := newFramer(r, nil, nil, 2)

	head := frameHeader{
		version: 2,
		op:      opReady,
		length:  r.Len() - 8,
	}

	err := framer.readFrame(&head)
	if err != ErrFrameTooBig {
		t.Fatalf("expected to get %v got %v", ErrFrameTooBig, err)
	}

	head, err = readHeader(r, make([]byte, 8))
	if err != nil {
		t.Fatal(err)
	}
	if head.op != opReady {
		t.Fatalf("expected to get header %v got %v", opReady, head.op)
	}
}
