package awsiam

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Name of the role against which the login is being attempted.
If 'role' is not specified, then the login endpoint looks for a role
bearing the name of the IAM principal that is trying to login in.
If a matching role is not found, the login fails.`,
			},
			"method": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `HTTP method to use for the AWS request. This must match what
has been signed in the presigned request. Currently, POST is the only supported value`,
			},
			"url": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Full URL against which to make the AWS request. If using a POST
request with the action specified in the body, this should just be
"/".`,
			},
			"body": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Base64-encoded request body. This must match the request
body included in the signature.`,
			},
			"headers": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Base64-encoded JSON representation of the request headers.
This must at a minimum include the headers over which AWS has included a
signature.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLoginUpdate,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLoginUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	// BEGIN boring data parsing
	method := data.Get("method").(string)
	if method == "" {
		return logical.ErrorResponse("missing method"), nil
	}

	// In the future, might consider supporting GET
	if method != "POST" {
		return logical.ErrorResponse("invalid method; currently only 'POST' is supported"), nil
	}

	rawUrl := data.Get("url").(string)
	if rawUrl == "" {
		return logical.ErrorResponse("missing url"), nil
	}
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return logical.ErrorResponse("error parsing url"), nil
	}

	// TODO: There are two potentially valid cases we're not yet supporting that would
	// necessitate this check being changed. First, if we support GET requests.
	// Second if we support presigned POST requests
	bodyB64 := data.Get("body").(string)
	if bodyB64 == "" {
		return logical.ErrorResponse("missing body"), nil
	}
	bodyRaw, err := base64.StdEncoding.DecodeString(bodyB64)
	if err != nil {
		return logical.ErrorResponse("body is invalid base64"), nil
	}
	body := string(bodyRaw)

	headersB64 := data.Get("headers").(string)
	if headersB64 == "" {
		return logical.ErrorResponse("missing headers"), nil
	}
	headersJson, err := base64.StdEncoding.DecodeString(headersB64)
	if err != nil {
		return logical.ErrorResponse("headers is invalid base64"), nil
	}
	var headers http.Header
	err = jsonutil.DecodeJSON(headersJson, &headers)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("headers '%s' is invalid JSON: %s", headersJson, err)), nil
	}
	// END boring data parsing

	config, err := b.lockedClientConfigEntry(req.Storage)
	if err != nil {
		return logical.ErrorResponse("error getting configuration"), nil
	}

	endpoint := "https://sts.amazonaws.com"

	if config != nil {
		if config.HeaderValue != "" {
			ok, msg := ensureVaultHeaderValue(headers, parsedUrl, config.HeaderValue)
			if !ok {
				return logical.ErrorResponse(fmt.Sprintf("error validating %s header: %s", magicVaultHeader, msg)), nil
			}
		}
		if config.Endpoint != "" {
			endpoint = config.Endpoint
		}
	}

	clientArn, err := submitCallerIdentityRequest(method, endpoint, parsedUrl, body, headers)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error making upstream request: %s", err)), nil
	}
	canonicalArn, principalName, err := parseIamArn(clientArn)
	if err != nil {
		return logical.ErrorResponse("unrecognized IAM principal type"), nil
	}

	roleName := data.Get("role").(string)
	if roleName == "" {
		roleName = principalName
	}

	roleEntry, err := b.lockedAWSRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("entry for role %s not found", roleName)), nil
	}

	if roleEntry.BoundIamPrincipal != canonicalArn {
		return logical.ErrorResponse(fmt.Sprintf("IAM Principal '%s' does not belong to the role '%s'", clientArn, roleName)), nil
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			Policies: roleEntry.Policies,
			Metadata: map[string]string{
				"client_arn":    clientArn,
				"canonical_arn": canonicalArn,
			},
			InternalData: map[string]interface{}{
				"role_name": roleName,
			},
			DisplayName: principalName,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       roleEntry.TTL,
			},
		},
	}

	shortestTTL := b.System().DefaultLeaseTTL()
	if roleEntry.TTL > time.Duration(0) && roleEntry.TTL < shortestTTL {
		shortestTTL = roleEntry.TTL
	}

	maxTTL := b.System().MaxLeaseTTL()
	if roleEntry.MaxTTL > time.Duration(0) && roleEntry.MaxTTL < maxTTL {
		maxTTL = roleEntry.MaxTTL
	}

	if shortestTTL > maxTTL {
		resp.AddWarning(fmt.Sprintf("Effective TTL of %q exceeded the effective max_ttl of %q; TTL value is capped accordingly", (shortestTTL / time.Second).String(), (maxTTL / time.Second).String()))
		shortestTTL = maxTTL
	}

	resp.Auth.TTL = shortestTTL

	return resp, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	canonicalArn := req.Auth.Metadata["canonical_arn"]
	if canonicalArn == "" {
		return nil, fmt.Errorf("unable to retrieve canonical ARN from metadata during renewal")
	}

	roleName := req.Auth.InternalData["role_name"].(string)
	if roleName == "" {
		return nil, fmt.Errorf("error retrieving role_name during renewal")
	}
	roleEntry, err := b.lockedAWSRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, fmt.Errorf("role entry not found")
	}

	if roleEntry.BoundIamPrincipal != canonicalArn {
		return nil, fmt.Errorf("role no longer bound to arn '%s'", canonicalArn)
	}

	return framework.LeaseExtend(roleEntry.TTL, roleEntry.MaxTTL, b.System())(req, data)

}

// takes in ARNs like either arn:aws:iam::123456789012:user/MyUserName or
// arn:aws:sts::123456789012:assumed-role/RoleName/SessionName
// returns two strings
// The first is the ARN transformed into a canonical form, i.e.,
// the latter example transformed into arn:aws:iam::123456789012:role/RoleName
// The second is the "principal" name, i.e., either "MyUserName" or "RoleName"
func parseIamArn(iamArn string) (string, string, error) {
	fullParts := strings.Split(iamArn, ":")
	principalFullName := fullParts[5]
	parts := strings.Split(principalFullName, "/")
	principalName := parts[1]
	transformedArn := iamArn
	if parts[0] == "assumed-role" {
		transformedArn = fmt.Sprintf("arn:aws:iam::%s:role/%s", fullParts[4], principalName)
	} else if parts[0] != "user" {
		return "", "", fmt.Errorf("unrecognized principal type: '%s'", parts[0])
	}
	return transformedArn, principalName, nil
}

func ensureVaultHeaderValue(headers http.Header, requestUrl *url.URL, requiredHeaderValue string) (bool, string) {
	providedValue := ""
	for k, v := range headers {
		if strings.ToLower(magicVaultHeader) == strings.ToLower(k) {
			providedValue = strings.Join(v, ",")
			break
		}
	}
	if providedValue == "" {
		return false, fmt.Sprintf("didn't find %s", magicVaultHeader)
	}

	// NOT doing a constant time compare here since the value is NOT intended to be secret
	if providedValue != requiredHeaderValue {
		return false, fmt.Sprintf("expected %s but got %s", requiredHeaderValue, providedValue)
	}

	if authzHeaders, ok := headers["Authorization"]; ok {
		// authzHeader looks like AWS4-HMAC-SHA256 Credential=AKI..., SignedHeaders=host;x-amz-date;x-vault-awsiam-id, Signature=...
		// We need to extract out the SignedHeaders
		re := regexp.MustCompile(".*SignedHeaders=([^,]+)")
		authzHeader := strings.Join(authzHeaders, ",")
		matches := re.FindSubmatch([]byte(authzHeader))
		if len(matches) < 1 {
			return false, "vault header wasn't signed"
		}
		if len(matches) > 2 {
			return false, "found multiple SignedHeaders components"
		}
		signedHeaders := string(matches[1])
		return ensureHeaderIsSigned(signedHeaders, magicVaultHeader)
	}
	// TODO: If we support GET requests, then we need to parse the X-Amz-SignedHeaders
	// argument out of the query string and search in there for the header value
	return false, "Missing Authorization header"
}

func buildHttpRequest(method, endpoint string, parsedUrl *url.URL, body string, headers http.Header) *http.Request {
	// This is all a bit complicated because the AWS signature algorithm requires that
	// the Host header be included in the signed headers. See
	// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
	// The use cases we want to support, in order of increasing complexity, are:
	// 1. All defaults (client assumes sts.amazonaws.com and server has no override)
	// 2. Alternate STS regions: client wants to go to a specific region, in which case
	//    Vault must be confiugred with that endpoint as well. The client's signed request
	//    will include a signature over what the client expects the Host header to be,
	//    so we cannot change that and must match.
	// 3. Alternate STS regions with a proxy that is transparent to Vault's clients.
	//    In this case, Vault is aware of the proxy, as the proxy is configured as the
	//    endpoint, but the clients should NOT be aware of the proxy (because STS will
	//    not be aware of the proxy)
	// It's also annoying because:
	// 1. The AWS Sigv4 algorithm requires the Host header to be defined
	// 2. Some of the official SDKs (at least botocore and aws-sdk-go) don't actually
	//    incude an explicit Host header in the HTTP requests they generate, relying on
	//    the underlying HTTP library to do that for them.
	// 3. To get a validly signed request, the SDKs check if a Host header has been set
	//    and, if not, add an inferred host header (based on the URI) to the internal
	//    data structure used for calculating the signature, but never actually expose
	//    that to clients. So then they just "hope" that the underlying library actually
	//    adds the right Host header which was included in the signature calculation.
	// We could either explicity require all Vault clients to explicitly add the Host header
	// in the encoded request, or we could also implicitly infer it from the URI.
	// We choose to support both -- allow you to explicitly set a Host header, but if not,
	// infer one from the URI.
	// HOWEVER, we have to preserve the request URI portion of the client's
	// URL because the GetCallerIdentity Action can be encoded in either the body
	// or the URL. So, we need to rebuild the URL sent to the http library to have the
	// custom, Vault-specified endpoint with the client-side request parameters.
	targetUrl := fmt.Sprintf("%s/%s", endpoint, parsedUrl.RequestURI())
	request, err := http.NewRequest(method, targetUrl, strings.NewReader(body))
	if err != nil {
		return nil
	}
	request.Host = parsedUrl.Host
	for k, vals := range headers {
		for _, val := range vals {
			request.Header.Add(k, val)
		}
	}
	return request
}

func ensureHeaderIsSigned(signedHeaders, headerToSign string) (bool, string) {
	// Not doing a constant time compare here, the values aren't secret
	for _, header := range strings.Split(signedHeaders, ";") {
		if header == strings.ToLower(headerToSign) {
			return true, ""
		}
	}
	return false, fmt.Sprintf("Vault header wasn't signed")
}

func parseGetCallerIdentityResponse(response string) (GetCallerIdentityResponse, error) {
	decoder := xml.NewDecoder(strings.NewReader(response))
	result := GetCallerIdentityResponse{}
	err := decoder.Decode(&result)
	return result, err
}

func submitCallerIdentityRequest(method, endpoint string, parsedUrl *url.URL, body string, headers http.Header) (string, error) {
	// NOTE: We need to ensure we're calling STS, instead of acting as an unintended network proxy
	// The protection against this is that this method will only call the endpoint specified in the
	// client config (defaulting to sts.amazonaws.com), so it would require a Vault admin to override
	// the endpoint to talk to alternate web addresses
	request := buildHttpRequest(method, endpoint, parsedUrl, body, headers)
	client := cleanhttp.DefaultClient()
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error making request: %s", err)
	}
	if response != nil {
		defer response.Body.Close()
	}
	// we check for status code afterwards to also print out response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return "", fmt.Errorf("received error code %s from STS: %s", response.StatusCode, string(responseBody))
	}
	callerIdentityResponse, err := parseGetCallerIdentityResponse(string(responseBody))
	if err != nil {
		return "", fmt.Errorf("error parsing STS response")
	}
	clientArn := callerIdentityResponse.GetCallerIdentityResult[0].Arn
	if clientArn == "" {
		return "", fmt.Errorf("no ARN validated")
	}
	return clientArn, nil
}

type GetCallerIdentityResponse struct {
	XMLName                 xml.Name                  `xml:"GetCallerIdentityResponse"`
	GetCallerIdentityResult []GetCallerIdentityResult `xml:"GetCallerIdentityResult"`
	ResponseMetadata        []ResponseMetadata        `xml:"ResponseMetadata"`
}

type GetCallerIdentityResult struct {
	Arn     string `xml:"Arn"`
	UserId  string `xml:"UserId"`
	Account string `xml:"Account"`
}

type ResponseMetadata struct {
	RequestId string `xml:"RequestId"`
}

const pathLoginSyn = `
Authenticates an AWS IAM principal with Vault.
`

const pathLoginDesc = `
An AWS IAM principal is authenticated using the AWS STS GetCallerIdentity API method.

`

const magicVaultHeader = "X-Vault-AWSIAM-Server-Id"
