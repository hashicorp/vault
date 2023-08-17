package approle

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
)

// testHTTPServer creates a test HTTP server that handles requests until
// the listener returned is closed.
func testHTTPServer(
	t *testing.T, handler http.Handler,
) (*api.Config, net.Listener) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	server := &http.Server{Handler: handler}
	go server.Serve(ln)

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("http://%s", ln.Addr())

	return config, ln
}

func init() {
	os.Setenv("VAULT_TOKEN", "")
}

func TestLogin(t *testing.T) {
	secretIDEnvVar := "APPROLE_SECRET_ID"
	allowedRoleID := "my-role-id"
	allowedSecretID := "my-secret-id"

	content := []byte(allowedSecretID)
	tmpfile, err := os.CreateTemp("", "file-containing-secret-id")
	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	err = os.Setenv(secretIDEnvVar, allowedSecretID)
	if err != nil {
		t.Fatalf("error writing secret ID to env var: %v", err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("error writing to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("error closing temp file: %v", err)
	}

	// a response to return if the correct values were passed to login
	authSecret := &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: "a-client-token",
		},
	}

	authBytes, err := json.Marshal(authSecret)
	if err != nil {
		t.Fatalf("error marshaling json: %v", err)
	}

	handler := func(w http.ResponseWriter, req *http.Request) {
		payload := make(map[string]interface{})
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			t.Fatalf("error decoding json: %v", err)
		}
		if payload["role_id"] == allowedRoleID && payload["secret_id"] == allowedSecretID {
			w.Write(authBytes)
		}
	}

	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	config.Address = strings.ReplaceAll(config.Address, "127.0.0.1", "localhost")
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatalf("error initializing Vault client: %v", err)
	}

	authFromFile, err := NewAppRoleAuth(allowedRoleID, &SecretID{FromFile: tmpfile.Name()})
	if err != nil {
		t.Fatalf("error initializing AppRoleAuth with secret ID file: %v", err)
	}

	loginRespFromFile, err := client.Auth().Login(context.TODO(), authFromFile)
	if err != nil {
		t.Fatalf("error logging in with secret ID from file: %v", err)
	}
	if loginRespFromFile.Auth == nil || loginRespFromFile.Auth.ClientToken == "" {
		t.Fatalf("no authentication info returned by login")
	}

	authFromEnv, err := NewAppRoleAuth(allowedRoleID, &SecretID{FromEnv: secretIDEnvVar})
	if err != nil {
		t.Fatalf("error initializing AppRoleAuth with secret ID env var: %v", err)
	}

	loginRespFromEnv, err := client.Auth().Login(context.TODO(), authFromEnv)
	if err != nil {
		t.Fatalf("error logging in with secret ID from env var: %v", err)
	}
	if loginRespFromEnv.Auth == nil || loginRespFromEnv.Auth.ClientToken == "" {
		t.Fatalf("no authentication info returned by login with secret ID from env var")
	}

	authFromStr, err := NewAppRoleAuth(allowedRoleID, &SecretID{FromString: allowedSecretID})
	if err != nil {
		t.Fatalf("error initializing AppRoleAuth with secret ID string: %v", err)
	}

	loginRespFromStr, err := client.Auth().Login(context.TODO(), authFromStr)
	if err != nil {
		t.Fatalf("error logging in with string: %v", err)
	}
	if loginRespFromStr.Auth == nil || loginRespFromStr.Auth.ClientToken == "" {
		t.Fatalf("no authentication info returned by login with secret ID from string")
	}
}
