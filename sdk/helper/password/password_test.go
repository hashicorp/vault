package password

import "testing"

type testCase struct {
	name     string
	input    string
	expected string
}

func TestRemoveiTermDelete(t *testing.T) {
	var tests = []testCase{
		{"NoDelete", "TestingStuff", "TestingStuff"},
		{"SingleDelete", "Testing\x7fStuff", "Testing\x7fStuff"},
		{"DeleteFirst", "\x7fTestingStuff", "\x7fTestingStuff"},
		{"DoubleDelete", "\x7f\x7fTestingStuff", "\x7f\x7fTestingStuff"},
		{"SpaceFirst", "\x20TestingStuff", "\x20TestingStuff"},
		{"iTermDelete", "\x20\x7fTestingStuff", "TestingStuff"},
	}

	for _, test := range tests {
		result := removeiTermDelete(test.input)
		if result != test.expected {
			t.Errorf("Test %s failed, input: '%s', expected: '%s', output: '%s'", test.name, test.input, test.expected, result)
		}
	}
}
