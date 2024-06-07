// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package agent

import (
	"context"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	vaultcert "github.com/hashicorp/vault/builtin/credential/cert"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	agentcert "github.com/hashicorp/vault/command/agentproxyshared/auth/cert"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/file"
	"github.com/hashicorp/vault/helper/dhutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestCertEndToEnd(t *testing.T) {
	cases := []struct {
		name             string
		withCertRoleName bool
		ahWrapping       bool
	}{
		{
			"with name with wrapping",
			true,
			true,
		},
		{
			"with name without wrapping",
			true,
			false,
		},
		{
			"without name with wrapping",
			false,
			true,
		},
		{
			"without name without wrapping",
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testCertEndToEnd(t, tc.withCertRoleName, tc.ahWrapping)
		})
	}
}

func testCertEndToEnd(t *testing.T, withCertRoleName, ahWrapping bool) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"cert": vaultcert.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	// Setup Vault
	err := client.Sys().EnableAuthWithOptions("cert", &api.EnableAuthOptions{
		Type: "cert",
	})
	if err != nil {
		t.Fatal(err)
	}

	certificatePEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cluster.CACert.Raw})

	certRoleName := "test"
	_, err = client.Logical().Write(fmt.Sprintf("auth/cert/certs/%s", certRoleName), map[string]interface{}{
		"certificate": string(certificatePEM),
		"policies":    "default",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Generate encryption params
	pub, pri, err := dhutil.GeneratePublicPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	ouf, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	out := ouf.Name()
	ouf.Close()
	os.Remove(out)
	t.Logf("output: %s", out)

	dhpathf, err := ioutil.TempFile("", "auth.dhpath.test.")
	if err != nil {
		t.Fatal(err)
	}
	dhpath := dhpathf.Name()
	dhpathf.Close()
	os.Remove(dhpath)

	// Write DH public key to file
	mPubKey, err := jsonutil.EncodeJSON(&dhutil.PublicKeyInfo{
		Curve25519PublicKey: pub,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(dhpath, mPubKey, 0o600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote dh param file", "path", dhpath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	aaConfig := map[string]interface{}{}

	if withCertRoleName {
		aaConfig["name"] = certRoleName
	}

	am, err := agentcert.NewCertAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.cert"),
		MountPath: "auth/cert",
		Config:    aaConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	ahConfig := &auth.AuthHandlerConfig{
		Logger:                       logger.Named("auth.handler"),
		Client:                       client,
		EnableReauthOnNewCredentials: true,
	}
	if ahWrapping {
		ahConfig.WrapTTL = 10 * time.Second
	}
	ah := auth.NewAuthHandler(ahConfig)
	errCh := make(chan error)
	go func() {
		errCh <- ah.Run(ctx, am)
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	config := &sink.SinkConfig{
		Logger:    logger.Named("sink.file"),
		AAD:       "foobar",
		DHType:    "curve25519",
		DHPath:    dhpath,
		DeriveKey: true,
		Config: map[string]interface{}{
			"path": out,
		},
	}
	if !ahWrapping {
		config.WrapTTL = 10 * time.Second
	}
	fs, err := file.NewFileSink(config)
	if err != nil {
		t.Fatal(err)
	}
	config.Sink = fs

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: logger.Named("sink.server"),
		Client: client,
	})
	go func() {
		errCh <- ss.Run(ctx, ah.OutputCh, []*sink.SinkConfig{config}, ah.AuthInProgress)
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	// This has to be after the other defers so it happens first. It allows
	// successful test runs to immediately cancel all of the runner goroutines
	// and unblock any of the blocking defer calls by the runner's DoneCh that
	// comes before this and avoid successful tests from taking the entire
	// timeout duration.
	defer cancel()

	cloned, err := client.Clone()
	if err != nil {
		t.Fatal(err)
	}

	checkToken := func() string {
		timeout := time.Now().Add(5 * time.Second)
		for {
			if time.Now().After(timeout) {
				t.Fatal("did not find a written token after timeout")
			}
			val, err := ioutil.ReadFile(out)
			if err == nil {
				os.Remove(out)
				if len(val) == 0 {
					t.Fatal("written token was empty")
				}

				// First decrypt it
				resp := new(dhutil.Envelope)
				if err := jsonutil.DecodeJSON(val, resp); err != nil {
					continue
				}

				shared, err := dhutil.GenerateSharedSecret(pri, resp.Curve25519PublicKey)
				if err != nil {
					t.Fatal(err)
				}
				aesKey, err := dhutil.DeriveSharedKey(shared, pub, resp.Curve25519PublicKey)
				if err != nil {
					t.Fatal(err)
				}
				if len(aesKey) == 0 {
					t.Fatal("got empty aes key")
				}

				val, err = dhutil.DecryptAES(aesKey, resp.EncryptedPayload, resp.Nonce, []byte("foobar"))
				if err != nil {
					t.Fatalf("error: %v\nresp: %v", err, string(val))
				}

				// Now unwrap it
				wrapInfo := new(api.SecretWrapInfo)
				if err := jsonutil.DecodeJSON(val, wrapInfo); err != nil {
					t.Fatal(err)
				}
				switch {
				case wrapInfo.TTL != 10:
					t.Fatalf("bad wrap info: %v", wrapInfo.TTL)
				case !ahWrapping && wrapInfo.CreationPath != "sys/wrapping/wrap":
					t.Fatalf("bad wrap path: %v", wrapInfo.CreationPath)
				case ahWrapping && wrapInfo.CreationPath != "auth/cert/login":
					t.Fatalf("bad wrap path: %v", wrapInfo.CreationPath)
				case wrapInfo.Token == "":
					t.Fatal("wrap token is empty")
				}
				cloned.SetToken(wrapInfo.Token)
				secret, err := cloned.Logical().Unwrap("")
				if err != nil {
					t.Fatal(err)
				}
				if ahWrapping {
					switch {
					case secret.Auth == nil:
						t.Fatal("unwrap secret auth is nil")
					case secret.Auth.ClientToken == "":
						t.Fatal("unwrap token is nil")
					}
					return secret.Auth.ClientToken
				} else {
					switch {
					case secret.Data == nil:
						t.Fatal("unwrap secret data is nil")
					case secret.Data["token"] == nil:
						t.Fatal("unwrap token is nil")
					}
					return secret.Data["token"].(string)
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
	checkToken()
}

func TestCertEndToEnd_CertsInConfig(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"cert": vaultcert.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"pki": pki.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	// /////////////
	// PKI setup
	// /////////////

	// Mount /pki as a root CA
	err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set the cluster's certificate as the root CA in /pki
	pemBundleRootCA := string(cluster.CACertPEM) + string(cluster.CAKeyPEM)
	_, err = client.Logical().Write("pki/config/ca", map[string]interface{}{
		"pem_bundle": pemBundleRootCA,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Mount /pki2 to operate as an intermediate CA
	err = client.Sys().Mount("pki2", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a CSR for the intermediate CA
	secret, err := client.Logical().Write("pki2/intermediate/generate/internal", nil)
	if err != nil {
		t.Fatal(err)
	}
	intermediateCSR := secret.Data["csr"].(string)

	// Sign the intermediate CSR using /pki
	secret, err = client.Logical().Write("pki/root/sign-intermediate", map[string]interface{}{
		"permitted_dns_domains": ".myvault.com",
		"csr":                   intermediateCSR,
	})
	if err != nil {
		t.Fatal(err)
	}
	intermediateCertPEM := secret.Data["certificate"].(string)

	// Configure the intermediate cert as the CA in /pki2
	_, err = client.Logical().Write("pki2/intermediate/set-signed", map[string]interface{}{
		"certificate": intermediateCertPEM,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a role on the intermediate CA mount
	_, err = client.Logical().Write("pki2/roles/myvault-dot-com", map[string]interface{}{
		"allowed_domains":  "myvault.com",
		"allow_subdomains": "true",
		"max_ttl":          "5m",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue a leaf cert using the intermediate CA
	secret, err = client.Logical().Write("pki2/issue/myvault-dot-com", map[string]interface{}{
		"common_name": "cert.myvault.com",
		"format":      "pem",
		"ip_sans":     "127.0.0.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	leafCertPEM := secret.Data["certificate"].(string)
	leafCertKeyPEM := secret.Data["private_key"].(string)

	// Create temporary files for CA cert, client cert and client cert key.
	// This is used to configure TLS in the api client.
	caCertFile, err := ioutil.TempFile("", "caCert")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(caCertFile.Name())
	if _, err := caCertFile.Write([]byte(cluster.CACertPEM)); err != nil {
		t.Fatal(err)
	}
	if err := caCertFile.Close(); err != nil {
		t.Fatal(err)
	}

	leafCertFile, err := ioutil.TempFile("", "leafCert")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(leafCertFile.Name())
	if _, err := leafCertFile.WriteString(leafCertPEM); err != nil {
		t.Fatal(err)
	}
	if err := leafCertFile.Close(); err != nil {
		t.Fatal(err)
	}

	leafCertKeyFile, err := ioutil.TempFile("", "leafCertKey")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(leafCertKeyFile.Name())
	if _, err := leafCertKeyFile.WriteString(leafCertKeyPEM); err != nil {
		t.Fatal(err)
	}
	if err := leafCertKeyFile.Close(); err != nil {
		t.Fatal(err)
	}

	// /////////////
	// Cert auth setup
	// /////////////

	// Enable the cert auth method
	err = client.Sys().EnableAuthWithOptions("cert", &api.EnableAuthOptions{
		Type: "cert",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set the intermediate CA cert as a trusted certificate in the backend
	_, err = client.Logical().Write("auth/cert/certs/myvault-dot-com", map[string]interface{}{
		"display_name": "myvault.com",
		"policies":     "default",
		"certificate":  intermediateCertPEM,
	})
	if err != nil {
		t.Fatal(err)
	}

	// /////////////
	// Auth handler (auto-auth) setup
	// /////////////

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	am, err := agentcert.NewCertAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.cert"),
		MountPath: "auth/cert",
		Config: map[string]interface{}{
			"ca_cert":     caCertFile.Name(),
			"client_cert": leafCertFile.Name(),
			"client_key":  leafCertKeyFile.Name(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ahConfig := &auth.AuthHandlerConfig{
		Logger:                       logger.Named("auth.handler"),
		Client:                       client,
		EnableReauthOnNewCredentials: true,
	}

	ah := auth.NewAuthHandler(ahConfig)
	errCh := make(chan error)
	go func() {
		errCh <- ah.Run(ctx, am)
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	// /////////////
	// Sink setup
	// /////////////

	// Use TempFile to get us a generated file name to use for the sink.
	ouf, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	ouf.Close()
	out := ouf.Name()
	os.Remove(out)
	t.Logf("output: %s", out)

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": out,
		},
	}
	fs, err := file.NewFileSink(config)
	if err != nil {
		t.Fatal(err)
	}
	config.Sink = fs

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: logger.Named("sink.server"),
		Client: client,
	})
	go func() {
		errCh <- ss.Run(ctx, ah.OutputCh, []*sink.SinkConfig{config}, ah.AuthInProgress)
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	// This has to be after the other defers so it happens first. It allows
	// successful test runs to immediately cancel all of the runner goroutines
	// and unblock any of the blocking defer calls by the runner's DoneCh that
	// comes before this and avoid successful tests from taking the entire
	// timeout duration.
	defer cancel()

	// Read the token from the sink
	timeout := time.Now().Add(5 * time.Second)
	for {
		if time.Now().After(timeout) {
			t.Fatal("did not find a written token after timeout")
		}

		// Attempt to read the sink file until we get a token or the timeout is
		// reached.
		val, err := ioutil.ReadFile(out)
		if err == nil {
			os.Remove(out)
			if len(val) == 0 {
				t.Fatal("written token was empty")
			}

			t.Logf("sink token: %s", val)

			break
		}

		time.Sleep(250 * time.Millisecond)
	}
}
