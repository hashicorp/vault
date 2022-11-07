package logging

import (
	"errors"
	"reflect"
	"testing"
)

func Test_ParseLogFormat(t *testing.T) {
	type testData struct {
		format      string
		expected    LogFormat
		expectedErr error
	}

	tests := []testData{
		{format: "", expected: UnspecifiedFormat, expectedErr: nil},
		{format: " ", expected: UnspecifiedFormat, expectedErr: nil},
		{format: "standard", expected: StandardFormat, expectedErr: nil},
		{format: "STANDARD", expected: StandardFormat, expectedErr: nil},
		{format: "json", expected: JSONFormat, expectedErr: nil},
		{format: " json ", expected: JSONFormat, expectedErr: nil},
		{format: "bogus", expected: UnspecifiedFormat, expectedErr: errors.New("Unknown log format: bogus")},
	}

	for _, test := range tests {
		result, err := ParseLogFormat(test.format)
		if test.expected != result {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
		if !reflect.DeepEqual(test.expectedErr, err) {
			t.Errorf("expected error %v, got %v", test.expectedErr, err)
		}
	}
}
