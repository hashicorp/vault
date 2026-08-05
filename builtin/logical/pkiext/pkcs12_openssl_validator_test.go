// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pkiext

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	pkihelper "github.com/hashicorp/vault/helper/testhelpers/pki"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/stretchr/testify/require"
)

const (
	opensslImageRepo = "docker.mirror.hashicorp.services/alpine/openssl"
	openssl3_5       = "3.5.6"
	openssl3_3       = "3.3.3"
	pkcs12Password   = "123-secure-password"
)

// TestPKCS12OpenSSLValidation validates PKCS#12 output from different encoders
// can be read by different OpenSSL versions.
// This test focuses on PKCS#12 format compatibility, not PKI backend functionality.
func TestPKCS12OpenSSLValidation(t *testing.T) {
	t.Parallel()

	// Cleanup containers after tests complete
	envs := map[string]*alpineEnv{}
	t.Cleanup(func() {
		for tag, env := range envs {
			if err := env.runner.Stop(context.Background(), env.container.Container.ID); err != nil {
				t.Logf("Warning: failed to stop container (tag: %s): %v", tag, err)
			}
		}
	})

	// Generate test CA, leaf and key once for all subtests
	result := pkihelper.GenerateCertWithRoot(t)
	leafKey, leafCert, caChain := result.Leaf.Key, result.Leaf.Cert, []*x509.Certificate{result.RootCa.Cert}

	// Test matrix:
	// 2 bundle types (with and without a private key) x 2 encoders (modern2026, modern2023) x 2 OpenSSL versions
	// - modern2026 encoder is incompatible with OpenSSL 3.3 (uses newer algorithms)
	// - modern2023 encoder works with both OpenSSL 3.3 and 3.5
	testCases := []struct {
		bundleType  string // "keystore" or "trust store"
		encoder     string
		version     string
		shouldError bool
	}{
		{bundleType: "keystore", encoder: "modern2026", version: openssl3_5},
		{bundleType: "keystore", encoder: "modern2023", version: openssl3_5},
		{bundleType: "keystore", encoder: "modern2026", version: openssl3_3, shouldError: true},
		{bundleType: "keystore", encoder: "modern2023", version: openssl3_3},

		{bundleType: "trust store", encoder: "modern2026", version: openssl3_5},
		{bundleType: "trust store", encoder: "modern2023", version: openssl3_5},
		{bundleType: "trust store", encoder: "modern2026", version: openssl3_3, shouldError: true},
		{bundleType: "trust store", encoder: "modern2023", version: openssl3_3},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("bundle=%s encoder=%s openssl=%s", tc.bundleType, tc.encoder, tc.version)
		if tc.shouldError {
			name += " (should error)"
		}
		t.Run(name, func(t *testing.T) {
			log := corehelpers.NewTestLogger(t)

			var pkcs12Bytes []byte
			var err error
			if tc.bundleType == "keystore" {
				// Only pass private key for keystores
				pkcs12Bytes, err = pki.EncodeToPKCS12(
					tc.encoder,
					leafKey,
					leafCert,
					caChain,
					pkcs12Password)
			} else {
				pkcs12Bytes, err = pki.EncodeToPKCS12(
					tc.encoder, nil,
					leafCert,
					caChain,
					pkcs12Password)
			}
			require.NoError(t, err, "EncodeToPKCS12 should succeed")
			require.NotEmpty(t, pkcs12Bytes)

			// Validate with openSSL
			opensslInfo, err := runOpenSSLInfo(t, envs, pkcs12Bytes, tc.version)

			if tc.shouldError {
				// Expect openssl to fail for incompatible PKCS12 encoders (modern2026 with OpenSSL 3.3)
				require.Error(t, err, "openssl should fail to read PKCS12 file with incompatible encoder")
				log.Info("OpenSSL failed as expected with incompatible encoder", "error", err)
				// Verify it's the expected algorithm incompatibility error and not some other failure
				errMsg := err.Error()
				isExpectedError := strings.Contains(errMsg, "unknown digest algorithm") || strings.Contains(errMsg, "PBMAC1") || strings.Contains(errMsg, "mac generation error")
				require.True(t, isExpectedError, "Expected algorithm parsing error but got: %s", errMsg)
				return
			}

			require.NoError(t, err, "OpenSSL should successfully read PKCS12 file")
			log.Info("OpenSSL info", "output", opensslInfo)
			require.Contains(t, opensslInfo, "PKCS7 Encrypted data: PBES2, PBKDF2, AES-256-CBC, Iteration 2048, PRF hmacWithSHA256", "Encoder should use expected encryption")

			// Verify encryption uses expected algorithms
			expectedMac := "MAC: PBMAC1 using PBKDF2"
			if tc.encoder == "modern2023" {
				expectedMac = "MAC: sha256"
			}
			require.Contains(t, opensslInfo, expectedMac, "%s encoder should use %s", tc.encoder, expectedMac)

			runner := envs[tc.version].runner
			containerID := envs[tc.version].container.Container.ID

			// Determine if this bundle includes a private key
			hasPrivateKey := tc.bundleType == "keystore"
			validateCertificateOrder(t, runner, containerID, hasPrivateKey)

			// Validate that OpenSSL can extract the private key
			// and exists for bundles with key and empty for trust stores
			output, err := runOpenSSLCmd(runner, containerID, []string{
				"pkcs12",
				"-in", "/tmp/bundle.p12",
				"-nocerts",
				"-nodes",
				"-passin", "pass:" + pkcs12Password,
				"-noenc",
			})
			require.NoError(t, err, "OpenSSL should successfully extract private key from PKCS12")
			if hasPrivateKey {
				require.Contains(t, opensslInfo, "Shrouded Keybag: PBES2, PBKDF2, AES-256-CBC, Iteration 2048, PRF hmacWithSHA256", "private key should use expected encryption")
				require.Contains(t, output, "-----BEGIN PRIVATE KEY-----", "OpenSSL should output private key")
			} else {
				require.True(t, strings.TrimSpace(output) == "", "Key should be empty")
			}
		})
	}

	//  validates PKCS#12 bundles work with different key types
	keyTypes := []struct {
		name    string
		keyType string
		keyBits int
	}{
		// RSA keys
		{name: "RSA-2048", keyType: "rsa", keyBits: 2048},
		{name: "RSA-4096", keyType: "rsa", keyBits: 4096},

		// ECDSA keys (ECDSA-P256 is used for interop tests above so not re-tested here)
		{name: "ECDSA-P384", keyType: "ec", keyBits: 384},
		{name: "ECDSA-P521", keyType: "ec", keyBits: 521},

		// Ed25519
		{name: "Ed25519", keyType: "ed25519", keyBits: 0},
	}

	for _, tc := range keyTypes {
		// Generate key and certificate for this key type
		privateKey, cert, caChain, err := generateKeyAndCert(t, tc.keyType, tc.keyBits)
		require.NoError(t, err, "Should generate cert and private key: %s, bits: %s", tc.keyType, tc.keyBits)

		// Test each encoder type
		for _, encoder := range []string{"modern2023", "modern2026"} {
			name := fmt.Sprintf("encoder=%s key=%s", encoder, tc.name)

			t.Run(name, func(t *testing.T) {
				// Encode to PKCS#12
				pkcs12Bytes, err := pki.EncodeToPKCS12(
					encoder,
					privateKey,
					cert,
					caChain,
					pkcs12Password,
				)
				require.NoError(t, err, "Failed to encode PKCS#12 for key type: %s, bits:", tc.keyType, tc.keyBits)

				// Validate with OpenSSL
				opensslInfo, err := runOpenSSLInfo(t, envs, pkcs12Bytes, openssl3_5)
				require.NoError(t, err, "OpenSSL should read PKCS#12 bundle")
				// Verify expected encryption algorithms
				require.Contains(t, opensslInfo, "PBES2, PBKDF2, AES-256-CBC")
				// Extract and verify private key
				runner := envs[openssl3_5].runner
				containerID := envs[openssl3_5].container.Container.ID
				keyOutput, err := runOpenSSLCmd(runner, containerID, []string{
					"pkcs12", "-in", "/tmp/bundle.p12",
					"-nocerts", "-nodes",
					"-passin", "pass:" + pkcs12Password,
					"-noenc",
				})
				require.NoError(t, err, "Should extract private key")
				require.Contains(t, keyOutput, "-----BEGIN PRIVATE KEY-----",
					"Should contain private key for %s", tc.keyType)
			})
		}
	}
}

