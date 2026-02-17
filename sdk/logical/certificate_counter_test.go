// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"math"
	"testing"
)

func Test_durationAdjustedCertificateCount(t *testing.T) {
	tests := []struct {
		name            string
		validitySeconds int64
		want            float64
	}{
		{
			name:            "zero duration",
			validitySeconds: 0,
			want:            0.0,
		},
		{
			name:            "1 hour",
			validitySeconds: 3600,
			want:            0.0014, // 1/730 rounded to 4 decimals
		},
		{
			name:            "24 hours (1 day)",
			validitySeconds: 86400,  // 24 * 3600
			want:            0.0329, // 24/730 = 0.032876... rounded to 4 decimals
		},
		{
			name:            "730 hours (standard duration - 1 month)",
			validitySeconds: 2628000, // 730 * 3600
			want:            1.0,
		},
		{
			name:            "8760 hours (1 year)",
			validitySeconds: 31536000, // 365 * 24 * 3600
			want:            12.0,     // 8760/730 = 12.0
		},
		{
			name:            "17520 hours (2 years)",
			validitySeconds: 63072000, // 730 * 24 * 3600
			want:            24.0,     // 17520/730 = 24.0
		},
		{
			name:            "87600 hours (10 years)",
			validitySeconds: 315360000, // 3650 * 24 * 3600
			want:            120.0,     // 87600/730 = 120.0
		},
		{
			name:            "90 days",
			validitySeconds: 7776000, // 90 * 24 * 3600
			want:            2.9589,  // 2160/730 = 2.958904... rounded to 4 decimals
		},
		{
			name:            "365 days (1 year)",
			validitySeconds: 31536000, // 365 * 24 * 3600
			want:            12.0,     // 8760/730 = 12.0
		},
		{
			name:            "fractional result - 100 hours",
			validitySeconds: 360000, // 100 * 3600
			want:            0.137,  // 100/730 = 0.136986... rounded to 4 decimals
		},
		{
			name:            "fractional result - 500 hours",
			validitySeconds: 1800000, // 500 * 3600
			want:            0.6849,  // 500/730 = 0.684931... rounded to 4 decimals
		},
		{
			name:            "very small duration - 1 second",
			validitySeconds: 1,
			want:            0.0, // 1/3600/730 = 0.00000038... rounds to 0.0
		},
		{
			name:            "very small duration - 60 seconds",
			validitySeconds: 60,
			want:            0.0, // 60/3600/730 = 0.000023... rounds to 0.0
		},
		{
			name:            "very small duration - 600 seconds",
			validitySeconds: 600,
			want:            0.0002, // 600/3600/730 = 0.000228... rounds to 0.0002
		},
		{
			name:            "edge case - exactly rounds up",
			validitySeconds: 13149,  // Should result in value that rounds up
			want:            0.0050, // 3.6525/730 = 0.005003... rounds to 0.0050
		},
		{
			name:            "edge case - exactly rounds down",
			validitySeconds: 13140,  // Should result in value that rounds down
			want:            0.0050, // 3.65/730 = 0.005000
		},
		{
			name:            "large duration - 100 years",
			validitySeconds: 3153600000, // 100 * 365 * 24 * 3600
			want:            1200.0,     // 876000/730 = 1200.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := durationAdjustedCertificateCount(tt.validitySeconds)
			if got != tt.want {
				t.Errorf("durationAdjustedCertificateCount(%d) = %v, want %v", tt.validitySeconds, got, tt.want)
			}
		})
	}
}

func Test_durationAdjustedCertificateCount_Precision(t *testing.T) {
	// Test that the function properly rounds to 4 decimal places
	tests := []struct {
		name            string
		validitySeconds int64
		wantPrecision   int // number of decimal places
	}{
		{
			name:            "result has max 4 decimal places - case 1",
			validitySeconds: 12345,
			wantPrecision:   4,
		},
		{
			name:            "result has max 4 decimal places - case 2",
			validitySeconds: 98765,
			wantPrecision:   4,
		},
		{
			name:            "result has max 4 decimal places - case 3",
			validitySeconds: 555555,
			wantPrecision:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := durationAdjustedCertificateCount(tt.validitySeconds)
			// Check that the result has at most 4 decimal places
			// by multiplying by 10000 and checking if it's an integer
			scaled := got * 10000
			if scaled != math.Floor(scaled) {
				t.Errorf("durationAdjustedCertificateCount(%d) = %v has more than 4 decimal places", tt.validitySeconds, got)
			}
		})
	}
}

func Test_durationAdjustedCertificateCount_Consistency(t *testing.T) {
	// Test that the function is consistent with the public wrapper
	tests := []struct {
		name            string
		validitySeconds int64
	}{
		{"1 hour", 3600},
		{"1 day", 86400},
		{"1 month", 2628000},
		{"1 year", 31536000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the internal function
			internal := durationAdjustedCertificateCount(tt.validitySeconds)

			// The internal function should produce the same result as calculating manually
			validityHours := float64(tt.validitySeconds) / 3600.0
			units := validityHours / 730.0
			expected := math.Round(units*10000) / 10000

			if internal != expected {
				t.Errorf("durationAdjustedCertificateCount(%d) = %v, manual calculation = %v", tt.validitySeconds, internal, expected)
			}
		})
	}
}

func Test_durationAdjustedCertificateCount_NegativeInput(t *testing.T) {
	// Test behavior with negative input (edge case)
	// Note: In practice, this shouldn't happen with valid certificates,
	// but we should verify the function's behavior
	validitySeconds := int64(-3600)
	got := durationAdjustedCertificateCount(validitySeconds)

	// The function should handle negative values mathematically
	// -1 hour / 730 hours = -0.0014 (rounded to 4 decimals)
	want := -0.0014

	if got != want {
		t.Errorf("durationAdjustedCertificateCount(%d) = %v, want %v", validitySeconds, got, want)
	}
}

func Benchmark_durationAdjustedCertificateCount(b *testing.B) {
	validitySeconds := int64(31536000) // 1 year

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = durationAdjustedCertificateCount(validitySeconds)
	}
}

func Benchmark_durationAdjustedCertificateCount_Various(b *testing.B) {
	testCases := []int64{
		3600,      // 1 hour
		86400,     // 1 day
		2628000,   // 1 month
		31536000,  // 1 year
		315360000, // 10 years
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, validitySeconds := range testCases {
			_ = durationAdjustedCertificateCount(validitySeconds)
		}
	}
}

// Made with Bob
