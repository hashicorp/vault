// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/docker"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	defaultPGImage   = "docker.mirror.hashicorp.services/postgres"
	defaultPGVersion = "16-alpine"
	defaultPGPass    = "secret"
)

func defaultRunOpts(t *testing.T) docker.RunOptions {
	return docker.RunOptions{
		ContainerName: fmt.Sprintf("postgres-%s", sanitize(t.Name())),
		ImageRepo:     defaultPGImage,
		ImageTag:      defaultPGVersion,
		Env: []string{
			"POSTGRES_PASSWORD=" + defaultPGPass,
			"POSTGRES_DB=database",
		},
		Ports:             []string{"5432/tcp"},
		DoNotAutoRemove:   false,
		OmitLogTimestamps: true,
		LogConsumer: func(s string) {
			if t.Failed() {
				t.Logf("container logs: %s", s)
			}
		},
	}
}

func requireVaultEnv(t *testing.T) {
	t.Helper()

	if os.Getenv("VAULT_ADDR") == "" || os.Getenv("VAULT_TOKEN") == "" {
		t.Skip("skipping blackbox test: VAULT_ADDR and VAULT_TOKEN are required")
	}
}

func PrepareTestContainer(t *testing.T) (func(), string) {
	_, cleanup, connURL, _ := prepareTestContainer(t, defaultRunOpts(t), defaultPGPass, false, false, false)

	return cleanup, connURL
}

func PrepareTestContainerMultiHost(t *testing.T) (func(), string) {
	_, cleanup, connURL, _ := prepareTestContainer(t, defaultRunOpts(t), defaultPGPass, false, false, true)

	return cleanup, connURL
}

func prepareTestContainer(t *testing.T, runOpts docker.RunOptions, password string, addSuffix, forceLocalAddr, useFallback bool,
) (*docker.Runner, func(), string, string) {
	requireVaultEnv(t)

	if os.Getenv("PG_URL") != "" {
		envPGURL := os.Getenv("PG_URL")
		if useFallback {
			envPGURL = withFallbackHost(envPGURL)
		}
		return nil, func() {}, envPGURL, ""
	}

	runner, err := docker.NewServiceRunner(runOpts)
	if err != nil {
		if strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
			t.Fatalf("skipping blackbox test: docker daemon not available: %v", err)
		}
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	svc, containerID, err := runner.StartNewService(context.Background(), addSuffix, forceLocalAddr, connectPostgres(password, useFallback))
	if err != nil {
		if strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
			t.Fatalf("skipping blackbox test: docker daemon not available: %v", err)
		}
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	return runner, svc.Cleanup, svc.Config.URL().String(), containerID
}

func withFallbackHost(connURL string) string {
	u, err := url.Parse(connURL)
	if err != nil || u.Host == "" {
		return connURL
	}

	if strings.Contains(u.Host, "localhost:55,") {
		return connURL
	}

	u.Host = "localhost:55," + u.Host
	return u.String()
}

func sanitize(name string) string {
	lower := strings.ToLower(name)
	var b strings.Builder
	b.Grow(len(lower))

	lastDash := false
	for _, r := range lower {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}

		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}

	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "test"
	}

	if len(out) > 54 {
		const hashLen = 8
		sum := sha256.Sum256([]byte(out))
		hash := hex.EncodeToString(sum[:])[:hashLen]
		prefixLen := 54 - 1 - hashLen
		out = out[:prefixLen] + "-" + hash
	}

	return out
}

func connectPostgres(password string, useFallback bool) docker.ServiceAdapter {
	return func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		hostAddr := fmt.Sprintf("%s:%d", host, port)
		if useFallback {
			// Use an unreachable first host to exercise driver fallback behavior.
			hostAddr = "localhost:55," + hostAddr
		}

		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword("postgres", password),
			Host:     hostAddr,
			Path:     "postgres",
			RawQuery: "sslmode=disable",
		}

		db, err := sql.Open("pgx", u.String())
		if err != nil {
			return nil, err
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			db.Close() // Explicit close on error
			return nil, err
		}

		return docker.NewServiceURL(u), nil
	}
}
