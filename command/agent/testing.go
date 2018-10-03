package agent

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const envVarRunAccTests = "VAULT_ACC"

var runAcceptanceTests = os.Getenv(envVarRunAccTests) == "1"

func GetTestJWT(t *testing.T) (string, *ecdsa.PrivateKey) {
	t.Helper()
	cl := jwt.Claims{
		Subject:   "r3qXcK2bix9eFECzsU3Sbmh0K16fatW6@clients",
		Issuer:    "https://team-vault.auth0.com/",
		NotBefore: jwt.NewNumericDate(time.Now().Add(-5 * time.Second)),
		Audience:  jwt.Audience{"https://vault.plugin.auth.jwt.test"},
	}

	privateCl := struct {
		User   string   `json:"https://vault/user"`
		Groups []string `json:"https://vault/groups"`
	}{
		"jeff",
		[]string{"foo", "bar"},
	}

	var key *ecdsa.PrivateKey
	block, _ := pem.Decode([]byte(TestECDSAPrivKey))
	if block != nil {
		var err error
		key, err = x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			t.Fatal(err)
		}
	}

	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.ES256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		t.Fatal(err)
	}

	raw, err := jwt.Signed(sig).Claims(cl).Claims(privateCl).CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}

	return raw, key
}

func readToken(fileName string) (*logical.HTTPWrapInfo, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	wrapper := &logical.HTTPWrapInfo{}
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(wrapper); err != nil {
		return nil, err
	}
	return wrapper, nil
}

const (
	TestECDSAPrivKey string = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKfldwWLPYsHjRL9EVTsjSbzTtcGRu6icohNfIqcb6A+oAoGCCqGSM49
AwEHoUQDQgAE4+SFvPwOy0miy/FiTT05HnwjpEbSq+7+1q9BFxAkzjgKnlkXk5qx
hzXQvRmS4w9ZsskoTZtuUI+XX7conJhzCQ==
-----END EC PRIVATE KEY-----`

	TestECDSAPubKey string = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE4+SFvPwOy0miy/FiTT05HnwjpEbS
q+7+1q9BFxAkzjgKnlkXk5qxhzXQvRmS4w9ZsskoTZtuUI+XX7conJhzCQ==
-----END PUBLIC KEY-----`
)
