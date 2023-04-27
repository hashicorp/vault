// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

func TestRenewer_Renew(t *testing.T) {
	t.Parallel()

	client, vaultDone := testVaultServer(t)
	defer vaultDone()

	t.Run("group", func(t *testing.T) {
		t.Run("kv", func(t *testing.T) {
			t.Parallel()

			if _, err := client.Logical().Write("secret/value", map[string]interface{}{
				"foo": "bar",
			}); err != nil {
				t.Fatal(err)
			}

			secret, err := client.Logical().Read("secret/value")
			if err != nil {
				t.Fatal(err)
			}

			v, err := client.NewLifetimeWatcher(&api.RenewerInput{
				Secret: secret,
			})
			if err != nil {
				t.Fatal(err)
			}
			go v.Renew()
			defer v.Stop()

			select {
			case err := <-v.DoneCh():
				if err != api.ErrLifetimeWatcherNotRenewable {
					t.Fatal(err)
				}
			case renew := <-v.RenewCh():
				t.Errorf("received renew, but should have been nil: %#v", renew)
			case <-time.After(500 * time.Millisecond):
				t.Error("should have been non-renewable")
			}
		})

		t.Run("transit", func(t *testing.T) {
			t.Parallel()

			if err := client.Sys().Mount("transit", &api.MountInput{
				Type: "transit",
			}); err != nil {
				t.Fatal(err)
			}

			secret, err := client.Logical().Write("transit/encrypt/my-app", map[string]interface{}{
				"plaintext": "Zm9vCg==",
			})
			if err != nil {
				t.Fatal(err)
			}

			v, err := client.NewLifetimeWatcher(&api.RenewerInput{
				Secret: secret,
			})
			if err != nil {
				t.Fatal(err)
			}
			go v.Renew()
			defer v.Stop()

			select {
			case err := <-v.DoneCh():
				if err != api.ErrLifetimeWatcherNotRenewable {
					t.Fatal(err)
				}
			case renew := <-v.RenewCh():
				t.Errorf("received renew, but should have been nil: %#v", renew)
			case <-time.After(500 * time.Millisecond):
				t.Error("should have been non-renewable")
			}
		})

		t.Run("auth", func(t *testing.T) {
			t.Parallel()

			secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
				Policies:       []string{"default"},
				TTL:            "5s",
				ExplicitMaxTTL: "10s",
			})
			if err != nil {
				t.Fatal(err)
			}

			v, err := client.NewLifetimeWatcher(&api.RenewerInput{
				Secret: secret,
			})
			if err != nil {
				t.Fatal(err)
			}
			go v.Renew()
			defer v.Stop()

			renewed, done := false, false
			timeout := time.After(10 * time.Second)
			for {
				if done {
					break
				}
				select {
				case err := <-v.DoneCh():
					if renewed {
						// If we renewed but there's an error, we fail
						if err != nil {
							t.Fatalf("renewal failed with an error: %v", err)
						}
						// We can break out early here
						done = true
					} else {
						t.Errorf("should have renewed once before returning: %s", err)
					}
					done = true
				case renew := <-v.RenewCh():
					if renew == nil {
						t.Fatal("renew is nil")
					}
					if renew.Secret.Auth == nil {
						t.Fatal("renew auth is nil")
					}
					if !renew.Secret.Auth.Renewable {
						t.Errorf("expected lease to be renewable: %#v", renew)
					}
					if renew.Secret.Auth.LeaseDuration > 5 {
						t.Errorf("expected lease to <= 5s: %#v", renew)
					}
					if renew.Secret.Auth.ClientToken == "" {
						t.Error("expected a client token")
					}
					if renew.Secret.Auth.Accessor == "" {
						t.Error("expected an accessor")
					}
					renewed = true
				case <-timeout:
					if !renewed {
						t.Errorf("no renewal")
					}
					done = true
				}
			}
		})
	})
}
