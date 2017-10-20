package util

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"regexp"
	"strings"
)

const (
	labelRegex string = "^(?P<key>[a-z]([\\w-]+)?):(?P<value>[\\w-]*)$"
)

// GcpCredentials represents a simplified version of the Google Cloud Platform credentials file format.
type GcpCredentials struct {
	ClientEmail  string `json:"client_email" structs:"client_email" mapstructure:"client_email"`
	ClientId     string `json:"client_id" structs:"client_id" mapstructure:"client_id"`
	PrivateKeyId string `json:"private_key_id" structs:"private_key_id" mapstructure:"private_key_id"`
	PrivateKey   string `json:"private_key" structs:"private_key" mapstructure:"private_key"`
	ProjectId    string `json:"project_id" structs:"project_id" mapstructure:"project_id"`
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
	pemBytes, err := base64.StdEncoding.DecodeString(pemString)
	if err != nil {
		return nil, err
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

// OAuth2RSAPublicKey returns the PEM key file string for Google Oauth2 public cert for the given 'kid' id.
func OAuth2RSAPublicKey(kid, oauth2BasePath string) (interface{}, error) {
	oauth2Client, err := googleoauth2.New(cleanhttp.DefaultClient())
	if err != nil {
		return "", err
	}

	if len(oauth2BasePath) > 0 {
		oauth2Client.BasePath = oauth2BasePath
	}

	jwks, err := oauth2Client.GetCertForOpenIdConnect().Do()
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid && jose.SignatureAlgorithm(key.Alg) == jose.RS256 {
			// Trim extra '=' from key so it can be parsed.
			key.N = strings.TrimRight(key.N, "=")
			js, err := key.MarshalJSON()
			if err != nil {
				return nil, fmt.Errorf("unable to marshal json %v", err)
			}
			key := &jose.JSONWebKey{}
			if err := key.UnmarshalJSON(js); err != nil {
				return nil, fmt.Errorf("unable to unmarshal json %v", err)
			}

			return key.Key, nil
		}
	}

	return nil, fmt.Errorf("could not find public key with kid '%s'", kid)
}

func ParseGcpLabels(labels []string) (parsed map[string]string, invalid []string) {
	parsed = map[string]string{}
	invalid = []string{}

	re := regexp.MustCompile(labelRegex)
	for _, labelStr := range labels {
		matches := re.FindStringSubmatch(labelStr)
		if len(matches) == 0 {
			invalid = append(invalid, labelStr)
			continue
		}

		captureNames := re.SubexpNames()
		var keyPtr, valPtr *string
		for i, name := range captureNames {
			if name == "key" {
				keyPtr = &matches[i]
			} else if name == "value" {
				valPtr = &matches[i]
			}
		}

		if keyPtr == nil || valPtr == nil || len(*keyPtr) < 1 {
			invalid = append(invalid, labelStr)
			continue
		} else {
			parsed[*keyPtr] = *valPtr
		}
	}

	return parsed, invalid
}
