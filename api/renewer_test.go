package api

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestRenewer_NewRenewer(t *testing.T) {
	t.Parallel()

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		i    *RenewerInput
		e    *Renewer
		err  bool
	}{
		{
			"nil",
			nil,
			nil,
			true,
		},
		{
			"missing_secret",
			&RenewerInput{
				Secret: nil,
			},
			nil,
			true,
		},
		{
			"default_grace",
			&RenewerInput{
				Secret: &Secret{},
			},
			&Renewer{
				secret: &Secret{},
				grace:  DefaultRenewerGrace,
			},
			false,
		},
		{
			"custom_grace",
			&RenewerInput{
				Secret: &Secret{},
				Grace:  30,
			},
			&Renewer{
				secret: &Secret{},
				grace:  30,
			},
			false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d_%s", i, tc.name), func(t *testing.T) {
			v, err := client.NewRenewer(tc.i)
			if (err != nil) != tc.err {
				t.Fatal(err)
			}

			if v == nil {
				return
			}

			// Zero-out channels because reflect
			v.client = nil
			v.doneCh = nil
			v.tickCh = nil
			v.stopCh = nil

			if !reflect.DeepEqual(tc.e, v) {
				t.Errorf("not equal\nexp: %#v\nact: %#v", tc.e, v)
			}
		})
	}
}

func TestRenewer_Renew(t *testing.T) {
	client, vaultDone := testVaultServer(t)
	defer vaultDone()

	pgURL, pgDone := testPostgresDatabase(t)
	defer pgDone()

	// Generic
	if _, err := client.Logical().Write("secret/value", map[string]interface{}{
		"foo": "bar",
	}); err != nil {
		t.Fatal(err)
	}

	// Transit
	if err := client.Sys().Mount("transit", &MountInput{
		Type: "transit",
	}); err != nil {
		t.Fatal(err)
	}

	// PostgreSQL
	if err := client.Sys().Mount("database", &MountInput{
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
		"default_ttl": "2s",
		"max_ttl":     "5s",
	}); err != nil {
		t.Fatal(err)
	}

	t.Run("generic", func(t *testing.T) {
		secret, err := client.Logical().Read("secret/value")
		if err != nil {
			t.Fatal(err)
		}

		v, err := client.NewRenewer(&RenewerInput{
			Secret: secret,
		})
		if err != nil {
			t.Fatal(err)
		}
		go v.Renew()
		defer v.Stop()

		select {
		case err := <-v.DoneCh():
			if err != ErrRenewerNotRenewable {
				t.Fatal(err)
			}
		}
	})

	t.Run("transit", func(t *testing.T) {
		secret, err := client.Logical().Write("transit/encrypt/my-app", map[string]interface{}{
			"plaintext": "Zm9vCg==",
		})
		if err != nil {
			t.Fatal(err)
		}

		v, err := client.NewRenewer(&RenewerInput{
			Secret: secret,
		})
		if err != nil {
			t.Fatal(err)
		}
		go v.Renew()
		defer v.Stop()

		select {
		case err := <-v.DoneCh():
			if err != ErrRenewerNotRenewable {
				t.Fatal(err)
			}
		}
	})

	t.Run("dynamic", func(t *testing.T) {
		secret, err := client.Logical().Read("database/creds/readonly")
		if err != nil {
			t.Fatal(err)
		}

		v, err := client.NewRenewer(&RenewerInput{
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
		case <-v.TickCh():
			// Received a renewal
		case <-time.After(5 * time.Second):
			t.Errorf("no data in 5s")
		}

		select {
		case err := <-v.DoneCh():
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(5 * time.Second):
			t.Errorf("no data in 5s")
		}
	})

	t.Run("auth", func(t *testing.T) {
		secret, err := client.Auth().Token().Create(&TokenCreateRequest{
			Policies:       []string{"default"},
			TTL:            "2s",
			ExplicitMaxTTL: "5s",
		})
		if err != nil {
			t.Fatal(err)
		}

		v, err := client.NewRenewer(&RenewerInput{
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
		case <-v.TickCh():
			// Received a renewal
		case <-time.After(5 * time.Second):
			t.Errorf("no data in 5s")
		}

		select {
		case err := <-v.DoneCh():
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(5 * time.Second):
			t.Errorf("no data in 5s")
		}
	})
}
