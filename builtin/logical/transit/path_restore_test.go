package transit

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/logical"
)

func TestTransit_Restore(t *testing.T) {
	// Test setup:
	// - Create a key
	// - Configure it to be exportable, allowing deletion, and backups
	// - Capture backup
	// - Delete key
	// - Run test cases
	//
	// Each test case should start with no key present. If the 'Seed' parameter is
	// in the struct, we'll start by restoring it (without force) to run that test
	// as if the key already existed

	keyType := "aes256-gcm96"
	b, s := createBackendWithStorage(t)
	keyName := testhelpers.RandomWithPrefix("my-key")

	// Create a key
	keyReq := &logical.Request{
		Path:      "keys/" + keyName,
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"type":       keyType,
			"exportable": true,
		},
	}
	resp, err := b.HandleRequest(context.Background(), keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	// Configure the key to allow its deletion and backup
	configReq := &logical.Request{
		Path:      fmt.Sprintf("keys/%s/config", keyName),
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"deletion_allowed":       true,
			"allow_plaintext_backup": true,
		},
	}
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	// Take a backup of the key
	backupReq := &logical.Request{
		Path:      "backup/" + keyName,
		Operation: logical.ReadOperation,
		Storage:   s,
	}
	resp, err = b.HandleRequest(context.Background(), backupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	backupKey := resp.Data["backup"].(string)
	if backupKey == "" {
		t.Fatal("failed to get a backup")
	}

	// Delete the key to start test cases with clean slate
	keyReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	// helper func to get a pointer value for a boolean
	boolPtr := func(b bool) *bool {
		return &b
	}

	keyExitsError := fmt.Errorf("key \"%s\" already exists", keyName)

	testCases := []struct {
		Name string
		// Seed dermines if we start the test by restoring the initial backup we
		// took, to test a restore operation based on the key existing or not
		Seed bool
		// Force is a pointer to differentiate between default false and given false
		Force *bool
		// The error we expect, if any
		ExpectedErr error
	}{
		{
			// key does not already exist
			Name: "Default restore",
		},
		{
			// key already exists
			Name:        "Restore-without-force",
			Seed:        true,
			ExpectedErr: keyExitsError,
		},
		{
			// key already exists, use force to force a restore
			Name:  "Restore-with-force",
			Seed:  true,
			Force: boolPtr(true),
		},
		{
			// using force shouldn't matter if the key doesn't exist
			Name:  "Restore-with-force-no-seed",
			Force: boolPtr(true),
		},
		{
			// using force shouldn't matter if the key doesn't exist
			Name:  "Restore-force-false",
			Force: boolPtr(false),
		},
		{
			// using false force should still error
			Name:        "Restore-force-false",
			Seed:        true,
			Force:       boolPtr(false),
			ExpectedErr: keyExitsError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var resp *logical.Response
			var err error
			if tc.Seed {
				// restore our key to test a pre-existing key
				seedRestoreReq := &logical.Request{
					Path:      "restore",
					Operation: logical.UpdateOperation,
					Storage:   s,
					Data: map[string]interface{}{
						"backup": backupKey,
					},
				}

				resp, err := b.HandleRequest(context.Background(), seedRestoreReq)
				if resp != nil && resp.IsError() {
					t.Fatalf("resp: %#v\nerr: %v", resp, err)
				}
				if err != nil && tc.ExpectedErr == nil {
					t.Fatalf("did not expect an error in SeedKey restore: %s", err)
				}
			}

			restoreReq := &logical.Request{
				Path:      "restore",
				Operation: logical.UpdateOperation,
				Storage:   s,
				Data: map[string]interface{}{
					"backup": backupKey,
				},
			}

			if tc.Force != nil {
				restoreReq.Data["force"] = *tc.Force
			}

			resp, err = b.HandleRequest(context.Background(), restoreReq)
			if resp != nil && resp.IsError() {
				t.Fatalf("resp: %#v\nerr: %v", resp, err)
			}
			if err == nil && tc.ExpectedErr != nil {
				t.Fatalf("expected an error, but got none")
			}
			if err != nil && tc.ExpectedErr == nil {
				t.Fatalf("unexpected error:%s", err)
			}

			if err != nil && tc.ExpectedErr != nil {
				if err.Error() != tc.ExpectedErr.Error() {
					t.Fatalf("expected error: (%s), got: (%s)", tc.ExpectedErr.Error(), err.Error())
				}
			}

			// cleanup / delete key after each run
			keyReq.Operation = logical.DeleteOperation
			resp, err = b.HandleRequest(context.Background(), keyReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("resp: %#v\nerr: %v", resp, err)
			}
		})
	}
}
