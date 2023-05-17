// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	docker.ServiceURL
	Token string
}

func (c *Config) APIConfig() *nomadapi.Config {
	apiConfig := nomadapi.DefaultConfig()
	apiConfig.Address = c.URL().String()
	apiConfig.SecretID = c.Token
	return apiConfig
}

func (c *Config) Client() (*nomadapi.Client, error) {
	apiConfig := c.APIConfig()

	return nomadapi.NewClient(apiConfig)
}

func prepareTestContainer(t *testing.T, bootstrap bool) (func(), *Config) {
	if retAddress := os.Getenv("NOMAD_ADDR"); retAddress != "" {
		s, err := docker.NewServiceURLParse(retAddress)
		if err != nil {
			t.Fatal(err)
		}
		return func() {}, &Config{*s, os.Getenv("NOMAD_TOKEN")}
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/multani/nomad",
		ImageTag:      "1.1.6",
		ContainerName: "nomad",
		Ports:         []string{"4646/tcp"},
		Cmd:           []string{"agent", "-dev"},
		Env:           []string{`NOMAD_LOCAL_CONFIG=bind_addr = "0.0.0.0" acl { enabled = true }`},
	})
	if err != nil {
		t.Fatalf("Could not start docker Nomad: %s", err)
	}

	var nomadToken string
	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		var err error
		nomadapiConfig := nomadapi.DefaultConfig()
		nomadapiConfig.Address = fmt.Sprintf("http://%s:%d/", host, port)
		nomad, err := nomadapi.NewClient(nomadapiConfig)
		if err != nil {
			return nil, err
		}

		_, err = nomad.Status().Leader()
		if err != nil {
			t.Logf("[DEBUG] Nomad is not ready yet: %s", err)
			return nil, err
		}

		if bootstrap {
			aclbootstrap, _, err := nomad.ACLTokens().Bootstrap(nil)
			if err != nil {
				return nil, err
			}
			nomadToken = aclbootstrap.SecretID
			t.Logf("[WARN] Generated Master token: %s", nomadToken)
		}

		nomadAuthConfig := nomadapi.DefaultConfig()
		nomadAuthConfig.Address = nomad.Address()

		if bootstrap {
			nomadAuthConfig.SecretID = nomadToken

			nomadAuth, err := nomadapi.NewClient(nomadAuthConfig)
			if err != nil {
				return nil, err
			}

			err = preprePolicies(nomadAuth)
			if err != nil {
				return nil, err
			}
		}

		u, _ := docker.NewServiceURLParse(nomadapiConfig.Address)
		return &Config{
			ServiceURL: *u,
			Token:      nomadToken,
		}, nil
	})
	if err != nil {
		t.Fatalf("Could not start docker Nomad: %s", err)
	}

	return svc.Cleanup, svc.Config.(*Config)
}

func preprePolicies(nomadClient *nomadapi.Client) error {
	policy := &nomadapi.ACLPolicy{
		Name:        "test",
		Description: "test",
		Rules: `namespace "default" {
        policy = "read"
      }
      `,
	}
	anonPolicy := &nomadapi.ACLPolicy{
		Name:        "anonymous",
		Description: "Deny all access for anonymous requests",
		Rules: `namespace "default" {
            policy = "deny"
        }
        agent {
            policy = "deny"
        }
        node {
            policy = "deny"
        }
        `,
	}

	_, err := nomadClient.ACLPolicies().Upsert(policy, nil)
	if err != nil {
		return err
	}

	_, err = nomadClient.ACLPolicies().Upsert(anonPolicy, nil)
	if err != nil {
		return err
	}

	return nil
}

