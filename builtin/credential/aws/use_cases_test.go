package awsauth

/*
	The AWS auth method backend acceptance tests are for testing high-level
	use cases the AWS auth engine has.
*/

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

// This is directly based on our docs here:
// https://www.vaultproject.io/docs/auth/aws
func TestEC2LoginRenewDefaultSettings(t *testing.T) {
	mockEC2Client = &fakeEC2Client{
		describeInstanceOutputToReturn: &ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{
						{
							InstanceId: aws.String("i-de0f1344"),
							State: &ec2.InstanceState{
								Name: aws.String("running"),
							},
							ImageId: aws.String("ami-fce3c696"),
						},
					},
				},
			},
		},
	}
	mockSTSClient = &fakeSTSClient{
		getCallerIdentityOutputToReturn: &sts.GetCallerIdentityOutput{
			Account: aws.String("241656615859"),
			Arn:     aws.String("arn:aws:iam::241656615859:tester/tester"),
			UserId:  aws.String("5678"),
		},
	}

	testEnv, err := newTestEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	{
		// This is the fake key and secret in our docs.
		// vault write auth/aws/config/client secret_key=vCtSM8ZUEQ3mOFVlYPBQkf2sO6F/W7a5TVzrl3Oj access_key=VKIAJBRHKH6EVTTNXDHA
		req := &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "config/client",
			Storage:   testEnv.conf.StorageView,
			Data: map[string]interface{}{
				"secret_key": "vCtSM8ZUEQ3mOFVlYPBQkf2sO6F/W7a5TVzrl3Oj",
				"access_key": "VKIAJBRHKH6EVTTNXDHA",
			},
		}
		resp, err := testEnv.backend.HandleRequest(testEnv.ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Fatalf("expected nil response but received %+v", resp)
		}
	}
	{
		// vault write auth/aws/role/dev-role auth_type=ec2 bound_ami_id=ami-fce3c696 policies=prod,dev max_ttl=500h
		req := &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "role/dev-role",
			Storage:   testEnv.conf.StorageView,
			Data: map[string]interface{}{
				"auth_type":    "ec2",
				"bound_ami_id": "ami-fce3c696",
				"policies":     []string{"prod", "dev"},
				"max_ttl":      "500h",
			},
		}
		resp, err := testEnv.backend.HandleRequest(testEnv.ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Fatalf("expected nil response but received %+v", resp)
		}
	}
	renewalReq := &logical.Request{}
	{
		// vault write auth/aws/login role=dev-role \
		//		pkcs7=MIAGCSqGSIb3DQEHAqCAMIACAQExCzAJBgUrDgMCGgUAMIAGCSqGSIb3DQEHAaCAJIAEggGmewogICJkZXZwYXlQcm9kdWN0Q29kZXMiIDogbnVsbCwKICAicHJpdmF0ZUlwIiA6ICIxNzIuMzEuNjMuNjAiLAogICJhdmFpbGFiaWxpdHlab25lIiA6ICJ1cy1lYXN0LTFjIiwKICAidmVyc2lvbiIgOiAiMjAxMC0wOC0zMSIsCiAgImluc3RhbmNlSWQiIDogImktZGUwZjEzNDQiLAogICJiaWxsaW5nUHJvZHVjdHMiIDogbnVsbCwKICAiaW5zdGFuY2VUeXBlIiA6ICJ0Mi5taWNybyIsCiAgImFjY291bnRJZCIgOiAiMjQxNjU2NjE1ODU5IiwKICAiaW1hZ2VJZCIgOiAiYW1pLWZjZTNjNjk2IiwKICAicGVuZGluZ1RpbWUiIDogIjIwMTYtMDQtMDVUMTY6MjY6NTVaIiwKICAiYXJjaGl0ZWN0dXJlIiA6ICJ4ODZfNjQiLAogICJrZXJuZWxJZCIgOiBudWxsLAogICJyYW1kaXNrSWQiIDogbnVsbCwKICAicmVnaW9uIiA6ICJ1cy1lYXN0LTEiCn0AAAAAAAAxggEXMIIBEwIBATBpMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQwIJAJa6SNnlXhpnMAkGBSsOAwIaBQCgXTAYBgkqhkiG9w0BCQMxCwYJKoZIhvcNAQcBMBwGCSqGSIb3DQEJBTEPFw0xNjA0MDUxNjI3MDBaMCMGCSqGSIb3DQEJBDEWBBRtiynzMTNfTw1TV/d8NvfgVw+XfTAJBgcqhkjOOAQDBC4wLAIUVfpVcNYoOKzN1c+h1Vsm/c5U0tQCFAK/K72idWrONIqMOVJ8Uen0wYg4AAAAAAAA \
		//		nonce=5defbf9e-a8f9-3063-bdfc-54b7a42a1f95
		req := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "login",
			Storage:   testEnv.conf.StorageView,
			Data: map[string]interface{}{
				"role":  "dev-role",
				"pkcs7": "MIAGCSqGSIb3DQEHAqCAMIACAQExCzAJBgUrDgMCGgUAMIAGCSqGSIb3DQEHAaCAJIAEggGmewogICJkZXZwYXlQcm9kdWN0Q29kZXMiIDogbnVsbCwKICAicHJpdmF0ZUlwIiA6ICIxNzIuMzEuNjMuNjAiLAogICJhdmFpbGFiaWxpdHlab25lIiA6ICJ1cy1lYXN0LTFjIiwKICAidmVyc2lvbiIgOiAiMjAxMC0wOC0zMSIsCiAgImluc3RhbmNlSWQiIDogImktZGUwZjEzNDQiLAogICJiaWxsaW5nUHJvZHVjdHMiIDogbnVsbCwKICAiaW5zdGFuY2VUeXBlIiA6ICJ0Mi5taWNybyIsCiAgImFjY291bnRJZCIgOiAiMjQxNjU2NjE1ODU5IiwKICAiaW1hZ2VJZCIgOiAiYW1pLWZjZTNjNjk2IiwKICAicGVuZGluZ1RpbWUiIDogIjIwMTYtMDQtMDVUMTY6MjY6NTVaIiwKICAiYXJjaGl0ZWN0dXJlIiA6ICJ4ODZfNjQiLAogICJrZXJuZWxJZCIgOiBudWxsLAogICJyYW1kaXNrSWQiIDogbnVsbCwKICAicmVnaW9uIiA6ICJ1cy1lYXN0LTEiCn0AAAAAAAAxggEXMIIBEwIBATBpMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQwIJAJa6SNnlXhpnMAkGBSsOAwIaBQCgXTAYBgkqhkiG9w0BCQMxCwYJKoZIhvcNAQcBMBwGCSqGSIb3DQEJBTEPFw0xNjA0MDUxNjI3MDBaMCMGCSqGSIb3DQEJBDEWBBRtiynzMTNfTw1TV/d8NvfgVw+XfTAJBgcqhkjOOAQDBC4wLAIUVfpVcNYoOKzN1c+h1Vsm/c5U0tQCFAK/K72idWrONIqMOVJ8Uen0wYg4AAAAAAAA",
				"nonce": "5defbf9e-a8f9-3063-bdfc-54b7a42a1f95",
			},
		}
		resp, err := testEnv.backend.HandleRequest(testEnv.ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Auth == nil {
			t.Fatal("expected to receive auth")
		}
		renewalReq.Auth = resp.Auth
	}
	{
		// Test renewal.
		renewalReq.Storage = testEnv.conf.StorageView
		resp, err := testEnv.backend.pathLoginRenew(testEnv.ctx, renewalReq, nil)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response but received")
		}
		if resp.Auth == nil {
			t.Fatal("expected to receive auth")
		}
	}
}

