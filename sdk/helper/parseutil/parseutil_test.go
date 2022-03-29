package parseutil

import (
	"encoding/json"
	"math/cmplx"
	"testing"
	"time"
)

func Test_ParseCapacityString(t *testing.T) {
	testCases := []struct {
		name     string
		inp      interface{}
		valid    bool
		expected uint64
	}{
		{
			"bare number as an int",
			5,
			true,
			uint64(5),
		},
		{
			"bare number as a float",
			5.0,
			true,
			uint64(5),
		},
		{
			"bare number as a string",
			"5",
			true,
			uint64(5),
		},
		{
			"string",
			"haha",
			false,
			uint64(0),
		},
		{
			"random data structure",
			struct{}{},
			false,
			uint64(0),
		},
		{
			"kb",
			"5kb",
			true,
			uint64(5000),
		},
		{
			"kib",
			"5kib",
			true,
			uint64(5120),
		},
		{
			"KB",
			"5KB",
			true,
			uint64(5000),
		},
		{
			"KIB",
			"5KIB",
			true,
			uint64(5120),
		},
		{
			"kB",
			"5kB",
			true,
			uint64(5000),
		},
		{
			"Kb",
			"5Kb",
			true,
			uint64(5000),
		},
		{
			"space kb",
			"5 kb",
			true,
			uint64(5000),
		},
		{
			"space KB",
			"5 KB",
			true,
			uint64(5000),
		},
		{
			"kb surrounding spaces",
			" 5 kb ",
			true,
			uint64(5000),
		},
		{
			"mb",
			"5mb",
			true,
			uint64(5000000),
		},
		{
			"mib",
			"5mib",
			true,
			uint64(5242880),
		},
		{
			"gb",
			"5gb",
			true,
			uint64(5000000000),
		},
		{
			"gib",
			"5gib",
			true,
			uint64(5368709120),
		},
		{
			"tb",
			"5tb",
			true,
			uint64(5000000000000),
		},
		{
			"tib",
			"5tib",
			true,
			uint64(5497558138880),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			outp, err := ParseCapacityString(tc.inp)
			if tc.valid && err != nil {
				t.Errorf("failed to parse: %v. err: %v", tc.inp, err)
			}
			if !tc.valid && err == nil {
				t.Errorf("no error for: %v", tc.inp)
			}
			if outp != tc.expected {
				t.Errorf("input %v parsed as %v, expected %v", tc.inp, outp, tc.expected)
			}
		})
	}
}

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

func Test_ParseAbsoluteTime(t *testing.T) {
	testCases := []struct {
		inp      interface{}
		valid    bool
		expected time.Time
	}{
		{
			"2020-12-11T09:08:07.654321Z",
			true,
			time.Date(2020, 12, 11, 9, 8, 7, 654321000, time.UTC),
		},
		{
			"2020-12-11T09:08:07+02:00",
			true,
			time.Date(2020, 12, 11, 7, 8, 7, 0, time.UTC),
		},
		{
			"2021-12-11T09:08:07Z",
			true,
			time.Date(2021, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			"2021-12-11T09:08:07",
			false,
			time.Time{},
		},
		{
			"1670749687",
			true,
			time.Date(2022, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			1670749687,
			true,
			time.Date(2022, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			uint32(1670749687),
			true,
			time.Date(2022, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			json.Number("1670749687"),
			true,
			time.Date(2022, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			nil,
			true,
			time.Time{},
		},
		{
			struct{}{},
			false,
			time.Time{},
		},
		{
			true,
			false,
			time.Time{},
		},
	}
	for _, tc := range testCases {
		outp, err := ParseAbsoluteTime(tc.inp)
		if err != nil {
			if tc.valid {
				t.Errorf("failed to parse: %v", tc.inp)
			}
			continue
		}
		if err == nil && !tc.valid {
			t.Errorf("no error for: %v", tc.inp)
			continue
		}
		if !outp.Equal(tc.expected) {
			t.Errorf("input %v parsed as %v, expected %v", tc.inp, outp, tc.expected)
		}
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

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func Test_ParseCommaStringSlice(t *testing.T) {
	cases := []struct {
		name     string
		inp      interface{}
		expected []string
		valid    bool
	}{
		{
			"nil",
			nil,
			[]string{},
			true,
		},
		{
			"empty string",
			"",
			[]string{},
			true,
		},
		{
			"string without commas",
			"foo",
			[]string{"foo"},
			true,
		},
		{
			"comma-separated string",
			"foo,bar,baz",
			[]string{"foo", "bar", "baz"},
			true,
		},
		{
			"comma-separated string with trim",
			"  foo ,    bar   ,baz  ",
			[]string{"foo", "bar", "baz"},
			true,
		},
		{
			"json number",
			json.Number("123"),
			[]string{"123"},
			true,
		},
		{
			"int",
			1,
			[]string{"1"},
			true,
		},
		{
			"float",
			5.5,
			[]string{"5.5"},
			true,
		},
		{
			"rune",
			'a',
			[]string{"97"},
			true,
		},
		{
			"bool",
			true,
			[]string{"1"},
			true,
		},
		{
			"byte",
			byte(10),
			[]string{"10"},
			true,
		},
		{
			"complex",
			cmplx.Sqrt(-1),
			nil,
			false,
		},
		{
			"time",
			time.Now(),
			nil,
			false,
		},
		{
			"string slice",
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "bar", "baz"},
			true,
		},
		{
			"json number slice",
			[]json.Number{json.Number("1"), json.Number("2")},
			[]string{"1", "2"},
			true,
		},
		{
			"int slice",
			[]int{1, 2, 3},
			[]string{"1", "2", "3"},
			true,
		},
		{
			"float slice",
			[]float64{1.1, 1.2, 1.3},
			[]string{"1.1", "1.2", "1.3"},
			true,
		},
		{
			"rune slice",
			[]rune{'a', 'b', 'c'},
			[]string{"97", "98", "99"},
			true,
		},
		{
			"bool slice",
			[]bool{true, false, true},
			[]string{"1", "0", "1"},
			true,
		},
		{
			"complex slice",
			[]complex128{cmplx.Sqrt(-1)},
			nil,
			false,
		},
		{
			"map",
			map[string]interface{}{"foo": "bar"},
			nil,
			false,
		},
		{
			"struct",
			struct{ name string }{"foo"},
			nil,
			false,
		},
	}

	for _, tc := range cases {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			outp, err := ParseCommaStringSlice(tc.inp)

			if tc.valid && err != nil {
				t.Errorf("failed to parse: %v. err: %v", tc.inp, err)
			}

			if !tc.valid && err == nil {
				t.Errorf("no error for: %v", tc.inp)
			}

			if !equalStringSlice(outp, tc.expected) {
				t.Errorf("input %v parsed as %v, expected %v", tc.inp, outp, tc.expected)
			}
		})
	}
}
