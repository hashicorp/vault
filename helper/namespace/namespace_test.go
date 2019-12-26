package namespace

import "testing"

func TestSplitIDFromString(t *testing.T) {
	tcases := []struct {
		input  string
		id     string
		prefix string
	}{
		{
			"foo",
			"",
			"foo",
		},
		{
			"foo.id",
			"id",
			"foo",
		},
		{
			"foo.foo.id",
			"id",
			"foo.foo",
		},
		{
			"foo.foo/foo.id",
			"id",
			"foo.foo/foo",
		},
		{
			"foo.foo/.id",
			"id",
			"foo.foo/",
		},
		{
			"foo.foo/foo",
			"",
			"foo.foo/foo",
		},
		{
			"foo.foo/f",
			"",
			"foo.foo/f",
		},
		{
			"foo.foo/",
			"",
			"foo.foo/",
		},
		{
			"b.foo",
			"",
			"b.foo",
		},
		{
			"s.foo",
			"",
			"s.foo",
		},
		{
			"t.foo",
			"foo",
			"t",
		},
	}

	for _, c := range tcases {
		pre, id := SplitIDFromString(c.input)
		if pre != c.prefix || id != c.id {
			t.Fatalf("bad test case: %s != %s, %s != %s", pre, c.prefix, id, c.id)
		}
	}
}
