package s3_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/internal/test/unit"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

type s3BucketTest struct {
	bucket string
	url    string
}

var (
	_ = unit.Imported

	sslTests = []s3BucketTest{
		{"abc", "https://abc.s3.mock-region.amazonaws.com/"},
		{"a$b$c", "https://s3.mock-region.amazonaws.com/a%24b%24c"},
		{"a.b.c", "https://s3.mock-region.amazonaws.com/a.b.c"},
		{"a..bc", "https://s3.mock-region.amazonaws.com/a..bc"},
	}

	nosslTests = []s3BucketTest{
		{"a.b.c", "http://a.b.c.s3.mock-region.amazonaws.com/"},
		{"a..bc", "http://s3.mock-region.amazonaws.com/a..bc"},
	}

	forcepathTests = []s3BucketTest{
		{"abc", "https://s3.mock-region.amazonaws.com/abc"},
		{"a$b$c", "https://s3.mock-region.amazonaws.com/a%24b%24c"},
		{"a.b.c", "https://s3.mock-region.amazonaws.com/a.b.c"},
		{"a..bc", "https://s3.mock-region.amazonaws.com/a..bc"},
	}
)

func runTests(t *testing.T, svc *s3.S3, tests []s3BucketTest) {
	for _, test := range tests {
		req, _ := svc.ListObjectsRequest(&s3.ListObjectsInput{Bucket: &test.bucket})
		req.Build()
		assert.Equal(t, test.url, req.HTTPRequest.URL.String())
	}
}

func TestHostStyleBucketBuild(t *testing.T) {
	s := s3.New(nil)
	runTests(t, s, sslTests)
}

func TestHostStyleBucketBuildNoSSL(t *testing.T) {
	s := s3.New(&aws.Config{DisableSSL: true})
	runTests(t, s, nosslTests)
}

func TestPathStyleBucketBuild(t *testing.T) {
	s := s3.New(&aws.Config{S3ForcePathStyle: true})
	runTests(t, s, forcepathTests)
}
