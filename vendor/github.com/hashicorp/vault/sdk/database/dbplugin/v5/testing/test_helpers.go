package dbtesting

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
)

func getRequestTimeout(t *testing.T) time.Duration {
	rawDur := os.Getenv("VAULT_TEST_DATABASE_REQUEST_TIMEOUT")
	if rawDur == "" {
		return 5 * time.Second
	}

	dur, err := time.ParseDuration(rawDur)
	if err != nil {
		t.Fatalf("Failed to parse custom request timeout %q: %s", rawDur, err)
	}
	return dur
}

func AssertInitialize(t *testing.T, db dbplugin.Database, req dbplugin.InitializeRequest) dbplugin.InitializeResponse {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout(t))
	defer cancel()

	resp, err := db.Initialize(ctx, req)
	if err != nil {
		t.Fatalf("Failed to initialize: %s", err)
	}
	return resp
}

func AssertNewUser(t *testing.T, db dbplugin.Database, req dbplugin.NewUserRequest) dbplugin.NewUserResponse {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout(t))
	defer cancel()

	resp, err := db.NewUser(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create new user: %s", err)
	}

	if resp.Username == "" {
		t.Fatalf("Missing username from NewUser response")
	}
	return resp
}

func AssertUpdateUser(t *testing.T, db dbplugin.Database, req dbplugin.UpdateUserRequest) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout(t))
	defer cancel()

	_, err := db.UpdateUser(ctx, req)
	if err != nil {
		t.Fatalf("Failed to update user: %s", err)
	}
}

func AssertDeleteUser(t *testing.T, db dbplugin.Database, req dbplugin.DeleteUserRequest) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout(t))
	defer cancel()

	_, err := db.DeleteUser(ctx, req)
	if err != nil {
		t.Fatalf("Failed to delete user %q: %s", req.Username, err)
	}
}

func AssertClose(t *testing.T, db dbplugin.Database) {
	t.Helper()
	err := db.Close()
	if err != nil {
		t.Fatalf("Failed to close database: %s", err)
	}
}
