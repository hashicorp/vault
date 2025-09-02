// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type SignCertInput interface {
	CreationBundleInput
	GetCSR() (*x509.CertificateRequest, error)
	IsCA() bool
	UseCSRValues() bool
	GetPermittedDomains() []string
	GetExcludedDomains() []string
	GetPermittedIpRanges() ([]*net.IPNet, error)
	GetExcludedIpRanges() ([]*net.IPNet, error)
	GetPermittedEmailAddresses() []string
	GetExcludedEmailAddresses() []string
	GetPermittedUriDomains() []string
	GetExcludedUriDomains() []string
}

func NewBasicSignCertInput(csr *x509.CertificateRequest, isCA, useCSRValues bool) BasicSignCertInput {
	return NewBasicSignCertInputWithIgnore(csr, isCA, useCSRValues, false)
}

func NewBasicSignCertInputWithIgnore(csr *x509.CertificateRequest, isCA, useCSRValues, ignoreCsrSignature bool) BasicSignCertInput {
	return BasicSignCertInput{
		isCA:               isCA,
		useCSRValues:       useCSRValues,
		csr:                csr,
		ignoreCsrSignature: ignoreCsrSignature,
	}
}

var _ SignCertInput = BasicSignCertInput{}

type BasicSignCertInput struct {
	isCA               bool
	useCSRValues       bool
	csr                *x509.CertificateRequest
	ignoreCsrSignature bool
}

func (b BasicSignCertInput) GetTTL() int {
	return 0
}

func (b BasicSignCertInput) GetOptionalNotAfter() (interface{}, bool) {
	return "", false
}

func (b BasicSignCertInput) GetCommonName() string {
	return ""
}

func (b BasicSignCertInput) GetSerialNumber() string {
	return ""
}

func (b BasicSignCertInput) GetExcludeCnFromSans() bool {
	return false
}

func (b BasicSignCertInput) GetOptionalAltNames() (interface{}, bool) {
	return []string{}, false
}

func (b BasicSignCertInput) GetOtherSans() []string {
	return []string{}
}

func (b BasicSignCertInput) GetIpSans() []string {
	return []string{}
}

func (b BasicSignCertInput) GetURISans() []string {
	return []string{}
}

func (b BasicSignCertInput) GetOptionalSkid() (interface{}, bool) {
	return "", false
}

func (b BasicSignCertInput) IsUserIdInSchema() (interface{}, bool) {
	return []string{}, false
}

func (b BasicSignCertInput) GetUserIds() []string {
	return []string{}
}

func (b BasicSignCertInput) GetCSR() (*x509.CertificateRequest, error) {
	return b.csr, nil
}

func (b BasicSignCertInput) IgnoreCSRSignature() bool {
	return b.ignoreCsrSignature
}

func (b BasicSignCertInput) IsCA() bool {
	return b.isCA
}

func (b BasicSignCertInput) UseCSRValues() bool {
	return b.useCSRValues
}

func (b BasicSignCertInput) GetPermittedDomains() []string {
	return []string{}
}

func (b BasicSignCertInput) GetExcludedDomains() []string {
	return []string{}
}

// GetPermittedIpRanges returns the permitted IP ranges for the name constraints extension.
// ignore-nil-nil-function-check
func (b BasicSignCertInput) GetPermittedIpRanges() ([]*net.IPNet, error) {
	return nil, nil
}

// GetExcludedIpRanges returns the excluded IP ranges for the name constraints extension.
// ignore-nil-nil-function-check
func (b BasicSignCertInput) GetExcludedIpRanges() ([]*net.IPNet, error) {
	return nil, nil
}

func (b BasicSignCertInput) GetPermittedEmailAddresses() []string {
	return []string{}
}

func (b BasicSignCertInput) GetExcludedEmailAddresses() []string {
	return []string{}
}

func (b BasicSignCertInput) GetPermittedUriDomains() []string {
	return []string{}
}

func (b BasicSignCertInput) GetExcludedUriDomains() []string {
	return []string{}
}

