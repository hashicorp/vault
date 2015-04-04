package gatedwriter

import (
	"bytes"
	"io"
	"testing"
)

func TestWriter_impl(t *testing.T) {
	var _ io.Writer = new(Writer)
}

func TestWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	w := &Writer{Writer: buf}
	w.Write([]byte("foo\n"))
	w.Write([]byte("bar\n"))

	if buf.String() != "" {
		t.Fatalf("bad: %s", buf.String())
	}

	w.Flush()

	if buf.String() != "foo\nbar\n" {
		t.Fatalf("bad: %s", buf.String())
	}

	w.Write([]byte("baz\n"))

	if buf.String() != "foo\nbar\nbaz\n" {
		t.Fatalf("bad: %s", buf.String())
	}
}
