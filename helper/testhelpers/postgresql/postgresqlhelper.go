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

func PrepareTestContainerWithSSL(t *testing.T, ctx context.Context, sslMode string) (func(), string) {
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
		certhelpers.CommonName("ca"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)
	serverCert := certhelpers.NewCert(t,
		certhelpers.CommonName("server"),
		certhelpers.DNS("localhost"),
		certhelpers.Parent(caCert),
	)
	clientCert := certhelpers.NewCert(t,
		certhelpers.CommonName("postgres"),
		certhelpers.DNS("localhost"),
		certhelpers.Parent(caCert),
	)

	bCtx := docker.NewBuildContext()
	bCtx["ca.crt"] = docker.PathContentsFromBytes(caCert.CombinedPEM())
	bCtx["server.crt"] = docker.PathContentsFromBytes(serverCert.CombinedPEM())
	bCtx["server.key"] = &docker.FileContents{
		Data: serverCert.PrivateKeyPEM(),
		Mode: 0o600,
		// postgres uid
		UID: 999,
	}

	// https://www.postgresql.org/docs/current/auth-pg-hba-conf.html
	clientAuthConfig := "echo 'hostssl all all all cert clientcert=verify-ca' > /var/lib/postgresql/data/pg_hba.conf"
	bCtx["ssl-conf.sh"] = docker.PathContentsFromString(clientAuthConfig)
	pgConfig := `
cat << EOF > /var/lib/postgresql/data/postgresql.conf
# PostgreSQL configuration file
listen_addresses = '*'
max_connections = 100
shared_buffers = 128MB
dynamic_shared_memory_type = posix
max_wal_size = 1GB
min_wal_size = 80MB
ssl = on
ssl_ca_file = '/var/lib/postgresql/ca.crt'
ssl_cert_file = '/var/lib/postgresql/server.crt'
ssl_key_file= '/var/lib/postgresql/server.key'
EOF
`
	bCtx["pg-conf.sh"] = docker.PathContentsFromString(pgConfig)

	err = runner.CopyTo(id, "/var/lib/postgresql/", bCtx)
	if err != nil {
		t.Fatalf("failed to copy to container: %v", err)
	}

	// overwrite the postgresql.conf config file with our ssl settings
	mustRunCommand(t, ctx, runner, id,
		[]string{"bash", "/var/lib/postgresql/pg-conf.sh"})

	// overwrite the pg_hba.conf file and set it to require SSL for each connection
	mustRunCommand(t, ctx, runner, id,
		[]string{"bash", "/var/lib/postgresql/ssl-conf.sh"})

	// reload so the config changes take effect and ssl is enabled
	mustRunCommand(t, ctx, runner, id,
		[]string{"psql", "-U", "postgres", "-c", "SELECT pg_reload_conf()"})

	if sslMode == "disable" {
		// return non-tls connection url
		return svc.Cleanup, svc.Config.URL().String()
	}

	sslConfig, err := connectPostgresSSL(
		t,
		svc.Config.URL().Host,
		sslMode,
		string(caCert.CombinedPEM()),
		string(clientCert.CombinedPEM()),
		string(clientCert.PrivateKeyPEM()),
	)
	if err != nil {
		svc.Cleanup()
		t.Fatalf("failed to connect to postgres container via SSL: %v", err)
	}
	return svc.Cleanup, sslConfig.URL().String()
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

func connectPostgresSSL(t *testing.T, host, sslMode, caCert, clientCert, clientKey string) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme: "postgres",
		User:   url.User("postgres"),
		Host:   host,
		Path:   "postgres",
		RawQuery: url.Values{
			"sslmode":     {sslMode},
			"sslinline":   {"true"},
			"sslrootcert": {caCert},
			"sslcert":     {clientCert},
			"sslkey":      {clientKey},
		}.Encode(),
	}

	db, err := connutil.OpenPostgres("pgx", u.String())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
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
	_, stderr, retcode, err := runner.RunCmdWithOutput(ctx, containerID, cmd)
	if err != nil {
		t.Fatalf("Could not run command (%v) in container: %v", cmd, err)
	}
	if retcode != 0 || len(stderr) != 0 {
		t.Fatalf("exit code: %v, stderr: %v", retcode, string(stderr))
	}
}
