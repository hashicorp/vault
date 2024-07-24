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

	"github.com/hashicorp/vault/helper/testhelpers/certhelpers"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/helper/docker"
)

const defaultPostgresImage = "docker.mirror.hashicorp.services/postgres"

const defaultPostgresVersion = "13.4-buster"

const defaultPostgresPassword = "secret"

func defaultRunOpts(t *testing.T) docker.RunOptions {
	return docker.RunOptions{
		ContainerName: "postgres",
		ImageRepo:     defaultPostgresImage,
		ImageTag:      defaultPostgresVersion,
		Env: []string{
			"POSTGRES_PASSWORD=" + defaultPostgresPassword,
			"POSTGRES_DB=database",
		},
		Ports:           []string{"5432/tcp"},
		DoNotAutoRemove: false,
		LogConsumer: func(s string) {
			if t.Failed() {
				t.Logf("container logs: %s", s)
			}
		},
	}
}

func PrepareTestContainer(t *testing.T) (func(), string) {
	_, cleanup, url, _ := prepareTestContainer(t, defaultRunOpts(t), defaultPostgresPassword, true, false)

	return cleanup, url
}

// PrepareTestContainerWithVaultUser will setup a test container with a Vault
// admin user configured so that we can safely call rotate-root without
// rotating the root DB credentials
func PrepareTestContainerWithVaultUser(t *testing.T, ctx context.Context) (func(), string) {
	runner, cleanup, url, id := prepareTestContainer(t, defaultRunOpts(t), defaultPostgresPassword, true, false)

	cmd := []string{"psql", "-U", "postgres", "-c", "CREATE USER vaultadmin WITH LOGIN PASSWORD 'vaultpass' SUPERUSER"}
	mustRunCommand(t, ctx, runner, id, cmd)

	return cleanup, url
}

func PrepareTestContainerWithSSL(t *testing.T, ctx context.Context) (func(), string, certhelpers.Certificate) {
	runOpts := defaultRunOpts(t)
	runner, err := docker.NewServiceRunner(runOpts)
	if err != nil {
		t.Fatalf("Could not provision docker service runner: %s", err)
	}

	// first we connect with username/password because ssl is not enabled yet
	svc, id, err := runner.StartNewService(context.Background(), false, false, connectPostgres(defaultPostgresPassword, runOpts.ImageRepo))
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	// Create certificates for postgres authentication
	caCert := certhelpers.NewCert(t,
		certhelpers.CommonName("test certificate authority"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)
	serverCert := certhelpers.NewCert(t,
		certhelpers.CommonName("postgres"),
		certhelpers.DNS("localhost"),
		certhelpers.Parent(caCert),
	)
	clientCert := certhelpers.NewCert(t,
		certhelpers.CommonName("client"),
		certhelpers.DNS("localhost"),
		certhelpers.Parent(caCert),
	)

	bCtx := docker.NewBuildContext()
	bCtx["ca.crt"] = docker.PathContentsFromBytes(caCert.CombinedPEM())
	bCtx["server.crt"] = docker.PathContentsFromBytes(serverCert.CombinedPEM())
	bCtx["server.key"] = docker.PathContentsFromBytes(serverCert.PrivateKeyPEM())
	// https://www.postgresql.org/docs/current/auth-pg-hba-conf.html
	clientAuthConfig := "echo 'hostssl all all all cert clientcert=verify-ca' > /var/lib/postgresql/data/pg_hba.conf"
	bCtx["ssl-conf.sh"] = docker.PathContentsFromString(clientAuthConfig)

	err = runner.CopyTo(id, "/var/lib/postgresql/", bCtx)
	if err != nil {
		t.Fatalf("failed to copy to container: %v", err)
	}

	// run the ssl init script to overwrite the pg_hba.conf file and set it to
	// require SSL for each connection
	mustRunCommand(t, ctx, runner, id,
		[]string{"bash", "/var/lib/postgresql/ssl-conf.sh"})

	// reload so the config changes take effect and ssl is enabled
	mustRunCommand(t, ctx, runner, id,
		[]string{"psql", "-U", "postgres", "-c", "SELECT pg_reload_conf()"})
	mustRunCommand(t, ctx, runner, id,
		[]string{"cat", "/var/lib/postgresql/data/pg_hba.conf"})

	svcConfig, err := connectPostgresSSL(t, svc.Config.URL(), string(caCert.CombinedPEM()), string(clientCert.CombinedPEM()), string(clientCert.PrivateKeyPEM()))
	if err != nil {
		// svc.Cleanup()
		t.Fatalf("failed to connect to postgres container via SSL: %v", err)
	}
	return svc.Cleanup, svcConfig.URL().String(), clientCert
}

func PrepareTestContainerWithPassword(t *testing.T, password string) (func(), string) {
	runOpts := defaultRunOpts(t)
	runOpts.Env = []string{
		"POSTGRES_PASSWORD=" + password,
		"POSTGRES_DB=database",
	}

	_, cleanup, url, _ := prepareTestContainer(t, runOpts, password, true, false)

	return cleanup, url
}

func prepareTestContainer(t *testing.T, runOpts docker.RunOptions, password string, addSuffix, forceLocalAddr bool,
) (*docker.Runner, func(), string, string) {
	if os.Getenv("PG_URL") != "" {
		return nil, func() {}, "", os.Getenv("PG_URL")
	}

	runner, err := docker.NewServiceRunner(runOpts)
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	svc, containerID, err := runner.StartNewService(context.Background(), addSuffix, forceLocalAddr, connectPostgres(password, runOpts.ImageRepo))
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	return runner, svc.Cleanup, svc.Config.URL().String(), containerID
}

func connectPostgresSSL(t *testing.T, connURL *url.URL, caCert, clientCert, clientKey string) (docker.ServiceConfig, error) {
	caCert = "foo"
	u := url.URL{
		Scheme:   "postgres",
		User:     url.User("client"),
		Host:     connURL.Host,
		Path:     "postgres",
		RawQuery: url.QueryEscape("sslmode=verify-full&sslinline=true&sslrootcert=" + caCert + "&sslcert=" + clientCert + "&sslkey=" + clientKey),
	}

	t.Logf("\nurl: %s\n", u.String())
	db, err := connutil.OpenPostgres("pgx", u.String())
	if err != nil {
		t.Fatalf("open err %s", err)
		return nil, err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatalf("ping err %s", err)
		return nil, err
	}
	return docker.NewServiceURL(u), nil
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

func mustRunCommand(t *testing.T, ctx context.Context, runner *docker.Runner, containerID string, cmd []string) {
	t.Helper()
	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, containerID, cmd)
	if err != nil {
		t.Fatalf("Could not run command (%v) in container: %v", cmd, err)
	}
	if retcode != 0 || len(stderr) != 0 {
		t.Fatalf("exit code: %v, stderr: %v", retcode, string(stderr))
	}
	t.Logf("stdout: %v", string(stdout))
}
