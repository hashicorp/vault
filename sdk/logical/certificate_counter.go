// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"crypto/x509"
	"math"
	"time"
)

// CertificateCounter is an interface for incrementing the count of issued and stored
// certificates.
type CertificateCounter interface {
	// IncrementCount increments the count of issued and stored certificates.
	AddCount(params CertCount)

	// Increment returns a CertCountIncrementer that can be used to add
	// to the count.
	Increment() CertCountIncrementer
}

// CertCount represents the parameters for incrementing certificate counts.
type CertCount struct {
	IssuedCerts uint64 // TODO(victorr): Rename to PkiIssuedCerts
	StoredCerts uint64

	// PkiDurationAdjustedCerts tracks the normalized certificate duration units for billing
	// purposes. Each certificate's billable units = (Validity Hours ÷ 730), rounded to 4 decimal
	// places.
	PkiDurationAdjustedCerts float64
	SSHIssuedCerts           float64
	SSHIssuedOTPs            float64

	// PkiMountAttributions contains per-mount/namespace attribution for PKI certificates.
	// Keyed by mount accessor.
	PkiMountAttributions map[string]MountAttribution

	// SshCertMountAttributions contains per-mount/namespace attribution for SSH certificates.
	// Keyed by mount accessor.
	SshCertMountAttributions map[string]MountAttribution

	// SshOtpMountAttributions contains per-mount/namespace attribution for SSH OTP credentials.
	// Keyed by mount accessor.
	SshOtpMountAttributions map[string]MountAttribution
}

func (i *CertCount) Add(other CertCount) {
	i.IssuedCerts += other.IssuedCerts
	i.StoredCerts += other.StoredCerts
	i.PkiDurationAdjustedCerts += other.PkiDurationAdjustedCerts
	i.SSHIssuedCerts += other.SSHIssuedCerts
	i.SSHIssuedOTPs += other.SSHIssuedOTPs

	mergeMountAttributions(&i.PkiMountAttributions, other.PkiMountAttributions)
	mergeMountAttributions(&i.SshCertMountAttributions, other.SshCertMountAttributions)
	mergeMountAttributions(&i.SshOtpMountAttributions, other.SshOtpMountAttributions)
}

// mergeMountAttributions merges src into *dst, accumulating numeric counts for
// mounts that already exist in *dst.
func mergeMountAttributions(dst *map[string]MountAttribution, src map[string]MountAttribution) {
	if len(src) == 0 {
		return
	}
	if *dst == nil {
		*dst = make(map[string]MountAttribution)
	}

	asFloat64 := func(v interface{}) float64 {
		switch n := v.(type) {
		case float64:
			return n
		case float32:
			return float64(n)
		case int:
			return float64(n)
		case int64:
			return float64(n)
		case uint:
			return float64(n)
		case uint64:
			return float64(n)
		case interface{ Float64() (float64, error) }:
			f, _ := n.Float64()
			return f
		default:
			return 0
		}
	}

	for accessor, attr := range src {
		if existing, ok := (*dst)[accessor]; ok {
			existing.Count = asFloat64(existing.Count) + asFloat64(attr.Count)
			(*dst)[accessor] = existing
			continue
		}
		(*dst)[accessor] = attr
	}
}

func (i *CertCount) IsZero() bool {
	return i.IssuedCerts == 0 &&
		i.StoredCerts == 0 &&
		i.PkiDurationAdjustedCerts == 0 &&
		i.SSHIssuedCerts == 0 &&
		i.SSHIssuedOTPs == 0 &&
		len(i.PkiMountAttributions) == 0 &&
		len(i.SshCertMountAttributions) == 0 &&
		len(i.SshOtpMountAttributions) == 0
}

