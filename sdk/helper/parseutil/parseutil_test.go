package parseutil

import (
	"encoding/json"
	"testing"
	"time"
)

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
