package parseip

import (
	"testing"
)

func Test_TrimLeadingZeroes(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"0", ""},
		{"no leading 0s", "no leading 0s"},
		{"0 but only one", " but only one"},
		{"00 two zeroes", " two zeroes"},
		{"0 0 should trim one", " 0 should trim one"},
	}
	for _, tt := range tests {
		if got := TrimLeadingZeroes(tt.in); got != tt.want {
			t.Errorf("TrimLeadingZeroes() = %v, want %v", got, tt.want)
		}
	}
}
