// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package template

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/config"
	"github.com/stretchr/testify/require"
)

func TestPKIRotationHandlerValidation(t *testing.T) {
	for _, tc := range []struct {
		name, method, body   string
		enabled, ready, once bool
		want                 int
	}{
		{name: "disabled", method: "POST", body: `{"destination":"/cert.pem"}`, want: 404},
		{name: "method", method: "GET", enabled: true, want: 405},
		{name: "missing destination", method: "POST", body: `{}`, enabled: true, want: 400},
		{name: "unknown field", method: "POST", body: `{"destination":"/cert.pem","force":true}`, enabled: true, want: 400},
		{name: "trailing JSON", method: "POST", body: `{"destination":"/cert.pem"}{}`, enabled: true, want: 400},
		{name: "oversized", method: "POST", body: `{"destination":"` + strings.Repeat("x", 4096) + `"}`, enabled: true, want: 400},
		{name: "unknown destination", method: "POST", body: `{"destination":"/other.pem"}`, enabled: true, want: 404},
		{name: "not ready", method: "POST", body: `{"destination":"/cert.pem"}`, enabled: true, want: 503},
		{name: "one shot", method: "POST", body: `{"destination":"/cert.pem"}`, enabled: true, ready: true, once: true, want: 409},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(&ServerConfig{Logger: hclog.NewNullLogger(), ExitAfterAuth: tc.once, AgentConfig: &config.Config{Templates: []*ctconfig.TemplateConfig{{Destination: ctconfig.String("/cert.pem")}}}})
			server.runnerStarted.Store(tc.ready)
			response := httptest.NewRecorder()
			server.PKIRotationHandler(tc.enabled).ServeHTTP(response, httptest.NewRequest(tc.method, "/agent/v1/templates/rotate", strings.NewReader(tc.body)))
			require.Equal(t, tc.want, response.Code)
		})
	}
}

func TestPKIRotationHandlerDispatch(t *testing.T) {
	for _, count := range []int{0, 1} {
		t.Run(string(rune('0'+count)), func(t *testing.T) {
			server := NewServer(&ServerConfig{Logger: hclog.NewNullLogger(), AgentConfig: &config.Config{Templates: []*ctconfig.TemplateConfig{{Destination: ctconfig.String("/cert.pem")}}}})
			server.runnerStarted.Store(true)
			done := make(chan string, 1)
			go func() { request := <-server.rotateRequests; done <- request.destination; request.result <- count }()
			response := httptest.NewRecorder()
			server.PKIRotationHandler(true).ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/agent/v1/templates/rotate", strings.NewReader(`{"destination":"/cert.pem"}`)))
			require.Equal(t, "/cert.pem", <-done)
			if count == 0 {
				require.Equal(t, 409, response.Code)
			} else {
				require.Equal(t, 202, response.Code)
				require.JSONEq(t, `{"destination":"/cert.pem","dependencies":1,"status":"accepted"}`, response.Body.String())
			}
		})
	}
}

func TestPKIRotationHandlerStoppedOrCanceled(t *testing.T) {
	for _, stopped := range []bool{false, true} {
		server := NewServer(&ServerConfig{Logger: hclog.NewNullLogger(), AgentConfig: &config.Config{Templates: []*ctconfig.TemplateConfig{{Destination: ctconfig.String("/cert.pem")}}}})
		server.runnerStarted.Store(true)
		request := httptest.NewRequest(http.MethodPost, "/agent/v1/templates/rotate", strings.NewReader(`{"destination":"/cert.pem"}`))
		if stopped {
			server.Stop()
		} else {
			ctx, cancel := context.WithCancel(request.Context())
			cancel()
			request = request.WithContext(ctx)
		}
		response := httptest.NewRecorder()
		server.PKIRotationHandler(true).ServeHTTP(response, request)
		require.Equal(t, 503, response.Code)
	}
	var server *Server
	response := httptest.NewRecorder()
	server.PKIRotationHandler(true).ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/agent/v1/templates/rotate", strings.NewReader(`{"destination":"/cert.pem"}`)))
	require.Equal(t, 503, response.Code)
}

