package logical

import (
	"reflect"
	"testing"
)

func TestWAL(t *testing.T) {
	s := new(InmemStorage)

	// WAL should be empty to start
	keys, err := ListWAL(s)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(keys) > 0 {
		t.Fatalf("bad: %#v", keys)
	}

	// Write an entry to the WAL
	id, err := PutWAL(s, "bar")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// The key should be in the WAL
	keys, err = ListWAL(s)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !reflect.DeepEqual(keys, []string{id}) {
		t.Fatalf("bad: %#v", keys)
	}

	// Should be able to get the value
	v, err := GetWAL(s, id)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Should be able to delete the value
	if err := DeleteWAL(s, id); err != nil {
		t.Fatalf("err: %s", err)
	}
	v, err = GetWAL(s, id)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v != nil {
		t.Fatalf("bad: %#v", v)
	}
}
