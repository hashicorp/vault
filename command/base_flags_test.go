package command

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_BoolPtr(t *testing.T) {
	var boolPtr BoolPtr
	value := newBoolPtrValue(nil, &boolPtr, false)

	require.False(t, boolPtr.IsSet())
	require.False(t, boolPtr.Get())

	err := value.Set("false")
	require.NoError(t, err)

	require.True(t, boolPtr.IsSet())
	require.False(t, boolPtr.Get())

	err = value.Set("true")
	require.NoError(t, err)

	require.True(t, boolPtr.IsSet())
	require.True(t, boolPtr.Get())

	var boolPtrFalseDefault BoolPtr
	_ = newBoolPtrValue(new(bool), &boolPtrFalseDefault, false)

	require.True(t, boolPtrFalseDefault.IsSet())
	require.False(t, boolPtrFalseDefault.Get())

	var boolPtrTrueDefault BoolPtr
	defTrue := true
	_ = newBoolPtrValue(&defTrue, &boolPtrTrueDefault, false)

	require.True(t, boolPtrTrueDefault.IsSet())
	require.True(t, boolPtrTrueDefault.Get())

	var boolPtrHidden BoolPtr
	value = newBoolPtrValue(nil, &boolPtrHidden, true)
	require.Equal(t, true, value.Hidden())
}

func Test_TimeParsing(t *testing.T) {
	var zeroTime time.Time

	testCases := []struct {
		Input    string
		Formats  TimeFormat
		Valid    bool
		Expected time.Time
	}{
		{
			"2020-08-24",
			TimeVar_TimeOrDay,
			true,
			time.Date(2020, 8, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			"2099-09",
			TimeVar_TimeOrDay,
			false,
			zeroTime,
		},
		{
			"2099-09",
			TimeVar_TimeOrDay | TimeVar_Month,
			true,
			time.Date(2099, 9, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			"2021-01-02T03:04:05-02:00",
			TimeVar_TimeOrDay,
			true,
			time.Date(2021, 1, 2, 5, 4, 5, 0, time.UTC),
		},
		{
			"2021-01-02T03:04:05",
			TimeVar_TimeOrDay,
			false, // Missing timezone not supported
			time.Date(2021, 1, 2, 3, 4, 5, 0, time.UTC),
		},
		{
			"2021-01-02T03:04:05+02:00",
			TimeVar_TimeOrDay,
			true,
			time.Date(2021, 1, 2, 1, 4, 5, 0, time.UTC),
		},
		{
			"1598313593",
			TimeVar_TimeOrDay,
			true,
			time.Date(2020, 8, 24, 23, 59, 53, 0, time.UTC),
		},
		{
			"2037",
			TimeVar_TimeOrDay,
			false,
			zeroTime,
		},
		{
			"20201212",
			TimeVar_TimeOrDay,
			false,
			zeroTime,
		},
		{
			"9999999999999999999999999999999999999999999999",
			TimeVar_TimeOrDay,
			false,
			zeroTime,
		},
		{
			"2021-13-02T03:04:05-02:00",
			TimeVar_TimeOrDay,
			false,
			zeroTime,
		},
		{
			"2021-12-02T24:04:05+00:00",
			TimeVar_TimeOrDay,
			false,
			zeroTime,
		},
		{
			"2021-01-02T03:04:05.234567890Z",
			TimeVar_TimeOrDay,
			true,
			time.Date(2021, 1, 2, 3, 4, 5, 234567890, time.UTC),
		},
	}

	for _, tc := range testCases {
		var result time.Time
		timeVal := newTimeValue(zeroTime, &result, false, tc.Formats)
		err := timeVal.Set(tc.Input)
		if err == nil && !tc.Valid {
			t.Errorf("Time %q parsed without error as %v, but is not valid", tc.Input, result)
			continue
		}
		if err != nil {
			if tc.Valid {
				t.Errorf("Time %q parsed as error, but is valid", tc.Input)
			}
			continue
		}
		if !tc.Expected.Equal(result) {
			t.Errorf("Time %q parsed incorrectly, expected %v but got %v", tc.Input, tc.Expected.UTC(), result.UTC())
		}
	}
}
