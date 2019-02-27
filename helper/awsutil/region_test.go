package awsutil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/user"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
)

const testConfigFile = `[default]
region=%s
output=json`

var (
	logger               = hclog.NewNullLogger()
	expectedTestRegion   = "us-west-2"
	unexpectedTestRegion = "us-east-2"
)

func TestGetOrDefaultRegion_UserConfigPreferredFirst(t *testing.T) {
	configuredRegion := expectedTestRegion

	cleanupEnv := setEnvRegion(t, unexpectedTestRegion)
	defer cleanupEnv()

	cleanupFile := setConfigFileRegion(t, unexpectedTestRegion)
	defer cleanupFile()

	cleanupMetadata := setInstanceMetadata(t, unexpectedTestRegion)
	defer cleanupMetadata()

	result := GetOrDefaultRegion(logger, configuredRegion)
	if result != expectedTestRegion {
		t.Fatalf("expected: %s; actual: %s", expectedTestRegion, result)
	}
}

func TestGetOrDefaultRegion_EnvVarsPreferredSecond(t *testing.T) {
	configuredRegion := ""

	cleanupEnv := setEnvRegion(t, expectedTestRegion)
	defer cleanupEnv()

	cleanupFile := setConfigFileRegion(t, unexpectedTestRegion)
	defer cleanupFile()

	cleanupMetadata := setInstanceMetadata(t, unexpectedTestRegion)
	defer cleanupMetadata()

	result := GetOrDefaultRegion(logger, configuredRegion)
	if result != expectedTestRegion {
		t.Fatalf("expected: %s; actual: %s", expectedTestRegion, result)
	}
}

func TestGetOrDefaultRegion_ConfigFilesPreferredThird(t *testing.T) {
	configuredRegion := ""

	cleanupEnv := setEnvRegion(t, "")
	defer cleanupEnv()

	cleanupFile := setConfigFileRegion(t, expectedTestRegion)
	defer cleanupFile()

	cleanupMetadata := setInstanceMetadata(t, unexpectedTestRegion)
	defer cleanupMetadata()

	result := GetOrDefaultRegion(logger, configuredRegion)
	if result != expectedTestRegion {
		t.Fatalf("expected: %s; actual: %s", expectedTestRegion, result)
	}
}

func TestGetOrDefaultRegion_ConfigFileUnfound(t *testing.T) {
	configuredRegion := ""

	cleanupEnv := setEnvRegion(t, "")
	defer cleanupEnv()

	if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "foo"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Unsetenv("AWS_SHARED_CREDENTIALS_FILE"); err != nil {
			t.Fatal(err)
		}
	}()

	result := GetOrDefaultRegion(logger, configuredRegion)
	if result != DefaultRegion {
		t.Fatalf("expected: %s; actual: %s", DefaultRegion, result)
	}
}

func TestGetOrDefaultRegion_EC2InstanceMetadataPreferredFourth(t *testing.T) {
	configuredRegion := ""

	cleanupEnv := setEnvRegion(t, "")
	defer cleanupEnv()

	cleanupFile := setConfigFileRegion(t, "")
	defer cleanupFile()

	cleanupMetadata := setInstanceMetadata(t, expectedTestRegion)
	defer cleanupMetadata()

	result := GetOrDefaultRegion(logger, configuredRegion)
	if result != expectedTestRegion {
		t.Fatalf("expected: %s; actual: %s", expectedTestRegion, result)
	}
}

func TestGetOrDefaultRegion_DefaultsToDefaultRegionWhenRegionUnavailable(t *testing.T) {
	configuredRegion := ""

	cleanupEnv := setEnvRegion(t, "")
	defer cleanupEnv()

	cleanupFile := setConfigFileRegion(t, "")
	defer cleanupFile()

	result := GetOrDefaultRegion(logger, configuredRegion)
	if result != DefaultRegion {
		t.Fatalf("expected: %s; actual: %s", DefaultRegion, result)
	}
}

func setEnvRegion(t *testing.T, region string) (cleanup func()) {
	for _, envKey := range RegionEnvKeys {
		if err := os.Setenv(envKey, region); err != nil {
			t.Fatal(err)
		}
	}
	cleanup = func() {
		for _, envKey := range RegionEnvKeys {
			if err := os.Unsetenv(envKey); err != nil {
				t.Fatal(err)
			}
		}
	}
	return
}

func setConfigFileRegion(t *testing.T, region string) (cleanup func()) {

	var cleanupFuncs []func()

	cleanup = func() {
		for _, f := range cleanupFuncs {
			f()
		}
	}

	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	pathToAWSDir := usr.HomeDir + "/.aws"
	pathToConfig := pathToAWSDir + "/config"

	preExistingConfig, err := ioutil.ReadFile(pathToConfig)
	if err != nil {
		// File simply doesn't exist.
		if err := os.Mkdir(pathToAWSDir, os.ModeDir); err != nil {
			t.Fatal(err)
		}
		cleanupFuncs = append(cleanupFuncs, func() {
			if err := os.RemoveAll(pathToAWSDir); err != nil {
				t.Fatal(err)
			}
		})
	} else {
		cleanupFuncs = append(cleanupFuncs, func() {
			if err := ioutil.WriteFile(pathToConfig, preExistingConfig, 0644); err != nil {
				t.Fatal(err)
			}
		})
	}
	fileBody := fmt.Sprintf(testConfigFile, region)
	if err := ioutil.WriteFile(pathToConfig, []byte(fileBody), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", pathToConfig); err != nil {
		t.Fatal(err)
	}
	cleanupFuncs = append(cleanupFuncs, func() {
		if err := os.Unsetenv("AWS_SHARED_CREDENTIALS_FILE"); err != nil {
			t.Fatal(err)
		}
	})

	return
}

func setInstanceMetadata(t *testing.T, region string) (cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.String()
		switch reqPath {
		case "/latest/meta-data/instance-id":
			w.Write([]byte("i-1234567890abcdef0"))
			return
		case "/latest/meta-data/placement/availability-zone":
			// add a letter suffix, as a normal response is formatted like "us-east-1a"
			w.Write([]byte(region + "a"))
			return
		default:
			t.Fatalf("received unexpected request path: %s", reqPath)
		}
	}))
	originalMetadataBaseURL := ec2MetadataBaseURL
	ec2MetadataBaseURL = ts.URL
	cleanup = func() {
		ts.Close()
		ec2MetadataBaseURL = originalMetadataBaseURL
	}
	return
}
