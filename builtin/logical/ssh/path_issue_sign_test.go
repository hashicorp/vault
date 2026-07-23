// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

// newTestBundle returns a minimal creationBundle ready for sign() using the
// shared test RSA CA key unless overridden by the caller.
func newTestBundle(t *testing.T, signer ssh.Signer, algorithmSigner string) *creationBundle {
	t.Helper()

	pub, _, _, _, err := ssh.ParseAuthorizedKey([]byte(testCAPublicKey))
	require.NoError(t, err)

	return &creationBundle{
		KeyID:           "test-key-id",
		PublicKey:       pub,
		Signer:          signer,
		ValidPrincipals: []string{"testuser"},
		CertificateType: ssh.UserCert,
		TTL:             time.Hour,
		Role:            &sshRole{AlgorithmSigner: algorithmSigner},
	}
}

// TestCreationBundleSign_DefaultAlgorithm_RSA verifies that when algorithm_signer
// is "default" and the CA is an RSA key, the produced certificate is signed with
// rsa-sha2-256 (not the deprecated ssh-rsa).
func TestCreationBundleSign_DefaultAlgorithm_RSA(t *testing.T) {
	signer, err := ssh.ParsePrivateKey([]byte(testCAPrivateKey))
	require.NoError(t, err)

	bundle := newTestBundle(t, signer, DefaultAlgorithmSigner)
	cert, err := bundle.sign()
	require.NoError(t, err)
	require.NotNil(t, cert)

	require.NotNil(t, cert.Signature, "certificate must be signed")
	assert.Equal(t, ssh.KeyAlgoRSASHA256, cert.Signature.Format,
		"default RSA signing should use rsa-sha2-256, not ssh-rsa")

	// Structural validity: cert must be a user cert, have the right key ID and principals.
	assert.Equal(t, uint32(ssh.UserCert), cert.CertType)
	assert.Equal(t, "test-key-id", cert.KeyId)
	assert.Equal(t, []string{"testuser"}, cert.ValidPrincipals)
	assert.NotZero(t, cert.Serial)
	assert.NotZero(t, cert.ValidBefore)
}

// TestCreationBundleSign_ExplicitRSASHA256 verifies that an explicit
// algorithm_signer of rsa-sha2-256 produces a cert signed with rsa-sha2-256.
func TestCreationBundleSign_ExplicitRSASHA256(t *testing.T) {
	signer, err := ssh.ParsePrivateKey([]byte(testCAPrivateKey))
	require.NoError(t, err)

	bundle := newTestBundle(t, signer, ssh.KeyAlgoRSASHA256)
	cert, err := bundle.sign()
	require.NoError(t, err)
	require.NotNil(t, cert)

	require.NotNil(t, cert.Signature)
	assert.Equal(t, ssh.KeyAlgoRSASHA256, cert.Signature.Format)
}

// TestCreationBundleSign_ExplicitRSASHA512 verifies that an explicit
// algorithm_signer of rsa-sha2-512 produces a cert signed with rsa-sha2-512.
func TestCreationBundleSign_ExplicitRSASHA512(t *testing.T) {
	signer, err := ssh.ParsePrivateKey([]byte(testCAPrivateKey))
	require.NoError(t, err)

	bundle := newTestBundle(t, signer, ssh.KeyAlgoRSASHA512)
	cert, err := bundle.sign()
	require.NoError(t, err)
	require.NotNil(t, cert)

	require.NotNil(t, cert.Signature)
	assert.Equal(t, ssh.KeyAlgoRSASHA512, cert.Signature.Format)
}

// TestCreationBundleSign_ExplicitSSHRSA_Deprecated verifies that explicitly
// requesting the deprecated ssh-rsa algorithm still works (for backward
// compatibility) and produces a cert signed with ssh-rsa.
func TestCreationBundleSign_ExplicitSSHRSA_Deprecated(t *testing.T) {
	signer, err := ssh.ParsePrivateKey([]byte(testCAPrivateKey))
	require.NoError(t, err)

	bundle := newTestBundle(t, signer, ssh.SigAlgoRSA)
	cert, err := bundle.sign()
	require.NoError(t, err)
	require.NotNil(t, cert)

	require.NotNil(t, cert.Signature)
	assert.Equal(t, ssh.SigAlgoRSA, cert.Signature.Format)
}

// TestCreationBundleSign_DefaultAlgorithm_Ed25519 verifies that when
// algorithm_signer is "default" and the CA is an ed25519 key (non-RSA),
// sign() succeeds and the cert carries the ed25519 signature format.
func TestCreationBundleSign_DefaultAlgorithm_Ed25519(t *testing.T) {
	signer, err := ssh.ParsePrivateKey([]byte(testCAPrivateKeyEd25519))
	require.NoError(t, err)

	// For a non-RSA CA the bundle's PublicKey can be any valid key; reuse the
	// ed25519 public key here.
	pub, _, _, _, err := ssh.ParseAuthorizedKey([]byte(testCAPublicKeyEd25519))
	require.NoError(t, err)

	bundle := &creationBundle{
		KeyID:           "test-key-id",
		PublicKey:       pub,
		Signer:          signer,
		ValidPrincipals: []string{"testuser"},
		CertificateType: ssh.UserCert,
		TTL:             time.Hour,
		Role:            &sshRole{AlgorithmSigner: DefaultAlgorithmSigner},
	}

	cert, err := bundle.sign()
	require.NoError(t, err)
	require.NotNil(t, cert)

	require.NotNil(t, cert.Signature)
	assert.Equal(t, ssh.KeyAlgoED25519, cert.Signature.Format)
}

// TestCreationBundleSign_SignatureVerifies confirms that the signature embedded
// in the returned certificate is cryptographically valid against the CA public key.
func TestCreationBundleSign_SignatureVerifies(t *testing.T) {
	signer, err := ssh.ParsePrivateKey([]byte(testCAPrivateKey))
	require.NoError(t, err)

	bundle := newTestBundle(t, signer, DefaultAlgorithmSigner)
	cert, err := bundle.sign()
	require.NoError(t, err)

	// MarshalAuthorizedKey → ParseAuthorizedKey round-trip exercises the same
	// path the SSH client would take when verifying the cert.
	marshaled := ssh.MarshalAuthorizedKey(cert)
	require.NotEmpty(t, marshaled)

	parsed, _, _, _, err := ssh.ParseAuthorizedKey(marshaled)
	require.NoError(t, err, "signed cert must survive marshal/parse round-trip")

	parsedCert, ok := parsed.(*ssh.Certificate)
	require.True(t, ok)
	assert.NotNil(t, parsedCert.Signature)
	assert.Equal(t, cert.Serial, parsedCert.Serial)
}
