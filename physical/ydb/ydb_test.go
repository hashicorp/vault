package ydb

import (
	"fmt"
	"testing"

	log "github.com/hashicorp/go-hclog"
	helper "github.com/hashicorp/vault/helper/testhelpers/ydb"
	ydbconsts "github.com/hashicorp/vault/physical/ydb/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestYDBBackend(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	cleanup, cfg := helper.PrepareTestContainer(t)
	defer cleanup()

	logger.Info(fmt.Sprintf("YDB DSN: %v", cfg.DSN))
	logger.Info(fmt.Sprintf("YDB VAULT TABLE: %v", cfg.Table))

	backend, err := NewYDBBackend(map[string]string{
		"dsn":        cfg.DSN,
		"table":      cfg.Table,
		"balancer":   "single",
		"ha_enabled": "true",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}
	backend2, err := NewYDBBackend(map[string]string{
		"dsn":        cfg.DSN,
		"table":      cfg.Table,
		"balancer":   "single",
		"ha_enabled": "true",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create second backend: %v", err)
	}

	logger.Info("Running basic backend tests")
	physical.ExerciseBackend(t, backend)

	logger.Info("Running list prefix backend tests")
	physical.ExerciseBackend_ListPrefix(t, backend)

	logger.Info("Running transactional backend tests")
	physical.ExerciseTransactionalBackend(t, backend)

	logger.Info("Running ha backend tests")
	physical.ExerciseHABackend(t, backend.(physical.HABackend), backend2.(physical.HABackend))
}

func TestQuoteYDBIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		expectErr bool
	}{
		{
			name:  "simple table",
			input: "vault_kv",
			want:  "`vault_kv`",
		},
		{
			name:  "absolute path",
			input: "/local/vault_kv",
			want:  "/`local`/`vault_kv`",
		},
		{
			name:  "escapes backticks",
			input: "vault`kv",
			want:  "`vault``kv`",
		},
		{
			name:      "rejects reserved segment",
			input:     "/local/../vault_kv",
			expectErr: true,
		},
		{
			name:      "rejects empty middle segment",
			input:     "local//vault_kv",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := quoteYDBIdentifier(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error for %q", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("quoteYDBIdentifier(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestResolveYDBAuth(t *testing.T) {
	clearYDBAuthEnv := func(t *testing.T) {
		t.Helper()
		for _, key := range []string{
			ydbconsts.EnvToken,
			ydbconsts.EnvStaticCredentialsUser,
			ydbconsts.EnvStaticCredentialsPassword,
			ydbconsts.EnvAnonymousCredentials,
		} {
			t.Setenv(key, "")
		}
	}

	tests := []struct {
		name      string
		conf      map[string]string
		env       map[string]string
		wantKind  string
		wantValue string
	}{
		{
			name: "config token supported directly",
			conf: map[string]string{
				"token": "config-token",
			},
			wantKind:  "token",
			wantValue: "config-token",
		},
		{
			name: "vault env token overrides config",
			conf: map[string]string{
				"token": "config-token",
			},
			env: map[string]string{
				ydbconsts.EnvToken: "vault-token",
			},
			wantKind:  "token",
			wantValue: "vault-token",
		},
		{
			name: "config static credentials supported directly",
			conf: map[string]string{
				"static_credentials_user":     "user",
				"static_credentials_password": "password",
			},
			wantKind:  "static",
			wantValue: "user",
		},
		{
			name: "vault env static credentials override config",
			conf: map[string]string{
				"static_credentials_user":     "user",
				"static_credentials_password": "password",
			},
			env: map[string]string{
				ydbconsts.EnvStaticCredentialsUser:     "env-user",
				ydbconsts.EnvStaticCredentialsPassword: "env-password",
			},
			wantKind:  "static",
			wantValue: "env-user",
		},
		{
			name: "anonymous config supported directly",
			conf: map[string]string{
				"anonymous_credentials": "true",
			},
			wantKind: "anonymous",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clearYDBAuthEnv(t)
			for key, value := range tc.env {
				t.Setenv(key, value)
			}

			got := resolveYDBAuth(tc.conf)
			if got.kind != tc.wantKind {
				t.Fatalf("resolveYDBAuth kind = %q, want %q", got.kind, tc.wantKind)
			}
			if got.value != tc.wantValue {
				t.Fatalf("resolveYDBAuth value = %q, want %q", got.value, tc.wantValue)
			}
			if tc.wantKind == "static" && got.value2 == "" {
				t.Fatalf("resolveYDBAuth static password is empty")
			}
		})
	}
}

func TestGetYDBTransactionLimits(t *testing.T) {
	tests := []struct {
		name        string
		conf        map[string]string
		env         map[string]string
		wantEntries int
		wantSize    int
		expectErr   bool
	}{
		{
			name:        "defaults",
			wantEntries: defaultYDBTransactionMaxEntries,
			wantSize:    defaultYDBTransactionMaxSize,
		},
		{
			name: "config overrides",
			conf: map[string]string{
				"transaction_max_entries": "100",
				"transaction_max_size":    "262144",
			},
			wantEntries: 100,
			wantSize:    262144,
		},
		{
			name: "env overrides config",
			conf: map[string]string{
				"transaction_max_entries": "100",
				"transaction_max_size":    "262144",
			},
			env: map[string]string{
				ydbconsts.EnvTransactionMaxEntries: "50",
				ydbconsts.EnvTransactionMaxSize:    "131072",
			},
			wantEntries: 50,
			wantSize:    131072,
		},
		{
			name: "invalid entries",
			conf: map[string]string{
				"transaction_max_entries": "0",
			},
			expectErr: true,
		},
		{
			name: "invalid size",
			env: map[string]string{
				ydbconsts.EnvTransactionMaxSize: "abc",
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(ydbconsts.EnvTransactionMaxEntries, "")
			t.Setenv(ydbconsts.EnvTransactionMaxSize, "")
			for key, value := range tc.env {
				t.Setenv(key, value)
			}

			gotEntries, gotSize, err := getYDBTransactionLimits(tc.conf)
			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotEntries != tc.wantEntries || gotSize != tc.wantSize {
				t.Fatalf("getYDBTransactionLimits() = (%d, %d), want (%d, %d)", gotEntries, gotSize, tc.wantEntries, tc.wantSize)
			}
		})
	}
}