// This is directly based on our docs here:
// https://www.vaultproject.io/docs/auth/aws
func TestIAMLoginRenewDefaultSettings(t *testing.T) {
	mockIAMClient = &fakeIAMClient{
		roleOutputToReturn: &iam.GetRoleOutput{
			Role: &iam.Role{
				RoleId: aws.String("AROADBQP57FF2AEXAMPLE"),
				Arn:    aws.String("arn:aws:iam::241656615859:role/MyRole"),
			},
		},
	}
	mockSTSClient = &fakeSTSClient{
		getCallerIdentityOutputToReturn: &sts.GetCallerIdentityOutput{
			Account: aws.String("241656615859"),
			Arn:     aws.String("arn:aws:iam::241656615859:tester/tester"),
			UserId:  aws.String("5678"),
		},
		getCallerIdentityRequestToReturn: &request.Request{
			HTTPRequest: &http.Request{
				Method: "POST",
				Header: map[string][]string{
					"Authorization": {"SignedHeaders=" + strings.ToLower("X-Vault-AWS-IAM-Server-ID")},
				},
				Body: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
				URL: &url.URL{
					Scheme: "https://",
					Host:   "www.foo.com",
				},
			},
		},
	}
	directGetCallerIdentityResponse := `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
   <Arn>arn:aws:iam::241656615859:role/MyRole</Arn>
    <UserId>AROADBQP57FF2AEXAMPLE</UserId>
    <Account>123456789012</Account>
  </GetCallerIdentityResult>
  <ResponseMetadata>
    <RequestId>01234567-89ab-cdef-0123-456789abcdef</RequestId>
  </ResponseMetadata>
</GetCallerIdentityResponse>`

	testEnv, err := newTestEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	{
		// vault write auth/aws/role/dev-role-iam auth_type=iam \
		//		bound_iam_principal_arn=arn:aws:iam::241656615859:role/MyRole policies=prod,dev max_ttl=500h
		req := &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "role/dev-role-iam",
			Storage:   testEnv.conf.StorageView,
			Data: map[string]interface{}{
				"auth_type":               "iam",
				"bound_iam_principal_arn": "arn:aws:iam::241656615859:role/MyRole",
				"policies":                []string{"prod", "dev"},
				"max_ttl":                 "500h",
			},
		}
		resp, err := testEnv.backend.HandleRequest(testEnv.ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Fatalf("expected nil response but received %+v", resp)
		}
	}
	// We need a test server to respond to the GetCallerIdentity call.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(directGetCallerIdentityResponse))
	}))
	defer ts.Close()
	{
		// vault write auth/aws/config/client iam_server_id_header_value=vault.example.com
		req := &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "config/client",
			Storage:   testEnv.conf.StorageView,
			Data: map[string]interface{}{
				"iam_server_id_header_value": "vault.example.com",
				// This is a slight edit from the default use case that we must do to
				// point our raw/direct STS call at our test server.
				"sts_endpoint": ts.URL,
			},
		}
		resp, err := testEnv.backend.HandleRequest(testEnv.ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Fatalf("expected nil response but received %+v", resp)
		}
	}
	renewalReq := &logical.Request{}
	{
		// vault login -method=aws header_value=vault.example.com role=dev-role-iam
		requestBody, err := GenerateLoginData(&credentials.Credentials{}, "vault.example.com", "us-east-1")
		if err != nil {
			t.Fatal(err)
		}
		requestBody["role"] = "dev-role-iam"
		req := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "login",
			Storage:   testEnv.conf.StorageView,
			Data:      requestBody,
		}
		resp, err := testEnv.backend.HandleRequest(testEnv.ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Auth == nil {
			t.Fatal("expected to receive auth")
		}
		renewalReq.Auth = resp.Auth
	}
	{
		// Test renewal.
		renewalReq.Storage = testEnv.conf.StorageView
		resp, err := testEnv.backend.pathLoginRenew(testEnv.ctx, renewalReq, nil)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response but received")
		}
		if resp.Auth == nil {
			t.Fatal("expected to receive auth")
		}
	}
}

