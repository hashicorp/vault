package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestRenewLeaseThroughCache validates that cached lease renew responses
// that include an "Age" header will correctly adjust the TTL age so that
// API consumers work with the true TTL and not the original cached TTL.
func TestRenewLeaseThroughCache(t *testing.T) {
	age := time.Hour * 10
	ttl := time.Hour * 24
	expectedDuration := int((ttl - age).Seconds())
	mockVaultAgentCache := httptest.NewServer(http.HandlerFunc(agedVaultCacheResponseHandler(age, ttl)))
	defer mockVaultAgentCache.Close()

	cfg := DefaultConfig()
	cfg.AgentAddress = mockVaultAgentCache.URL

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Sys().Renew("112321-bd6b-818g-edbb-e462338bb0aa", 1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.LeaseDuration != expectedDuration {
		t.Fatalf("expected lease duration to be %d seconds not %d  seconds", expectedDuration, resp.LeaseDuration)
	}
}

func agedVaultCacheResponseHandler(age, ttl time.Duration) func(http.ResponseWriter, *http.Request) {
	ageStr := fmt.Sprintf("%.0f", (age).Seconds())
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Age", ageStr)

		renewResponseTemplate := `{
			"request_id": "82601a91-cd7a-718f-feca-f573449cc1bb",
			"lease_id": "112321-bd6b-818g-edbb-e462338bb0aa",
			"renewable": true,
			"lease_duration": %.0f,
			"data": {
			},
			"warp_info": null,
			"warnings": null,
			"auth": null
		}`

		renewResponse := fmt.Sprintf(renewResponseTemplate, (ttl).Seconds())
		_, _ = w.Write([]byte(renewResponse))
	}
}
