package audit

import (
	"reflect"
	"testing"
)

func TestHashWalker(t *testing.T) {
	replaceText := "foo"

	cases := []struct {
		Input  interface{}
		Output interface{}
	}{
		{
			map[string]interface{}{
				"hello": "foo",
			},
			map[string]interface{}{
				"hello": replaceText,
			},
		},

		{
			map[string]interface{}{
				"hello": []interface{}{"world"},
			},
			map[string]interface{}{
				"hello": []interface{}{replaceText},
			},
		},
	}

	for _, tc := range cases {
		output, err := HashStructure(tc.Input, func(string) (string, error) {
			return replaceText, nil
		})
		if err != nil {
			t.Fatalf("err: %s\n\n%#v", err, tc.Input)
		}
		if !reflect.DeepEqual(output, tc.Output) {
			t.Fatalf("bad:\n\n%#v\n\n%#v", tc.Input, output)
		}
	}
}

func TestHashSHA1(t *testing.T) {
	fn := HashSHA1("")
	result, err := fn("foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if result != "sha1:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33" {
		t.Fatalf("bad: %#v", result)
	}
}
