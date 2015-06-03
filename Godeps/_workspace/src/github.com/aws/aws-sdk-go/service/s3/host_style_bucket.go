package s3

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

var reDomain = regexp.MustCompile(`^[a-z0-9][a-z0-9\.\-]{1,61}[a-z0-9]$`)
var reIPAddress = regexp.MustCompile(`^(\d+\.){3}\d+$`)

// dnsCompatibleBucketName returns true if the bucket name is DNS compatible.
// Buckets created outside of the classic region MUST be DNS compatible.
func dnsCompatibleBucketName(bucket string) bool {
	return reDomain.MatchString(bucket) &&
		!reIPAddress.MatchString(bucket) &&
		!strings.Contains(bucket, "..")
}

// hostStyleBucketName returns true if the request should put the bucket in
// the host. This is false if S3ForcePathStyle is explicitly set or if the
// bucket is not DNS compatible.
func hostStyleBucketName(r *aws.Request, bucket string) bool {
	if r.Config.S3ForcePathStyle {
		return false
	}

	// Bucket might be DNS compatible but dots in the hostname will fail
	// certificate validation, so do not use host-style.
	if r.HTTPRequest.URL.Scheme == "https" && strings.Contains(bucket, ".") {
		return false
	}

	// Use host-style if the bucket is DNS compatible
	return dnsCompatibleBucketName(bucket)
}

func updateHostWithBucket(r *aws.Request) {
	b := awsutil.ValuesAtPath(r.Params, "Bucket")
	if len(b) == 0 {
		return
	}

	if bucket := b[0].(string); bucket != "" && hostStyleBucketName(r, bucket) {
		r.HTTPRequest.URL.Host = bucket + "." + r.HTTPRequest.URL.Host
		r.HTTPRequest.URL.Path = strings.Replace(r.HTTPRequest.URL.Path, "/{Bucket}", "", -1)
		if r.HTTPRequest.URL.Path == "" {
			r.HTTPRequest.URL.Path = "/"
		}
	}
}
