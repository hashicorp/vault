// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pkiext

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	pkihelper "github.com/hashicorp/vault/helper/testhelpers/pki"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/stretchr/testify/require"
)

const (
	javaImageRepo          = "docker.mirror.hashicorp.services/ibm-semeru-runtimes"
	java21Image            = "open-21-jdk"
	java26Image            = "open-26-jdk"
	javaPKCS12DefaultAlias = "1"
)

// TestPKCS12JavaValidation validates PKCS#12 bundles with different encoders
// can be read by Java applications using keytool.
// This test focuses on PKCS12 format compatibility, not PKI backend functionality.
func TestPKCS12JavaValidation(t *testing.T) {
	t.Parallel()

	// Cleanup containers after tests complete
	envs := map[string]*javaEnv{}
	t.Cleanup(func() {
		for tag, env := range envs {
			if err := env.runner.Stop(context.Background(), env.container.Container.ID); err != nil {
				t.Logf("Warning: failed to stop container (tag: %s): %v", tag, err)
			}
		}
	})

	// Generate test CA, leaf and key once for all subtests
	leafKey, leafCert, caChain := generateChainAndLeafCert(t)

	// Test matrix:
	// 2 bundle types (with and without a private key) x 2 encoders (modern2026, modern2023) x 2 Java versions
	// - modern2026 encoder is incompatible with Java 21 (uses newer algorithms)
	// - modern2023 encoder works with both Java 21 and 26
	testCases := []struct {
		bundleType  string // "keystore" or "trust store"
		encoder     string
		version     string
		shouldError bool
	}{
		{bundleType: "keystore", encoder: "modern2026", version: java26Image},
		{bundleType: "keystore", encoder: "modern2023", version: java26Image},
		{bundleType: "keystore", encoder: "modern2026", version: java21Image, shouldError: true},
		{bundleType: "keystore", encoder: "modern2023", version: java21Image},

		{bundleType: "trust store", encoder: "modern2026", version: java26Image},
		{bundleType: "trust store", encoder: "modern2023", version: java26Image},
		{bundleType: "trust store", encoder: "modern2026", version: java21Image, shouldError: true},
		{bundleType: "trust store", encoder: "modern2023", version: java21Image},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("bundle=%s encoder=%s java=%s", tc.bundleType, tc.encoder, tc.version)
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

			// Validate with Java keytool
			keytoolOutput, err := runJavaKeytoolInspect(t, envs, pkcs12Bytes, tc.version)

			if tc.shouldError {
				// Expect keytool to fail for incompatible PKCS12 encoders (modern2026 with Java 21)
				require.Error(t, err, "keytool should fail to read PKCS12 file with incompatible encoder")
				// Verify it's the expected algorithm incompatibility error and not some other failure
				errMsg := err.Error()
				require.Contains(t, errMsg, "NoSuchAlgorithmException", "should fail due to algorithm not available")
				require.Contains(t, errMsg, "HmacPBE", "should fail due to HMAC-based PKCS12 algorithm")
				log.Info("Keytool failed as expected with algorithm incompatibility", "error", err)
				return
			}

			require.NoError(t, err, "keytool should successfully read PKCS12 file")
			log.Info("Java keytool output", "output", keytoolOutput)

			// Verify keytool can read the PKCS12 file
			require.Contains(t, keytoolOutput, "Keystore type: PKCS12", "keytool should recognize PKCS12 format")

			// Verify certificate chain and entry type based on bundle type
			if tc.bundleType == "keystore" {
				require.Contains(t, keytoolOutput, "Alias name: "+javaPKCS12DefaultAlias, "bundle with key should use numeric alias")
				require.Contains(t, keytoolOutput, "Certificate chain length: 2", "keytool should show complete certificate chain")
				require.Contains(t, keytoolOutput, "Entry type: PrivateKeyEntry", "bundle with key should create PrivateKeyEntry")
			} else {
				// Trust stores use CN-based aliases, not numeric ones
				require.Contains(t, keytoolOutput, "Entry type: trustedCertEntry", "trust store should create trustedCertEntry without private key")
				// Trust stores contain separate entries for each cert, not a chain
				require.Contains(t, keytoolOutput, "Your keystore contains 2 entries", "trust store should contain leaf cert and CA cert as separate entries")
			}

			// Verify CA certificate is present
			require.Contains(t, keytoolOutput, "Owner: CN=Root CA", "keytool should show CA certificate")
		})
	}
}

