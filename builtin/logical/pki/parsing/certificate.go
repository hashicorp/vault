// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package parsing

import (
	"crypto/x509"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/certutil"
)

func SerialFromCert(cert *x509.Certificate) string {
	return SerialFromBigInt(cert.SerialNumber)
}

func SerialFromBigInt(serial *big.Int) string {
	return strings.TrimSpace(certutil.GetHexFormatted(serial.Bytes(), ":"))
}

// NormalizeSerialForStorageFromBigInt given a serial number, format it as a string
// that is safe to store within a filesystem
func NormalizeSerialForStorageFromBigInt(serial *big.Int) string {
	return strings.TrimSpace(certutil.GetHexFormatted(serial.Bytes(), "-"))
}

// NormalizeSerialForStorage given a serial number with ':' characters, convert
// them to '-' which is safe to store within filesystems
func NormalizeSerialForStorage(serial string) string {
	return strings.ReplaceAll(strings.ToLower(serial), ":", "-")
}

func ParseCertificateFromString(pemCert string) (*x509.Certificate, error) {
	return ParseCertificateFromBytes([]byte(pemCert))
}

func ParseCertificateFromBytes(certBytes []byte) (*x509.Certificate, error) {
	block, err := DecodePem(certBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %w", err)
	}

	return cert, nil
}

func ParseCertificatesFromString(pemCerts string) ([]*x509.Certificate, error) {
	return ParseCertificatesFromBytes([]byte(pemCerts))
}

func ParseCertificatesFromBytes(certBytes []byte) ([]*x509.Certificate, error) {
	block, err := DecodePem(certBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %w", err)
	}

	cert, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %w", err)
	}

	return cert, nil
}

var (
	// keyUsageToString maps a x509.KeyUsage bitmask to its name.
	keyUsageToString = map[x509.KeyUsage]string{
		x509.KeyUsageDigitalSignature:  "DigitalSignature",
		x509.KeyUsageContentCommitment: "ContentCommitment",
		x509.KeyUsageKeyEncipherment:   "KeyEncipherment",
		x509.KeyUsageDataEncipherment:  "DataEncipherment",
		x509.KeyUsageKeyAgreement:      "KeyAgreement",
		x509.KeyUsageCertSign:          "CertSign",
		x509.KeyUsageCRLSign:           "CRLSign",
		x509.KeyUsageEncipherOnly:      "EncipherOnly",
		x509.KeyUsageDecipherOnly:      "DecipherOnly",
	}

	lowerStringToKeyUsage = flipAndLowerCaseKeyUsageMap(keyUsageToString)
)

func flipAndLowerCaseKeyUsageMap(aMap map[x509.KeyUsage]string) map[string]x509.KeyUsage {
	flipped := make(map[string]x509.KeyUsage, len(aMap))
	for k, v := range aMap {
		flipped[strings.ToLower(v)] = k
	}
	return flipped
}

// ParseKeyUsages returns a bitmap of all the key usage strings that we
// can map back to a known key usages. Unknown values are ignored.
func ParseKeyUsages(input []string) x509.KeyUsage {
	var parsedKeyUsages x509.KeyUsage

	for _, k := range input {
		lk := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(k)), "keyusage")
		if ku, ok := lowerStringToKeyUsage[lk]; ok {
			parsedKeyUsages |= ku
		}
	}

	return parsedKeyUsages
}

// KeyUsageIsPresent checks the provided bitmap (keyUsages) for presence of the provided x509.KeyUsage.
func KeyUsageIsPresent(keyUsages x509.KeyUsage, usage x509.KeyUsage) bool {
	if _, ok := keyUsageToString[usage]; !ok {
		return false
	}
	return keyUsages&usage != 0
}

// KeyUsagesArePresent checks that all the requested key usages are present within the usages provided
func KeyUsagesArePresent(usages x509.KeyUsage, reqUsages []x509.KeyUsage) bool {
	for _, reqUsage := range reqUsages {
		if !KeyUsageIsPresent(usages, reqUsage) {
			return false
		}
	}

	return true
}

// KeyUsageStrings return all the known values represented by strings that are set
// within the passed in keyUsages argument.
func KeyUsageStrings(keyUsages x509.KeyUsage) []string {
	var keyUsageStrings []string
	for ku, name := range keyUsageToString {
		if KeyUsageIsPresent(keyUsages, ku) {
			keyUsageStrings = append(keyUsageStrings, name)
		}
	}
	return keyUsageStrings
}

// KeyUsageToString convert the individual key usage value into a string. If unknown
// a string saying the value is unknown is returned.
func KeyUsageToString(keyUsage x509.KeyUsage) string {
	if kuName, ok := keyUsageToString[keyUsage]; ok {
		return kuName
	}

	return fmt.Sprintf("unknown key usage: %d", keyUsage)
}