// durationAdjustedCertificateCount calculates the billable units for a certificate based on its
// validity duration. WARNING: Beware the maximum value for time.Duration (approximately 290 years).
//
// The calculation follows the billing specification:
// - Standard duration is 730 hours (1 month)
// - Units = (Validity Hours ÷ 730), rounded to 4 decimal places
// - Example: 1-year cert (8760 hours) = 12.0000 units
// - Example: 1-day cert (24 hours) = 0.0329 units
func durationAdjustedCertificateCount(validitySeconds int64) float64 {
	const standardDuration = 730.0
	validityHours := float64(validitySeconds) / 3600.0
	units := validityHours / standardDuration
	// Round to 4 decimal places
	ret := math.Round(units*10000) / 10000
	if ret == 0.0 && validitySeconds > 0 {
		// Ensure we don't return 0.0, which would be interpreted as no billable units.
		return 0.0001
	}
	return ret
}

type CertCountIncrementer interface {
	// WithMountInfo attaches mount and namespace metadata to the AddIssuedCertificate,
	// AddSSHCertificate, or AddSSHOTP call that follows, to record per-mount attribution
	// for that operation. MountAttribution.Count is not set here, it will be set to the
	// billable units calculated in the following Add* call.
	WithMountInfo(mount MountAttribution) CertCountIncrementer
	AddIssuedCertificate(stored bool, cert *x509.Certificate) CertCountIncrementer
	AddSSHCertificate(ttl time.Duration) CertCountIncrementer
	AddSSHOTP() CertCountIncrementer
}

type certCountIncrementer struct {
	counter   CertificateCounter
	mountInfo *MountAttribution
}

var _ CertCountIncrementer = (*certCountIncrementer)(nil)

// NewCertCountIncrementer creates a new CertCountIncrementer for the given counter.
func NewCertCountIncrementer(counter CertificateCounter) CertCountIncrementer {
	return &certCountIncrementer{counter: counter}
}

// WithMountInfo attaches mount/namespace metadata that will be recorded as attribution
// on the next Add* call. Count on the provided MountAttribution is ignored.
func (c *certCountIncrementer) WithMountInfo(mount MountAttribution) CertCountIncrementer {
	c.mountInfo = &mount
	return c
}

// AddIssuedCertificate increments the issued certificate count by 1, the stored certificate
// count if stored is true, and adds the calculated billable units based on the certificate's
// validity duration.
// cert: The X.509 certificate to extract validity duration from.
func (c *certCountIncrementer) AddIssuedCertificate(stored bool, cert *x509.Certificate) CertCountIncrementer {
	validity := int64(cert.NotAfter.Unix() - cert.NotBefore.Unix())
	units := durationAdjustedCertificateCount(validity)
	count := CertCount{
		IssuedCerts:              1,
		PkiDurationAdjustedCerts: units,
	}
	if stored {
		count.StoredCerts = 1
	}
	if c.mountInfo != nil {
		if c.mountInfo.MountAccessor != "" {
			attr := *c.mountInfo
			attr.Count = units
			count.PkiMountAttributions = map[string]MountAttribution{attr.MountAccessor: attr}
		}
		c.mountInfo = nil // always consume, even when accessor was blank
	}
	c.counter.AddCount(count)
	return c
}

func (c *certCountIncrementer) AddSSHCertificate(ttl time.Duration) CertCountIncrementer {
	units := durationAdjustedCertificateCount(int64(ttl.Seconds()))
	count := CertCount{
		SSHIssuedCerts: units,
	}
	if c.mountInfo != nil {
		if c.mountInfo.MountAccessor != "" {
			attr := *c.mountInfo
			attr.Count = units
			count.SshCertMountAttributions = map[string]MountAttribution{attr.MountAccessor: attr}
		}
		c.mountInfo = nil // always consume
	}
	c.counter.AddCount(count)
	return c
}

func (c *certCountIncrementer) AddSSHOTP() CertCountIncrementer {
	const otpUnits = 0.0014
	count := CertCount{
		SSHIssuedOTPs: otpUnits,
	}
	if c.mountInfo != nil {
		if c.mountInfo.MountAccessor != "" {
			attr := *c.mountInfo
			attr.Count = otpUnits
			count.SshOtpMountAttributions = map[string]MountAttribution{attr.MountAccessor: attr}
		}
		c.mountInfo = nil // always consume
	}
	c.counter.AddCount(count)
	return c
}
