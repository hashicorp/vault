// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pkiext

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	pkihelper "github.com/hashicorp/vault/helper/testhelpers/pki"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/stretchr/testify/require"
)

const (
	javaImageRepo   = "docker.mirror.hashicorp.services/ibm-semeru-runtimes"
	java8Image      = "open-8-jdk"
	java26Image     = "open-26-jdk"
	jksDefaultAlias = "1"
	bundlePassword  = "my-very-secure-password"
)

// TestPKCS12AndJKSJavaValidation verifies encoded PKCS#12 and JKS bundles, with different
// key types, can inspected by Java keytool.
// PKCS#12 coverage checks encoder compatibility across Java versions and bundle contents.
// JKS coverage checks alias handling and the expected store entries exist.
// This test focuses on format compatibility, not PKI backend functionality.
func TestPKCS12AndJKSJavaValidation(t *testing.T) {
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
	result := pkihelper.GenerateCertWithRoot(t)
	leafKey, leafCert, caChain := result.Leaf.Key, result.Leaf.Cert, []*x509.Certificate{result.RootCa.Cert}

	// Test matrix:
	// 2 bundle types (with and without a private key) x 2 encoders (modern2026, modern2023) x 2 Java versions
	// - modern2026 encoder is incompatible with Java 21 (uses newer algorithms)
	// - modern2023 encoder works with both Java 21 and 26
	testCasesPKCS12 := []struct {
		bundleType  string // "keystore" or "trust store"
		encoder     string
		version     string
		shouldError bool
	}{
		{bundleType: "keystore", encoder: "modern2026", version: java26Image},
		{bundleType: "keystore", encoder: "modern2023", version: java26Image},
		{bundleType: "keystore", encoder: "modern2026", version: java8Image, shouldError: true},
		{bundleType: "keystore", encoder: "modern2023", version: java8Image},

		{bundleType: "trust store", encoder: "modern2026", version: java26Image},
		{bundleType: "trust store", encoder: "modern2023", version: java26Image},
		{bundleType: "trust store", encoder: "modern2026", version: java8Image, shouldError: true},
		{bundleType: "trust store", encoder: "modern2023", version: java8Image},
	}

	for _, tc := range testCasesPKCS12 {
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
					bundlePassword)
			} else {
				pkcs12Bytes, err = pki.EncodeToPKCS12(
					tc.encoder, nil,
					leafCert,
					caChain,
					bundlePassword)
			}
			require.NoError(t, err, "EncodeToPKCS12 should succeed")
			require.NotEmpty(t, pkcs12Bytes)

			// Validate with Java keytool
			keytoolOutput, err := runJavaKeytoolInspect(t, envs, pkcs12Bytes, tc.version, "PKCS12")

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
				require.Contains(t, keytoolOutput, "Alias name: "+jksDefaultAlias, "bundle with key should use numeric alias")
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

	// Test matrix: 2 bundle types (with and without a private key) x different aliases x 2 Java versions
	testCasesJKS := []struct {
		name       string
		bundleType string // "keystore" or "trust store"
		version    string
		alias      string
	}{
		{name: "keystore default alias", bundleType: "keystore", version: java26Image, alias: jksDefaultAlias},
		{name: "keystore with java 8", bundleType: "keystore", version: java8Image, alias: jksDefaultAlias},
		{name: "keystore custom alias", bundleType: "keystore", version: java26Image, alias: "myapp"},

		{name: "trust store default alias", bundleType: "trust store", version: java26Image, alias: jksDefaultAlias},
		{name: "trust store with java 8", bundleType: "trust store", version: java8Image, alias: jksDefaultAlias},
		{name: "trust store non-numeric alias", bundleType: "trust store", version: java26Image, alias: "myapp"},
		{name: "trust store non-default numeric alias", bundleType: "trust store", version: java26Image, alias: "5"},
	}

	for _, tc := range testCasesJKS {
		t.Run(tc.name, func(t *testing.T) {
			log := corehelpers.NewTestLogger(t)

			var jksBytes []byte
			var err error
			if tc.bundleType == "keystore" {
				// Keystore with private key
				jksBytes, err = pki.EncodeToJKS(
					leafKey,
					leafCert,
					caChain,
					tc.alias,
					bundlePassword)
			} else {
				// Trust store without private key
				jksBytes, err = pki.EncodeToJKS(
					nil,
					leafCert,
					caChain,
					tc.alias,
					bundlePassword)
			}
			require.NoError(t, err, "EncodeToJKS should succeed")
			require.NotEmpty(t, jksBytes)

			// Validate with Java keytool
			keytoolOutput, err := runJavaKeytoolInspect(t, envs, jksBytes, tc.version, "JKS")
			require.NoError(t, err, "keytool should successfully read JKS file")
			log.Info("Java keytool output", "output", keytoolOutput)
			aliases := extractAliasNames(t, keytoolOutput)

			// Verify keytool can read the JKS file
			require.Contains(t, keytoolOutput, "Keystore type: JKS", "keytool should recognize JKS format")

			// Verify certificate chain and entry type based on bundle type
			if tc.bundleType == "keystore" {
				require.Equal(t, []string{tc.alias}, aliases, "keystore should contain exactly the provided alias")
				require.Contains(t, keytoolOutput, "Entry type: PrivateKeyEntry", "keystore should create PrivateKeyEntry")
				require.Contains(t, keytoolOutput, "Certificate chain length: 2", "keytool should show complete certificate chain (leaf + CA)")
			} else {
				// Trust stores contain separate entries for each cert
				require.Contains(t, keytoolOutput, "Entry type: trustedCertEntry", "trust store should create trustedCertEntry without private key")
				require.Contains(t, keytoolOutput, "Your keystore contains 2 entries", "trust store should contain leaf cert and CA cert as separate entries")
				require.NotContains(t, keytoolOutput, "Entry type: PrivateKeyEntry", "keystore should not create PrivateKeyEntry")
				require.ElementsMatch(t, expectedNumericAliases(2), aliases, "trust store should contain generated numeric entry aliases for each entry")
			}

			// Verify CA certificate is present
			require.Contains(t, keytoolOutput, "Owner: CN=Root CA", "keytool should show CA certificate")
			// Verify leaf certificate is present
			require.Contains(t, keytoolOutput, "Owner: CN=localhost", "keytool should show leaf certificate")
		})
	}

	// validates PKCS12/JKS bundles with different key types can be read by java
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
		// Generate key and certificate for key type
		privateKey, cert, caChain, err := generateKeyAndCert(t, tc.keyType, tc.keyBits)
		require.NoError(t, err, "Should generate cert and private key: %s, bits: %s", tc.keyType, tc.keyBits)

		// Test each encoder type
		for _, encoder := range []string{"modern2023", "modern2026"} {
			name := fmt.Sprintf("PKCS#12 encoder=%s key=%s", encoder, tc.name)

			t.Run(name, func(t *testing.T) {
				pkcs12Bytes, err := pki.EncodeToPKCS12(
					encoder,
					privateKey,
					cert,
					caChain,
					bundlePassword,
				)
				require.NoError(t, err, "Failed to encode PKCS#12 for key type: %s, bits:", tc.keyType, tc.keyBits)

				// Validate with Java keytool
				keytoolOutput, err := runJavaKeytoolInspect(t, envs, pkcs12Bytes, java26Image, "PKCS12")
				require.NoError(t, err, "keytool should successfully read PKCS12 file")
				// Verify keytool can read the PKCS12 file
				require.Contains(t, keytoolOutput, "Keystore type: PKCS12", "keytool should recognize PKCS12 format")
			})
		}

		// Test JKS
		t.Run("JKS "+tc.name, func(t *testing.T) {
			jksBytes, err := pki.EncodeToJKS(
				leafKey,
				leafCert,
				caChain,
				"1",
				bundlePassword)
			require.NoError(t, err, "EncodeToJKS should succeed")
			require.NotEmpty(t, jksBytes)

			// Validate with Java keytool
			keytoolOutput, err := runJavaKeytoolInspect(t, envs, jksBytes, java26Image, "JKS")
			require.NoError(t, err, "keytool should successfully read JKS file")
			// Verify keytool can read the JKS file
			require.Contains(t, keytoolOutput, "Keystore type: JKS", "keytool should recognize JKS format")
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

func runJavaKeytoolInspect(t *testing.T, envs map[string]*javaEnv, pkcs12Bytes []byte, imageTag string, storeType string) (string, error) {
	env := getOrBuildJavaEnv(t, envs, imageTag)

	_, err := runJavaCmd(t, env.runner, env.container.Container.ID, []string{"java", "-version"})
	if err != nil {
		return "", fmt.Errorf("failed to get java version: %w", err)
	}

	_, err = runJavaCmd(t, env.runner, env.container.Container.ID, []string{"keytool", "-J-version"})
	if err != nil {
		return "", fmt.Errorf("failed to get keytool version: %w", err)
	}

	file := "bundle.p12"
	if storeType == "JKS" {
		file = "bundle.jks"
	}
	pfxCtx := docker.NewBuildContext()
	pfxCtx[file] = docker.PathContentsFromBytes(pkcs12Bytes)
	if err := env.runner.CopyTo(env.container.Container.ID, "/tmp/", pfxCtx); err != nil {
		return "", fmt.Errorf("could not copy bundle into container for store type: %s, Error: %w", storeType, err)
	}

	return runJavaCmd(t, env.runner, env.container.Container.ID, []string{
		"keytool",
		"-list",
		"-v",
		"-storetype", storeType,
		"-keystore", "/tmp/" + file,
		"-storepass", bundlePassword,
	})
}

func extractAliasNames(t *testing.T, keytoolOutput string) []string {
	t.Helper()

	aliasPattern := regexp.MustCompile(`(?m)^Alias name: (.+)$`)
	matches := aliasPattern.FindAllStringSubmatch(keytoolOutput, -1)
	aliases := make([]string, 0, len(matches))
	for _, match := range matches {
		require.Len(t, match, 2, "alias match should include the full line and captured alias")
		aliases = append(aliases, match[1])
	}

	return aliases
}

func expectedNumericAliases(count int) []string {
	aliases := make([]string, 0, count)
	for index := 1; index <= count; index++ {
		aliases = append(aliases, strconv.Itoa(index))
	}

	return aliases
}
