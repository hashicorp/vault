package ldap

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
	t *testing.T, handler http.Handler) (*api.Config, net.Listener) {
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
	allowedPassword := "6hrtL!*bro!ywbQbvDwW"

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
		if payload["password"] == allowedPassword {
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

	authFromStr, err := NewLDAPAuth("my-ldap-username", &Password{FromString: allowedPassword})
	if err != nil {
		t.Fatalf("error initializing LDAPAuth with password string: %v", err)
	}
	loginRespFromStr, err := client.Auth().Login(context.TODO(), authFromStr)
	if err != nil {
		t.Fatalf("error logging in with string: %v", err)
	}
	if loginRespFromStr.Auth == nil || loginRespFromStr.Auth.ClientToken == "" {
		t.Fatalf("no authentication info returned by login with password from string")
	}

	_, err = NewLDAPAuth("", &Password{FromString: allowedPassword})
	if err.Error() != "no user name provided for login" {
		t.Fatalf("Auth object created for empty username: %v", err)
	}

	_, err = NewLDAPAuth("my-ldap-username", nil)
	if err.Error() != "no password provided for login" {
		t.Fatalf("Auth object created when passing a nil Password struct: %v", err)
	}
	ldapMount := WithMountPath("customMount")
	_, err = NewLDAPAuth("my-ldap-username", &Password{FromString: allowedPassword}, ldapMount)
	if err != nil {
		t.Fatalf("error initializing LDAPAuth with custom mountpath: %v", err)
	}
}