func TestBackend_config_Bootstrap(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, svccfg := prepareTestContainer(t, false)
	defer cleanup()

	connData := map[string]interface{}{
		"address": svccfg.URL().String(),
		"token":   "",
	}

	confReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Storage:   config.StorageView,
		Data:      connData,
	}

	resp, err := b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	confReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	expected := map[string]interface{}{
		"address":               connData["address"].(string),
		"max_token_name_length": 0,
		"ca_cert":               "",
		"client_cert":           "",
	}
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
	}

	nomadClient, err := svccfg.Client()
	if err != nil {
		t.Fatalf("failed to construct nomaad client, %v", err)
	}

	token, _, err := nomadClient.ACLTokens().Bootstrap(nil)
	if err == nil {
		t.Fatalf("expected acl system to be bootstrapped already, but was able to get the bootstrap token : %v", token)
	}
	// NOTE: fragile test, but it's the only way, AFAIK, to check that nomad is
	// bootstrapped
	if !strings.Contains(err.Error(), "bootstrap already done") {
		t.Fatalf("expected acl system to be bootstrapped already: err: %v", err)
	}
}

func TestBackend_config_access(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, svccfg := prepareTestContainer(t, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": svccfg.URL().String(),
		"token":   svccfg.Token,
	}

	confReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Storage:   config.StorageView,
		Data:      connData,
	}

	resp, err := b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	confReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	expected := map[string]interface{}{
		"address":               connData["address"].(string),
		"max_token_name_length": 0,
		"ca_cert":               "",
		"client_cert":           "",
	}
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
	}
	if resp.Data["token"] != nil {
		t.Fatalf("token should not be set in the response")
	}
}

func TestBackend_config_access_with_certs(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, svccfg := prepareTestContainer(t, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address":     svccfg.URL().String(),
		"token":       svccfg.Token,
		"ca_cert":     caCert,
		"client_cert": clientCert,
		"client_key":  clientKey,
	}

	confReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Storage:   config.StorageView,
		Data:      connData,
	}

	resp, err := b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	confReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	expected := map[string]interface{}{
		"address":               connData["address"].(string),
		"max_token_name_length": 0,
		"ca_cert":               caCert,
		"client_cert":           clientCert,
	}
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
	}
	if resp.Data["token"] != nil {
		t.Fatalf("token should not be set in the response")
	}
}

