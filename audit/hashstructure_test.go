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