func newTestEnvironment() (*testEnvironment, error) {
	ctx := context.Background()
	conf := &logical.BackendConfig{
		StorageView: &logical.InmemStorage{},
		Logger:      hclog.Default(),
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: 24 * time.Hour * 32,
			MaxLeaseTTLVal:     24 * time.Hour * 32,
		},
		BackendUUID: "1234-5678-9012-3456",
	}
	b, err := Factory(ctx, conf)
	if err != nil {
		return nil, err
	}
	return &testEnvironment{
		ctx:     ctx,
		conf:    conf,
		backend: b.(*backend),
	}, nil
}

type testEnvironment struct {
	ctx     context.Context
	conf    *logical.BackendConfig
	backend *backend
}

type fakeSTSClient struct {
	stsiface.STSAPI
	getCallerIdentityOutputToReturn  *sts.GetCallerIdentityOutput
	getCallerIdentityRequestToReturn *request.Request
}

func (f *fakeSTSClient) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return f.getCallerIdentityOutputToReturn, nil
}

func (f *fakeSTSClient) GetCallerIdentityRequest(*sts.GetCallerIdentityInput) (req *request.Request, output *sts.GetCallerIdentityOutput) {
	return f.getCallerIdentityRequestToReturn, nil
}

type fakeEC2Client struct {
	ec2iface.EC2API
	describeInstanceOutputToReturn *ec2.DescribeInstancesOutput
}

func (f *fakeEC2Client) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return f.describeInstanceOutputToReturn, nil
}

type fakeIAMClient struct {
	iamiface.IAMAPI
	roleOutputToReturn *iam.GetRoleOutput
}

func (f *fakeIAMClient) GetInstanceProfile(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
	return nil, nil
}

func (f *fakeIAMClient) GetRole(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
	return f.roleOutputToReturn, nil
}

func (f *fakeIAMClient) GetUser(*iam.GetUserInput) (*iam.GetUserOutput, error) {
	return nil, nil
}