func TestBackend_renew_revoke(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, svccfg := prepareTestContainer(t, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": svccfg.URL().String(),
		"token":   svccfg.Token,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "role/test"
	req.Data = map[string]interface{}{
		"policies": []string{"policy"},
		"lease":    "6h",
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}

	generatedSecret := resp.Secret
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token    string `mapstructure:"secret_id"`
		Accessor string `mapstructure:"accessor_id"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}
	t.Logf("[WARN] Generated token: %s with accessor %s", d.Token, d.Accessor)

	// Build a client and verify that the credentials work
	nomadapiConfig := nomadapi.DefaultConfig()
	nomadapiConfig.Address = connData["address"].(string)
	nomadapiConfig.SecretID = d.Token
	client, err := nomadapi.NewClient(nomadapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[WARN] Verifying that the generated token works...")
	_, err = client.Agent().Members, nil
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Build a management client and verify that the token does not exist anymore
	nomadmgmtConfig := nomadapi.DefaultConfig()
	nomadmgmtConfig.Address = connData["address"].(string)
	nomadmgmtConfig.SecretID = connData["token"].(string)
	mgmtclient, err := nomadapi.NewClient(nomadmgmtConfig)
	if err != nil {
		t.Fatal(err)
	}

	q := &nomadapi.QueryOptions{
		Namespace: "default",
	}

	t.Log("[WARN] Verifying that the generated token does not exist...")
	_, _, err = mgmtclient.ACLTokens().Info(d.Accessor, q)
	if err == nil {
		t.Fatal("err: expected error")
	}
}

func TestBackend_CredsCreateEnvVar(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, svccfg := prepareTestContainer(t, true)
	defer cleanup()

	req := logical.TestRequest(t, logical.UpdateOperation, "role/test")
	req.Data = map[string]interface{}{
		"policies": []string{"policy"},
		"lease":    "6h",
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("NOMAD_TOKEN", svccfg.Token)
	defer os.Unsetenv("NOMAD_TOKEN")
	os.Setenv("NOMAD_ADDR", svccfg.URL().String())
	defer os.Unsetenv("NOMAD_ADDR")

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}
}

func TestBackend_max_token_name_length(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, svccfg := prepareTestContainer(t, true)
	defer cleanup()

	testCases := []struct {
		title       string
		roleName    string
		tokenLength int
	}{
		{
			title: "Default",
		},
		{
			title:       "ConfigOverride",
			tokenLength: 64,
		},
		{
			title:       "ConfigOverride-LongName",
			roleName:    "testlongerrolenametoexceed64charsdddddddddddddddddddddddd",
			tokenLength: 64,
		},
		{
			title:    "Notrim",
			roleName: "testlongersubrolenametoexceed64charsdddddddddddddddddddddddd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// setup config/access
			connData := map[string]interface{}{
				"address":               svccfg.URL().String(),
				"token":                 svccfg.Token,
				"max_token_name_length": tc.tokenLength,
			}
			expected := map[string]interface{}{
				"address":               svccfg.URL().String(),
				"max_token_name_length": tc.tokenLength,
				"ca_cert":               "",
				"client_cert":           "",
			}

			expectedMaxTokenNameLength := maxTokenNameLength
			if tc.tokenLength != 0 {
				expectedMaxTokenNameLength = tc.tokenLength
			}

			confReq := logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "config/access",
				Storage:   config.StorageView,
				Data:      connData,
			}

			resp, err := b.HandleRequest(context.Background(), &confReq)
			if err != nil || (resp != nil && resp.IsError()) || resp != nil {
				t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
			}
			confReq.Operation = logical.ReadOperation
			resp, err = b.HandleRequest(context.Background(), &confReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
			}

			// verify token length is returned in the config/access query
			if !reflect.DeepEqual(expected, resp.Data) {
				t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
			}
			// verify token is not returned
			if resp.Data["token"] != nil {
				t.Fatalf("token should not be set in the response")
			}

			// create a role to create nomad credentials with
			// Seeds random with current timestamp

			if tc.roleName == "" {
				tc.roleName = "test"
			}
			roleTokenName := testhelpers.RandomWithPrefix(tc.roleName)

			confReq.Path = "role/" + roleTokenName
			confReq.Operation = logical.UpdateOperation
			confReq.Data = map[string]interface{}{
				"policies": []string{"policy"},
				"lease":    "6h",
			}
			resp, err = b.HandleRequest(context.Background(), &confReq)
			if err != nil {
				t.Fatal(err)
			}

			confReq.Operation = logical.ReadOperation
			confReq.Path = "creds/" + roleTokenName
			resp, err = b.HandleRequest(context.Background(), &confReq)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("resp nil")
			}
			if resp.IsError() {
				t.Fatalf("resp is error: %v", resp.Error())
			}

			// extract the secret, so we can query nomad directly
			generatedSecret := resp.Secret
			generatedSecret.TTL = 6 * time.Hour

			var d struct {
				Token    string `mapstructure:"secret_id"`
				Accessor string `mapstructure:"accessor_id"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				t.Fatal(err)
			}

			// Build a client and verify that the credentials work
			nomadapiConfig := nomadapi.DefaultConfig()
			nomadapiConfig.Address = connData["address"].(string)
			nomadapiConfig.SecretID = d.Token
			client, err := nomadapi.NewClient(nomadapiConfig)
			if err != nil {
				t.Fatal(err)
			}

			// default query options for Nomad queries ... not sure if needed
			qOpts := &nomadapi.QueryOptions{
				Namespace: "default",
			}

			// connect to Nomad and verify the token name does not exceed the
			// max_token_name_length
			token, _, err := client.ACLTokens().Self(qOpts)
			if err != nil {
				t.Fatal(err)
			}

			if len(token.Name) > expectedMaxTokenNameLength {
				t.Fatalf("token name exceeds max length (%d): %s (%d)", expectedMaxTokenNameLength, token.Name, len(token.Name))
			}
		})
	}
}

