package sliceflag

import (
	"flag"
	"reflect"
	"testing"
)

func TestStringFlag_implements(t *testing.T) {
	var raw interface{}
	raw = new(StringFlag)
	if _, ok := raw.(flag.Value); !ok {
		t.Fatalf("StringFlag should be a Value")
	}
}

func TestStringFlagSet(t *testing.T) {
	sv := new(StringFlag)
	err := sv.Set("foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = sv.Set("bar")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual([]string(*sv), expected) {
		t.Fatalf("Bad: %#v", sv)
	}
}
