package api_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

func TestParseSecret(t *testing.T) {
	raw := strings.TrimSpace(`
{
	"lease_id": "foo",
	"renewable": true,
	"lease_duration": 10,
	"data": {
		"key": "value"
	},
	"warnings": [
		"a warning!"
	],
	"wrap_info": {
		"token": "token",
		"ttl": 60,
		"creation_time": "2016-06-07T15:52:10-04:00",
		"wrapped_accessor": "abcd1234"
	}
}`)

	rawTime, _ := time.Parse(time.RFC3339, "2016-06-07T15:52:10-04:00")

	secret, err := api.ParseSecret(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &api.Secret{
		LeaseID:       "foo",
		Renewable:     true,
		LeaseDuration: 10,
		Data: map[string]interface{}{
			"key": "value",
		},
		Warnings: []string{
			"a warning!",
		},
		WrapInfo: &api.SecretWrapInfo{
			Token:           "token",
			TTL:             60,
			CreationTime:    rawTime,
			WrappedAccessor: "abcd1234",
		},
	}
	if !reflect.DeepEqual(secret, expected) {
		t.Fatalf("bad:\ngot\n%#v\nexpected\n%#v\n", secret, expected)
	}
}

func TestSecret_TokenID(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		secret *api.Secret
		exp    string
	}{
		{
			"nil",
			nil,
			"",
		},
		{
			"nil_auth",
			&api.Secret{
				Auth: nil,
			},
			"",
		},
		{
			"empty_auth_client_token",
			&api.Secret{
				Auth: &api.SecretAuth{
					ClientToken: "",
				},
			},
			"",
		},
		{
			"real_auth_client_token",
			&api.Secret{
				Auth: &api.SecretAuth{
					ClientToken: "my-token",
				},
			},
			"my-token",
		},
		{
			"nil_data",
			&api.Secret{
				Data: nil,
			},
			"",
		},
		{
			"empty_data",
			&api.Secret{
				Data: map[string]interface{}{},
			},
			"",
		},
		{
			"data_not_string",
			&api.Secret{
				Data: map[string]interface{}{
					"id": 123,
				},
			},
			"",
		},
		{
			"data_string",
			&api.Secret{
				Data: map[string]interface{}{
					"id": "my-token",
				},
			},
			"my-token",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			act := tc.secret.TokenID()
			if act != tc.exp {
				t.Errorf("expected %q to be %q", act, tc.exp)
			}
		})
	}

	t.Run("auth", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		secret, err := client.Logical().Write("auth/userpass/login/test", map[string]interface{}{
			"password": "test",
		})
		if err != nil || secret == nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		if secret.TokenID() != token {
			t.Errorf("expected %q to be %q", secret.TokenID(), token)
		}
	})

	t.Run("token-create", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		if secret.TokenID() != token {
			t.Errorf("expected %q to be %q", secret.TokenID(), token)
		}
	})

	t.Run("token-lookup", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Lookup(token)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenID() != token {
			t.Errorf("expected %q to be %q", secret.TokenID(), token)
		}
	})

	t.Run("token-lookup-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenID() != token {
			t.Errorf("expected %q to be %q", secret.TokenID(), token)
		}
	})

	t.Run("token-renew", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Renew(token, 0)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenID() != token {
			t.Errorf("expected %q to be %q", secret.TokenID(), token)
		}
	})

	t.Run("token-renew-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().RenewSelf(0)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenID() != token {
			t.Errorf("expected %q to be %q", secret.TokenID(), token)
		}
	})
}

