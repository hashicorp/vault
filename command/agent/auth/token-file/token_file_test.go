package token_file

import (
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"

	"github.com/hashicorp/vault/command/agent/auth"
)

func TestNewTokenFileAuthMethodEmptyConfig(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)
	_, err := NewTokenFileAuthMethod(&auth.AuthConfig{
		Logger: logger.Named("auth.method"),
		Config: map[string]interface{}{},
	})
	if err == nil {
		t.Fatal("Expected error due to empty config")
	}
}

func TestNewTokenFileRemoveTokenString(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)
	_, err := NewTokenFileAuthMethod(&auth.AuthConfig{
		Logger: logger.Named("auth.method"),
		Config: map[string]interface{}{
			"token_file_path":                 "/some/path",
			"remove_token_file_after_reading": "string",
		},
	})
	if err == nil {
		t.Fatalf("Expected error when giving string for bool")
	}
}

func TestNewTokenFileEmptyFilePath(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)
	_, err := NewTokenFileAuthMethod(&auth.AuthConfig{
		Logger: logger.Named("auth.method"),
		Config: map[string]interface{}{
			"token_file_path": "",
		},
	})
	if err == nil {
		t.Fatalf("Expected error when giving empty file path")
	}
}

func TestNewTokenFileAuthenticate(t *testing.T) {
	tokenFile, err := os.CreateTemp("", "token_file")
	tokenFileContents := "super-secret-token"
	if err != nil {
		t.Fatal(err)
	}
	tokenFileName := tokenFile.Name()
	tokenFile.Close() // WriteFile doesn't need it open
	os.WriteFile(tokenFileName, []byte(tokenFileContents), 0o666)
	defer os.Remove(tokenFileName)

	logger := logging.NewVaultLogger(log.Trace)
	am, err := NewTokenFileAuthMethod(&auth.AuthConfig{
		Logger: logger.Named("auth.method"),
		Config: map[string]interface{}{
			"token_file_path":                 tokenFileName,
			"remove_token_file_after_reading": false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	path, headers, data, err := am.Authenticate(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if path != "auth/token/lookup" {
		t.Fatalf("Incorrect path, was %s", path)
	}
	if headers != nil {
		t.Fatalf("Expected no headers, instead got %v", headers)
	}
	if data == nil {
		t.Fatal("Data was nil")
	}
	tokenDataFromAuthMethod := data["token"].(string)
	if tokenDataFromAuthMethod != tokenFileContents {
		t.Fatalf("Incorrect token file contents return by auth method, expected %s, got %s", tokenFileContents, tokenDataFromAuthMethod)
	}

	_, err = os.Stat(tokenFileName)
	if err != nil {
		t.Fatal("Token file removed despite remove token file being set to false")
	}
}

func TestNewTokenFileAuthenticateRemoveAfterReading(t *testing.T) {
	tokenFile, err := os.CreateTemp("", "token_file")
	tokenFileContents := "super-secret-token"
	if err != nil {
		t.Fatal(err)
	}
	tokenFileName := tokenFile.Name()
	tokenFile.Close() // WriteFile doesn't need it open
	os.WriteFile(tokenFileName, []byte(tokenFileContents), 0o666)
	defer os.Remove(tokenFileName)

	logger := logging.NewVaultLogger(log.Trace)
	am, err := NewTokenFileAuthMethod(&auth.AuthConfig{
		Logger: logger.Named("auth.method"),
		Config: map[string]interface{}{
			"token_file_path":                 tokenFileName,
			"remove_token_file_after_reading": true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	path, headers, data, err := am.Authenticate(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if path != "auth/token/lookup" {
		t.Fatalf("Incorrect path, was %s", path)
	}
	if headers != nil {
		t.Fatalf("Expected no headers, instead got %v", headers)
	}
	if data == nil {
		t.Fatal("Data was nil")
	}
	tokenDataFromAuthMethod := data["token"].(string)
	if tokenDataFromAuthMethod != tokenFileContents {
		t.Fatalf("Incorrect token file contents return by auth method, expected %s, got %s", tokenFileContents, tokenDataFromAuthMethod)
	}

	_, err = os.Stat(tokenFileName)
	if err == nil {
		t.Fatal("no error returned from stat, indicating the file is still present")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("unexpected error: %v", err)
	}
}
