package password

import "testing"

type testCase struct {
	name     string
	input    string
	expected string
}

func TestRemoveDeletes(t *testing.T) {
	var tests = []testCase{
		{"NoDelete", "TestingStuff", "TestingStuff"},
		{"SingleDelete", "Testing\x7fStuff", "TestinStuff"},
		{"DeleteFirst", "\x7fTestingStuff", "TestingStuff"},
		{"DoubleDelete", "Testing\x7f\x7fStuff", "TestiStuff"},
		{"LastDelete", "TestingStuff\x7f", "TestingStuf"},
	}

	for _, test := range tests {
		result := removeDeletes(test.input)
		if result != test.expected {
			t.Errorf("Test %s failed, input: '%s', expected: '%s', output: '%s'", test.name, test.input, test.expected, result)
		}
	}
}
