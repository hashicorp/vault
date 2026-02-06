// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/stretchr/testify/require"
)

// AssertUIAvailable performs a raw HTTP GET request to the Vault UI
// to ensure it returns a 200 OK and serves HTML.
func (s *Session) AssertUIAvailable() {
	s.t.Helper()

	// client.Address() returns the API address (e.g. http://127.0.0.1:8200)
	// The UI is usually at /ui/
	uiURL := fmt.Sprintf("%s/ui/", s.Client.Address())

	resp, err := http.Get(uiURL)
	require.NoError(s.t, err)
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	require.Equal(s.t, http.StatusOK, resp.StatusCode, "UI endpoint returned non-200 status")

	// Optional: Check Content-Type
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		s.t.Fatalf("Expected text/html content type for UI, got %s", ct)
	}
}

// AssertFileContainsSecret scans a file on the local disk (audit log)
// and ensures the provided secret string is NOT present.
func (s *Session) AssertFileDoesNotContainSecret(filePath, secretValue string) {
	s.t.Helper()

	if secretValue == "" {
		return
	}

	content, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		s.t.Fatalf("Audit log file not found: %s", filePath)
	}
	require.NoError(s.t, err)

	fileBody := string(content)
	if strings.Contains(fileBody, secretValue) {
		s.t.Fatalf("Security Violation: Found secret value in file %s", filePath)
	}
}
