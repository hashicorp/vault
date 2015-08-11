package jwt

import (
	"fmt"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_basic(t *testing.T) {
	tokenClaims := map[string]interface{}{
		"iss": "Test Issuer", 
		"sub": "Test Subject",
		"aud": "Test Audience",
		"iat": 1438898720,
		"nbf": 1438898720,
		"exp": 1538898720,
		"jti": "jti",
		"ran": "random",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepWriteRole(t, "test", "HS256", "test"),
			testAccStepReadRole(t, "test", "HS256", "test"),
			testAccStepSignToken(t, "test", tokenClaims, false),
			testAccStepDeleteRole(t, "test"),
			testAccStepReadRole(t, "test", "HS256", "test"),
		},
	})
}

func TestBackend_defaults(t *testing.T) {
	tokenClaims := map[string]interface{}{
		"iat": 1438898720,
		"nbf": 1438898720,
		"exp": 1538898720,
		"jti": "9fe94d93-7bb4-434c-b197-731b4b4c70d3",
		"ran": "random",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepWriteRole(t, "test", "HS256", "test"),
			testAccStepReadRole(t, "test", "HS256", "test"),
			testAccStepSignToken(t, "test", tokenClaims, true),
			testAccStepDeleteRole(t, "test"),
			testAccStepReadRole(t, "test", "HS256", "test"),
		},
	})
}

func testAccStepWriteRole(t *testing.T, name string, algorithm string, key string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"algorithm": algorithm,
			"key": key,
			"default_issuer": "Test Default Issuer",
			"default_subject": "Test Default Subject",
			"default_audience": "Test Default Audience",
		},
	}
}

func testAccStepDeleteRole(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + name,
	}
}

func testAccStepReadRole(t *testing.T, name string, algorithm string, key string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				return fmt.Errorf("missing response")
			}
			var d struct {
				Name           string        `json:"name" mapstructure:"name"`
				Algorithm      string        `json:"algorithm" structs:"algorithm" mapstructure:"algorithm"`
				Key            string        `json:"key" structs:"key" mapstructure:"key"`
				Issuer         string        `json:"iss" structs:"iss" mapstructure:"iss"`
				Subject        string        `json:"sub" structs:"sub" mapstructure:"sub"`
				Audience       string        `json:"aud" structs:"aud" mapstructure:"aud"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Name != name {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.Algorithm != algorithm {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.Key != key {
				return fmt.Errorf("bad: %#v", d)
			}
			
			return nil
		},
	}
}

func testAccStepSignToken(t *testing.T, name string, tokenClaims map[string]interface{}, defaults bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "issue/" + name,
		Data:      tokenClaims,
		Check: func(resp *logical.Response) error {
			var d struct {
				JTI        string `mapstructure:"jti"`
				Token      string `mapstructure:"token"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if d.Token == "" {
				return fmt.Errorf("missing token")
			}

			token, err := jwt.Parse(d.Token, func(token *jwt.Token) (interface{}, error) {
				return token, nil
		    })
			if err != nil {
				return fmt.Errorf("error parsing token")
			}

			if d.JTI != token.Claims["jti"] {
				return fmt.Errorf("bad: %#v", d)
			}
			
			if token.Claims["ran"] != "random" {
				return fmt.Errorf("bad: %#v", d)
			}

			if defaults == true {
				if token.Claims["sub"] != "Test Default Subject" {
					return fmt.Errorf("bad: %#v", d)
				}
				if token.Claims["aud"] != "Test Default Audience" {
					return fmt.Errorf("bad: %#v", d)
				}
				if token.Claims["iss"] != "Test Default Issuer" {
					return fmt.Errorf("bad: %#v", d)
				}
			}

			return nil
		},
	}
}
