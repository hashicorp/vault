package duration

import (
	"testing"
	"time"
)

func Test_ParseDurationSecond(t *testing.T) {
	outp, err := ParseDurationSecond("9876s")
	if err != nil {
		t.Fatal(err)
	}
	if outp != time.Duration(9876)*time.Second {
		t.Fatal("not equivalent")
	}
	outp, err = ParseDurationSecond("9876")
	if err != nil {
		t.Fatal(err)
	}
	if outp != time.Duration(9876)*time.Second {
		t.Fatal("not equivalent")
	}
}
