package duration

import "testing"

func Test_ParseDurationSecond(t *testing.T) {
	outp, err := ParseDurationSecond("9876s")
	if err != nil {
		t.Fatal(err)
	}
	if outp != 9876 {
		t.Fatal("not equivalent")
	}
	outp, err = ParseDurationSecond("9876")
	if err != nil {
		t.Fatal(err)
	}
	if outp != 9876 {
		t.Fatal("not equivalent")
	}
}