func SignCert(b logical.SystemView, role *RoleEntry, entityInfo EntityInfo, caSign *certutil.CAInfoBundle, signInput SignCertInput) (*certutil.ParsedCertBundle, []string, error) {
	if role == nil {
		return nil, nil, errutil.InternalError{Err: "no role found in data bundle"}
	}

	csr, err := signInput.GetCSR()
	if err != nil {
		return nil, nil, err
	}

	if csr.PublicKeyAlgorithm == x509.UnknownPublicKeyAlgorithm || csr.PublicKey == nil {
		return nil, nil, errutil.UserError{Err: "Refusing to sign CSR with empty PublicKey. This usually means the SubjectPublicKeyInfo field has an OID not recognized by Go, such as 1.2.840.113549.1.1.10 for rsaPSS."}
	}

	// This switch validates that the CSR key type matches the role and sets
	// the Value in the actualKeyType/actualKeyBits values.
	actualKeyType := ""
	actualKeyBits := 0

	switch role.KeyType {
	case "rsa":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.RSA {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf("role requires keys of type %s", role.KeyType)}
		}

		pubKey, ok := csr.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}

		actualKeyType = "rsa"
		actualKeyBits = pubKey.N.BitLen()
	case "ec":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.ECDSA {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires keys of type %s",
				role.KeyType)}
		}
		pubKey, ok := csr.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}

		actualKeyType = "ec"
		actualKeyBits = pubKey.Params().BitSize
	case "ed25519":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.Ed25519 {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires keys of type %s",
				role.KeyType)}
		}

		_, ok := csr.PublicKey.(ed25519.PublicKey)
		if !ok {
			return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}

		actualKeyType = "ed25519"
		actualKeyBits = 0
	case "any":
		// We need to compute the actual key type and key bits, to correctly
		// validate minimums and SignatureBits below.
		switch csr.PublicKeyAlgorithm {
		case x509.RSA:
			pubKey, ok := csr.PublicKey.(*rsa.PublicKey)
			if !ok {
				return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
			}
			if pubKey.N.BitLen() < 2048 {
				return nil, nil, errutil.UserError{Err: "RSA keys < 2048 bits are unsafe and not supported"}
			}

			actualKeyType = "rsa"
			actualKeyBits = pubKey.N.BitLen()
		case x509.ECDSA:
			pubKey, ok := csr.PublicKey.(*ecdsa.PublicKey)
			if !ok {
				return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
			}

			actualKeyType = "ec"
			actualKeyBits = pubKey.Params().BitSize
		case x509.Ed25519:
			_, ok := csr.PublicKey.(ed25519.PublicKey)
			if !ok {
				return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
			}

			actualKeyType = "ed25519"
			actualKeyBits = 0
		default:
			return nil, nil, errutil.UserError{Err: "Unknown key type in CSR: " + csr.PublicKeyAlgorithm.String()}
		}
	default:
		return nil, nil, errutil.InternalError{Err: fmt.Sprintf("unsupported key type Value: %s", role.KeyType)}
	}

	// Before validating key lengths, update our KeyBits/SignatureBits based
	// on the actual CSR key type.
	if role.KeyType == "any" {
		// We update the Value of KeyBits and SignatureBits here (from the
		// role), using the specified key type. This allows us to convert
		// the default Value (0) for SignatureBits and KeyBits to a
		// meaningful Value.
		//
		// We ignore the role's original KeyBits Value if the KeyType is any
		// as legacy (pre-1.10) roles had default values that made sense only
		// for RSA keys (key_bits=2048) and the older code paths ignored the role Value
		// set for KeyBits when KeyType was set to any. This also enforces the
		// docs saying when key_type=any, we only enforce our specified minimums
		// for signing operations
		var err error
		if role.KeyBits, role.SignatureBits, err = certutil.ValidateDefaultOrValueKeyTypeSignatureLength(
			actualKeyType, 0, role.SignatureBits); err != nil {
			return nil, nil, errutil.InternalError{Err: fmt.Sprintf("unknown internal error updating default values: %v", err)}
		}

		// We're using the KeyBits field as a minimum Value below, and P-224 is safe
		// and a previously allowed Value. However, the above call defaults
		// to P-256 as that's a saner default than P-224 (w.r.t. generation), so
		// override it here to allow 224 as the smallest size we permit.
		if actualKeyType == "ec" {
			role.KeyBits = 224
		}
	}

	// At this point, role.KeyBits and role.SignatureBits should both
	// be non-zero, for RSA and ECDSA keys. Validate the actualKeyBits based on
	// the role's values. If the KeyType was any, and KeyBits was set to 0,
	// KeyBits should be updated to 2048 unless some other Value was chosen
	// explicitly.
	//
	// This validation needs to occur regardless of the role's key type, so
	// that we always validate both RSA and ECDSA key sizes.
	if actualKeyType == "rsa" {
		if actualKeyBits < role.KeyBits {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires a minimum of a %d-bit key, but CSR's key is %d bits",
				role.KeyBits, actualKeyBits)}
		}

		if actualKeyBits < 2048 {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"Vault requires a minimum of a 2048-bit key, but CSR's key is %d bits",
				actualKeyBits)}
		}
	} else if actualKeyType == "ec" {
		if actualKeyBits < role.KeyBits {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires a minimum of a %d-bit key, but CSR's key is %d bits",
				role.KeyBits,
				actualKeyBits)}
		}
	}

	creation, warnings, err := GenerateCreationBundle(b, role, entityInfo, signInput, caSign, csr)
	if err != nil {
		return nil, nil, err
	}
	if creation.Params == nil {
		return nil, nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	creation.Params.IsCA = signInput.IsCA()
	creation.Params.UseCSRValues = signInput.UseCSRValues()

	if signInput.IsCA() {
		creation.Params.PermittedDNSDomains = signInput.GetPermittedDomains()
		creation.Params.ExcludedDNSDomains = signInput.GetExcludedDomains()
		creation.Params.PermittedIPRanges, err = signInput.GetPermittedIpRanges()
		if err != nil {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf("error parsinng permitted IP ranges: %v", err)}
		}
		creation.Params.ExcludedIPRanges, err = signInput.GetExcludedIpRanges()
		if err != nil {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf("error parsinng excluded IP ranges: %v", err)}
		}
		creation.Params.PermittedEmailAddresses = signInput.GetPermittedEmailAddresses()
		creation.Params.ExcludedEmailAddresses = signInput.GetExcludedEmailAddresses()
		creation.Params.PermittedURIDomains = signInput.GetPermittedUriDomains()
		creation.Params.ExcludedURIDomains = signInput.GetExcludedUriDomains()
	} else {
		for _, ext := range csr.Extensions {
			if ext.Id.Equal(certutil.ExtensionBasicConstraintsOID) {
				warnings = append(warnings, "specified CSR contained a Basic Constraints extension that was ignored during issuance")
			}
		}
	}

	parsedBundle, err := certutil.SignCertificate(creation)
	if err != nil {
		return nil, nil, err
	}

	return parsedBundle, warnings, nil
}
