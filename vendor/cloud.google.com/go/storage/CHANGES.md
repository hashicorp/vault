# Changes

## v1.10.0
- Bump dependency on google.golang.org/api to capture changes to retry logic
  which will make retries on writes more resilient.
- Improve documentation for Writer.ChunkSize.
- Fix a bug in lifecycle to allow callers to clear lifecycle rules on a bucket.

## v1.9.0
- Add retry for transient network errors on most operations (with the exception
  of writes).
- Bump dependency for google.golang.org/api to capture a change in the default
  HTTP transport which will improve performance for reads under heavy load.
- Add CRC32C checksum validation option to Composer.

## v1.8.0
- Add support for V4 signed post policies.

## v1.7.0
- V4 signed URL support:
  - Add support for bucket-bound domains and virtual hosted style URLs.
  - Add support for query parameters in the signature.
  - Fix text encoding to align with standards.
- Add the object name to query parameters for write calls.
- Fix retry behavior when reading files with Content-Encoding gzip.
- Fix response header in reader.
- New code examples:
   - Error handling for `ObjectHandle` preconditions.
   - Existence checks for buckets and objects.

## v1.6.0

- Updated option handling:
  - Don't drop custom scopes (#1756)
  - Don't drop port in provided endpoint (#1737)

## v1.5.0

- Honor WithEndpoint client option for reads as well as writes.
- Add archive storage class to docs.
- Make fixes to storage benchwrapper.

## v1.4.0

- When listing objects in a bucket, allow callers to specify which attributes
  are queried. This allows for performance optimization.

## v1.3.0

- Use `storage.googleapis.com/storage/v1` by default for GCS requests
  instead of `www.googleapis.com/storage/v1`.

## v1.2.1

- Fixed a bug where UniformBucketLevelAccess and BucketPolicyOnly were not
  being sent in all cases.

## v1.2.0

- Add support for UniformBucketLevelAccess. This configures access checks
  to use only bucket-level IAM policies.
  See: https://godoc.org/cloud.google.com/go/storage#UniformBucketLevelAccess.
- Fix userAgent to use correct version.

## v1.1.2

- Fix memory leak in BucketIterator and ObjectIterator.

## v1.1.1

- Send BucketPolicyOnly even when it's disabled.

## v1.1.0

- Performance improvements for ObjectIterator and BucketIterator.
- Fix Bucket.ObjectIterator size calculation checks.
- Added HMACKeyOptions to all the methods which allows for options such as
  UserProject to be set per invocation and optionally be used.

## v1.0.0

This is the first tag to carve out storage as its own module. See:
https://github.com/golang/go/wiki/Modules#is-it-possible-to-add-a-module-to-a-multi-module-repository.
