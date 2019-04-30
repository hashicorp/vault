package util

import (
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

	credsStr := os.Getenv("GOOGLE_CREDENTIALS")
	if credsStr == "" {
		tb.Fatal("set GOOGLE_CREDENTIALS to the path to JSON creds on disk to run integration tests")
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