type alpineEnv struct {
	runner    *docker.Runner
	container *docker.StartResult
}

func getOrBuildOpenSSLEnv(t *testing.T, envs map[string]*alpineEnv, imageTag string) *alpineEnv {
	// Return cached runner and container for a given image tag if it exists
	if env, ok := envs[imageTag]; ok {
		return env
	}

	// Otherwise create the runner and container
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     opensslImageRepo,
		ImageTag:      imageTag,
		ContainerName: "openssl_pkcs12_" + uuid.New().String()[:8], // 8 chars is not guaranteed unique but should be fine for test containers
		Entrypoint:    []string{"sleep", "infinity"},               // Containers are cleaned up after subtests run
		LogConsumer: func(s string) {
			if t.Failed() {
				t.Logf("container logs: %s", s)
			}
		},
	})
	if err != nil {
		t.Fatalf("Could not provision docker service runner: %s", err)
	}

	result, err := runner.Start(context.Background(), true, false)
	if err != nil {
		t.Fatalf("Could not start container: %s (repo:%s, tag:%s)", err, opensslImageRepo, imageTag)
	}

	// Install OpenSSL in the container
	ctx := context.Background()
	installCmd := []string{"apk", "add", "--no-cache", "openssl"}
	_, stderr, retcode, err := runner.RunCmdWithOutput(ctx, result.Container.ID, installCmd)
	if err != nil || retcode != 0 {
		t.Fatalf("Failed to install OpenSSL version: %v, stderr: %s, (repo: %s, tag: %s)", err, string(stderr), opensslImageRepo, imageTag)
	}

	envs[imageTag] = &alpineEnv{runner: runner, container: result}
	return envs[imageTag]
}

