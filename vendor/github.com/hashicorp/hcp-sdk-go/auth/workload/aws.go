// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workload

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	// Environment Variable Reference:
	// https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html

	// awsEnvRegion is the region to send requests to. It takes precedence of
	// default region.
	awsEnvRegion = "AWS_REGION"

	// awsEnvDefaultRegion is where requests will be sent to by default, if not
	// overridden.
	awsEnvDefaultRegion = "AWS_DEFAULT_REGION"

	// awsEnvAccessKeyID stores the AWS access key.
	awsEnvAccessKeyID = "AWS_ACCESS_KEY_ID"

	// awsEnvSecretAccessKeyId stores the secret key associated with the access key.
	awsEnvSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	// awsEnvSessionToken stores session token value that is required if you are
	// using temporary security credentials that you retrieved directly from AWS
	// STS operations.
	awsEnvSessionToken = "AWS_SESSION_TOKEN"

	// awsRegionURL is the metadata URL to discover the instances current
	// region.
	awsRegionURL = "http://169.254.169.254/latest/meta-data/placement/region"

	// awsSecurityCredentialsURL is used to find the assigned role and retrieve
	// its credentials.
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html#instance-metadata-security-credentials
	awsSecurityCredentialsURL = "http://169.254.169.254/latest/meta-data/iam/security-credentials"

	// awsSessionTokenURL retrieves a short lived session token.
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-metadata-v2-how-it-works.html
	awsSessionTokenURL = "http://169.254.169.254/latest/api/token"

	// awsIMDSv2SessionTTLHeader is used to configure the session token TTL.
	awsIMDSv2SessionTTLHeader = "X-aws-ec2-metadata-token-ttl-seconds"

	// awsIMDSv2SessionTTL is the session ttl we request.
	awsIMDSv2SessionTTL = "300"

	// awsIMDSv2SessionTokenHeader is used to pass the short lived session
	// token to an IMDSv2 endpoint.
	awsIMDSv2SessionTokenHeader = "X-aws-ec2-metadata-token"

	// awsRegionalSTSEndpoint is the regional STS endpoint. The region must be
	// formatted in.
	awsRegionalSTSEndpoint = "https://sts.%s.amazonaws.com?Action=GetCallerIdentity&Version=2011-06-15"

	// awsSecurityTokenHeader is the header to pass the short lived security
	// token.
	awsSecurityTokenHeader = "x-amz-security-token"

	// awsDateHeader is the header to pass the timestamp of the request.
	awsDateHeader = "x-amz-date"

	// awsDateFormat is the how to format the date timestamp
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/signing-elements.html#date
	awsDateFormat = "20060102T150405Z"

	// awsDateShortFormat is the how to format the date in the credential
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/signing-elements.html#authentication
	awsDateShortFormat = "20060102"

	// awsAlgorithm indicates how the request was signed
	awsAlgorithm = "AWS4-HMAC-SHA256"

	// awsRequestType is used to terminate the signed credential.
	// https://docs.aws.amazon.com/general/latest/gr/sigv4-create-string-to-sign.html
	awsRequestType = "aws4_request"

	// hcpWorkloadIdentityProviderHeader is the header used to reference the
	// workload identity provider the token exchange is meant for. The header
	// must be included in the signed request.
	hcpWorkloadIdentityProviderHeader = "x-hcp-workload-identity-provider"
)

// AWSCredentialSource sources credentials by interacting with the AWS IMDS
// endpoint to sign an AWS GetCallerIdentity request. The signed request can
// then be used by HCP to return HCP Service Principal credentials based on the
// identity of the AWS workload.
type AWSCredentialSource struct {
	// IMDSv2 indicates that IMDSv2 endpoint should be used.
	IMDSv2 bool `json:"imds_v2,omitempty"`

	// client is the http client used to make requests
	client *http.Client

	// now is a function that returns the current time
	now func() time.Time
}

// getCallerIdentityReq returns the signed AWS GetCallerIdentity request.
func (ac *AWSCredentialSource) getCallerIdentityReq(wipResourceName string) (*callerIdentityRequest, error) {
	if ac.client == nil {
		ac.client = cleanhttp.DefaultPooledClient()
	}

	if ac.now == nil {
		ac.now = time.Now
	}

	// get the request signer
	s, err := newAWSRequestSigner(ac.IMDSv2, ac.client, ac.now)
	if err != nil {
		return nil, err
	}

	// Create the request to the regional AWS STS GetCallerIdentity API.
	req, err := http.NewRequest("POST", fmt.Sprintf(awsRegionalSTSEndpoint, s.region), nil)
	if err != nil {
		return nil, err
	}

	// Add the workload identity provider resource name as a signed header. This
	// ensures that the token exchange is for only the specified resource.
	req.Header.Add(hcpWorkloadIdentityProviderHeader, wipResourceName)

	// Sign the request
	if err := s.SignRequest(req); err != nil {
		return nil, fmt.Errorf("failed to sign GetCallerIdentity request: %v", err)
	}

	// Convert the request
	idReq := &callerIdentityRequest{
		Headers: map[string]string{},
		Method:  req.Method,
		URL:     req.URL.String(),
	}
	for k, values := range req.Header {
		value := ""
		if len(values) > 0 {
			value = values[0]
		}
		idReq.Headers[k] = value
	}

	return idReq, nil
}

