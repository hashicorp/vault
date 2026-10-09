// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package tpm

import "github.com/google/go-attestation/attest"

// AKParams is the stable wire representation of an Attestation Key's public
// parameters sent by the client during gentpmcert/begin.
//
// All byte fields carry TPM 2.0 spec-defined wire-format structures, so their
// meaning is governed by the TCG specification rather than any particular
// version of the go-attestation library.
type AKParams struct {
	// Public is the AK's canonical encoding (TPM2B_PUBLIC / TPMT_PUBLIC).
	Public []byte `json:"public"`
	// UseTCSDActivationFormat is set when tcsd acts as an intermediary.
	UseTCSDActivationFormat bool `json:"use_tcsd_activation_format,omitempty"`
	// CreateData is the TPMS_CREATION_DATA structure for the AK.
	CreateData []byte `json:"create_data"`
	// CreateAttestation is the TPMS_ATTEST structure for the AK.
	CreateAttestation []byte `json:"create_attestation"`
	// CreateSignature is the TPMT_SIGNATURE over CreateAttestation.
	CreateSignature []byte `json:"create_signature"`
}

// akParamsFromAttest converts from the go-attestation library type to AKParams.
func akParamsFromAttest(p attest.AttestationParameters) AKParams {
	return AKParams{
		Public:                  p.Public,
		UseTCSDActivationFormat: p.UseTCSDActivationFormat,
		CreateData:              p.CreateData,
		CreateAttestation:       p.CreateAttestation,
		CreateSignature:         p.CreateSignature,
	}
}

// KeyCertification is the stable wire representation of an application key's
// certification parameters sent by the client during gentpmcert/finish.
type KeyCertification struct {
	// Public is the key's canonical encoding (TPMT_PUBLIC structure).
	Public []byte `json:"public"`
	// CreateData is the TPMS_CREATION_DATA structure for the key.
	CreateData []byte `json:"create_data"`
	// CreateAttestation is the TPMS_ATTEST structure for the key.
	CreateAttestation []byte `json:"create_attestation"`
	// CreateSignature is the TPMT_SIGNATURE over CreateAttestation.
	CreateSignature []byte `json:"create_signature"`
}

// keyCertificationFromAttest converts from the go-attestation library type to
// KeyCertification.
func keyCertificationFromAttest(c attest.CertificationParameters) KeyCertification {
	return KeyCertification{
		Public:            c.Public,
		CreateData:        c.CreateData,
		CreateAttestation: c.CreateAttestation,
		CreateSignature:   c.CreateSignature,
	}
}
