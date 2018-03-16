package api_test

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

func TestRenewer_Renew(t *testing.T) {
	t.Parallel()

	client, vaultDone := testVaultServer(t)
	defer vaultDone()

	pgURL, pgDone := testPostgresDB(t)
	defer pgDone()

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

			v, err := client.NewRenewer(&api.RenewerInput{
				Secret: secret,
			})
			if err != nil {
				t.Fatal(err)
			}
			go v.Renew()
			defer v.Stop()

			select {
			case err := <-v.DoneCh():
				if err != api.ErrRenewerNotRenewable {
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

			v, err := client.NewRenewer(&api.RenewerInput{
				Secret: secret,
			})
			if err != nil {
				t.Fatal(err)
			}
			go v.Renew()
			defer v.Stop()

			select {
			case err := <-v.DoneCh():
				if err != api.ErrRenewerNotRenewable {
					t.Fatal(err)
				}
			case renew := <-v.RenewCh():
				t.Errorf("received renew, but should have been nil: %#v", renew)
			case <-time.After(500 * time.Millisecond):
				t.Error("should have been non-renewable")
			}
		})

		t.Run("database", func(t *testing.T) {
			t.Parallel()

			if err := client.Sys().Mount("database", &api.MountInput{
				Type: "database",
			}); err != nil {
				t.Fatal(err)
			}
			if _, err := client.Logical().Write("database/config/postgresql", map[string]interface{}{
				"plugin_name":    "postgresql-database-plugin",
				"connection_url": pgURL,
				"allowed_roles":  "readonly",
			}); err != nil {
				t.Fatal(err)
			}
			if _, err := client.Logical().Write("database/roles/readonly", map[string]interface{}{
				"db_name": "postgresql",
				"creation_statements": `` +
					`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';` +
					`GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";`,
				"default_ttl": "5s",
				"max_ttl":     "10s",
			}); err != nil {
				t.Fatal(err)
			}

			secret, err := client.Logical().Read("database/creds/readonly")
			if err != nil {
				t.Fatal(err)
			}

			v, err := client.NewRenewer(&api.RenewerInput{
				Secret: secret,
			})
			if err != nil {
				t.Fatal(err)
			}
			go v.Renew()
			defer v.Stop()

			select {
			case err := <-v.DoneCh():
				t.Errorf("should have renewed once before returning: %s", err)
			case renew := <-v.RenewCh():
				if renew == nil {
					t.Fatal("renew is nil")
				}
				if !renew.Secret.Renewable {
					t.Errorf("expected lease to be renewable: %#v", renew)
				}
				if renew.Secret.LeaseDuration > 5 {
					t.Errorf("expected lease to <= 5s: %#v", renew)
				}
			case <-time.After(5 * time.Second):
				t.Errorf("no renewal")
			}

		outer:
			for {
				select {
				case err := <-v.DoneCh():
					if err != nil {
						t.Fatal(err)
					}
					break outer
				case renew := <-v.RenewCh():
					t.Logf("renew called, remaining lease duration: %d", renew.Secret.LeaseDuration)
					continue outer
				case <-time.After(5 * time.Second):
					t.Errorf("no data")
					break outer
				}
			}
		})

		t.Run("auth", func(t *testing.T) {
			t.Parallel()

			secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
				Policies:       []string{"default"},
				TTL:            "1s",
				ExplicitMaxTTL: "3s",
			})
			if err != nil {
				t.Fatal(err)
			}

			v, err := client.NewRenewer(&api.RenewerInput{
				Secret: secret,
			})
			if err != nil {
				t.Fatal(err)
			}
			go v.Renew()
			defer v.Stop()

			select {
			case err := <-v.DoneCh():
				t.Errorf("should have renewed once before returning: %s", err)
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
				if renew.Secret.Auth.LeaseDuration > 2 {
					t.Errorf("expected lease to < 2s: %#v", renew)
				}
				if renew.Secret.Auth.ClientToken == "" {
					t.Error("expected a client token")
				}
				if renew.Secret.Auth.Accessor == "" {
					t.Error("expected an accessor")
				}
			case <-time.After(3 * time.Second):
				t.Errorf("no renewal")
			}

		outer:
			for {
				select {
				case err := <-v.DoneCh():
					if err != nil {
						t.Fatal(err)
					}
					break outer
				case <-v.RenewCh():
					continue outer
				case <-time.After(3 * time.Second):
					t.Errorf("no data")
					break outer
				}
			}
		})
	})
}
