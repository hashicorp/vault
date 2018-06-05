package util

import (
	"os"
	"testing"

	"github.com/hashicorp/go-gcp-common/gcputil"
)

const googleCredentialsEnv = "TEST_GOOGLE_CREDENTIALS"
const googleProjectEnv = "TEST_GOOGLE_PROJECT"

func GetTestCredentials(t *testing.T) (string, *gcputil.GcpCredentials) {
	credentialsJSON := os.Getenv(googleCredentialsEnv)
	if credentialsJSON == "" {
		t.Fatalf("%s must be set to JSON string of valid Google credentials file", googleCredentialsEnv)
	}

	credentials, err := gcputil.Credentials(credentialsJSON)
	if err != nil {
		t.Fatalf("valid Google credentials JSON could not be read from %s env variable: %v", googleCredentialsEnv, err)
	}
	return credentialsJSON, credentials
}

func GetTestProject(t *testing.T) string {
	project := os.Getenv(googleProjectEnv)
	if project == "" {
		t.Fatalf("%s must be set to JSON string of valid Google credentials file", googleProjectEnv)
	}
	return project
}