// callerIdentityRequest is the signed request for the GetCallerIdentity
// endpoint.
type callerIdentityRequest struct {
	// headers is the HTTP request headers. This must include:
	//
	// x-amz-date: the date of the request
	//  host: the host of the request, e.g. sts.amazonaws.com
	//  x-hcp-workload-identity-provider: the resource_name of the workload
	//    identity provider the token exchange will be conducted against.
	//  Authorization: The AWS Signature for the request.
	//  X-Amz-Security-Token: The temporary security credentials' session used
	//    to sign the request. Described here:
	//    https://docs.aws.amazon.com/IAM/latest/UserGuide/create-signed-request.html#temporary-security-credentials
	Headers map[string]string `json:"headers,omitempty"`

	// method is the method of the HTTP request.
	Method string `json:"method,omitempty"`

	// url is the URL of the AWS endpoint being called.
	URL string `json:"url,omitempty"`
}

type awsRequestSigner struct {
	// client is the http client used to make requests
	client *http.Client

	// imdsSessionToken is the session token to send to IMDSv2 requests
	imdsSessionToken string

	// region is the AWS region being operated in
	region string

	// AWS credentials
	creds awsSecurityCredentials

	// now is a function that returns the current time
	now func() time.Time
}

type awsSecurityCredentials struct {
	AccessKeyID     string `json:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SecurityToken   string `json:"Token"`
}

func newAWSRequestSigner(imdsV2 bool, client *http.Client, now func() time.Time) (*awsRequestSigner, error) {
	s := &awsRequestSigner{
		client: client,
		now:    now,
	}
	s.sourceEnvVars()

	// Create a request context with a deadline
	reqCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If we are configured to use IMDSv2 and we haven't sourced everything we
	// need from environment variables then get a session token.
	if imdsV2 && (s.region == "" || s.creds.AccessKeyID == "") {
		if err := s.getSessionToken(reqCtx); err != nil {
			return nil, err
		}
	}

	if err := s.getRegion(reqCtx); err != nil {
		return nil, err
	}

	if err := s.getCredentials(reqCtx); err != nil {
		return nil, err
	}

	return s, nil
}

// SignRequest adds the appropriate headers to an http.Request
// or returns an error if something prevented this.
func (s *awsRequestSigner) SignRequest(req *http.Request) error {
	req.Header.Add("host", req.Host)

	if s.creds.SecurityToken != "" {
		req.Header.Add(awsSecurityTokenHeader, s.creds.SecurityToken)
	}

	now := s.now()
	if req.Header.Get("date") == "" {
		req.Header.Add(awsDateHeader, now.Format(awsDateFormat))
	}

	authorizationCode, err := s.generateAuthentication(req, now)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", authorizationCode)
	return nil
}

func (s *awsRequestSigner) sourceEnvVars() {
	// Try to determine the region
	region, regionOk := os.LookupEnv(awsEnvRegion)
	if regionOk {
		s.region = region
	}

	defaultRegion, defaultOk := os.LookupEnv(awsEnvDefaultRegion)
	if !regionOk && defaultOk {
		s.region = defaultRegion
	}

	// Try to get the AWS credentials
	accessKey, accessKeyOk := os.LookupEnv(awsEnvAccessKeyID)
	secretKey, secretKeyOk := os.LookupEnv(awsEnvSecretAccessKey)
	sessionToken := os.Getenv(awsEnvSessionToken)
	if accessKeyOk && secretKeyOk {
		s.creds.AccessKeyID = accessKey
		s.creds.SecretAccessKey = secretKey
		s.creds.SecurityToken = sessionToken
	}
}

func (s *awsRequestSigner) getSessionToken(ctx context.Context) error {
	// Create the request to retrieve the session token.
	req, err := http.NewRequestWithContext(ctx, "PUT", awsSessionTokenURL, nil)
	if err != nil {
		return err
	}

	// Configure the requested token TTL
	req.Header.Add(awsIMDSv2SessionTTLHeader, awsIMDSv2SessionTTL)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed retrieving AWS session token from metadata endpoint: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("failed reading AWS session token response from metadata endpoint: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("AWS session token endpoint returned status code %d: %s", resp.StatusCode, string(respBody))
	}

	// Store the token in the header.
	s.imdsSessionToken = string(respBody)

	return nil
}

func (s *awsRequestSigner) getRegion(ctx context.Context) error {
	// Check if we retrieved this from the environment already
	if s.region != "" {
		return nil
	}

	// Create the request to retrieve the region
	req, err := http.NewRequestWithContext(ctx, "GET", awsRegionURL, nil)
	if err != nil {
		return err
	}

	if s.imdsSessionToken != "" {
		req.Header.Add(awsIMDSv2SessionTokenHeader, s.imdsSessionToken)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed retrieving AWS region from metadata endpoint: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("failed reading AWS region response from metadata endpoint: %v", err)

	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("metadata region endpoint returned status %d: %v", resp.StatusCode, string(respBody))
	}

	// The only value returned is the region
	s.region = string(respBody)
	return nil
}

func (s *awsRequestSigner) getCredentials(ctx context.Context) error {
	// Check if we retrieved these from the environment already
	if s.creds.AccessKeyID != "" {
		return nil
	}

	role, err := s.getRoleName(ctx)
	if err != nil {
		return err
	}

	// Create the request to retrieve the role
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", awsSecurityCredentialsURL, role), nil)
	if err != nil {
		return err
	}

	if s.imdsSessionToken != "" {
		req.Header.Add(awsIMDSv2SessionTokenHeader, s.imdsSessionToken)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed retrieving security credentials from metadata endpoint: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("failed reading AWS security credential response from metadata endpoint: %v", err)

	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("metadata security credential endpoint returned status %d: %v", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, &s.creds); err != nil {
		return fmt.Errorf("failed to unmarshall security credential response: %v", err)
	}

	return nil
}

func (s *awsRequestSigner) getRoleName(ctx context.Context) (string, error) {
	// Create the request to retrieve the role
	req, err := http.NewRequestWithContext(ctx, "GET", awsSecurityCredentialsURL, nil)
	if err != nil {
		return "", err
	}

	if s.imdsSessionToken != "" {
		req.Header.Add(awsIMDSv2SessionTokenHeader, s.imdsSessionToken)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed retrieving role name from metadata endpoint: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("failed reading AWS security credential response from metadata endpoint: %v", err)

	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("metadata security credential endpoint returned status %d: %v", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

// generateAuthentication generates the authentication header for the given
// request.
// https://docs.aws.amazon.com/IAM/latest/UserGuide/create-signed-request.html
func (s *awsRequestSigner) generateAuthentication(req *http.Request, timestamp time.Time) (string, error) {
	canonicalHeaderColumns, canonicalHeaderData := canonicalHeaders(req)

	dateStamp := timestamp.Format(awsDateShortFormat)
	credentialScope := fmt.Sprintf("%s/%s/%s/%s", dateStamp, s.region, "sts", awsRequestType)

	requestString, err := canonicalRequest(req, canonicalHeaderColumns, canonicalHeaderData)
	if err != nil {
		return "", err
	}
	requestHash, err := getSha256([]byte(requestString))
	if err != nil {
		return "", err
	}

	signingKey := []byte("AWS4" + s.creds.SecretAccessKey)
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s", awsAlgorithm, timestamp.Format(awsDateFormat), credentialScope, requestHash)
	for _, signingInput := range []string{
		dateStamp, s.region, "sts", awsRequestType, stringToSign,
	} {
		signingKey, err = getHmacSha256(signingKey, []byte(signingInput))
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", awsAlgorithm, s.creds.AccessKeyID, credentialScope, canonicalHeaderColumns, hex.EncodeToString(signingKey)), nil
}

func getSha256(input []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(input); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func getHmacSha256(key, input []byte) ([]byte, error) {
	hash := hmac.New(sha256.New, key)
	if _, err := hash.Write(input); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func canonicalHeaders(req *http.Request) (string, string) {
	// Header keys need to be sorted alphabetically.
	var headers []string
	lowerCaseHeaders := make(http.Header)
	for k, v := range req.Header {
		k := strings.ToLower(k)
		if _, ok := lowerCaseHeaders[k]; ok {
			// include additional values
			lowerCaseHeaders[k] = append(lowerCaseHeaders[k], v...)
		} else {
			headers = append(headers, k)
			lowerCaseHeaders[k] = v
		}
	}
	sort.Strings(headers)

	var fullHeaders bytes.Buffer
	for _, header := range headers {
		headerValue := strings.Join(lowerCaseHeaders[header], ",")
		fullHeaders.WriteString(header)
		fullHeaders.WriteRune(':')
		fullHeaders.WriteString(headerValue)
		fullHeaders.WriteRune('\n')
	}

	return strings.Join(headers, ";"), fullHeaders.String()
}

func requestDataHash(req *http.Request) (string, error) {
	var requestData []byte
	if req.Body != nil {
		requestBody, err := req.GetBody()
		if err != nil {
			return "", err
		}
		defer requestBody.Close()

		requestData, err = io.ReadAll(io.LimitReader(requestBody, 1<<20))
		if err != nil {
			return "", err
		}
	}

	return getSha256(requestData)
}

func canonicalRequest(req *http.Request, canonicalHeaderColumns, canonicalHeaderData string) (string, error) {
	dataHash, err := requestDataHash(req)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", req.Method, canonicalPath(req), canonicalQuery(req), canonicalHeaderData, canonicalHeaderColumns, dataHash), nil
}

func canonicalPath(req *http.Request) string {
	result := req.URL.EscapedPath()
	if result == "" {
		return "/"
	}
	return path.Clean(result)
}

func canonicalQuery(req *http.Request) string {
	queryValues := req.URL.Query()
	for queryKey := range queryValues {
		sort.Strings(queryValues[queryKey])
	}
	return queryValues.Encode()
}
