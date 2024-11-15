# Breaking Changes

## v1.6.0

### Behavioral change to `DefaultAzureCredential` in IMDS managed identity scenarios

As of `azidentity` v1.6.0, `DefaultAzureCredential` makes a minor behavioral change when it uses IMDS managed
identity. It sends its first request to IMDS without the "Metadata" header, to expedite validating whether the endpoint
is available. This precedes the credential's first token request and is guaranteed to fail with a 400 error. This error
response can appear in logs but doesn't indicate authentication failed.