func TestSecret_TokenAccessor(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		secret *api.Secret
		exp    string
	}{
		{
			"nil",
			nil,
			"",
		},
		{
			"nil_auth",
			&api.Secret{
				Auth: nil,
			},
			"",
		},
		{
			"empty_auth_accessor",
			&api.Secret{
				Auth: &api.SecretAuth{
					Accessor: "",
				},
			},
			"",
		},
		{
			"real_auth_accessor",
			&api.Secret{
				Auth: &api.SecretAuth{
					Accessor: "my-accessor",
				},
			},
			"my-accessor",
		},
		{
			"nil_data",
			&api.Secret{
				Data: nil,
			},
			"",
		},
		{
			"empty_data",
			&api.Secret{
				Data: map[string]interface{}{},
			},
			"",
		},
		{
			"data_not_string",
			&api.Secret{
				Data: map[string]interface{}{
					"accessor": 123,
				},
			},
			"",
		},
		{
			"data_string",
			&api.Secret{
				Data: map[string]interface{}{
					"accessor": "my-accessor",
				},
			},
			"my-accessor",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			act := tc.secret.TokenAccessor()
			if act != tc.exp {
				t.Errorf("expected %q to be %q", act, tc.exp)
			}
		})
	}

	t.Run("auth", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		secret, err := client.Logical().Write("auth/userpass/login/test", map[string]interface{}{
			"password": "test",
		})
		if err != nil || secret == nil {
			t.Fatal(err)
		}
		_, accessor := secret.Auth.ClientToken, secret.Auth.Accessor

		if secret.TokenAccessor() != accessor {
			t.Errorf("expected %q to be %q", secret.TokenAccessor(), accessor)
		}
	})

	t.Run("token-create", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		_, accessor := secret.Auth.ClientToken, secret.Auth.Accessor

		if secret.TokenAccessor() != accessor {
			t.Errorf("expected %q to be %q", secret.TokenAccessor(), accessor)
		}
	})

	t.Run("token-lookup", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		token, accessor := secret.Auth.ClientToken, secret.Auth.Accessor

		secret, err = client.Auth().Token().Lookup(token)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenAccessor() != accessor {
			t.Errorf("expected %q to be %q", secret.TokenAccessor(), accessor)
		}
	})

	t.Run("token-lookup-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		token, accessor := secret.Auth.ClientToken, secret.Auth.Accessor

		client.SetToken(token)
		secret, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenAccessor() != accessor {
			t.Errorf("expected %q to be %q", secret.TokenAccessor(), accessor)
		}
	})

	t.Run("token-renew", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		token, accessor := secret.Auth.ClientToken, secret.Auth.Accessor

		secret, err = client.Auth().Token().Renew(token, 0)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenAccessor() != accessor {
			t.Errorf("expected %q to be %q", secret.TokenAccessor(), accessor)
		}
	})

	t.Run("token-renew-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		token, accessor := secret.Auth.ClientToken, secret.Auth.Accessor

		client.SetToken(token)
		secret, err = client.Auth().Token().RenewSelf(0)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenAccessor() != accessor {
			t.Errorf("expected %q to be %q", secret.TokenAccessor(), accessor)
		}
	})
}

func TestSecret_TokenMeta(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		secret *api.Secret
		exp    map[string]string
	}{
		{
			"nil",
			nil,
			nil,
		},
		{
			"nil_auth",
			&api.Secret{
				Auth: nil,
			},
			nil,
		},
		{
			"nil_auth_metadata",
			&api.Secret{
				Auth: &api.SecretAuth{
					Metadata: nil,
				},
			},
			nil,
		},
		{
			"empty_auth_metadata",
			&api.Secret{
				Auth: &api.SecretAuth{
					Metadata: map[string]string{},
				},
			},
			nil,
		},
		{
			"real_auth_metadata",
			&api.Secret{
				Auth: &api.SecretAuth{
					Metadata: map[string]string{"foo": "bar"},
				},
			},
			map[string]string{"foo": "bar"},
		},
		{
			"nil_data",
			&api.Secret{
				Data: nil,
			},
			nil,
		},
		{
			"empty_data",
			&api.Secret{
				Data: map[string]interface{}{},
			},
			nil,
		},
		{
			"data_not_map",
			&api.Secret{
				Data: map[string]interface{}{
					"meta": 123,
				},
			},
			nil,
		},
		{
			"data_map",
			&api.Secret{
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			map[string]string{"foo": "bar"},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			act := tc.secret.TokenMeta()
			if !reflect.DeepEqual(act, tc.exp) {
				t.Errorf("expected %#v to be %#v", act, tc.exp)
			}
		})
	}

	t.Run("auth", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		secret, err := client.Logical().Write("auth/userpass/login/test", map[string]interface{}{
			"password": "test",
		})
		if err != nil || secret == nil {
			t.Fatal(err)
		}

		meta := map[string]string{"username": "test"}
		if !reflect.DeepEqual(secret.TokenMeta(), meta) {
			t.Errorf("expected %#v to be %#v", secret.TokenMeta(), meta)
		}
	})

	t.Run("token-create", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		meta := map[string]string{"foo": "bar"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Metadata: meta,
		})
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenMeta(), meta) {
			t.Errorf("expected %#v to be %#v", secret.TokenMeta(), meta)
		}
	})

	t.Run("token-lookup", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		meta := map[string]string{"foo": "bar"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Metadata: meta,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Lookup(token)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenMeta(), meta) {
			t.Errorf("expected %#v to be %#v", secret.TokenMeta(), meta)
		}
	})

	t.Run("token-lookup-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		meta := map[string]string{"foo": "bar"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Metadata: meta,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenMeta(), meta) {
			t.Errorf("expected %#v to be %#v", secret.TokenMeta(), meta)
		}
	})

	t.Run("token-renew", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		meta := map[string]string{"foo": "bar"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Metadata: meta,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Renew(token, 0)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenMeta(), meta) {
			t.Errorf("expected %#v to be %#v", secret.TokenMeta(), meta)
		}
	})

	t.Run("token-renew-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		meta := map[string]string{"foo": "bar"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Metadata: meta,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().RenewSelf(0)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenMeta(), meta) {
			t.Errorf("expected %#v to be %#v", secret.TokenMeta(), meta)
		}
	})
}

