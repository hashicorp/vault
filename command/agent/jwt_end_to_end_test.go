// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	vaultjwt "github.com/hashicorp/vault-plugin-auth-jwt"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	agentjwt "github.com/hashicorp/vault/command/agentproxyshared/auth/jwt"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/file"
	"github.com/hashicorp/vault/helper/dhutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestJWTEndToEnd(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		ahWrapping            bool
		useSymlink            bool
		removeJWTAfterReading bool
	}{
		{false, false, false},
		{true, false, false},
		{false, true, false},
		{true, true, false},
		{false, false, true},
		{true, false, true},
		{false, true, true},
		{true, true, true},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(fmt.Sprintf("ahWrapping=%v, useSymlink=%v, removeJWTAfterReading=%v", tc.ahWrapping, tc.useSymlink, tc.removeJWTAfterReading), func(t *testing.T) {
			t.Parallel()
			testJWTEndToEnd(t, tc.ahWrapping, tc.useSymlink, tc.removeJWTAfterReading)
		})
	}
}

func testJWTEndToEnd(t *testing.T, ahWrapping, useSymlink, removeJWTAfterReading bool) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"jwt": vaultjwt.Factory,
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
	err := client.Sys().EnableAuthWithOptions("jwt", &api.EnableAuthOptions{
		Type: "jwt",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/jwt/config", map[string]interface{}{
		"bound_issuer":           "https://team-vault.auth0.com/",
		"jwt_validation_pubkeys": TestECDSAPubKey,
		"jwt_supported_algs":     "ES256",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/jwt/role/test", map[string]interface{}{
		"role_type":       "jwt",
		"bound_subject":   "r3qXcK2bix9eFECzsU3Sbmh0K16fatW6@clients",
		"bound_audiences": "https://vault.plugin.auth.jwt.test",
		"user_claim":      "https://vault/user",
		"groups_claim":    "https://vault/groups",
		"policies":        "test",
		"period":          "3s",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Generate encryption params
	pub, pri, err := dhutil.GeneratePublicPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	// We close these right away because we're just basically testing
	// permissions and finding a usable file name
	inf, err := os.CreateTemp("", "auth.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	in := inf.Name()
	inf.Close()
	os.Remove(in)
	symlink, err := os.CreateTemp("", "auth.jwt.symlink.test.")
	if err != nil {
		t.Fatal(err)
	}
	symlinkName := symlink.Name()
	symlink.Close()
	os.Remove(symlinkName)
	os.Symlink(in, symlinkName)
	t.Logf("input: %s", in)

	ouf, err := os.CreateTemp("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	out := ouf.Name()
	ouf.Close()
	os.Remove(out)
	t.Logf("output: %s", out)

	dhpathf, err := os.CreateTemp("", "auth.dhpath.test.")
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
	if err := os.WriteFile(dhpath, mPubKey, 0o600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote dh param file", "path", dhpath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	var fileNameToUseAsPath string
	if useSymlink {
		fileNameToUseAsPath = symlinkName
	} else {
		fileNameToUseAsPath = in
	}
	am, err := agentjwt.NewJWTAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.jwt"),
		MountPath: "auth/jwt",
		Config: map[string]interface{}{
			"path":                        fileNameToUseAsPath,
			"role":                        "test",
			"remove_jwt_after_reading":    removeJWTAfterReading,
			"remove_jwt_follows_symlinks": true,
			"jwt_read_period":             "0.5s",
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

	// Check that no jwt file exists
	_, err = os.Lstat(in)
	if err == nil {
		t.Fatal("expected err")
	}
	if !os.IsNotExist(err) {
		t.Fatal("expected notexist err")
	}
	_, err = os.Lstat(out)
	if err == nil {
		t.Fatal("expected err")
	}
	if !os.IsNotExist(err) {
		t.Fatal("expected notexist err")
	}

	cloned, err := client.Clone()
	if err != nil {
		t.Fatal(err)
	}

	// Get a token
	jwtToken, _ := GetTestJWT(t)

	if err := os.WriteFile(in, []byte(jwtToken), 0o600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test jwt", "path", in)
	}

	checkToken := func() string {
		timeout := time.Now().Add(5 * time.Second)
		for {
			if time.Now().After(timeout) {
				t.Fatal("did not find a written token after timeout")
			}
			val, err := os.ReadFile(out)
			if err == nil {
				os.Remove(out)
				if len(val) == 0 {
					t.Fatal("written token was empty")
				}

				// First, ensure JWT has been removed
				if removeJWTAfterReading {
					_, err = os.Stat(in)
					if err == nil {
						t.Fatal("no error returned from stat, indicating the jwt is still present")
					}
					if !os.IsNotExist(err) {
						t.Fatalf("unexpected error: %v", err)
					}
				} else {
					_, err := os.Stat(in)
					if err != nil {
						t.Fatal("JWT file removed despite removeJWTAfterReading being set to false")
					}
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
				case ahWrapping && wrapInfo.CreationPath != "auth/jwt/login":
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
	origToken := checkToken()

	// We only check this if the renewer is actually renewing for us
	if !ahWrapping {
		// Period of 3 seconds, so should still be alive after 7
		timeout := time.Now().Add(7 * time.Second)
		cloned.SetToken(origToken)
		for {
			if time.Now().After(timeout) {
				break
			}
			secret, err := cloned.Auth().Token().LookupSelf()
			if err != nil {
				t.Fatal(err)
			}
			ttl, err := secret.Data["ttl"].(json.Number).Int64()
			if err != nil {
				t.Fatal(err)
			}
			if ttl > 3 {
				t.Fatalf("unexpected ttl: %v", secret.Data["ttl"])
			}
		}
	}

	// Get another token to test the backend pushing the need to authenticate
	// to the handler
	jwtToken, _ = GetTestJWT(t)
	if err := os.WriteFile(in, []byte(jwtToken), 0o600); err != nil {
		t.Fatal(err)
	}

	newToken := checkToken()
	if newToken == origToken {
		t.Fatal("found same token written")
	}

	if !ahWrapping {
		// Repeat the period test. At the end the old token should have expired and
		// the new token should still be alive after 7
		timeout := time.Now().Add(7 * time.Second)
		cloned.SetToken(newToken)
		for {
			if time.Now().After(timeout) {
				break
			}
			secret, err := cloned.Auth().Token().LookupSelf()
			if err != nil {
				t.Fatal(err)
			}
			ttl, err := secret.Data["ttl"].(json.Number).Int64()
			if err != nil {
				t.Fatal(err)
			}
			if ttl > 3 {
				t.Fatalf("unexpected ttl: %v", secret.Data["ttl"])
			}
		}

		cloned.SetToken(origToken)
		_, err = cloned.Auth().Token().LookupSelf()
		if err == nil {
			t.Fatal("expected error")
		}
	}
}
