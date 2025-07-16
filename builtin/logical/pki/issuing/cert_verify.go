// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"fmt"
	"os"
	"strconv"

	ctx509 "github.com/google/certificate-transparency-go/x509"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// disableVerifyCertificateEnvVar is an environment variable that can be used to disable the
// verification done when issuing or signing certificates that was added by VAULT-22013. It
// is meant as a scape hatch to avoid breaking deployments that the new verification would
// break.
const disableVerifyCertificateEnvVar = "VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION"

func isCertificateVerificationDisabled() (bool, error) {
	disableRaw, ok := os.LookupEnv(disableVerifyCertificateEnvVar)
	if !ok {
		return false, nil
	}

	disable, err := strconv.ParseBool(disableRaw)
	if err != nil {
		return false, fmt.Errorf("failed parsing environment variable %s: %w", disableVerifyCertificateEnvVar, err)
	}

	return disable, nil
}

func VerifyCertificate(ctx context.Context, storage logical.Storage, issuerId IssuerID, parsedBundle *certutil.ParsedCertBundle) error {
	if verificationDisabled, err := isCertificateVerificationDisabled(); err != nil {
		return err
	} else if verificationDisabled {
		return nil
	}

	// Note that we use github.com/google/certificate-transparency-go/x509 to perform certificate verification,
	// since that library provides options to disable checks that the standard library does not.
	options := ctx509.VerifyOptions{
		KeyUsages:                      nil,
		MaxConstraintComparisions:      0, // Use the library's 'sensible default'
		DisableTimeChecks:              true,
		DisableEKUChecks:               true,
		DisableCriticalExtensionChecks: false,
		DisableNameChecks:              false,
		DisablePathLenChecks:           false,
		DisableNameConstraintChecks:    false,
	}
	if err := entSetCertVerifyOptions(ctx, storage, issuerId, &options); err != nil {
		return err
	}

	return certutil.VerifyCertificate(parsedBundle, options)
}