func TestSecret_TokenRemainingUses(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		secret *api.Secret
		exp    int
	}{
		{
			"nil",
			nil,
			0,
		},
		{
			"nil_data",
			&api.Secret{
				Data: nil,
			},
			0,
		},
		{
			"empty_data",
			&api.Secret{
				Data: map[string]interface{}{},
			},
			0,
		},
		{
			"data_not_json_number",
			&api.Secret{
				Data: map[string]interface{}{
					"num_uses": 123,
				},
			},
			0,
		},
		{
			"data_json_number",
			&api.Secret{
				Data: map[string]interface{}{
					"num_uses": json.Number("123"),
				},
			},
			123,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			act := tc.secret.TokenRemainingUses()
			if act != tc.exp {
				t.Errorf("expected %d to be %d", act, tc.exp)
			}
		})
	}

	t.Run("auth", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		uses := 5

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
			"num_uses": uses,
		}); err != nil {
			t.Fatal(err)
		}

		secret, err := client.Logical().Write("auth/userpass/login/test", map[string]interface{}{
			"password": "test",
		})
		if err != nil || secret == nil {
			t.Fatal(err)
		}

		// Remaining uses is not returned from this API
		uses = 0
		if secret.TokenRemainingUses() != uses {
			t.Errorf("expected %d to be %d", secret.TokenRemainingUses(), uses)
		}
	})

	t.Run("token-create", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		uses := 5

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			NumUses:  uses,
		})
		if err != nil {
			t.Fatal(err)
		}

		// /auth/token/create does not return the number of uses
		uses = 0
		if secret.TokenRemainingUses() != uses {
			t.Errorf("expected %d to be %d", secret.TokenRemainingUses(), uses)
		}
	})

	t.Run("token-lookup", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		uses := 5

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			NumUses:  uses,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Lookup(token)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenRemainingUses() != uses {
			t.Errorf("expected %d to be %d", secret.TokenRemainingUses(), uses)
		}
	})

	t.Run("token-lookup-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		uses := 5

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			NumUses:  uses,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}

		uses = uses - 1 // we just used it
		if secret.TokenRemainingUses() != uses {
			t.Errorf("expected %d to be %d", secret.TokenRemainingUses(), uses)
		}
	})

	t.Run("token-renew", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		uses := 5

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			NumUses:  uses,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Renew(token, 0)
		if err != nil {
			t.Fatal(err)
		}

		// /auth/token/renew does not return the number of uses
		uses = 0
		if secret.TokenRemainingUses() != uses {
			t.Errorf("expected %d to be %d", secret.TokenRemainingUses(), uses)
		}
	})

	t.Run("token-renew-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		uses := 5

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			NumUses:  uses,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().RenewSelf(0)
		if err != nil {
			t.Fatal(err)
		}

		// /auth/token/renew-self does not return the number of uses
		uses = 0
		if secret.TokenRemainingUses() != uses {
			t.Errorf("expected %d to be %d", secret.TokenRemainingUses(), uses)
		}
	})
}

