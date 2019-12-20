package awsutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	// These variables will be used to toggle the test server's
	// behavior for each test.
	returnErr bool
	v2Enabled bool
)

func TestMain(m *testing.M) {
	// Start a test server for us to use.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if returnErr {
			// Bad request.
			w.WriteHeader(400)
			return
		}
		if v2Enabled {
			w.WriteHeader(200)
			w.Write([]byte(instanceMetadataServiceTokenResponse))
			return
		}
		// The instance metadata service is turned off.
		w.WriteHeader(403)
	}))
	defer ts.Close()

	// Override the instance metadata service's URL with the test server's
	// for the duration of our tests.
	original := InstanceMetadataService.BaseURL
	InstanceMetadataService.BaseURL = ts.URL
	defer func() {
		InstanceMetadataService.BaseURL = original
	}()
	os.Exit(m.Run())
}

// Note, it's possible to disable instance metadata altogether
// but we'll let the caller discover that on their own.
func TestPrepareInstanceMetadataReq_V2Enabled(t *testing.T) {
	v2Enabled = true
	returnErr = false

	req, err := http.NewRequest(http.MethodGet, "http://localhost:1000", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstanceMetadataService.PrepareRequest(req); err != nil {
		t.Fatal(err)
	}
	result := req.Header.Get(InstanceMetadataService.TokenHeader)
	if result != instanceMetadataServiceTokenResponse {
		t.Fatalf("expected %q but received %q", instanceMetadataServiceTokenResponse, result)
	}
}

func TestPrepareInstanceMetadataReq_V2Disabled(t *testing.T) {
	v2Enabled = false
	returnErr = false

	req, err := http.NewRequest(http.MethodGet, "http://localhost:1000", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstanceMetadataService.PrepareRequest(req); err != nil {
		t.Fatal(err)
	}
	result := req.Header.Get(InstanceMetadataService.TokenHeader)
	if result != "" {
		t.Fatalf("expected %q but received %q", "", result)
	}
}

func TestPrepareInstanceMetadataReq_MiscErr(t *testing.T) {
	v2Enabled = true
	returnErr = true

	req, err := http.NewRequest(http.MethodGet, "http://localhost:1000", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstanceMetadataService.PrepareRequest(req); err == nil {
		t.Fatal("expected err")
	}
}

const instanceMetadataServiceTokenResponse = "AQAAAHjB6Liyq02K_l4K9DqomS5iztS2Jc0Mpma12-s16qLUk75Amw=="
