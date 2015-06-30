package aws

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

var testCredentials = credentials.NewChainCredentials([]credentials.Provider{
	&credentials.EnvProvider{},
	&credentials.SharedCredentialsProvider{
		Filename: "TestFilename",
		Profile:  "TestProfile"},
	&credentials.EC2RoleProvider{ExpiryWindow: 5 * time.Minute},
})

var copyTestConfig = Config{
	Credentials:             testCredentials,
	Endpoint:                "CopyTestEndpoint",
	Region:                  "COPY_TEST_AWS_REGION",
	DisableSSL:              true,
	ManualSend:              true,
	HTTPClient:              http.DefaultClient,
	LogHTTPBody:             true,
	LogLevel:                2,
	Logger:                  os.Stdout,
	MaxRetries:              DefaultRetries,
	DisableParamValidation:  true,
	DisableComputeChecksums: true,
	S3ForcePathStyle:        true,
}

func TestCopy(t *testing.T) {
	want := copyTestConfig
	got := copyTestConfig.Copy()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Copy() = %+v", got)
		t.Errorf("    want %+v", want)
	}
}

func TestCopyReturnsNewInstance(t *testing.T) {
	want := copyTestConfig
	got := copyTestConfig.Copy()
	if &got == &want {
		t.Errorf("Copy() = %p; want different instance as source %p", &got, &want)
	}
}

var mergeTestZeroValueConfig = Config{MaxRetries: DefaultRetries}

var mergeTestConfig = Config{
	Credentials:             testCredentials,
	Endpoint:                "MergeTestEndpoint",
	Region:                  "MERGE_TEST_AWS_REGION",
	DisableSSL:              true,
	ManualSend:              true,
	HTTPClient:              http.DefaultClient,
	LogHTTPBody:             true,
	LogLevel:                2,
	Logger:                  os.Stdout,
	MaxRetries:              10,
	DisableParamValidation:  true,
	DisableComputeChecksums: true,
	S3ForcePathStyle:        true,
}

var mergeTests = []struct {
	cfg  *Config
	in   *Config
	want *Config
}{
	{&Config{}, nil, &Config{}},
	{&Config{}, &mergeTestZeroValueConfig, &Config{}},
	{&Config{}, &mergeTestConfig, &mergeTestConfig},
}

func TestMerge(t *testing.T) {
	for _, tt := range mergeTests {
		got := tt.cfg.Merge(tt.in)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Config %+v", tt.cfg)
			t.Errorf(" Merge(%+v)", tt.in)
			t.Errorf("   got %+v", got)
			t.Errorf("  want %+v", tt.want)
		}
	}
}