func TestSecret_TokenPolicies(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		secret *api.Secret
		exp    []string
	}{
		{
			"nil",
			nil,
			nil,
		},
		{
			"nil_auth",
			&api.Secret{
				Auth: nil,
			},
			nil,
		},
		{
			"nil_auth_policies",
			&api.Secret{
				Auth: &api.SecretAuth{
					Policies: nil,
				},
			},
			nil,
		},
		{
			"empty_auth_policies",
			&api.Secret{
				Auth: &api.SecretAuth{
					Policies: []string{},
				},
			},
			nil,
		},
		{
			"real_auth_policies",
			&api.Secret{
				Auth: &api.SecretAuth{
					Policies: []string{"foo"},
				},
			},
			[]string{"foo"},
		},
		{
			"nil_data",
			&api.Secret{
				Data: nil,
			},
			nil,
		},
		{
			"empty_data",
			&api.Secret{
				Data: map[string]interface{}{},
			},
			nil,
		},
		{
			"data_not_slice",
			&api.Secret{
				Data: map[string]interface{}{
					"policies": 123,
				},
			},
			nil,
		},
		{
			"data_slice",
			&api.Secret{
				Data: map[string]interface{}{
					"policies": []interface{}{"foo"},
				},
			},
			[]string{"foo"},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			act := tc.secret.TokenPolicies()
			if !reflect.DeepEqual(act, tc.exp) {
				t.Errorf("expected %#v to be %#v", act, tc.exp)
			}
		})
	}

	t.Run("auth", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		policies := []string{"bar", "default", "foo"}

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": strings.Join(policies, ","),
		}); err != nil {
			t.Fatal(err)
		}

		secret, err := client.Logical().Write("auth/userpass/login/test", map[string]interface{}{
			"password": "test",
		})
		if err != nil || secret == nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenPolicies(), policies) {
			t.Errorf("expected %#v to be %#v", secret.TokenPolicies(), policies)
		}
	})

	t.Run("token-create", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		policies := []string{"bar", "default", "foo"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: policies,
		})
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenPolicies(), policies) {
			t.Errorf("expected %#v to be %#v", secret.TokenPolicies(), policies)
		}
	})

	t.Run("token-lookup", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		policies := []string{"bar", "default", "foo"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: policies,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Lookup(token)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenPolicies(), policies) {
			t.Errorf("expected %#v to be %#v", secret.TokenPolicies(), policies)
		}
	})

	t.Run("token-lookup-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		policies := []string{"bar", "default", "foo"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: policies,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenPolicies(), policies) {
			t.Errorf("expected %#v to be %#v", secret.TokenPolicies(), policies)
		}
	})

	t.Run("token-renew", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		policies := []string{"bar", "default", "foo"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: policies,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Renew(token, 0)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenPolicies(), policies) {
			t.Errorf("expected %#v to be %#v", secret.TokenPolicies(), policies)
		}
	})

	t.Run("token-renew-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		policies := []string{"bar", "default", "foo"}

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: policies,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().RenewSelf(0)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(secret.TokenPolicies(), policies) {
			t.Errorf("expected %#v to be %#v", secret.TokenPolicies(), policies)
		}
	})
}

func TestSecret_TokenIsRenewable(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		secret *api.Secret
		exp    bool
	}{
		{
			"nil",
			nil,
			false,
		},
		{
			"nil_auth",
			&api.Secret{
				Auth: nil,
			},
			false,
		},
		{
			"auth_renewable_false",
			&api.Secret{
				Auth: &api.SecretAuth{
					Renewable: false,
				},
			},
			false,
		},
		{
			"auth_renewable_true",
			&api.Secret{
				Auth: &api.SecretAuth{
					Renewable: true,
				},
			},
			true,
		},
		{
			"nil_data",
			&api.Secret{
				Data: nil,
			},
			false,
		},
		{
			"empty_data",
			&api.Secret{
				Data: map[string]interface{}{},
			},
			false,
		},
		{
			"data_not_bool",
			&api.Secret{
				Data: map[string]interface{}{
					"renewable": 123,
				},
			},
			false,
		},
		{
			"data_bool_true",
			&api.Secret{
				Data: map[string]interface{}{
					"renewable": true,
				},
			},
			true,
		},
		{
			"data_bool_false",
			&api.Secret{
				Data: map[string]interface{}{
					"renewable": true,
				},
			},
			true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			act := tc.secret.TokenIsRenewable()
			if act != tc.exp {
				t.Errorf("expected %t to be %t", act, tc.exp)
			}
		})
	}

	t.Run("auth", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		renewable := true

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		secret, err := client.Logical().Write("auth/userpass/login/test", map[string]interface{}{
			"password": "test",
		})
		if err != nil || secret == nil {
			t.Fatal(err)
		}

		if secret.TokenIsRenewable() != renewable {
			t.Errorf("expected %t to be %t", secret.TokenIsRenewable(), renewable)
		}
	})

	t.Run("token-create", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		renewable := true

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:  []string{"default"},
			Renewable: &renewable,
		})
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenIsRenewable() != renewable {
			t.Errorf("expected %t to be %t", secret.TokenIsRenewable(), renewable)
		}
	})

	t.Run("token-lookup", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		renewable := true

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:  []string{"default"},
			Renewable: &renewable,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Lookup(token)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenIsRenewable() != renewable {
			t.Errorf("expected %t to be %t", secret.TokenIsRenewable(), renewable)
		}
	})

	t.Run("token-lookup-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		renewable := true

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:  []string{"default"},
			Renewable: &renewable,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenIsRenewable() != renewable {
			t.Errorf("expected %t to be %t", secret.TokenIsRenewable(), renewable)
		}
	})

	t.Run("token-renew", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		renewable := true

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:  []string{"default"},
			Renewable: &renewable,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Renew(token, 0)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenIsRenewable() != renewable {
			t.Errorf("expected %t to be %t", secret.TokenIsRenewable(), renewable)
		}
	})

	t.Run("token-renew-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		renewable := true

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:  []string{"default"},
			Renewable: &renewable,
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().RenewSelf(0)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenIsRenewable() != renewable {
			t.Errorf("expected %t to be %t", secret.TokenIsRenewable(), renewable)
		}
	})
}

