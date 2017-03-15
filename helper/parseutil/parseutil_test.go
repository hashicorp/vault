package parseutil

import (
	"encoding/json"
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
	outp, err = ParseDurationSecond(json.Number("4352"))
	if err != nil {
		t.Fatal(err)
	}
	if outp != time.Duration(4352)*time.Second {
		t.Fatal("not equivalent")
	}
}

func Test_ParseBool(t *testing.T) {
	outp, err := ParseBool("true")
	if err != nil {
		t.Fatal(err)
	}
	if !outp {
		t.Fatal("wrong output")
	}
	outp, err = ParseBool(1)
	if err != nil {
		t.Fatal(err)
	}
	if !outp {
		t.Fatal("wrong output")
	}
	outp, err = ParseBool(true)
	if err != nil {
		t.Fatal(err)
	}
	if !outp {
		t.Fatal("wrong output")
	}
}