func runOpenSSLCmd(runner *docker.Runner, containerID string, args []string) (string, error) {
	opensslCmd := append([]string{"openssl"}, args...)
	ctx := context.Background()
	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, containerID, opensslCmd)

	var errorMsg string
	if err != nil {
		errorMsg += fmt.Sprintf("error: %s", err)
	}
	if retcode != 0 {
		errorMsg += fmt.Sprintf("non-zero retcode: %v", retcode)
		if len(stderr) > 0 {
			errorMsg += fmt.Sprintf("stderr: %s", string(stderr))
		}
	}
	if errorMsg != "" {
		return "", fmt.Errorf("failed running command: %v, %s", opensslCmd, errorMsg)
	}
	// -info output is returned in stderr
	combined := string(stdout) + "\n" + string(stderr)
	return combined, nil
}

func runOpenSSLInfo(t *testing.T, envs map[string]*alpineEnv, pkcs12Bytes []byte, imageTag string) (string, error) {
	env := getOrBuildOpenSSLEnv(t, envs, imageTag)

	_, err := runOpenSSLCmd(env.runner, env.container.Container.ID, []string{"version"})
	if err != nil {
		return "", fmt.Errorf("failed to get openssl version: %w", err)
	}

	pfxCtx := docker.NewBuildContext()
	pfxCtx["bundle.p12"] = docker.PathContentsFromBytes(pkcs12Bytes)
	if err := env.runner.CopyTo(env.container.Container.ID, "/tmp/", pfxCtx); err != nil {
		return "", fmt.Errorf("could not copy pkcs12 bundle into container: %w", err)
	}

	return runOpenSSLCmd(env.runner, env.container.Container.ID, []string{
		"pkcs12",
		"-in", "/tmp/bundle.p12",
		"-info",
		"-noout",
		"-passin", "pass:" + pkcs12Password,
	})
}

