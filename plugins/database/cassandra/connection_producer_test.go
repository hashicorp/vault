// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cassandra

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/testhelpers/cassandra"
	pkihelper "github.com/hashicorp/vault/helper/testhelpers/pki"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/stretchr/testify/require"
)

var insecureFileMounts = map[string]string{
	"test-fixtures/no_tls/cassandra.yaml": "/etc/cassandra/cassandra.yaml",
}

func TestSelfSignedCA(t *testing.T) {
	// Generate certificates dynamically to avoid expiration issues
	certBundle := generateCassandraCerts(t)

	copyFromTo := map[string]string{
		certBundle.StoresDir:   "/bitnami/cassandra/secrets/",
		certBundle.CqlshrcFile: "/.cassandra/cqlshrc",
	}

	tlsConfig := &tls.Config{
		RootCAs: certBundle.CertPool,
	}
	// The test container is reached through a loopback port mapping, and the generated leaf cert is only valid
	// for localhost and 127.0.0.1. Set ServerName explicitly so hostname verification matches the test cert.
	tlsConfig.ServerName = "localhost"
	sslOpts := &gocql.SslOptions{
		Config:                 tlsConfig,
		EnableHostVerification: true,
	}

	host, cleanup := cassandra.PrepareTestContainer(t,
		cassandra.ContainerName("cassandra"),
		cassandra.Image("bitnamilegacy/cassandra", "3.11.11"),
		cassandra.CopyFromTo(copyFromTo),
		cassandra.SslOpts(sslOpts),
		cassandra.Env("CASSANDRA_KEYSTORE_PASSWORD=cassandra"),
		cassandra.Env("CASSANDRA_TRUSTSTORE_PASSWORD=cassandra"),
		cassandra.Env("CASSANDRA_INTERNODE_ENCRYPTION=none"),
		cassandra.Env("CASSANDRA_CLIENT_ENCRYPTION=true"),
	)
	t.Cleanup(cleanup)

	type testCase struct {
		config    map[string]interface{}
		expectErr bool
	}

	caPEM := certBundle.CaPEM
	badCAPEM := certBundle.BadCaPEM

	tests := map[string]testCase{
		// ///////////////////////
		// pem_json tests
		"pem_json/ca only": {
			config: map[string]interface{}{
				"pem_json": toJSON(t, certutil.CertBundle{
					CAChain: []string{caPEM},
				}),
			},
			expectErr: false,
		},
		"pem_json/bad ca": {
			config: map[string]interface{}{
				"pem_json": toJSON(t, certutil.CertBundle{
					CAChain: []string{badCAPEM},
				}),
			},
			expectErr: true,
		},
		"pem_json/missing ca": {
			config: map[string]interface{}{
				"pem_json": "",
			},
			expectErr: true,
		},

		// ///////////////////////
		// pem_bundle tests
		"pem_bundle/ca only": {
			config: map[string]interface{}{
				"pem_bundle": caPEM,
			},
			expectErr: false,
		},
		"pem_bundle/unrecognized CA": {
			config: map[string]interface{}{
				"pem_bundle": badCAPEM,
			},
			expectErr: true,
		},
		"pem_bundle/missing ca": {
			config: map[string]interface{}{
				"pem_bundle": "",
			},
			expectErr: true,
		},

		// ///////////////////////
		// no cert data provided
		"no cert data/tls=true": {
			config: map[string]interface{}{
				"tls": "true",
			},
			expectErr: true,
		},
		"no cert data/tls=false": {
			config: map[string]interface{}{
				"tls": "false",
			},
			expectErr: true,
		},
		"no cert data/insecure_tls": {
			config: map[string]interface{}{
				"insecure_tls": "true",
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Set values that we don't know until the cassandra container is started
			config := map[string]interface{}{
				"hosts":            host.Name,
				"port":             host.Port,
				"username":         "cassandra",
				"password":         "cassandra",
				"protocol_version": "4",
				"connect_timeout":  "30s",
				"tls":              "true",

				// Match the generated test certificate, which is valid for localhost and 127.0.0.1.
				"tls_server_name": "localhost",
			}

			// Apply the generated & common fields to the config to be sent to the DB
			for k, v := range test.config {
				config[k] = v
			}

			db := new()
			initReq := dbplugin.InitializeRequest{
				Config:           config,
				VerifyConnection: true,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := db.Initialize(ctx, initReq)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			// If no error expected, run a NewUser query to make sure the connection
			// actually works in case Initialize doesn't catch it
			if !test.expectErr {
				assertNewUser(t, db, sslOpts)
			}
		})
	}
}

func assertNewUser(t *testing.T, db *Cassandra, sslOpts *gocql.SslOptions) {
	newUserReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "dispname",
			RoleName:    "rolename",
		},
		Statements: dbplugin.Statements{
			Commands: []string{
				"create user '{{username}}' with password '{{password}}'",
			},
		},
		RollbackStatements: dbplugin.Statements{},
		Password:           "gh8eruajASDFAsgy89svn",
		Expiration:         time.Now().Add(5 * time.Second),
	}

	newUserResp := dbtesting.AssertNewUser(t, db, newUserReq)
	t.Logf("Username: %s", newUserResp.Username)

	assertCreds(t, db.Hosts, db.Port, newUserResp.Username, newUserReq.Password, sslOpts, 5*time.Second)
}

