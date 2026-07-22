// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"crypto/x509"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
			want:            0, // If the duration is zero, the normalized unit should be zero
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
			want:            0.0001, // 1/3600/730 = 0.00000038... rounds to 0.0 but should return default minimum 0.0001
		},
		{
			name:            "very small duration - 60 seconds",
			validitySeconds: 60,
			want:            0.0001, // 60/3600/730 = 0.000023... rounds to 0.0 and should return default minimum 0.0001
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

// makeCert creates a test certificate with the given validity.
func makeCert(t *testing.T, validity time.Duration) *x509.Certificate {
	t.Helper()

	// AddIssuedCertificate only uses NotBefore/NotAfter, so a minimal x509.Certificate is sufficient.
	notBefore := time.Now()
	return &x509.Certificate{
		NotBefore: notBefore,
		NotAfter:  notBefore.Add(validity),
	}
}

// TestWithMountInfo_PKI verifies that WithMountInfo followed by AddIssuedCertificate
// populates PkiMountAttributions with the correct count (duration-adjusted units) and
// metadata, while leaving SshCertMountAttributions and SshOtpMountAttributions empty.
func TestWithMountInfo_PKI(t *testing.T) {
	recording := &TestCertificateCounter{}

	mount := MountAttribution{
		MountAccessor:    "pki_abc123",
		MountPath:        "pki/",
		MountType:        "pki",
		NamespaceID:      "root",
		NamespacePath:    "",
		BackendAwareUUID: "pki-backend-uuid-001",
	}

	// 730 h = 1 standard duration = 1.0000 billing unit
	cert := makeCert(t, 730*time.Hour)
	NewCertCountIncrementer(recording).WithMountInfo(mount).AddIssuedCertificate(true, cert)

	require.Equal(t, uint64(1), recording.Record.IssuedCerts)
	require.Equal(t, uint64(1), recording.Record.StoredCerts)
	require.InDelta(t, 1.0, recording.Record.PkiDurationAdjustedCerts, 0.0001)

	require.Len(t, recording.Record.PkiMountAttributions, 1)
	attr := recording.Record.PkiMountAttributions[mount.MountAccessor]
	require.Equal(t, mount.MountAccessor, attr.MountAccessor)
	require.Equal(t, mount.MountPath, attr.MountPath)
	require.Equal(t, mount.MountType, attr.MountType)
	require.Equal(t, mount.NamespaceID, attr.NamespaceID)
	require.Equal(t, mount.BackendAwareUUID, attr.BackendAwareUUID)
	require.InDelta(t, 1.0, attr.Count.(float64), 0.0001, "per-mount count should equal billing units")

	require.Empty(t, recording.Record.SshCertMountAttributions)
	require.Empty(t, recording.Record.SshOtpMountAttributions)
}

// TestWithMountInfo_SSH verifies that WithMountInfo followed by AddSSHCertificate
// populates SshCertMountAttributions.
func TestWithMountInfo_SSH(t *testing.T) {
	recording := &TestCertificateCounter{}

	mount := MountAttribution{
		MountAccessor:    "ssh_xyz789",
		MountPath:        "ssh/",
		MountType:        "ssh",
		NamespaceID:      "ns1",
		NamespacePath:    "ns1/",
		BackendAwareUUID: "ssh-backend-uuid-001",
	}

	// 730 h = 1.0 billing unit
	NewCertCountIncrementer(recording).WithMountInfo(mount).AddSSHCertificate(730 * time.Hour)

	require.InDelta(t, 1.0, recording.Record.SSHIssuedCerts, 0.0001)
	require.Len(t, recording.Record.SshCertMountAttributions, 1)
	attr := recording.Record.SshCertMountAttributions[mount.MountAccessor]
	require.Equal(t, mount.MountAccessor, attr.MountAccessor)
	require.Equal(t, mount.MountPath, attr.MountPath)
	require.Equal(t, mount.MountType, attr.MountType)
	require.Equal(t, mount.NamespaceID, attr.NamespaceID)
	require.Equal(t, mount.NamespacePath, attr.NamespacePath)
	require.Equal(t, mount.BackendAwareUUID, attr.BackendAwareUUID)
	require.InDelta(t, 1.0, attr.Count.(float64), 0.0001)

	require.Empty(t, recording.Record.PkiMountAttributions)
	require.Empty(t, recording.Record.SshOtpMountAttributions)
}

// TestWithMountInfo_SSHOTP verifies that WithMountInfo followed by AddSSHOTP
// populates SshOtpMountAttributions with the fixed OTP unit value.
func TestWithMountInfo_SSHOTP(t *testing.T) {
	recording := &TestCertificateCounter{}

	mount := MountAttribution{
		MountAccessor:    "ssh_otp_001",
		MountPath:        "ssh/",
		MountType:        "ssh",
		NamespaceID:      "root",
		BackendAwareUUID: "ssh-backend-uuid-002",
	}

	NewCertCountIncrementer(recording).WithMountInfo(mount).AddSSHOTP()

	require.InDelta(t, 0.0014, recording.Record.SSHIssuedOTPs, 0.00001)
	require.Len(t, recording.Record.SshOtpMountAttributions, 1)
	attr := recording.Record.SshOtpMountAttributions[mount.MountAccessor]
	require.Equal(t, mount.MountAccessor, attr.MountAccessor)
	require.Equal(t, mount.MountPath, attr.MountPath)
	require.Equal(t, mount.MountType, attr.MountType)
	require.Equal(t, mount.NamespaceID, attr.NamespaceID)
	require.Equal(t, mount.BackendAwareUUID, attr.BackendAwareUUID)
	require.InDelta(t, 0.0014, attr.Count.(float64), 0.00001)

	require.Empty(t, recording.Record.PkiMountAttributions)
	require.Empty(t, recording.Record.SshCertMountAttributions)
}

// TestWithMountInfo_EmptyAccessorSkipped verifies that when MountAccessor is empty
// no attribution entry is recorded (guards against a blank map key), and that
// mountInfo is consumed so a subsequent Add* on the same incrementer also has no attribution.
func TestWithMountInfo_EmptyAccessorSkipped(t *testing.T) {
	recording := &TestCertificateCounter{}

	emptyMount := MountAttribution{
		MountAccessor: "", // intentionally blank
		MountPath:     "pki/",
		MountType:     "pki",
	}

	cert := makeCert(t, 730*time.Hour)
	inc := NewCertCountIncrementer(recording)

	// First call: WithMountInfo(empty) then Add — must not produce attribution.
	inc.WithMountInfo(emptyMount).AddIssuedCertificate(false, cert)
	require.Equal(t, uint64(1), recording.Record.IssuedCerts)
	require.Empty(t, recording.Record.PkiMountAttributions, "blank accessor must not produce attribution entry")

	// Reset recorder and call Add again on the same incrementer without WithMountInfo.
	// mountInfo should have been consumed (set to nil) by the first call.
	recording.Record = CertCount{}
	inc.AddIssuedCertificate(false, cert)
	require.Equal(t, uint64(1), recording.Record.IssuedCerts)
	require.Empty(t, recording.Record.PkiMountAttributions, "stale mountInfo must not leak into a subsequent Add call")
}

// TestWithMountInfo_NotCalledNoAttribution verifies that omitting WithMountInfo
// produces no attribution maps, maintaining backward compatibility.
func TestWithMountInfo_NotCalledNoAttribution(t *testing.T) {
	recording := &TestCertificateCounter{}

	cert := makeCert(t, 730*time.Hour)
	NewCertCountIncrementer(recording).AddIssuedCertificate(true, cert)

	require.Equal(t, uint64(1), recording.Record.IssuedCerts)
	require.Empty(t, recording.Record.PkiMountAttributions)
	require.Empty(t, recording.Record.SshCertMountAttributions)
	require.Empty(t, recording.Record.SshOtpMountAttributions)
}

// TestCertCount_Add_MergesAttributions verifies that CertCount.Add correctly
// accumulates per-mount counts across multiple Add calls for all three metric types.
func TestCertCount_Add_MergesAttributions(t *testing.T) {
	a := CertCount{
		PkiDurationAdjustedCerts: 1.0,
		SSHIssuedCerts:           0.5,
		SSHIssuedOTPs:            0.0014,
		PkiMountAttributions: map[string]MountAttribution{
			"pki_aaa": {MountAccessor: "pki_aaa", Count: 1.0},
		},
		SshCertMountAttributions: map[string]MountAttribution{
			"ssh_aaa": {MountAccessor: "ssh_aaa", Count: 0.5},
		},
		SshOtpMountAttributions: map[string]MountAttribution{
			"otp_aaa": {MountAccessor: "otp_aaa", Count: 0.0014},
		},
	}
	b := CertCount{
		PkiDurationAdjustedCerts: 2.0,
		SSHIssuedCerts:           1.0,
		SSHIssuedOTPs:            0.0014,
		PkiMountAttributions: map[string]MountAttribution{
			"pki_aaa": {MountAccessor: "pki_aaa", Count: 2.0}, // same mount — should sum
			"pki_bbb": {MountAccessor: "pki_bbb", Count: 2.0}, // new mount — should be added
		},
		SshCertMountAttributions: map[string]MountAttribution{
			"ssh_aaa": {MountAccessor: "ssh_aaa", Count: 1.0}, // accumulate
		},
		SshOtpMountAttributions: map[string]MountAttribution{
			"otp_bbb": {MountAccessor: "otp_bbb", Count: 0.0014}, // new
		},
	}

	a.Add(b)

	// PKI
	require.InDelta(t, 3.0, a.PkiDurationAdjustedCerts, 0.0001)
	require.Len(t, a.PkiMountAttributions, 2)
	require.InDelta(t, 3.0, a.PkiMountAttributions["pki_aaa"].Count.(float64), 0.0001, "same-mount PKI counts should accumulate")
	require.InDelta(t, 2.0, a.PkiMountAttributions["pki_bbb"].Count.(float64), 0.0001, "new PKI mount should be present")

	// SSH cert
	require.InDelta(t, 1.5, a.SSHIssuedCerts, 0.0001)
	require.Len(t, a.SshCertMountAttributions, 1)
	require.InDelta(t, 1.5, a.SshCertMountAttributions["ssh_aaa"].Count.(float64), 0.0001, "same-mount SSH counts should accumulate")

	// SSH OTP
	require.InDelta(t, 0.0028, a.SSHIssuedOTPs, 0.00001)
	require.Len(t, a.SshOtpMountAttributions, 2)
	require.InDelta(t, 0.0014, a.SshOtpMountAttributions["otp_aaa"].Count.(float64), 0.00001)
	require.InDelta(t, 0.0014, a.SshOtpMountAttributions["otp_bbb"].Count.(float64), 0.00001)
}

// TestCertCount_IsZero_WithAttributions verifies that IsZero returns false
// when only attribution maps are populated.
func TestCertCount_IsZero_WithAttributions(t *testing.T) {
	c := CertCount{
		PkiMountAttributions: map[string]MountAttribution{
			"pki_aaa": {MountAccessor: "pki_aaa", Count: 1.0},
		},
	}
	require.False(t, c.IsZero(), "CertCount with only attribution should not be zero")

	empty := CertCount{}
	require.True(t, empty.IsZero())
}