func TestSecret_TokenTTL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		secret *api.Secret
		exp    time.Duration
	}{
		{
			"nil",
			nil,
			0,
		},
		{
			"nil_auth",
			&api.Secret{
				Auth: nil,
			},
			0,
		},
		{
			"nil_auth_lease_duration",
			&api.Secret{
				Auth: &api.SecretAuth{
					LeaseDuration: 0,
				},
			},
			0,
		},
		{
			"real_auth_lease_duration",
			&api.Secret{
				Auth: &api.SecretAuth{
					LeaseDuration: 3600,
				},
			},
			1 * time.Hour,
		},
		{
			"nil_data",
			&api.Secret{
				Data: nil,
			},
			0,
		},
		{
			"empty_data",
			&api.Secret{
				Data: map[string]interface{}{},
			},
			0,
		},
		{
			"data_not_json_number",
			&api.Secret{
				Data: map[string]interface{}{
					"ttl": 123,
				},
			},
			0,
		},
		{
			"data_json_number",
			&api.Secret{
				Data: map[string]interface{}{
					"ttl": json.Number("3600"),
				},
			},
			1 * time.Hour,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			act := tc.secret.TokenTTL()
			if act != tc.exp {
				t.Errorf("expected %q to be %q", act, tc.exp)
			}
		})
	}

	t.Run("auth", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ttl := 30 * time.Minute

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password":         "test",
			"policies":         "default",
			"ttl":              ttl.String(),
			"explicit_max_ttl": ttl.String(),
		}); err != nil {
			t.Fatal(err)
		}

		secret, err := client.Logical().Write("auth/userpass/login/test", map[string]interface{}{
			"password": "test",
		})
		if err != nil || secret == nil {
			t.Fatal(err)
		}

		if secret.TokenTTL() == 0 || secret.TokenTTL() > ttl {
			t.Errorf("expected %q to non-zero and less than %q", secret.TokenTTL(), ttl)
		}
	})

	t.Run("token-create", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ttl := 30 * time.Minute

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:       []string{"default"},
			TTL:            ttl.String(),
			ExplicitMaxTTL: ttl.String(),
		})
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenTTL() == 0 || secret.TokenTTL() > ttl {
			t.Errorf("expected %q to non-zero and less than %q", secret.TokenTTL(), ttl)
		}
	})

	t.Run("token-lookup", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ttl := 30 * time.Minute

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:       []string{"default"},
			TTL:            ttl.String(),
			ExplicitMaxTTL: ttl.String(),
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Lookup(token)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenTTL() == 0 || secret.TokenTTL() > ttl {
			t.Errorf("expected %q to non-zero and less than %q", secret.TokenTTL(), ttl)
		}
	})

	t.Run("token-lookup-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ttl := 30 * time.Minute

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:       []string{"default"},
			TTL:            ttl.String(),
			ExplicitMaxTTL: ttl.String(),
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenTTL() == 0 || secret.TokenTTL() > ttl {
			t.Errorf("expected %q to non-zero and less than %q", secret.TokenTTL(), ttl)
		}
	})

	t.Run("token-renew", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ttl := 30 * time.Minute

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:       []string{"default"},
			TTL:            ttl.String(),
			ExplicitMaxTTL: ttl.String(),
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		secret, err = client.Auth().Token().Renew(token, 0)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenTTL() == 0 || secret.TokenTTL() > ttl {
			t.Errorf("expected %q to non-zero and less than %q", secret.TokenTTL(), ttl)
		}
	})

	t.Run("token-renew-self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ttl := 30 * time.Minute

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies:       []string{"default"},
			TTL:            ttl.String(),
			ExplicitMaxTTL: ttl.String(),
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		client.SetToken(token)
		secret, err = client.Auth().Token().RenewSelf(0)
		if err != nil {
			t.Fatal(err)
		}

		if secret.TokenTTL() == 0 || secret.TokenTTL() > ttl {
			t.Errorf("expected %q to non-zero and less than %q", secret.TokenTTL(), ttl)
		}
	})
}
