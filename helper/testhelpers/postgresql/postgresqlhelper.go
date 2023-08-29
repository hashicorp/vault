// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/docker"
)

func PrepareTestContainer(t *testing.T, version string) (func(), string) {
	_, cleanup, url, _ := PrepareTestContainerRunner(t, version)

	return cleanup, url
}

func PrepareTestContainerRunner(t *testing.T, version string) (*docker.Runner, func(), string, string) {
	env := []string{
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=database",
	}

	return prepareTestContainer(t, "postgres", "docker.mirror.hashicorp.services/postgres", version, "secret", true, false, false, env)
}

// PrepareTestContainerWithVaultUser will setup a test container with a Vault
// admin user configured so that we can safely call rotate-root without
// rotating the root DB credentials
func PrepareTestContainerWithVaultUser(t *testing.T, ctx context.Context, version string) (func(), string) {
	env := []string{
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=database",
	}

	runner, cleanup, url, id := prepareTestContainer(t, "postgres", "docker.mirror.hashicorp.services/postgres", version, "secret", true, false, false, env)

	cmd := []string{"psql", "-U", "postgres", "-c", "CREATE USER vaultadmin WITH LOGIN PASSWORD 'vaultpass' SUPERUSER"}
	_, err := runner.RunCmdInBackground(ctx, id, cmd)
	if err != nil {
		t.Fatalf("Could not run command (%v) in container: %v", cmd, err)
	}

	return cleanup, url
}

func PrepareTestContainerWithPassword(t *testing.T, version, password string) (func(), string) {
	env := []string{
		"POSTGRES_PASSWORD=" + password,
		"POSTGRES_DB=database",
	}

	_, cleanup, url, _ := prepareTestContainer(t, "postgres", "docker.mirror.hashicorp.services/postgres", version, password, true, false, false, env)

	return cleanup, url
}

func PrepareTestContainerRepmgr(t *testing.T, name, version string, envVars []string) (*docker.Runner, func(), string, string) {
	env := append(envVars,
		"REPMGR_PARTNER_NODES=psql-repl-node-0,psql-repl-node-1",
		"REPMGR_PRIMARY_HOST=psql-repl-node-0",
		"REPMGR_PASSWORD=repmgrpass",
		"POSTGRESQL_PASSWORD=secret")

	return prepareTestContainer(t, name, "docker.mirror.hashicorp.services/bitnami/postgresql-repmgr", version, "secret", false, true, true, env)
}

func prepareTestContainer(t *testing.T, name, repo, version, password string,
	addSuffix, forceLocalAddr, doNotAutoRemove bool, envVars []string,
) (*docker.Runner, func(), string, string) {
	if os.Getenv("PG_URL") != "" {
		return nil, func() {}, "", os.Getenv("PG_URL")
	}

	if version == "" {
		version = "11"
	}

	runOpts := docker.RunOptions{
		ContainerName:   name,
		ImageRepo:       repo,
		ImageTag:        version,
		Env:             envVars,
		Ports:           []string{"5432/tcp"},
		DoNotAutoRemove: doNotAutoRemove,
	}
	if repo == "bitnami/postgresql-repmgr" {
		runOpts.NetworkID = os.Getenv("POSTGRES_MULTIHOST_NET")
	}

	runner, err := docker.NewServiceRunner(runOpts)
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	svc, containerID, err := runner.StartNewService(context.Background(), addSuffix, forceLocalAddr, connectPostgres(password, repo))
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	return runner, svc.Cleanup, svc.Config.URL().String(), containerID
}

func connectPostgres(password, repo string) docker.ServiceAdapter {
	return func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword("postgres", password),
			Host:     fmt.Sprintf("%s:%d", host, port),
			Path:     "postgres",
			RawQuery: "sslmode=disable",
		}

		db, err := sql.Open("pgx", u.String())
		if err != nil {
			return nil, err
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			return nil, err
		}
		return docker.NewServiceURL(u), nil
	}
}

func StopContainer(t *testing.T, ctx context.Context, runner *docker.Runner, containerID string) {
	if err := runner.Stop(ctx, containerID); err != nil {
		t.Fatalf("Could not stop docker Postgres: %s", err)
	}
}

func RestartContainer(t *testing.T, ctx context.Context, runner *docker.Runner, containerID string) {
	if err := runner.Restart(ctx, containerID); err != nil {
		t.Fatalf("Could not restart docker Postgres: %s", err)
	}
}
