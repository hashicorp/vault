package gcputil

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/googleapi"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	defaultHomeCredentialsFile                = ".gcp/credentials"

	// Global URL: https://cloud.google.com/iam/docs/creating-managing-service-account-keys
	serviceAccountPublicKeyUrlTemplate        = "https://www.googleapis.com/service_accounts/v1/metadata/x509/%s?alt=json"

	// Global URL: Base URL from golang.org/x/oauth2, v1 returns x509 keys
	googleOauthProviderX509CertUrl            = "https://www.googleapis.com/oauth2/v1/certs"
)

// GcpCredentials represents a simplified version of the Google Cloud Platform credentials file format.
type GcpCredentials struct {
	ClientEmail  string `json:"client_email" structs:"client_email" mapstructure:"client_email"`
	ClientId     string `json:"client_id" structs:"client_id" mapstructure:"client_id"`
	PrivateKeyId string `json:"private_key_id" structs:"private_key_id" mapstructure:"private_key_id"`
	PrivateKey   string `json:"private_key" structs:"private_key" mapstructure:"private_key"`
	ProjectId    string `json:"project_id" structs:"project_id" mapstructure:"project_id"`
}

// FindCredentials attempts to obtain GCP credentials in the
// following ways:
// * Parse JSON from provided credentialsJson
// * Parse JSON from the environment variables GOOGLE_CREDENTIALS or GOOGLE_CLOUD_KEYFILE_JSON
// * Parse JSON file ~/.gcp/credentials
// * Google Application Default Credentials (see https://developers.google.com/identity/protocols/application-default-credentials)
func FindCredentials(credsJson string, ctx context.Context, scopes ...string) (*GcpCredentials, oauth2.TokenSource, error) {
	var creds *GcpCredentials
	var err error
	// 1. Parse JSON from provided credentialsJson
	if credsJson == "" {
		// 2. JSON from env var GOOGLE_CREDENTIALS
		credsJson = os.Getenv("GOOGLE_CREDENTIALS")
	}

	if credsJson == "" {
		// 3. JSON from env var GOOGLE_CLOUD_KEYFILE_JSON
		credsJson = os.Getenv("GOOGLE_CLOUD_KEYFILE_JSON")
	}

	if credsJson == "" {
		// 4. JSON from ~/.gcp/credentials
		home, err := homedir.Dir()
		if err != nil {
			return nil, nil, errors.New("could not find home directory")
		}
		credBytes, err := ioutil.ReadFile(filepath.Join(home, defaultHomeCredentialsFile))
		if err == nil {
			credsJson = string(credBytes)
		}
	}

	// Parse JSON into credentials.
	if credsJson != "" {
		creds, err = Credentials(credsJson)
		if err == nil {
			conf := jwt.Config{
				Email:      creds.ClientEmail,
				PrivateKey: []byte(creds.PrivateKey),
				Scopes:     scopes,
				TokenURL:   "https://accounts.google.com/o/oauth2/token",
			}
			return creds, conf.TokenSource(ctx), nil
		}
	}

	// 5. Use Application default credentials.
	defaultCreds, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		return nil, nil, err
	}

	if defaultCreds.JSON != nil {
		creds, err = Credentials(string(defaultCreds.JSON))
		if err != nil {
			return nil, nil, errors.New("could not read credentials from application default credential JSON")
		}
	}

	return creds, defaultCreds.TokenSource, nil
}

// Credentials attempts to parse GcpCredentials from a JSON string.
func Credentials(credentialsJson string) (*GcpCredentials, error) {
	credentials := &GcpCredentials{}
	if err := json.Unmarshal([]byte(credentialsJson), &credentials); err != nil {
		return nil, err
	}
	return credentials, nil
}

// GetHttpClient creates an HTTP client from the given Google credentials and scopes.
func GetHttpClient(credentials *GcpCredentials, clientScopes ...string) (*http.Client, error) {
	conf := jwt.Config{
		Email:      credentials.ClientEmail,
		PrivateKey: []byte(credentials.PrivateKey),
		Scopes:     clientScopes,
		TokenURL:   "https://accounts.google.com/o/oauth2/token",
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cleanhttp.DefaultClient())
	client := conf.Client(ctx)
	return client, nil
}

// PublicKey returns a public key from a Google PEM key file (type TYPE_X509_PEM_FILE).
func PublicKey(pemString string) (interface{}, error) {
	// Attempt to base64 decode
	pemBytes := []byte(pemString)
	if b64decoded, err := base64.StdEncoding.DecodeString(pemString); err == nil {
		pemBytes = b64decoded
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("Unable to find pem block in key")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert.PublicKey, nil
}

func ServiceAccountPublicKey(serviceAccount string, keyId string) (interface{}, error) {
	keyUrl := fmt.Sprintf(serviceAccountPublicKeyUrlTemplate, url.PathEscape(serviceAccount))
	res, err := cleanhttp.DefaultClient().Get(keyUrl)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}

	jwks := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("unable to decode JSON response: %v", err)
	}
	kRaw, ok := jwks[keyId]
	if !ok {
		return nil, fmt.Errorf("service account %q key %q not found at GET %q", keyId, serviceAccount, keyUrl)
	}

	kStr, ok := kRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error - decoded JSON key value %v is not string", kRaw)
	}
	return PublicKey(kStr)
}

// OAuth2RSAPublicKey returns the PEM key file string for Google Oauth2 public cert for the given 'kid' id.
func OAuth2RSAPublicKey(ctx context.Context, keyId string) (interface{}, error) {
	certUrl := googleOauthProviderX509CertUrl
	res, err := cleanhttp.DefaultClient().Get(certUrl)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}

	jwks := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("unable to decode JSON response: %v", err)
	}
	kRaw, ok := jwks[keyId]
	if !ok {
		return nil, fmt.Errorf("key %q not found (GET %q)", keyId, certUrl)
	}

	kStr, ok := kRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error - decoded JSON key value %v is not string", kRaw)
	}
	return PublicKey(kStr)
}
