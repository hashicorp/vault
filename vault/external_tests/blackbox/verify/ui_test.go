// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package verify

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestVaultUIAvailability verifies that the Vault UI is accessible and properly configured.
// It checks that:
// 1. The root URL redirects to /ui/
// 2. The UI page loads successfully and doesn't show an error message
func TestVaultUIAvailability(t *testing.T) {
	v := blackbox.New(t)

	// Get the Vault address from the client config
	vaultAddr := v.Client.Address()

	// Create HTTP client that doesn't follow redirects automatically
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Test 1: Verify root URL redirects to /ui/
	t.Run("redirect_to_ui", func(t *testing.T) {
		resp, err := client.Get(vaultAddr)
		if err != nil {
			t.Fatalf("Failed to request Vault root URL: %v", err)
		}
		defer resp.Body.Close()

		// Check for redirect status code
		if resp.StatusCode != http.StatusMovedPermanently && resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusSeeOther && resp.StatusCode != http.StatusTemporaryRedirect {
			t.Fatalf("Expected redirect status code (301, 302, 303, 307), got: %d", resp.StatusCode)
		}

		// Check redirect location
		location := resp.Header.Get("Location")
		expectedSuffix := "/ui/"
		if !strings.HasSuffix(location, expectedSuffix) {
			t.Fatalf("Expected redirect to end with %s, got: %s", expectedSuffix, location)
		}

		t.Logf("Successfully verified redirect from %s to %s", vaultAddr, location)
	})

	// Test 2: Verify UI page loads and is available
	t.Run("ui_page_loads", func(t *testing.T) {
		uiURL := fmt.Sprintf("%s/ui/", vaultAddr)

		// Use default client that follows redirects
		resp, err := http.Get(uiURL)
		if err != nil {
			t.Fatalf("Failed to request Vault UI: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for UI page, got: %d", resp.StatusCode)
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read UI response body: %v", err)
		}

		bodyStr := string(body)

		// Check that the UI is not showing an error message
		if strings.Contains(bodyStr, "Vault UI is not available") {
			t.Fatal("Vault UI is not available - error message found in response")
		}

		// Verify we got some HTML content (basic sanity check)
		if !strings.Contains(bodyStr, "<html") && !strings.Contains(bodyStr, "<!DOCTYPE") {
			t.Fatal("UI response doesn't appear to be HTML content")
		}

		t.Logf("Successfully verified UI is available at %s", uiURL)
	})
}