// validateCertificateOrder validates that certificates are in the correct order (leaf first, then CA chain)
// hasPrivateKey indicates whether the PKCS12 bundle contains a private key (true for keystores, false for trust stores)
func validateCertificateOrder(t *testing.T, runner *docker.Runner, containerID string, hasPrivateKey bool) {
	// Count total certificates (should be exactly 2: leaf + CA)
	allCertsOutput, err := runOpenSSLCmd(runner, containerID, []string{
		"pkcs12",
		"-in", "/tmp/bundle.p12",
		"-nokeys",
		"-passin", "pass:" + pkcs12Password,
	})
	require.NoError(t, err, "OpenSSL should extract all certificates")
	certCount := strings.Count(allCertsOutput, "-----BEGIN CERTIFICATE-----")
	require.Equal(t, 2, certCount, "Should have exactly 2 certificates (leaf + CA)")
	require.Contains(t, allCertsOutput, "subject=CN=Root CA", "Should contain CA certificate")

	// For trust stores: No private key, so all certs are treated as CA certs
	// only validate leaf vs CA certs for bundles with private key.
	if hasPrivateKey {
		leafOutput, err := runOpenSSLCmd(runner, containerID, []string{
			"pkcs12",
			"-in", "/tmp/bundle.p12",
			"-clcerts", // extracts certs that have corresponding private keys (i.e. leaf certs)
			"-nokeys",
			"-passin", "pass:" + pkcs12Password,
		})
		require.NoError(t, err, "OpenSSL should extract leaf certificate")
		require.Contains(t, leafOutput, "subject=CN=localhost", "leaf cert should have correct subject")

		caOutput, err := runOpenSSLCmd(runner, containerID, []string{
			"pkcs12",
			"-in", "/tmp/bundle.p12",
			"-cacerts", // extract CA certificates only using
			"-nokeys",
			"-passin", "pass:" + pkcs12Password,
		})
		require.NoError(t, err, "OpenSSL should extract CA certificates")
		require.Contains(t, caOutput, "subject=CN=Root CA", "CA cert should have correct subject")
	}
}

// generateKeyAndCert creates a private key and certificate for the specified key type
func generateKeyAndCert(t *testing.T, keyType string, keyBits int) (crypto.Signer, *x509.Certificate, []*x509.Certificate, error) {
	container := &privateKeyContainer{}
	if err := certutil.GeneratePrivateKey(keyType, keyBits, container); err != nil {
		return nil, nil, nil, err
	}
	privateKey := container.key

	// Generate CA certificate
	caSerialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate CA serial: %w", err)
	}

	caTemplate := &x509.Certificate{
		SerialNumber: caSerialNumber,
		Subject: pkix.Name{
			CommonName: "Root CA",
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCertDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, privateKey.Public(), privateKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create CA certificate: %w", err)
	}

	caCert, err := x509.ParseCertificate(caCertDER)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Generate leaf certificate
	leafSerialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate leaf serial: %w", err)
	}

	leafTemplate := &x509.Certificate{
		SerialNumber: leafSerialNumber,
		Subject: pkix.Name{
			CommonName: "test.example.com",
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(2 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	leafCertDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, caCert, privateKey.Public(), privateKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create leaf certificate: %w", err)
	}

	leafCert, err := x509.ParseCertificate(leafCertDER)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse leaf certificate: %w", err)
	}

	return privateKey, leafCert, []*x509.Certificate{caCert}, nil
}

type privateKeyContainer struct {
	key           crypto.Signer
	keyType       certutil.PrivateKeyType
	serializedKey []byte
}

func (c *privateKeyContainer) SetParsedPrivateKey(key crypto.Signer, keyType certutil.PrivateKeyType, serializedKey []byte) {
	c.key = key
	c.keyType = keyType
	c.serializedKey = serializedKey
}