type javaEnv struct {
	runner    *docker.Runner
	container *docker.StartResult
}

func getOrBuildJavaEnv(t *testing.T, envs map[string]*javaEnv, imageTag string) *javaEnv {
	// Return cached runner and container for a given image tag if it exists
	if env, ok := envs[imageTag]; ok {
		return env
	}

	// Otherwise create the runner and container
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     javaImageRepo,
		ImageTag:      imageTag,
		ContainerName: "java_pkcs12_" + uuid.New().String()[:8], // 8 chars is not guaranteed unique but should be fine for test containers
		Entrypoint:    []string{"sleep", "infinity"},            // Containers are cleaned up after subtests run
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
		t.Fatalf("Could not start container: %s (repo:%s, tag:%s)", err, javaImageRepo, imageTag)
	}

	envs[imageTag] = &javaEnv{runner: runner, container: result}
	return envs[imageTag]
}

func runJavaCmd(t *testing.T, runner *docker.Runner, containerID string, cmd []string) (string, error) {
	ctx := context.Background()
	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("could not run command %v in container: %w", cmd, err)
	}

	if len(stderr) != 0 {
		t.Logf("Got stderr from command %v:%v", cmd, string(stderr))
	}

	if retcode != 0 {
		t.Logf("Got stdout from command %v:%v", cmd, string(stdout))
		// Return stdout as error because it contains the Java exception details which some tests expect
		return "", errors.New(string(stdout))
	}

	return string(stdout), nil
}

func runJavaKeytoolInspect(t *testing.T, envs map[string]*javaEnv, pkcs12Bytes []byte, imageTag string) (string, error) {
	env := getOrBuildJavaEnv(t, envs, imageTag)

	_, err := runJavaCmd(t, env.runner, env.container.Container.ID, []string{"java", "-version"})
	if err != nil {
		return "", fmt.Errorf("failed to get java version: %w", err)
	}

	_, err = runJavaCmd(t, env.runner, env.container.Container.ID, []string{"keytool", "-J-version"})
	if err != nil {
		return "", fmt.Errorf("failed to get keytool version: %w", err)
	}

	pfxCtx := docker.NewBuildContext()
	pfxCtx["bundle.p12"] = docker.PathContentsFromBytes(pkcs12Bytes)
	if err := env.runner.CopyTo(env.container.Container.ID, "/tmp/", pfxCtx); err != nil {
		return "", fmt.Errorf("could not copy pkcs12 bundle into container: %w", err)
	}

	return runJavaCmd(t, env.runner, env.container.Container.ID, []string{
		"keytool",
		"-list",
		"-v",
		"-storetype", "PKCS12",
		"-keystore", "/tmp/bundle.p12",
		"-storepass", pkcs12Password,
	})
}

// generateChainAndLeafCert creates a leaf certificate signed by the provided CA
func generateChainAndLeafCert(t *testing.T) (*ecdsa.PrivateKey, *x509.Certificate, []*x509.Certificate) {
	ca := pkihelper.GenerateRootCa(t)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	require.NoError(t, err)

	leafTemplate := &x509.Certificate{
		SerialNumber: serialNumber,
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

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	leafCertDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, ca.Cert, &leafKey.PublicKey, ca.Key)
	require.NoError(t, err)

	leafCert, err := x509.ParseCertificate(leafCertDER)
	require.NoError(t, err)

	return leafKey, leafCert, []*x509.Certificate{ca.Cert}
}
