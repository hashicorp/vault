package util

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-gcp-common/gcputil"
)

func GetTestCredentials(tb testing.TB) (string, *gcputil.GcpCredentials) {
	tb.Helper()

	if testing.Short() {
		tb.Skip("skipping integration test (short)")
	}

	var credsStr string
	credsEnv := os.Getenv("GOOGLE_CREDENTIALS")
	if credsEnv == "" {
		tb.Fatal("set GOOGLE_CREDENTIALS to JSON or path to JSON creds on disk to run integration tests")
	}

	// Attempt to read as file path; if invalid, assume given JSON value directly
	if _, err := os.Stat(credsEnv); err == nil {
		credsBytes, err := ioutil.ReadFile(credsEnv)
		if err != nil {
			tb.Fatalf("unable to read credentials file %s: %v", credsStr, err)
		}
		credsStr = string(credsBytes)
	} else {
		credsStr = credsEnv
	}

	creds, err := gcputil.Credentials(credsStr)
	if err != nil {
		tb.Fatalf("failed to parse GOOGLE_CREDENTIALS as JSON: %s", err)
	}
	return credsStr, creds
}

func GetTestProject(tb testing.TB) string {
	tb.Helper()

	if testing.Short() {
		tb.Skip("skipping integration test (short)")
	}

	project := strings.TrimSpace(os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if project == "" {
		tb.Fatal("set GOOGLE_CLOUD_PROJECT to the ID of a GCP project to run integration tests")
	}
	return project
}