const caCert = `-----BEGIN CERTIFICATE-----
MIIF7zCCA9egAwIBAgIINVVQic4bju8wDQYJKoZIhvcNAQELBQAwaDELMAkGA1UE
BhMCVVMxFDASBgNVBAoMC1Vuc3BlY2lmaWVkMR8wHQYDVQQLDBZjYS0zODQzMDY2
NDA5ODI5MjQwNTU5MSIwIAYDVQQDDBl4cHMxNS5sb2NhbC5jaXBoZXJib3kuY29t
MB4XDTIyMDYwMjIxMTgxN1oXDTIzMDcwNTIxMTgxN1owaDELMAkGA1UEBhMCVVMx
FDASBgNVBAoMC1Vuc3BlY2lmaWVkMR8wHQYDVQQLDBZjYS0zODQzMDY2NDA5ODI5
MjQwNTU5MSIwIAYDVQQDDBl4cHMxNS5sb2NhbC5jaXBoZXJib3kuY29tMIICIjAN
BgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA35VilgfqMUKhword7wORXRFyPbpz
8uqO7eRaylMnkAkbk5eoQB/iYfXjJ6ZBs5mJGQVz5ZNvh9EzZsk1J6wqYgbwVKUx
fh4kvW6sXtDirtb4ZQAK7OTLEoapUQGnGcvm+aEYfvC1sTBl4fbex7yyN5FYMJTM
TAUumhdq2pwujaj2xkN9DwZa89Tk7tbj9HE9DTRji7bnciEtrmTAOIOfOrT/1l3x
YW1BwYXpQ0TamJ58pC/iNgEp5FAxKt9d3RggesMA7pvG/f8fNgsa/Tku/PeEXNPA
+Yx4CcAipujmqpBKiKwJ6TOzp80m2zrZ7Da4Av5vVS5GsNJxhFYD1h8hU1ptK9BS
2CaTwBpV421C9BfEmtSAksGDIWYujfiHb6XNaQrt8Hu85GBuPUudVn0lpoXLn2xD
rGK8WEK2gWZ4eez3ZDLbpLui6c1m7AVlMtj374s+LHcD7JIxY475Na7pXmEWReqM
RUyCEq1spOOn70fOdhphhmpY6DoklOTOriPawCLNmkPWRnhrIwqyP1gse9YMqQ2n
LhWUkv/08m/0pb4e5ijVhsZNzv+1PXPWCk968nzt0BMDgJT+0ZiXsaU7FILXuo7Y
Ijgrj7dpXWx2MBdMGPFQdveog7Pa80Yb7r4ERW0DL78TxYC6m/S1p14PHwZpDZzQ
LrPrBcpI5XzI7osCAwEAAaOBnDCBmTAOBgNVHQ8BAf8EBAMCAqQwDAYDVR0TBAUw
AwEB/zA0BgNVHR4ELTAroCkwG4IZeHBzMTUubG9jYWwuY2lwaGVyYm95LmNvbTAK
hwh/AAAB/wAAADAkBgNVHREEHTAbghl4cHMxNS5sb2NhbC5jaXBoZXJib3kuY29t
MB0GA1UdDgQWBBR3bHgDp5RpzerMKRkaGDFN/ZeImjANBgkqhkiG9w0BAQsFAAOC
AgEArkuDYYWYHYxIoTeZkQz5L1y0H27ZWPJx5jBBuktPonDLQxBGAwkl6NZbJGLU
v+usII+eyjPKIgjhCiTXJAmeIngwWoN3PHOLMIPe9axuNt6qVoP4dQtzfpPR3buK
CWj9i3H0ixK73klk7QWZiBUDinYfEMSNRpU3G7NsqmqCXD4s5gB+8y9c7+zIiJyN
IaJBWpzI4eQBi/4cBhtM7Xa+CMB/8whhWYR6H+GXGZdNcP5f7bwneMstWKceTadk
IEzFucJHDySpEkIA2A9t33pV54FmEp+JVwvxAH4FABCnjPmhg0j1IonWV5pySWpG
hhEZpnRRH1XfpTA5i6dlyUA5DJjL8X1lYrgOK+LaoR52mQh5JBsMoVHFzN50DiMA
RTsbq4Qzozf23hU1BqW4NOzPTukgSGEcbT/DhXKPPPLL8JD0rPelJPq76X3TJjgZ
C9uMnZaDnxjppDXp5oBIXqC05FDxJ5sSODNOpKGyuzOU2qQLMau33yYOgaSAttBk
r29+LNFJ+0QzMuPjYXPznpxbsI+lrlZ3F2tDGGs8+JVceC1YX+cBEsEOiqNGTIip
/DY3b9gu5oiTwhcFyQW8+WFsirRS/g5t+M40WLKVPdK09z96krFXQMkL6a7LHLY1
n9ivwj+sTG1XmJYXp8naLg4wdzIUf2fJxaFNI5Yq4elZ8sY=
-----END CERTIFICATE-----`