type cassandraCertBundle struct {
	CaPEM       string
	BadCaPEM    string
	CertPool    *x509.CertPool
	StoresDir   string
	CqlshrcFile string
}

// generateCassandraCerts generates fresh certificates for Cassandra testing
func generateCassandraCerts(t *testing.T) *cassandraCertBundle {
	t.Helper()

	// Generate CA and server certificate using pkihelper
	result := pkihelper.GenerateCertWithRoot(t)

	// Create a temporary directory for all certificate files
	tempDir := t.TempDir()
	storesDir := filepath.Join(tempDir, "stores")
	err := os.MkdirAll(storesDir, 0o755)
	require.NoError(t, err)

	// Write CA PEM
	caPEM := string(pem.EncodeToMemory(result.RootCa.CertPem))

	// Generate a "bad" CA for negative testing
	badCA := pkihelper.GenerateRootCa(t)
	badCaPEM := string(pem.EncodeToMemory(badCA.CertPem))

	// Create cert pool for TLS config
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM([]byte(caPEM))

	// Create PKCS12 keystore and truststore for Cassandra
	// Cassandra needs: keystore (server cert + key) and truststore (CA cert)
	createCassandraStores(t, storesDir, result)

	// Create cqlshrc file
	cqlshrcFile := filepath.Join(tempDir, "cqlshrc")
	cqlshrcContent := `[ssl]
validate = false
version = SSLv23
`
	err = os.WriteFile(cqlshrcFile, []byte(cqlshrcContent), 0o644)
	require.NoError(t, err)

	return &cassandraCertBundle{
		CaPEM:       caPEM,
		BadCaPEM:    badCaPEM,
		CertPool:    certPool,
		StoresDir:   storesDir,
		CqlshrcFile: cqlshrcFile,
	}
}

// createCassandraStores creates Java keystore and truststore files for Cassandra
func createCassandraStores(t *testing.T, storesDir string, result pkihelper.LeafWithRoot) {
	t.Helper()

	// Check if keytool is available
	_, err := exec.LookPath("keytool")
	if err != nil {
		t.Skip("keytool not found in PATH, skipping keystore generation")
	}

	// Create PKCS12 file with server certificate and key
	p12File := filepath.Join(storesDir, "server.p12")

	// First, create a temporary PEM file with both cert and key
	certKeyPEM := filepath.Join(storesDir, "server-cert-key.pem")
	keyBytes, err := x509.MarshalECPrivateKey(result.Leaf.Key)
	require.NoError(t, err)
	certKeyData := append(pem.EncodeToMemory(result.Leaf.CertPem), pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	})...)
	err = os.WriteFile(certKeyPEM, certKeyData, 0o644)
	require.NoError(t, err)

	// Convert to PKCS12 using openssl
	_, err = exec.LookPath("openssl")
	if err != nil {
		t.Skip("openssl not found in PATH, skipping keystore generation")
	}

	cmd := exec.Command("openssl", "pkcs12", "-export",
		"-in", certKeyPEM,
		"-name", "cassandra",
		"-out", p12File,
		"-passout", "pass:cassandra")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create PKCS12 file: %v\nOutput: %s", err, output)
	}

	// Create keystore from PKCS12
	keystoreFile := filepath.Join(storesDir, "keystore")
	cmd = exec.Command("keytool", "-importkeystore",
		"-srckeystore", p12File,
		"-srcstoretype", "PKCS12",
		"-srcstorepass", "cassandra",
		"-destkeystore", keystoreFile,
		"-deststoretype", "JKS",
		"-deststorepass", "cassandra",
		"-noprompt")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create keystore: %v\nOutput: %s", err, output)
	}

	// Create truststore with CA certificate
	truststoreFile := filepath.Join(storesDir, "truststore")
	caCertFile := filepath.Join(storesDir, "ca.pem")
	err = os.WriteFile(caCertFile, pem.EncodeToMemory(result.RootCa.CertPem), 0o644)
	require.NoError(t, err)

	cmd = exec.Command("keytool", "-import",
		"-file", caCertFile,
		"-alias", "cassandra-ca",
		"-keystore", truststoreFile,
		"-storetype", "JKS",
		"-storepass", "cassandra",
		"-noprompt")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create truststore: %v\nOutput: %s", err, output)
	}
}

func toJSON(t *testing.T, val interface{}) string {
	t.Helper()
	b, err := json.Marshal(val)
	require.NoError(t, err)
	return string(b)
}