func rotationCertificate(t *testing.T, serial int64) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	cert := &x509.Certificate{SerialNumber: big.NewInt(serial), NotBefore: time.Now().Add(-time.Minute), NotAfter: time.Now().Add(90 * 24 * time.Hour)}
	der, err := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	require.NoError(t, err)
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

func TestPKIRotationThroughRunningServer(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("VAULT_NAMESPACE", "")
	old, replacement := rotationCertificate(t, 1), rotationCertificate(t, 2)
	var calls atomic.Int32
	vault := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/pki/issue/test" || r.Header.Get("X-Vault-Token") != "same-token" {
			t.Errorf("unexpected Vault request: %s", r.URL.Path)
			w.WriteHeader(403)
			return
		}
		calls.Add(1)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"certificate": replacement}})
	}))
	defer vault.Close()
	dir := t.TempDir()
	dest, hook := filepath.Join(dir, "cert.pem"), filepath.Join(dir, "hook.pem")
	require.NoError(t, os.WriteFile(dest, []byte(old), 0600))
	templates := []*ctconfig.TemplateConfig{{
		Contents:    ctconfig.String(`{{ with pkiCert "pki/issue/test" }}{{ .Cert }}{{ end }}`),
		Destination: ctconfig.String(dest),
		Exec:        &ctconfig.ExecConfig{Command: []string{os.Args[0], "-test.run=^TestPKIRotationExecHelper$", "--", "pki-rotation-helper", dest, hook}},
	}}
	agentConfig := &config.Config{Vault: &config.Vault{Address: vault.URL}, Templates: templates}
	server := NewServer(&ServerConfig{Logger: hclog.NewNullLogger(), LogWriter: io.Discard, AgentConfig: agentConfig})
	ctx, cancel := context.WithCancel(context.Background())
	incoming := make(chan string, 1)
	incoming <- "same-token"
	done := make(chan error, 1)
	go func() { done <- server.Run(ctx, incoming, templates, new(atomic.Bool), make(chan error, 1)) }()
	defer func() {
		cancel()
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Error("template server did not stop")
		}
		server.Stop()
	}()
	require.Eventually(t, server.runnerStarted.Load, 5*time.Second, 10*time.Millisecond)
	initialRunner := server.runner
	require.Eventually(t, func() bool { return len(initialRunner.RenderEvents()) == 1 }, 5*time.Second, 10*time.Millisecond)
	require.Zero(t, calls.Load())
	body, err := json.Marshal(map[string]string{"destination": dest})
	require.NoError(t, err)
	response := httptest.NewRecorder()
	server.PKIRotationHandler(true).ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/agent/v1/templates/rotate", strings.NewReader(string(body))))
	require.Equal(t, http.StatusAccepted, response.Code, response.Body.String())
	require.Eventually(t, func() bool { b, _ := os.ReadFile(hook); return string(b) == replacement }, 5*time.Second, 10*time.Millisecond)
	require.Same(t, initialRunner, server.runner)
	require.EqualValues(t, 1, calls.Load())
}

func TestPKIRotationExecHelper(t *testing.T) {
	args := os.Args
	if len(args) < 4 || args[len(args)-3] != "pki-rotation-helper" {
		return
	}
	data, err := os.ReadFile(args[len(args)-2])
	if err != nil {
		os.Exit(1)
	}
	if os.WriteFile(args[len(args)-1], data, 0600) != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func TestPKIRotationListenerConfiguration(t *testing.T) {
	t.Setenv("TEST_AAD_ENV", "aad")
	data, err := os.ReadFile("../config/test-fixtures/config-auto-auth-aws.hcl")
	require.NoError(t, err)
	for _, enabled := range []bool{false, true} {
		contents := string(data)
		if enabled {
			contents = strings.ReplaceAll(contents, "enable_quit = true", "enable_quit = true\n    enable_pki_rotate = true")
		}
		path := filepath.Join(t.TempDir(), "agent.hcl")
		require.NoError(t, os.WriteFile(path, []byte(contents), 0600))
		cfg, err := config.LoadConfigFile(path)
		require.NoError(t, err)
		require.Equal(t, enabled, cfg.Listeners[0].AgentAPI.EnablePKIRotate)
	}
}