const clientCert = `-----BEGIN CERTIFICATE-----
MIIEsDCCApigAwIBAgIIRY1JBRIynFYwDQYJKoZIhvcNAQELBQAwaDELMAkGA1UE
BhMCVVMxFDASBgNVBAoMC1Vuc3BlY2lmaWVkMR8wHQYDVQQLDBZjYS0zODQzMDY2
NDA5ODI5MjQwNTU5MSIwIAYDVQQDDBl4cHMxNS5sb2NhbC5jaXBoZXJib3kuY29t
MB4XDTIyMDYwMjIxMTgxOFoXDTIzMDcwNTIxMTgxOFowRzELMAkGA1UEBhMCVVMx
FDASBgNVBAoMC1Vuc3BlY2lmaWVkMSIwIAYDVQQDDBl4cHMxNS5sb2NhbC5jaXBo
ZXJib3kuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs+XYhsW2
vTwN7gY3xMxgbNN8d3aoeqCswOp05BBf0Vgv3febahm422ubXXd5Mg2UGiU7sJVe
4tUpDeupVVRX5Qr/hpiXgEyfRDAAAJKqrl65KSS62TCbT/eJZ0ah25HV1evI4uM2
0kl5QWhtQjDyaVlTS38YFqXXQvpOuU5DG6UbKnpMcpsCPTyUKEJvJ95ZLcz0HJ8I
kIHrnX0Lt0pOhkllj5Nk4cXhU8CFk8IGNz7SVAycrUsffAUMNNEbrIOIfOTPHR1c
q3X9hO4/5pt80uIDMFwwumoA7nQR0AhlKkw9SskCIzJhKwKwssQY7fmovNG0fOEd
/+vSHK7OsYW+gwIDAQABo38wfTAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYI
KwYBBQUHAwIwCQYDVR0TBAIwADAqBgNVHREEIzAhghl4cHMxNS5sb2NhbC5jaXBo
ZXJib3kuY29thwR/AAABMB8GA1UdIwQYMBaAFHdseAOnlGnN6swpGRoYMU39l4ia
MA0GCSqGSIb3DQEBCwUAA4ICAQBUSP4ZJglCCrYkM5Le7McdvfkM5uYv1aQn0sM4
gbyDEWO0fnv50vLpD3y4ckgHLoD52pAZ0hN8a7rwAUae21GA6DvEchSH5x/yvJiS
7FBlq39sAafe03ZlzDErNYJRkLcnPAqG74lJ1SSsMcs9gCPHM8R7HtNnhAga06L7
K8/G43dsGZCmEb+xcX2B9McCt8jBG6TJPTGafb3BJ0JTmR/tHdoLFIiNwI+qzd2U
lMnGlkIApULX8tmIMsWO0rjdiFkPWGcmfn9ChC0iDpQOAcKSDBcZlWrDNpzKk0mK
l0TbE6cxcmCUUpiwaXFrbkwVWQw4W0c4b3sWFtWifFbiR1qZ/OT2Y2sHbkbxwvPl
PjjXMDBAdRRwtNcTP1E55I5zvwzzBxUpxOob0miorhTJrZR9So0rgv7Roce4ED6M
WETYa/mGhe+Q7gBQygIVoryfQLgGBsHC+7V4RDvYTazwZkz9nLQxHLI/TAZU5ofM
WqdoUkMd68rxTTEUoMfGbftxjKA0raxGcO7/PjLR3O743EwCqeqYJ7OKWgGRLnui
kIKNUJlZ9umURUFzL++Bx4Pr95jWXb2WYqYYQxhDz0oR5q5smnFm5+/1/MLDMvDU
TrgBK6pey4QF33B/I55H1+7tGdv85Q57Z8UrNi/IQxR2sFlsOTeCwStpBQ56sdZk
Wi4+cQ==
-----END CERTIFICATE-----`

const clientKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCz5diGxba9PA3u
BjfEzGBs03x3dqh6oKzA6nTkEF/RWC/d95tqGbjba5tdd3kyDZQaJTuwlV7i1SkN
66lVVFflCv+GmJeATJ9EMAAAkqquXrkpJLrZMJtP94lnRqHbkdXV68ji4zbSSXlB
aG1CMPJpWVNLfxgWpddC+k65TkMbpRsqekxymwI9PJQoQm8n3lktzPQcnwiQgeud
fQu3Sk6GSWWPk2ThxeFTwIWTwgY3PtJUDJytSx98BQw00Rusg4h85M8dHVyrdf2E
7j/mm3zS4gMwXDC6agDudBHQCGUqTD1KyQIjMmErArCyxBjt+ai80bR84R3/69Ic
rs6xhb6DAgMBAAECggEAPBcja2kxcCZWNNKo4DiwYMmHwtPE1SlEazAlmWSKzP+b
BZbGt/sdj1VzURYuSnTUqqMTPBm41yYCj57PMix5K42v6sKfoIB3lqw94/MZxiLn
0IFvVErzJhP2NqQWPqSI++rFcFwbHMTkFuAN1tVIs73dn9M1NaNxsvKvRyCIM/wz
5YQSDyTkdW4jQM2RvUFOoqwmeyAlQoBRMgQ4bHfLHxmPEjFgw1MAmmG8bJdkupin
MVzhZyKj4Fh80Xa2MU4KokijjG41hmYbg/sjNHaHJFDA92Rwq13dhWytrauJDxa/
3yj8pHWc23Y3hXvRAf/cibDVzXmmLj49W1i06KuUCQKBgQDj5yF/DJV0IOkhfbol
+f5AGH4ZrEXA/JwA5SxHU+aKhUuPEqK/LeUWqiy3szFjOz2JOnCC0LMN42nsmMyK
sdQEKHp2SPd2wCxsAKZAuxrEi6yBt1mEPFFU5yzvZbdMqYChKJjm9fbRHtuc63s8
PyVw67Ii9o4ij+PxfTobIs18xwKBgQDKE59w3uUDt2uoqNC8x4m5onL2p2vtcTHC
CxU57mu1+9CRM8N2BEp2VI5JaXjqt6W4u9ISrmOqmsPgTwosAquKpA/nu3bVvR9g
WlN9dh2Xgza0/AFaA9CB++ier8RJq5xFlcasMUmgkhYt3zgKNgRDfjfREWM0yamm
P++hAYRcZQKBgHEuYQk6k6J3ka/rQ54GmEj2oPFZB88+5K7hIWtO9IhIiGzGYYK2
ZTYrT0fvuxA/5GCZYDTnNnUoQnuYqsQaamOiQqcpt5QG/kiozegJw9JmV0aYauFs
HyweHsfJaQ2uhE4E3mKdNnVGcORuYeZaqdp5gx8v+QibEyXj/g5p60kTAoGBALKp
TMOHXmW9yqKwtvThWoRU+13WQlcJSFvuXpL8mCCrBgkLAhqaypb6RV7ksLKdMhk1
fhNkOdxBv0LXvv+QUMhgK2vP084/yrjuw3hecOVfboPvduZ2DuiNp2p9rocQAjeH
p8LgRN+Bqbhe7fYhMf3WX1UqEVM/pQ3G43+vjq39AoGAOyD2/hFSIx6BMddUNTHG
BEsMUc/DHYslZebbF1zAWnkKdTt+URhtHAFB2tYRDgkZfwW+wr/w12dJTIkX965o
HO7tI4FgpU9b0i8FTuwYkBfjwp2j0Xd2/VBR8Qpd17qKl3I6NXDsf3ykjGZAvldH
Tll+qwEZpXSRa5OWWTpGV8I=
-----END PRIVATE KEY-----`
